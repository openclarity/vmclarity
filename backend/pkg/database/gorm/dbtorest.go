// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package gorm

import (
	"encoding/json"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func ConvertToRestTarget(target Target) (models.Target, error) {
	ret := models.Target{
		TargetInfo: &models.TargetType{},
	}

	switch target.Type {
	case targetTypeVM:
		var cloudProvider *models.CloudProvider
		if target.InstanceProvider != nil {
			cp := models.CloudProvider(*target.InstanceProvider)
			cloudProvider = &cp
		}
		if err := ret.TargetInfo.FromVMInfo(models.VMInfo{
			InstanceID:       *target.InstanceID,
			InstanceProvider: cloudProvider,
			Location:         *target.Location,
		}); err != nil {
			return ret, fmt.Errorf("FromVMInfo failed: %w", err)
		}
	case targetTypeDir, targetTypePod:
		fallthrough
	default:
		return ret, fmt.Errorf("unknown target type: %v", target.Type)
	}
	ret.Id = utils.StringPtr(target.ID.String())
	ret.ScansCount = &scanCount
	if len(scanFindingsSummaryB) > 0 {
		scanFindingsSummary := models.ScanFindingsSummary{}
		err := json.Unmarshal(scanFindingsSummaryB, &scanFindingsSummary)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal scan findings summary: %w", err)
		}
		ret.Summary = &scanFindingsSummary
	}

	return ret, nil
}

func ConvertToRestTargets(targets []*Target, scanCount map[string]int, summary map[string][]byte, total int64) (*models.Targets, error) {
	ret := models.Targets{
		Items: &[]models.Target{},
	}

	for _, target := range targets {
		tr, err := ConvertToRestTarget(target, scanCount[target.ID.String()], summary[target.ID.String()])
		if err != nil {
			return ret, fmt.Errorf("failed to convert target: %w", err)
		}
		*ret.Items = append(*ret.Items, tr)
	}

	ret.Total = utils.IntPtr(len(targets))

	return ret, nil
}

