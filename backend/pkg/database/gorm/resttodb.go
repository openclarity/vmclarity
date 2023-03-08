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

	uuid "github.com/satori/go.uuid"

	"github.com/openclarity/vmclarity/api/models"
)

// nolint:cyclop
func ConvertToDBScanResult(result models.TargetScanResult) (ScanResult, error) {
	var ret ScanResult
	var err error
	var scanResultUUID uuid.UUID

	if result.Id != nil {
		scanResultUUID, err = uuid.FromString(*result.Id)
		if err != nil {
			return ret, fmt.Errorf("failed to convert scanResultID %v to uuid: %v", *result.Id, err)
		}
	}
	ret.ScanID = *result.Scan.Id
	ret.TargetID = *result.Target.Id

	if result.Exploits != nil {
		ret.Exploits, err = json.Marshal(result.Exploits)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Malware != nil {
		ret.Malware, err = json.Marshal(result.Malware)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Misconfigurations != nil {
		ret.Misconfigurations, err = json.Marshal(result.Misconfigurations)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Rootkits != nil {
		ret.Rootkits, err = json.Marshal(result.Rootkits)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Sboms != nil {
		ret.Sboms, err = json.Marshal(result.Sboms)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if result.Secrets != nil {
		ret.Secrets, err = json.Marshal(result.Secrets)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Status != nil {
		ret.Status, err = json.Marshal(result.Status)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Vulnerabilities != nil {
		ret.Vulnerabilities, err = json.Marshal(result.Vulnerabilities)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	ret.Base = Base{ID: scanResultUUID}

	return ret, nil
}

func ConvertToDBScan(scan models.Scan) (Scan, error) {
	var ret Scan
	var err error
	var scanUUID uuid.UUID

	if scan.Id != nil {
		scanUUID, err = uuid.FromString(*scan.Id)
		if err != nil {
			return ret, fmt.Errorf("failed to convert scanID %v to uuid: %v", scan.Id, err)
		}
	}

	if scan.ScanConfig != nil {
		ret.ScanConfigID = &scan.ScanConfig.Id
	}

	ret.ScanEndTime = scan.EndTime

	ret.ScanStartTime = scan.StartTime

	if scan.ScanConfigSnapshot != nil {
		ret.ScanConfigSnapshot, err = json.Marshal(scan.ScanConfigSnapshot)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if scan.State != nil {
		ret.State = string(*scan.State)
	}
	if scan.StateMessage != nil {
		ret.StateMessage = *scan.StateMessage
	}
	if scan.StateReason != nil {
		ret.StateReason = string(*scan.StateReason)
	}

	if scan.Summary != nil {
		ret.Summary, err = json.Marshal(scan.Summary)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if scan.TargetIDs != nil {
		ret.TargetIDs, err = json.Marshal(scan.TargetIDs)
		if err != nil {
			return ret, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	ret.Base = Base{ID: scanUUID}

	return ret, nil
}

func ConvertToDbScopes(scope *models.ScopeType) (*Scopes, error) {
	var ret Scopes
	disc, err := scope.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("failed to get discriminator: %w", err)
	}
	switch disc {
	case "AwsScope":
		awsScope, err := scope.AsAwsScope()
		if err != nil {
			return nil, fmt.Errorf("failed to convert scope to aws scope: %w", err)
		}
		ret.Type = awsScope.ObjectType
		ret.AwsScopesRegions = convertToDbRegions(awsScope.Regions)
		return &ret, nil
	default:
		return nil, fmt.Errorf("unknown scope type: %v", disc)
	}
}

func convertToDbRegions(regions *[]models.AwsRegion) []AwsScopesRegion {
	var ret []AwsScopesRegion
	if regions != nil {
		for _, region := range *regions {
			ret = append(ret, convertToDbRegion(region))
		}
	}

	return ret
}

func convertToDbRegion(region models.AwsRegion) AwsScopesRegion {
	return AwsScopesRegion{
		RegionID:      *region.Id,
		AwsRegionVpcs: convertToDbVPCs(region.Vpcs),
	}
}

func convertToDbVPCs(vpcs *[]models.AwsVPC) []AwsRegionVpc {
	var ret []AwsRegionVpc
	if vpcs != nil {
		for _, vpc := range *vpcs {
			ret = append(ret, AwsRegionVpc{
				VpcID:                *vpc.Id,
				AwsVpcSecurityGroups: convertToDbSecurityGroups(vpc.SecurityGroups),
			})
		}
	}

	return ret
}

func convertToDbSecurityGroups(groups *[]models.AwsSecurityGroup) []AwsVpcSecurityGroup {
	var ret []AwsVpcSecurityGroup
	if groups != nil {
		for _, group := range *groups {
			ret = append(ret, AwsVpcSecurityGroup{
				GroupID: *group.Id,
			})
		}
	}

	return ret
}
