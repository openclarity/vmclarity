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

package sbom

import (
	"context"
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/core/version"
	"github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/sbom/types"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
	"github.com/openclarity/vmclarity/scanner/utils"
	"github.com/openclarity/vmclarity/scanner/utils/converter"
)

type SBOM struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.SBOM] {
	return &SBOM{
		conf: conf,
	}
}

func (s SBOM) GetType() families.FamilyType {
	return families.SBOM
}

// nolint:cyclop
func (s SBOM) Run(ctx context.Context, _ *families.Results) (*types.SBOM, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "sbom")
	logger.Info("SBOM Run...")

	if len(s.conf.Inputs) == 0 {
		return nil, errors.New("inputs list is empty")
	}

	// TODO: move the logic from cli utils to shared utils
	// TODO: now that we support multiple inputs,
	//  we need to change the fact the mergedResults assumes it is only for 1 input?
	hash, err := utils.GenerateHash(s.conf.Inputs[0].InputType, s.conf.Inputs[0].Input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash for source %s: %w", s.conf.Inputs[0].Input, err)
	}

	// Run all scanners using scan manager
	manager := scan_manager.New(s.conf.AnalyzersList, s.conf.AnalyzersConfig, Factory)
	results, err := manager.Scan(ctx, s.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for sbom: %w", err)
	}

	mergedResults := newMergedResults(s.conf.Inputs[0].InputType, hash)

	// Merge results
	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata)
		mergedResults = mergedResults.Merge(result.ScanResult)
	}

	for i, with := range s.conf.MergeWith {
		name := fmt.Sprintf("merge_with_%d", i)
		cdxBOMBytes, err := converter.GetCycloneDXSBOMFromFile(with.SbomPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get CDX SBOM from path=%s: %w", with.SbomPath, err)
		}
		results := types.CreateScannerResult(cdxBOMBytes, name, with.SbomPath, common.SBOM)
		logger.Infof("Merging result from %q", with.SbomPath)
		mergedResults = mergedResults.Merge(results)
	}

	// TODO(sambetts) Expose CreateMergedSBOM as well as
	// CreateMergedSBOMBytes so that we don't need to re-convert it
	mergedSBOMBytes, err := mergedResults.CreateMergedSBOMBytes("cyclonedx-json", version.CommitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create merged output: %w", err)
	}

	cdxBom, err := converter.GetCycloneDXSBOMFromBytes(mergedSBOMBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load merged output to CDX bom: %w", err)
	}

	// Create result from merged data
	sbomResults := types.NewSBOM()
	sbomResults.Data = cdxBom

	logger.Info("SBOM Done...")

	return sbomResults, nil
}
