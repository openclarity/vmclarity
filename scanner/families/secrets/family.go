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

package secrets

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/secrets/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
)

type Secrets struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.Findings] {
	return &Secrets{
		conf: conf,
	}
}

func (s Secrets) GetType() families.FamilyType {
	return families.Secrets
}

func (s Secrets) Run(ctx context.Context, _ *families.Results) (*types.Findings, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "secrets")
	logger.Info("Secrets Run...")

	manager := scan_manager.New[types.ScannersConfig, *types.ScannerResult](s.conf.ScannersList, s.conf.ScannersConfig, logger, types.Factory)
	results, err := manager.Scan(ctx, s.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for secrets: %w", err)
	}

	secretsResults := types.NewFindings()

	// Merge results.
	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata)
		if familiesutils.ShouldStripInputPath(result.ScanInput.StripPathFromResult, s.conf.StripInputPaths) {
			result.ScanResult = stripPathFromResult(result.ScanResult, result.ScanInput.Input)
		}
		secretsResults.Merge(result.ScanResult)
	}

	logger.Info("Secrets Done...")

	return secretsResults, nil
}

// StripPathFromResult strip input path from results wherever it is found.
func stripPathFromResult(result *types.ScannerResult, path string) *types.ScannerResult {
	for i := range result.Findings {
		result.Findings[i].File = familiesutils.TrimMountPath(result.Findings[i].File, path)
		result.Findings[i].Fingerprint = familiesutils.RemoveMountPathSubStringIfNeeded(result.Findings[i].Fingerprint, path)
	}

	return result
}
