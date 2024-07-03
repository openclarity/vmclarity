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

package misconfiguration

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
)

type Misconfiguration struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.Misconfigurations] {
	return &Misconfiguration{
		conf: conf,
	}
}

func (m Misconfiguration) GetType() families.FamilyType {
	return families.Misconfiguration
}

func (m Misconfiguration) Run(ctx context.Context, _ *families.Results) (*types.Misconfigurations, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "misconfiguration")
	logger.Info("Misconfiguration Run...")

	manager := scan_manager.New[types.ScannersConfig, *types.ScannerResult](m.conf.ScannersList, m.conf.ScannersConfig, logger, types.Factory)
	results, err := manager.Scan(ctx, m.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for misconfigurations: %w", err)
	}

	misConfigResults := types.NewMisconfigurations()

	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata)
		if familiesutils.ShouldStripInputPath(result.ScanInput.StripPathFromResult, m.conf.StripInputPaths) {
			result.ScanResult = stripPathFromResult(result.ScanResult, result.ScanInput.Input)
		}
		misConfigResults.Merge(result.ScanResult)
	}

	logger.Info("Misconfiguration Done...")

	return misConfigResults, nil
}

// StripPathFromResult strip input path from results wherever it is found.
func stripPathFromResult(result *types.ScannerResult, path string) *types.ScannerResult {
	for i := range result.Misconfigurations {
		result.Misconfigurations[i].Location = familiesutils.TrimMountPath(result.Misconfigurations[i].Location, path)
	}

	return result
}
