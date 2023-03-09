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
			return ret, fmt.Errorf("failed to unmarshal json: %w", err)
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

func ConvertToRestScopes(scopes *Scopes) (*models.ScopeType, error) {
	ret := models.ScopeType{}

	switch scopes.Type {
	case "AwsScope":
		awsScope := models.AwsScope{
			Regions: convertToRestRegions(scopes.AwsScopesRegions),
		}
		if err := ret.FromAwsScope(awsScope); err != nil {
			return nil, fmt.Errorf("FromAwsScope failed: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown scope type: %v", scopes.Type)
	}

	return &ret, nil
}

func convertToRestRegions(regions []AwsScopesRegion) *[]models.AwsRegion {
	var ret []models.AwsRegion
	for _, region := range regions {
		ret = append(ret, convertToRestRegion(region))
	}

	return &ret
}

func convertToRestRegion(region AwsScopesRegion) models.AwsRegion {
	return models.AwsRegion{
		Id:   &region.RegionID,
		Vpcs: convertToRestVPCs(region.AwsRegionVpcs),
	}
}

func convertToRestVPCs(vpcs []AwsRegionVpc) *[]models.AwsVPC {
	var ret []models.AwsVPC
	for i, _ := range vpcs {
		ret = append(ret, models.AwsVPC{
			Id:             &vpcs[i].VpcID,
			SecurityGroups: convertToRestSecurityGroups(vpcs[i].AwsVpcSecurityGroups),
		})
	}

	return &ret
}

func convertToRestSecurityGroups(groups []AwsVpcSecurityGroup) *[]models.AwsSecurityGroup {
	var ret []models.AwsSecurityGroup

	for i, _ := range groups {
		ret = append(ret, models.AwsSecurityGroup{
			Id: &groups[i].GroupID,
		})
	}

	return &ret
}
