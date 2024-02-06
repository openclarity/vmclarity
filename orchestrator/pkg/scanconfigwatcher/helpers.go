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

package scanconfigwatcher

import (
	"fmt"
	"time"

	"github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/cli/pkg/utils"
)

func newScanFromScanConfig(scanConfig *types.ScanConfig) *types.Scan {
	return &types.Scan{
		Name: utils.PointerTo(fmt.Sprintf("%s-%s", *scanConfig.Name, scanConfig.Scheduled.OperationTime.Format(time.RFC3339))),
		ScanConfig: &types.ScanConfigRelationship{
			Id: *scanConfig.Id,
		},
		AssetScanTemplate:   scanConfig.ScanTemplate.AssetScanTemplate,
		Scope:               scanConfig.ScanTemplate.Scope,
		MaxParallelScanners: scanConfig.ScanTemplate.MaxParallelScanners,
		TimeoutSeconds:      scanConfig.ScanTemplate.TimeoutSeconds,
		Status: types.NewScanStatus(
			types.ScanStatusStatePending,
			types.ScanStatusReasonCreated,
			nil,
		),
		Summary: &types.ScanSummary{
			JobsCompleted:          utils.PointerTo(0),
			JobsLeftToRun:          utils.PointerTo(0),
			TotalExploits:          utils.PointerTo(0),
			TotalMalware:           utils.PointerTo(0),
			TotalMisconfigurations: utils.PointerTo(0),
			TotalPackages:          utils.PointerTo(0),
			TotalRootkits:          utils.PointerTo(0),
			TotalSecrets:           utils.PointerTo(0),
			TotalInfoFinder:        utils.PointerTo(0),
			TotalVulnerabilities: &types.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo(0),
				TotalHighVulnerabilities:       utils.PointerTo(0),
				TotalLowVulnerabilities:        utils.PointerTo(0),
				TotalMediumVulnerabilities:     utils.PointerTo(0),
				TotalNegligibleVulnerabilities: utils.PointerTo(0),
			},
		},
	}
}
