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

package windows

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/sbom/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"

	log "github.com/sirupsen/logrus"
)

const AnalyzerName = "windows"

func init() {
	types.FactoryRegister(AnalyzerName, New)
}

type Analyzer struct {
	logger *log.Entry
}

func New(_ string, _ types.AnalyzersConfig, logger *log.Entry) familiestypes.Scanner[*types.ScannerResult] {
	return &Analyzer{
		logger: logger.Dup().WithField("analyzer", AnalyzerName),
	}
}

// nolint:cyclop
func (a *Analyzer) Scan(ctx context.Context, sourceType scannertypes.InputType, userInput string) (*types.ScannerResult, error) {
	a.logger.Infof("Called %s analyzer on source %v %v", AnalyzerName, sourceType, userInput)

	// Create Windows registry based on supported input types
	var err error
	var registry *Registry
	switch sourceType {
	case scannertypes.FILE: // Use file location to the registry
		registry, err = NewRegistry(userInput, a.logger)
	case scannertypes.ROOTFS, scannertypes.DIR: // Use mount drive as input
		registry, err = NewRegistryForMount(userInput, a.logger)
	case scannertypes.SBOM, scannertypes.IMAGE, scannertypes.DOCKERARCHIVE, scannertypes.OCIARCHIVE, scannertypes.OCIDIR: // Unsupported
		fallthrough
	default:
		return nil, fmt.Errorf("skipping analyzing unsupported source type: %s", sourceType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open registry: %w", err)
	}
	defer registry.Close()

	// Fetch BOM from registry details
	bom, err := registry.GetBOM()
	if err != nil {
		return nil, fmt.Errorf("failed to get bom from registry: %w", err)
	}

	// Return sbom
	result := types.CreateScannerResult(bom, AnalyzerName, userInput, sourceType)

	a.logger.Infof("Sending successful results")

	return result, nil
}
