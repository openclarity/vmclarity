// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job_manager // nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/openclarity/vmclarity/scanner/families/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"time"
)

type ScanResult[RT types.Result[RT]] struct {
	Metadata scannertypes.ScanInputMetadata
	Input    scannertypes.ScanInput
	Result   RT
	Error    error
}

type Manager[CT any, RT types.Result[RT]] struct {
	jobNames   []string
	config     CT
	logger     *logrus.Entry
	jobFactory *Factory[CT, RT]
}

func New[CT any, RT types.Result[RT]](jobNames []string, config CT, logger *logrus.Entry, factory *Factory[CT, RT]) *Manager[CT, RT] {
	return &Manager[CT, RT]{
		jobNames:   jobNames,
		config:     config,
		logger:     logger,
		jobFactory: factory,
	}
}

func (m *Manager[CT, RT]) Process(ctx context.Context, inputs []scannertypes.ScanInput) ([]ScanResult[RT], error) {
	mainResultCh := make(chan ScanResult[RT])
	workerPool := pool.New().WithContext(ctx).WithFirstError().WithCancelOnError()

	// Create processing jobs
	for _, jobName := range m.jobNames {
		job, err := m.jobFactory.CreateJob(jobName, m.config, m.logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create job %s: %w", jobName, err)
		}

		// schedule each {job}, {input} input pair to parallel worker
		for _, input := range inputs {
			workerPool.Go(func(ctx context.Context) error {
				m.logger.Infof("Started running job %s for input %s...", jobName, input.Input)

				// Process
				startTime := time.Now()
				scanResult, scanErr := job.Scan(ctx, input.InputType, input.Input)
				inputSize, _ := familiesutils.GetInputSize(input) // in megabytes

				// Wait for job to finish by waiting for the result. Once done, forward the
				// result formatted to main result channel
				mainResultCh <- ScanResult[RT]{
					Metadata: scannertypes.NewScanInputMetadata(
						jobName,
						startTime,
						time.Now(),
						inputSize,
						input,
					),
					Input:  input,
					Result: scanResult,
					Error:  scanErr,
				}

				m.logger.Infof("Finished running job %s for input %s", jobName, input.Input)

				return nil
			})
		}
	}

	// Wait for workers to finish and close main result channel to allow proper listening.
	// Write wait error to a separate worker channel to handle it at the end.
	workerErrCh := make(chan error, 1)
	go func() {
		workerErrCh <- workerPool.Wait()
		close(mainResultCh)
	}()

	// Read results from the main channel
	var resultError error
	var totalSuccessfulResultsCount int
	results := make([]ScanResult[RT], 0, len(m.jobNames))
	for processResult := range mainResultCh {
		if err := processResult.Error; err != nil {
			errStr := fmt.Errorf("%q scanner job failed: %w", processResult.Metadata, err)
			m.logger.Warning(errStr)
			resultError = multierror.Append(resultError, errStr)
		} else {
			m.logger.Infof("Got result for scanner job %q", processResult.Metadata)
			results = append(results, processResult)
			totalSuccessfulResultsCount++
		}
	}

	// Return error if all jobs failed to return results.
	// TODO: should it be configurable? allow the user to decide failure threshold?
	if totalSuccessfulResultsCount == 0 {
		return nil, resultError // nolint:wrapcheck
	}

	// Check if any of the workers failed
	if err := <-workerErrCh; err != nil {
		return nil, fmt.Errorf("failed to process inputs: %w", err)
	}

	return results, nil
}
