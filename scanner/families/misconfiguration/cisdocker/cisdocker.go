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
	types2 "github.com/openclarity/vmclarity/scanner/families/plugins/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"

	dockle_run "github.com/Portshift/dockle/pkg"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	log "github.com/sirupsen/logrus"
)

const ScannerName = "cisdocker"

func init() {
	types2.FactoryRegister(ScannerName, New)
}

type Scanner struct {
	logger *logrus.Entry
	config config.Config
}

func New(_ string, config types.ScannersConfig, logger *log.Entry) familiestypes.Scanner[*types.ScannerResult] {
	return &Scanner{
		logger: logger.Dup().WithField("scanner", ScannerName),
		config: config.CISDocker,
	}
}

func (a *Scanner) Scan(ctx context.Context, sourceType scannertypes.InputType, userInput string) (*types.ScannerResult, error) {
	// Validate this is an input type supported by the scanner,
	// otherwise return skipped.
	if !a.isValidInputType(sourceType) {
		return nil, fmt.Errorf("unsupported source type for %s: %s", ScannerName, sourceType)
	}

	a.logger.Infof("Running %s scan on %s...", ScannerName, userInput)

	dockleCfg := createDockleConfig(a.logger, sourceType, userInput, a.config)
	ctx, cancel := context.WithTimeout(ctx, dockleCfg.Timeout)
	defer cancel()

	assessmentMap, err := dockle_run.RunWithContext(ctx, dockleCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to run dockle: %w", err)
	}

	a.logger.Infof("Successfully scanned %s %s", sourceType, userInput)

	misconfigurations := parseDockleReport(sourceType, userInput, assessmentMap)

	return &types.ScannerResult{
		ScannerName:       ScannerName,
		Misconfigurations: misconfigurations,
	}, nil
}

func (a *Scanner) isValidInputType(sourceType scannertypes.InputType) bool {
	switch sourceType {
	case scannertypes.IMAGE, scannertypes.DOCKERARCHIVE, scannertypes.ROOTFS, scannertypes.DIR:
		return true
	case scannertypes.FILE, scannertypes.SBOM, scannertypes.OCIARCHIVE, scannertypes.OCIDIR:
		a.logger.Infof("source type %v is not supported for CIS Docker Benchmark scanner, skipping.", sourceType)
	default:
		a.logger.Infof("unknown source type %v, skipping.", sourceType)
	}
	return false
}
