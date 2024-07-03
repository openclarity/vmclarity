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

package rootkits

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
)

type Rootkits struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.Rootkits] {
	return &Rootkits{
		conf: conf,
	}
}

func (r Rootkits) GetType() families.FamilyType {
	return families.Rootkits
}

func (r Rootkits) Run(ctx context.Context, _ *families.Results) (*types.Rootkits, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "rootkits")
	logger.Info("Rootkits Run...")

	manager := scan_manager.New[types.ScannersConfig, *types.ScannerResult](r.conf.ScannersList, r.conf.ScannersConfig, logger, types.Factory)
	results, err := manager.Scan(ctx, r.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for rootkits: %w", err)
	}

	rootkitsResults := types.NewRootkits()

	// Merge results.
	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata)
		if familiesutils.ShouldStripInputPath(result.ScanInput.StripPathFromResult, r.conf.StripInputPaths) {
			result.ScanResult = stripPathFromResult(result.ScanResult, result.ScanInput.Input)
		}
		rootkitsResults.Merge(result.ScanResult)
	}

	logger.Info("Rootkits Done...")

	return rootkitsResults, nil
}

// StripPathFromResult strip input path from results wherever it is found.
func stripPathFromResult(result *types.ScannerResult, path string) *types.ScannerResult {
	for i := range result.Rootkits {
		result.Rootkits[i].Message = familiesutils.RemoveMountPathSubStringIfNeeded(result.Rootkits[i].Message, path)
	}

	return result
}
