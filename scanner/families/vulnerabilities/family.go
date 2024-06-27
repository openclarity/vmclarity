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
	"github.com/openclarity/vmclarity/scanner/scanner/job"
	"os"

	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/config"
	"github.com/openclarity/vmclarity/scanner/families/interfaces"
	"github.com/openclarity/vmclarity/scanner/families/results"
	"github.com/openclarity/vmclarity/scanner/families/sbom"
	"github.com/openclarity/vmclarity/scanner/families/types"
	"github.com/openclarity/vmclarity/scanner/job_manager"
	"github.com/openclarity/vmclarity/scanner/scanner"
)

const (
	sbomTempFilePath = "/tmp/sbom"
)

type Vulnerabilities struct {
	conf           Config
	ScannersConfig config.Config
}

func (v Vulnerabilities) Run(ctx context.Context, res *results.Results) (interfaces.IsResults, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "vulnerabilities")
	logger.Info("Vulnerabilities Run...")

	if v.conf.InputFromSbom {
		logger.Infof("Using input from SBOM results")

		sbomResults, err := results.GetResult[*sbom.Results](res)
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

		v.conf.Inputs = append(v.conf.Inputs, types.Input{
			Input:     sbomTempFilePath,
			InputType: "sbom",
		})
	}

	if len(v.conf.Inputs) == 0 {
		return nil, errors.New("inputs list is empty")
	}

	manager := job_manager.New(v.conf.ScannersList, v.conf.ScannersConfig, logger, job.Factory)
	processResults, err := manager.Process(ctx, v.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for vulnerabilities: %w", err)
	}

	mergedResults := scanner.NewMergedResults()

	var vulResults Results

	// Merge results.
	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.ScannerName)
		data, ok := result.Result.(*scanner.Results)
		if !ok {
			return nil, fmt.Errorf("received results of a wrong type: %T", result)
		}
		mergedResults = mergedResults.Merge(data)
		vulResults.Metadata.InputScans = append(vulResults.Metadata.InputScans, result.InputScanMetadata)
	}

	// TODO:
	// // Set source values.
	// mergedResults.SetSource(sharedscanner.Source{
	//	Type: "image",
	//	Name: config.ImageIDToScan,
	//	Hash: config.ImageHashToScan,
	// })

	vulResults.MergedResults = mergedResults
	
	logger.Info("Vulnerabilities Done...")

	return &vulResults, nil
}

func (v Vulnerabilities) GetType() types.FamilyType {
	return types.Vulnerabilities
}

// ensure types implement the requisite interfaces.
var _ interfaces.Family = &Vulnerabilities{}

func New(conf Config) *Vulnerabilities {
	return &Vulnerabilities{
		conf: conf,
	}
}
