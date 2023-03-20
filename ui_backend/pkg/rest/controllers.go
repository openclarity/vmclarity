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

package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	backendmodels "github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
	"github.com/openclarity/vmclarity/ui_backend/api/models"
)

func (s *ServerImpl) GetDashboardRiskiestRegions(ctx echo.Context) error {
	targets, err := s.BackendClient.GetTargets(context.TODO(), backendmodels.GetTargetsParams{})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get targets: %v", err))
	}

	// map region to findings count per finding type
	var findingsPerRegion = make(map[string]map[backendmodels.ScanType]int)

	// get all targets, and add their findings count from the last time
	// they were scanned to the total region findings count.
	// target/ScanFindingsSummary should contain the latest results per family.
	// TODO ScanFindingsSummary is currently not containing the latest results per family, but all results ever found.
	for _, target := range *targets.Items {
		location, err := getTargetLocation(target)
		if err != nil {
			log.Errorf("Failed to get target location, skipping target: %v", err)
			continue
		}
		if _, ok := findingsPerRegion[location]; !ok {
			findingsPerRegion[location] = make(map[backendmodels.ScanType]int)
		}
		addTargetFindingsCount(findingsPerRegion[location], target.Summary)
	}

	items := []models.RegionFindings{}
	for region, findings := range findingsPerRegion {
		items = append(items, models.RegionFindings{
			FindingsCount: createFindingsCount(findings),
			RegionName:    &region,
		})
	}

	return sendResponse(ctx, http.StatusOK, &models.RiskiestRegions{
		Count: utils.PointerTo(len(items)),
		Items: &items,
	})
}

func createFindingsCount(findings map[backendmodels.ScanType]int) *models.FindingsCount {
	var ret models.FindingsCount
	ret.Malware = utils.PointerTo(findings["Malware"])
	ret.Exploits = utils.PointerTo(findings["Exploits"])
	ret.Vulnerabilities = utils.PointerTo(findings["Vulnerabilities"])
	ret.Rootkits = utils.PointerTo(findings["Rootkits"])
	ret.Misconfigurations = utils.PointerTo(findings["Misconfigurations"])
	ret.Secrets = utils.PointerTo(findings["Secrets"])

	return &ret
}

func getTargetLocation(target backendmodels.Target) (string, error) {
	discriminator, err := target.TargetInfo.ValueByDiscriminator()
	if err != nil {
		return "", fmt.Errorf("failed to get value by discriminator: %w", err)
	}

	switch info := discriminator.(type) {
	case backendmodels.VMInfo:
		return info.Location, nil
	default:
		return "", fmt.Errorf("target type is not supported (%T): %w", discriminator, err)
	}
}

func addTargetFindingsCount(findingsCount map[backendmodels.ScanType]int, summary *backendmodels.ScanFindingsSummary) {
	if summary == nil {
		return
	}
	if summary.TotalExploits != nil {
		findingsCount["Exploits"] += *summary.TotalExploits
	}
	if summary.TotalMisconfigurations != nil {
		findingsCount["Misconfigurations"] += *summary.TotalMisconfigurations
	}
	if summary.TotalRootkits != nil {
		findingsCount["Rootkits"] += *summary.TotalRootkits
	}
	if summary.TotalSecrets != nil {
		findingsCount["Secrets"] += *summary.TotalSecrets
	}
	if summary.TotalMalware != nil {
		findingsCount["Malware"] += *summary.TotalMalware
	}
	if summary.TotalVulnerabilities != nil {
		findingsCount["Vulnerabilities"] += getTotalVulnerabilities(summary.TotalVulnerabilities)
	}
}

func getTotalVulnerabilities(summary *backendmodels.VulnerabilityScanSummary) int {
	total := 0
	if summary.TotalCriticalVulnerabilities != nil {
		total += *summary.TotalCriticalVulnerabilities
	}
	if summary.TotalHighVulnerabilities != nil {
		total += *summary.TotalHighVulnerabilities
	}
	if summary.TotalMediumVulnerabilities != nil {
		total += *summary.TotalMediumVulnerabilities
	}
	if summary.TotalLowVulnerabilities != nil {
		total += *summary.TotalLowVulnerabilities
	}
	// TODO add also negligible?
	if summary.TotalNegligibleVulnerabilities != nil {
		total += *summary.TotalNegligibleVulnerabilities
	}
	return total
}
