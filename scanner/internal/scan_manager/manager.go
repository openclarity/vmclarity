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
	"math/rand"
	"runtime"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/sourcegraph/conc/pool"

	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
)

type InputScanResult[T any] struct {
	Metadata   families.ScanInputMetadata
	ScanInput  common.ScanInput
	ScanResult T
}

type InputScanResultWithError[T any] struct {
	InputScanResult[T]
	Error error
}

// Manager allows parallelized scan of inputs for a single family scanner factory.
type Manager[CT, RT any] struct {
	config   CT
	scanners []string
	factory  *Factory[CT, RT]
}

func New[CT, RT any](scanners []string, config CT, factory *Factory[CT, RT]) *Manager[CT, RT] {
	return &Manager[CT, RT]{
		config:   config,
		scanners: scanners,
		factory:  factory,
	}
}

func (m *Manager[CT, RT]) Scan(ctx context.Context, inputs []common.ScanInput) ([]InputScanResult[RT], error) {
	logger := log.GetLoggerFromContextOrDefault(ctx)

	// Create processing jobs, do not cancel on error
	resultCh := make(chan InputScanResultWithError[RT])
	workerPool := pool.New().WithContext(ctx).WithMaxGoroutines(runtime.NumCPU())

	for _, scannerName := range m.scanners {
		// Do not continue processing further if we cannot create a specific scanner
		scanner, err := m.factory.createScanner(scannerName, m.config)
		if err != nil {
			logger.WithError(err).Errorf("Failed to create scanner %s", scannerName)
			continue
		}

		// Schedule each ({scanner}, {input}) pair to parallel worker
		for _, input := range inputs {
			workerPool.Go(func(ctx context.Context) error {
				// Set logger for a given family scanner and input pair
				logger := logger.WithFields(map[string]interface{}{
					"scanner": scannerName,
					"input":   input,
				})
				ctx = log.SetLoggerForContext(ctx, logger)

				// Fuzzy start processing to prevent spike requests for each input
				time.Sleep(time.Duration(rand.Int63n(int64(20 * time.Millisecond)))) // nolint:mnd,gosec,wrapcheck

				// Run scan
				logger.Infof("Scanning input = '%s' started...", input)

				startTime := time.Now()
				inputScanResult, inputScanErr := scanner.Scan(ctx, input.InputType, input.Input)
				inputSize, _ := familiesutils.GetInputSize(input) // in megabytes
				endTime := time.Now()

				// Log scan result details
				if inputScanErr != nil {
					logger.WithError(inputScanErr).Warnf("Scanning input = '%s' finished with error", input)
				} else {
					logger.Infof("Scanning input = '%s' finished successfully", input)
				}

				// Forward the result in custom format with error to the main result channel so
				// that we can handle the results in the main loop
				resultCh <- InputScanResultWithError[RT]{
					InputScanResult: InputScanResult[RT]{
						Metadata: families.NewScanInputMetadata(
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

	// Read from the main channel and handle all the forwarded results and errors
	var resultErrs error
	var results []InputScanResult[RT]

	for result := range resultCh {
		// NOTE(ramizpolic): We don't check for nil results. A scan is considered
		// successful if the family scanner returned a result with nil error, unrelated
		// of the actual result.
		if err := result.Error; err != nil {
			scanErr := fmt.Errorf("%s scanner job failed: %w", result.Metadata, err)
			logger.Warning(scanErr)

			resultErrs = multierror.Append(resultErrs, scanErr)
		} else {
			logger.Infof("Got result for scanner job %s", result.Metadata)

			results = append(results, result.InputScanResult)
		}
	}

	// Return error if all jobs failed to return results.
	// TODO: should it be configurable? allow the user to decide failure threshold?
	if len(results) == 0 {
		return nil, resultErrs // nolint:wrapcheck
	}

	return results, nil
}
