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

package cisdocker

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/cisdocker/config"
	job_manager2 "github.com/openclarity/vmclarity/scanner/internal/job_manager"
	types2 "github.com/openclarity/vmclarity/scanner/types"

	dockle_run "github.com/Portshift/dockle/pkg"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
)

const (
	ScannerName = "cisdocker"
)

type Scanner struct {
	name       string
	logger     *logrus.Entry
	config     config.Config
	resultChan chan job_manager2.Result
}

func New(_ string, c job_manager2.IsConfig, logger *logrus.Entry, resultChan chan job_manager2.Result) job_manager2.Job {
	conf := c.(types.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("scanner", ScannerName),
		config:     conf.CISDocker,
		resultChan: resultChan,
	}
}

func (a *Scanner) Run(ctx context.Context, sourceType types2.InputType, userInput string) error {
	go func(ctx context.Context) {
		retResults := types.ScannerResult{
			ScannerName: ScannerName,
		}

		// Validate this is an input type supported by the scanner,
		// otherwise return skipped.
		if !a.isValidInputType(sourceType) {
			a.sendResults(retResults, nil)
			return
		}

		a.logger.Infof("Running %s scan...", a.name)
		config := createDockleConfig(a.logger, sourceType, userInput, a.config)
		ctx, cancel := context.WithTimeout(ctx, config.Timeout)
		defer cancel()

		assessmentMap, err := dockle_run.RunWithContext(ctx, config)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to run %s scan: %w", a.name, err))
			return
		}

		a.logger.Infof("Successfully scanned %s %s", sourceType, userInput)

		retResults.Misconfigurations = parseDockleReport(sourceType, userInput, assessmentMap)

		a.sendResults(retResults, nil)
	}(ctx)

	return nil
}

func (a *Scanner) isValidInputType(sourceType types2.InputType) bool {
	switch sourceType {
	case types2.IMAGE, types2.DOCKERARCHIVE, types2.ROOTFS, types2.DIR:
		return true
	case types2.FILE, types2.SBOM, types2.OCIARCHIVE, types2.OCIDIR:
		a.logger.Infof("source type %v is not supported for CIS Docker Benchmark scanner, skipping.", sourceType)
	default:
		a.logger.Infof("unknown source type %v, skipping.", sourceType)
	}
	return false
}

func (a *Scanner) sendResults(results types.ScannerResult, err error) {
	if err != nil {
		a.logger.Error(err)
		results.Error = err
	}
	select {
	case a.resultChan <- results:
	default:
		a.logger.Error("Failed to send results on channel")
	}
}
