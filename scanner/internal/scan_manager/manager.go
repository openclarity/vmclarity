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

package scan_manager // nolint:revive,stylecheck

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"time"
)

type InputScanResult[T any] struct {
	Metadata   scannertypes.ScanInputMetadata
	ScanInput  scannertypes.ScanInput
	ScanResult T
}

type InputScanResultWithError[T any] struct {
	InputScanResult[T]
	Error error
}

type Manager[CT, RT any] struct {
	config   CT
	logger   *logrus.Entry
	scanners []string
	factory  *Factory[CT, RT]
}

func New[CT, RT any](scanners []string, config CT, logger *logrus.Entry, factory *Factory[CT, RT]) *Manager[CT, RT] {
	return &Manager[CT, RT]{
		config:   config,
		logger:   logger,
		scanners: scanners,
		factory:  factory,
	}
}

func (m *Manager[CT, RT]) Scan(ctx context.Context, inputs []scannertypes.ScanInput) ([]InputScanResult[RT], error) {
	resultCh := make(chan InputScanResultWithError[RT])

	// Create processing jobs, do not cancel on error
	workerPool := pool.New().WithContext(ctx)

	for _, scannerName := range m.scanners {
		// Do not continue processing further if we cannot create a specific scanner
		scanner, err := m.factory.CreateJob(scannerName, m.config, m.logger)
		if err != nil {
			m.logger.Errorf("Failed to create scanner %s, reason=%v", scannerName, err)
			continue
		}

		// schedule each {job}, {input} input pair to parallel worker
		for _, input := range inputs {
			workerPool.Go(func(ctx context.Context) error {
				m.logger.Infof("Started running scanner %s for input %s:%s...", scannerName, input.InputType, input.Input)

				// Run scan
				startTime := time.Now()
				inputScanResult, inputScanErr := scanner.Scan(ctx, input.InputType, input.Input)
				inputSize, _ := familiesutils.GetInputSize(input) // in megabytes
				endTime := time.Now()

				// Skip doing anything in case the scanner returned nil result and nil error
				if inputScanResult == nil && inputScanErr == nil {
					return nil
				}

				// Forward the result in custom format to main result channel
				resultCh <- InputScanResultWithError[RT]{
					InputScanResult: InputScanResult[RT]{
						Metadata: scannertypes.NewScanInputMetadata(
							scannerName,
							startTime,
							endTime,
							inputSize,
							input,
						),
						ScanInput:  input,
						ScanResult: inputScanResult,
					},
					Error: inputScanErr,
				}

				state := "SUCCESS"
				if inputScanErr != nil {
					state = "FAILURE"
				}

				m.logger.Infof("Finished running scanner %s for input %s with state=%s", scannerName, input.Input, state)

				return nil
			})
		}
	}

	// Wait for workers to finish and close main result channel to allow proper
	// listening. We don't return any errors from the processing loop.
	go func() {
		_ = workerPool.Wait()
		close(resultCh)
	}()

	// Read results from the main channel
	var resultError error
	var results []InputScanResult[RT]

	for result := range resultCh {
		if err := result.Error; err != nil {
			scanErr := fmt.Errorf("%q scanner job failed: %w", result.Metadata, err)
			m.logger.Warning(scanErr)

			resultError = multierror.Append(resultError, scanErr)
		} else {
			m.logger.Infof("Got result for scanner job %q", result.Metadata)

			results = append(results, result.InputScanResult)
		}
	}

	// Return error if all jobs failed to return results.
	// TODO: should it be configurable? allow the user to decide failure threshold?
	if len(results) == 0 {
		return nil, resultError // nolint:wrapcheck
	}

	return results, nil
}
