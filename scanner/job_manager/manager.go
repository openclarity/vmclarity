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
	"github.com/openclarity/vmclarity/scanner/utils"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"time"
)

type Manager struct {
	jobNames   []string
	config     IsConfig
	logger     *logrus.Entry
	jobFactory *Factory
}

func New(jobNames []string, config IsConfig, logger *logrus.Entry, factory *Factory) *Manager {
	return &Manager{
		jobNames:   jobNames,
		config:     config,
		logger:     logger,
		jobFactory: factory,
	}
}

type ProcessResult struct {
	types.InputScanMetadata
	Result Result
	Input  types.Input
}

func (m *Manager) Process(ctx context.Context, inputs []types.Input) ([]ProcessResult, error) {
	mainResultCh := make(chan ProcessResult)
	workerPool := pool.New().WithContext(ctx).WithFirstError().WithCancelOnError()

	// Create processing jobs
	for _, jobName := range m.jobNames {
		jobResultCh := make(chan Result)
		job, err := m.jobFactory.CreateJob(jobName, m.config, m.logger, jobResultCh)
		if err != nil {
			return nil, fmt.Errorf("failed to create job %s: %w", jobName, err)
		}

		// schedule each {job}, {input} input pair to parallel worker
		for _, input := range inputs {
			workerPool.Go(func(ctx context.Context) error {
				// Process
				startTime := time.Now()
				err := job.Run(ctx, utils.SourceType(input.InputType), input.Input)
				if err != nil {
					return fmt.Errorf("failed to run job %s for input %s: %w", jobName, input.Input, err)
				}
				endTime := time.Now()
				inputSize, _ := familiesutils.GetInputSize(input)

				// Read result from job channel and write to main result channel
				jobResult := <-jobResultCh
				mainResultCh <- ProcessResult{
					InputScanMetadata: types.CreateInputScanMetadata(
						jobName,
						startTime,
						endTime,
						inputSize,
						input,
					),
					Input:  input,
					Result: jobResult,
				}

				return nil
			})
		}
	}

	// Wait for workers to finish and close main result channel to allow proper listening.
	// Write wait error to a separate worker channel to handle it at the end.
	workerErrCh := make(chan error)
	go func() {
		workerErrCh <- workerPool.Wait()
		close(mainResultCh)
	}()

	// Read results from the main channel
	var resultError error
	var totalSuccessfulResultsCount int
	results := make([]ProcessResult, 0, len(m.jobNames))
	for processResult := range mainResultCh {
		if err := processResult.Result.GetError(); err != nil {
			errStr := fmt.Errorf("%q scanner job failed: %w", processResult.ScannerName, err)
			m.logger.Warning(errStr)
			resultError = multierror.Append(resultError, errStr)
		} else {
			m.logger.Infof("Got result for scanner job %q", processResult.ScannerName)
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
