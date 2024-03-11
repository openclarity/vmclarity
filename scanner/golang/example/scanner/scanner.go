// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package scanner

import (
	"context"
	"fmt"
	dockle_run "github.com/Portshift/dockle/pkg"
	"github.com/openclarity/vmclarity/scanner/types"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type scanner interface {
	GetScanResult() (*types.ScanResult, error)
	StartScan(template types.ScanTemplate) (*types.Scan, error)
	GetScan() (*types.Scan, error)
	StopScan() error
}

type Scanner struct {
	mu        sync.RWMutex
	result    *types.ScanResult
	resultErr error
	cancel    func()
}

func (s *Scanner) GetScanResult() (*types.ScanResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.result == nil {
		return nil, fmt.Errorf("scan not started")
	}
	if s.resultErr != nil {
		return nil, s.resultErr
	}
	return s.result, nil
}

func (s *Scanner) StartScan(template types.ScanTemplate) (*types.Scan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.result != nil {
		return nil, fmt.Errorf("scan already running")
	}

	// create scan
	s.result = &types.ScanResult{
		Findings: Ptr([]types.Finding{}),
		Scan: &types.Scan{
			JobsCompleted: Ptr(0),
			JobsLeftToRun: Ptr(len(template.AssetScanInputs)),
			StartTime:     Ptr(time.Now()),
			Template:      &template,
		},
	}

	// start scan
	go s.scan(template.AssetScanInputs)

	// return
	return s.result.Scan, nil
}

func (s *Scanner) GetScan() (*types.Scan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.result == nil || s.result.Scan == nil {
		return nil, fmt.Errorf("scan not started")
	}
	return s.result.Scan, nil
}

func (s *Scanner) StopScan() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.result == nil {
		return fmt.Errorf("scan not started")
	}
	s.cancel()
	return nil
}

func (s *Scanner) scan(inputs []types.AssetScanInput) {
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()
	procGroup, _ := errgroup.WithContext(ctx)

	for _, input := range inputs {
		input := input
		procGroup.Go(func() error {
			// Validate this is an input type supported by the scanner,
			// otherwise return skipped.
			if err := s.isValidInputType(input.Type); err != nil {
				return err
			}

			assessmentMap, err := dockle_run.RunFromConfig(createDockleConfig(input.Type, input.Path))
			if err != nil {
				return err
			}
			findings := parseDockleReport(input, assessmentMap)

			s.mu.Lock()
			*s.result.Findings = append(*s.result.Findings, findings...)
			s.mu.Unlock()

			return nil
		})
	}

	go func() {
		err := procGroup.Wait()
		if err != nil {
			s.mu.Lock()
			s.resultErr = err
			s.mu.Unlock()
		}
	}()
}

func (s *Scanner) isValidInputType(sourceType types.AssetScanInputType) error {
	switch sourceType {
	case types.AssetScanInputImage, types.AssetScanInputDockerArchive:
		return nil
	default:
		return fmt.Errorf("unsupported input type %v", sourceType)
	}
}
