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

package infofinder

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families/infofinder/job"
	"github.com/openclarity/vmclarity/scanner/families/infofinder/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/job_manager"
)

type InfoFinder struct {
	conf Config
}

func New(conf Config) *InfoFinder {
	return &InfoFinder{
		conf: conf,
	}
}

func (i InfoFinder) GetType() familiestypes.FamilyType {
	return familiestypes.InfoFinder
}

func (i InfoFinder) Run(ctx context.Context, _ *familiestypes.FamiliesResults) (familiestypes.FamilyResult, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "info finder")
	logger.Info("InfoFinder Run...")

	manager := job_manager.New(i.conf.ScannersList, i.conf.ScannersConfig, logger, job.Factory)
	processResults, err := manager.Process(ctx, i.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for infofinders: %w", err)
	}

	infoFinderResults := types.NewFamilyResult()

	// Merge results.
	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.ScannerName)
		if assetScan, ok := result.Result.(*types.ScannerResult); ok {
			if familiesutils.ShouldStripInputPath(result.Input.StripPathFromResult, i.conf.StripInputPaths) {
				assetScan = stripPathFromResult(assetScan, result.InputPath)
			}
			infoFinderResults.AddScannerResult(assetScan)
		} else {
			return nil, fmt.Errorf("received bad scanner result type %T, expected infofinderTypes.ScannerResult", result)
		}
		infoFinderResults.Metadata.InputScans = append(infoFinderResults.Metadata.InputScans, result.InputScanMetadata)
	}

	logger.Info("InfoFinder Done...")

	return infoFinderResults, nil
}

// stripPathFromResult strip input path from results wherever it is found.
func stripPathFromResult(result *types.ScannerResult, path string) *types.ScannerResult {
	for i := range result.Infos {
		result.Infos[i].Path = familiesutils.TrimMountPath(result.Infos[i].Path, path)
	}
	return result
}
