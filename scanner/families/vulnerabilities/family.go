// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
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

package vulnerabilities

import (
	"context"
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/vulnerabilities/types"
	"github.com/openclarity/vmclarity/scanner/internal/job_manager"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
	"os"

	"github.com/openclarity/vmclarity/core/log"
	sbomtypes "github.com/openclarity/vmclarity/scanner/families/sbom/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
)

const (
	sbomTempFilePath = "/tmp/sbom"
)

type Vulnerabilities struct {
	conf types.Config
}

func New(conf types.Config) familiestypes.Family[*types.Vulnerabilities] {
	return &Vulnerabilities{
		conf: conf,
	}
}

func (v Vulnerabilities) GetType() familiestypes.FamilyType {
	return familiestypes.Vulnerabilities
}

func (v Vulnerabilities) Run(ctx context.Context, res *familiestypes.Results) (*types.Vulnerabilities, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "vulnerabilities")
	logger.Info("Vulnerabilities Run...")

	if v.conf.InputFromSbom {
		logger.Infof("Using input from SBOM results")

		sbomResults, err := familiestypes.GetFamilyResult[*sbomtypes.FamilyResult](res)
		if err != nil {
			return nil, fmt.Errorf("failed to get sbom results: %w", err)
		}

		sbomBytes, err := sbomResults.EncodeToBytes("cyclonedx-json")
		if err != nil {
			return nil, fmt.Errorf("failed to encode sbom results to bytes: %w", err)
		}

		// TODO: need to avoid writing sbom to file
		if err := os.WriteFile(sbomTempFilePath, sbomBytes, 0o600 /* read & write */); err != nil { // nolint:mnd,gofumpt
			return nil, fmt.Errorf("failed to write sbom to file: %w", err)
		}

		v.conf.Inputs = append(v.conf.Inputs, scannertypes.ScanInput{
			Input:     sbomTempFilePath,
			InputType: scannertypes.SBOM,
		})
	}

	if len(v.conf.Inputs) == 0 {
		return nil, errors.New("inputs list is empty")
	}

	manager := job_manager.New[types.ScannersConfig, *types.ScannerResult](v.conf.ScannersList, v.conf.ScannersConfig, logger, types.Factory)
	processResults, err := manager.Process(ctx, v.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for vulnerabilities: %w", err)
	}

	vulResults := types.NewVulnerabilities()

	// Merge results.
	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.Result.ScannerInfo)
		vulResults.Merge(result.Result)
	}

	// TODO:
	// // Set source values.
	// mergedResults.SetSource(sharedscanner.Source{
	//	Type: "image",
	//	Name: config.ImageIDToScan,
	//	Hash: config.ImageHashToScan,
	// })

	logger.Info("Vulnerabilities Done...")

	return vulResults, nil
}