// nolint:cyclop
func ConvertToRestScanResult(scanResult *database.ScanResult, scan *database.Scan, target *database.Target) (*models.TargetScanResult, error) {
	var ret models.TargetScanResult

	if scanResult.Secrets != nil {
		ret.Secrets = &models.SecretScan{}
		if err := json.Unmarshal(scanResult.Secrets, ret.Secrets); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Vulnerabilities != nil {
		ret.Vulnerabilities = &models.VulnerabilityScan{}
		if err := json.Unmarshal(scanResult.Vulnerabilities, ret.Vulnerabilities); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}

	if scanResult.Exploits != nil {
		ret.Exploits = &models.ExploitScan{}
		if err := json.Unmarshal(scanResult.Exploits, ret.Exploits); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Malware != nil {
		ret.Malware = &models.MalwareScan{}
		if err := json.Unmarshal(scanResult.Malware, ret.Malware); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Misconfigurations != nil {
		ret.Misconfigurations = &models.MisconfigurationScan{}
		if err := json.Unmarshal(scanResult.Misconfigurations, ret.Misconfigurations); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Rootkits != nil {
		ret.Rootkits = &models.RootkitScan{}
		if err := json.Unmarshal(scanResult.Rootkits, ret.Rootkits); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Sboms != nil {
		ret.Sboms = &models.SbomScan{}
		if err := json.Unmarshal(scanResult.Sboms, ret.Sboms); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	if scanResult.Status != nil {
		ret.Status = &models.TargetScanStatus{}
		if err := json.Unmarshal(scanResult.Status, ret.Status); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}
	ret.Id = utils.StringPtr(scanResult.ID.String())
	if scan != nil {
		convertScan, err := ConvertScan(scan)
		if err != nil {
			return nil, fmt.Errorf("failed to convert scan: %w", err)
		}
		ret.Scan = convertScan
	} else {
		ret.Scan = &models.Scan{Id: &scanResult.ScanID}
	}

	if target != nil {
		convertTarget, err := ConvertTarget(target, 0, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to convert scan: %w", err)
		}
		ret.Target = convertTarget
	} else {
		ret.Target = &models.Target{Id: &scanResult.TargetID}
	}

	return ret, nil
}

func ConvertToRestScanResults(scanResults []*database.ScanResult, scans []*database.Scan, targets []*database.Target, total int64) (*models.TargetScanResults, error) {
	ret := models.TargetScanResults{
		Items: &[]models.TargetScanResult{},
	}

	for i := range scanResults {
		sr, err := ConvertToRestScanResult(scanResults[i], scans[i], targets[i])
		if err != nil {
			return ret, fmt.Errorf("failed to convert scan result: %w", err)
		}
		*ret.Items = append(*ret.Items, sr)
	}

	ret.Total = utils.IntPtr(len(scanResults))

	return ret, nil
}

func ConvertToRestScan(scan Scan) (models.Scan, error) {
	var ret models.Scan

	if scan.ScanConfigSnapshot != nil {
		ret.ScanConfigSnapshot = &models.ScanConfigData{}
		if err := json.Unmarshal(scan.ScanConfigSnapshot, ret.ScanConfigSnapshot); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}

	if scan.Summary != nil {
		ret.Summary = &models.ScanSummary{}
		if err := json.Unmarshal(scan.Summary, ret.Summary); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}

	if scan.TargetIDs != nil {
		ret.TargetIDs = &[]string{}
		if err := json.Unmarshal(scan.TargetIDs, ret.TargetIDs); err != nil {
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}

	ret.Id = utils.StringPtr(scan.ID.String())
	ret.StartTime = scan.ScanStartTime
	ret.EndTime = scan.ScanEndTime
	ret.ScanConfig = &models.ScanConfigRelationship{Id: *scan.ScanConfigID}
	ret.State = utils.PointerTo[models.ScanState](models.ScanState(scan.State))
	ret.StateMessage = utils.PointerTo[string](scan.StateMessage)
	ret.StateReason = utils.PointerTo[models.ScanStateReason](models.ScanStateReason(scan.StateReason))

	return ret, nil
}

func ConvertToRestScans(scans []Scan) (models.Scans, error) {
	ret := models.Scans{
		Items: &[]models.Scan{},
	}

	for _, scan := range scans {
		sc, err := ConvertToRestScan(scan)
		if err != nil {
			return models.Scans{}, fmt.Errorf("failed to convert scan: %w", err)
		}
		*ret.Items = append(*ret.Items, sc)
	}

	ret.Total = utils.IntPtr(len(scans))

	return ret, nil
}

func ConvertScopes(scopes *database.Scopes) (*models.ScopeType, error) {
	ret := models.ScopeType{}

	switch scopes.Type {
	case "AwsScope":
		awsScope := models.AwsScope{
			Regions: convertRegions(scopes.AwsScopesRegions),
		}
		if err := ret.FromAwsScope(awsScope); err != nil {
			return nil, fmt.Errorf("FromAwsScope failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown scope type: %v", scopes.Type)
	}

	return &ret, nil
}

func convertRegions(regions []database.AwsScopesRegion) *[]models.AwsRegion {
	var ret []models.AwsRegion
	for _, region := range regions {
		ret = append(ret, convertRegion(region))
	}

	return &ret
}

func convertRegion(region database.AwsScopesRegion) models.AwsRegion {
	return models.AwsRegion{
		Id:   &region.RegionID,
		Vpcs: convertVPCs(region.AwsRegionVpcs),
	}
}

func convertVPCs(vpcs []database.AwsRegionVpc) *[]models.AwsVPC {
	var ret []models.AwsVPC
	for i, _ := range vpcs {
		ret = append(ret, models.AwsVPC{
			Id:             &vpcs[i].VpcID,
			SecurityGroups: convertSecurityGroups(vpcs[i].AwsVpcSecurityGroups),
		})
	}

	return &ret
}

func convertSecurityGroups(groups []database.AwsVpcSecurityGroup) *[]models.AwsSecurityGroup {
	var ret []models.AwsSecurityGroup

	for i, _ := range groups {
		ret = append(ret, models.AwsSecurityGroup{
			Id: &groups[i].GroupID,
		})
	}

	return &ret
}
