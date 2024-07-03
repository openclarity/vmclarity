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
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/infofinder/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
)

type InfoFinder struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.Infos] {
	return &InfoFinder{
		conf: conf,
	}
}

func (i InfoFinder) GetType() families.FamilyType {
	return families.InfoFinder
}

func (i InfoFinder) Run(ctx context.Context, _ *families.Results) (*types.Infos, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "info finder")
	logger.Info("InfoFinder Run...")

	manager := scan_manager.New[types.ScannersConfig, *types.ScannerResult](i.conf.ScannersList, i.conf.ScannersConfig, logger, types.Factory)
	results, err := manager.Scan(ctx, i.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for infofinders: %w", err)
	}

	infoFinderResults := types.NewInfos()

	// Merge results.
	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata)
		if familiesutils.ShouldStripInputPath(result.ScanInput.StripPathFromResult, i.conf.StripInputPaths) {
			result.ScanResult = stripPathFromResult(result.ScanResult, result.ScanInput.Input)
		}
		infoFinderResults.Merge(result.ScanResult)
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
