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
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/internal/job_manager"
)

type Misconfiguration struct {
	conf types.Config
}

func New(conf types.Config) familiestypes.Family[*types.Misconfigurations] {
	return &Misconfiguration{
		conf: conf,
	}
}

func (m Misconfiguration) GetType() familiestypes.FamilyType {
	return familiestypes.Misconfiguration
}

func (m Misconfiguration) Run(ctx context.Context, _ *familiestypes.Results) (*types.Misconfigurations, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "misconfiguration")
	logger.Info("Misconfiguration Run...")

	manager := job_manager.New[types.ScannersConfig, *types.ScannerResult](m.conf.ScannersList, m.conf.ScannersConfig, logger, types.Factory)
	processResults, err := manager.Process(ctx, m.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for misconfigurations: %w", err)
	}

	misConfigResults := types.NewMisconfigurations()

	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.Result.ScannerName)
		if familiesutils.ShouldStripInputPath(result.Input.StripPathFromResult, m.conf.StripInputPaths) {
			result.Result = stripPathFromResult(result.Result, result.Input.Input)
		}
		misConfigResults.Merge(result.Result)
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
