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

package database

import (
	"fmt"
	"sync"

	"github.com/openclarity/vmclarity/backend/pkg/database/gorm"
	"github.com/openclarity/vmclarity/backend/pkg/database/types"
)

type DBDriver func(config types.DBConfig) (types.Database, error)

var (
	DBDrivers           map[string]DBDriver
	RegisterDriversOnce sync.Once
)

func RegisterDrivers() {
	RegisterDriversOnce.Do(func() {
		// If DBDrivers is initialised before this function is called don't
		// reset it, this is useful for testing.
		if DBDrivers == nil {
			DBDrivers = map[string]DBDriver{}
		}
		DBDrivers[types.DBDriverTypeLocal] = gorm.NewDatabase
	})
}

//
//func (db *Handler) CreateFakeData() {
//	// Create scopes
//
//	var scopes = &Scopes{
//		Type: "AwsScope",
//		AwsScopesRegions: []AwsScopesRegion{
//			{
//				RegionID: "eu-central-1",
//				AwsRegionVpcs: []AwsRegionVpc{
//					{
//						VpcID: "vpc-1-from-eu-central-1",
//						AwsVpcSecurityGroups: []AwsVpcSecurityGroup{
//							{
//								GroupID: "sg-1-from-vpc-1-from-eu-central-1",
//							},
//						},
//					},
//					{
//						VpcID: "vpc-2-from-eu-central-1",
//						AwsVpcSecurityGroups: []AwsVpcSecurityGroup{
//							{
//								GroupID: "sg-2-from-vpc-1-from-eu-central-1",
//							},
//						},
//					},
//				},
//			},
//			{
//				RegionID: "us-east-1",
//				AwsRegionVpcs: []AwsRegionVpc{
//					{
//						VpcID: "vpc-1-from-us-east-1",
//						AwsVpcSecurityGroups: []AwsVpcSecurityGroup{
//							{
//								GroupID: "sg-1-from-vpc-1-from-us-east-1",
//							},
//						},
//					},
//					{
//						VpcID: "vpc-2-from-us-east-1",
//						AwsVpcSecurityGroups: []AwsVpcSecurityGroup{
//							{
//								GroupID: "sg-1-from-vpc-2-from-us-east-1",
//							},
//							{
//								GroupID: "sg-2-from-vpc-2-from-us-east-1",
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//	if _, err := db.ScopesTable().SetScopes(scopes); err != nil {
//		log.Fatalf("failed to set scopes: %v", err)
//	}
//
//	// Create scan configs
//
//	// Scan config 1
//	scanConfig1Families := &models.ScanFamiliesConfig{
//		Exploits: &models.ExploitsConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Malware: &models.MalwareConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Misconfigurations: &models.MisconfigurationsConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Rootkits: &models.RootkitsConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Sbom: &models.SBOMConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Secrets: &models.SecretsConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Vulnerabilities: &models.VulnerabilitiesConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//	}
//	scanConfig1FamiliesB, err := json.Marshal(scanConfig1Families)
//	if err != nil {
//		log.Fatalf("failed marshal scanConfig1Families: %v", err)
//	}
//	tag1 := models.Tag{
//		Key:   utils.StringPtr("app"),
//		Value: utils.StringPtr("my-app1"),
//	}
//	tag2 := models.Tag{
//		Key:   utils.StringPtr("app"),
//		Value: utils.StringPtr("my-app2"),
//	}
//	tag3 := models.Tag{
//		Key:   utils.StringPtr("system"),
//		Value: utils.StringPtr("sys1"),
//	}
//	tag4 := models.Tag{
//		Key:   utils.StringPtr("system"),
//		Value: utils.StringPtr("sys2"),
//	}
//	ScanConfig1SecurityGroups := []models.AwsSecurityGroup{
//		{
//			Id: utils.StringPtr("sg-1-from-vpc-1-from-eu-central-1"),
//		},
//	}
//	ScanConfig1VPCs := []models.AwsVPC{
//		{
//			Id:             utils.StringPtr("vpc-1-from-eu-central-1"),
//			SecurityGroups: &ScanConfig1SecurityGroups,
//		},
//	}
//	ScanConfig1Regions := []models.AwsRegion{
//		{
//			Id:   utils.StringPtr("eu-central-1"),
//			Vpcs: &ScanConfig1VPCs,
//		},
//	}
//	scanConfig1SelectorTags := []models.Tag{tag1, tag2}
//	scanConfig1ExclusionTags := []models.Tag{tag3, tag4}
//	scanConfig1Scope := models.AwsScanScope{
//		All:                        utils.BoolPtr(false),
//		InstanceTagExclusion:       &scanConfig1ExclusionTags,
//		InstanceTagSelector:        &scanConfig1SelectorTags,
//		ObjectType:                 "AwsScanScope",
//		Regions:                    &ScanConfig1Regions,
//		ShouldScanStoppedInstances: utils.BoolPtr(false),
//	}
//
//	var scanConfig1ScopeType models.ScanScopeType
//
//	err = scanConfig1ScopeType.FromAwsScanScope(scanConfig1Scope)
//	if err != nil {
//		log.Fatalf("failed to convert scanConfig1Scope: %v", err)
//	}
//
//	scanConfig1ScopeB, err := scanConfig1ScopeType.MarshalJSON()
//	if err != nil {
//		log.Fatalf("failed to marshal scanConfig1ScopeType: %v", err)
//	}
//
//	single1 := models.SingleScheduleScanConfig{
//		OperationTime: time.Now(),
//	}
//	var scanConfig1Scheduled models.RuntimeScheduleScanConfigType
//	err = scanConfig1Scheduled.FromSingleScheduleScanConfig(single1)
//	if err != nil {
//		log.Fatalf("failed to create FromSingleScheduleScanConfig: %v", err)
//	}
//	scanConfig1ScheduledB, err := scanConfig1Scheduled.MarshalJSON()
//
//	// Scan config 2
//	scanConfig2Families := &models.ScanFamiliesConfig{
//		Exploits: &models.ExploitsConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Malware: &models.MalwareConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Misconfigurations: &models.MisconfigurationsConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Rootkits: &models.RootkitsConfig{
//			Enabled: utils.BoolPtr(true),
//		},
//		Sbom: &models.SBOMConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Secrets: &models.SecretsConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//		Vulnerabilities: &models.VulnerabilitiesConfig{
//			Enabled: utils.BoolPtr(false),
//		},
//	}
//	scanConfig2FamiliesB, err := json.Marshal(scanConfig2Families)
//	if err != nil {
//		log.Fatalf("failed marshal scanConfig2Families: %v", err)
//	}
//	ScanConfig2SecurityGroups := []models.AwsSecurityGroup{
//		{
//			Id: utils.StringPtr("sg-1-from-vpc-1-from-us-east-1"),
//		},
//	}
//	ScanConfig2VPCs := []models.AwsVPC{
//		{
//			Id:             utils.StringPtr("vpc-1-from-us-east-1"),
//			SecurityGroups: &ScanConfig2SecurityGroups,
//		},
//	}
//	ScanConfig2Regions := []models.AwsRegion{
//		{
//			Id:   utils.StringPtr("us-east-1"),
//			Vpcs: &ScanConfig2VPCs,
//		},
//	}
//	scanConfig2SelectorTags := []models.Tag{tag2}
//	scanConfig2ExclusionTags := []models.Tag{tag4}
//	scanConfig2Scope := models.AwsScanScope{
//		All:                        utils.BoolPtr(false),
//		InstanceTagExclusion:       &scanConfig2ExclusionTags,
//		InstanceTagSelector:        &scanConfig2SelectorTags,
//		ObjectType:                 "AwsScanScope",
//		Regions:                    &ScanConfig2Regions,
//		ShouldScanStoppedInstances: utils.BoolPtr(true),
//	}
//
//	var scanConfig2ScopeType models.ScanScopeType
//
//	err = scanConfig2ScopeType.FromAwsScanScope(scanConfig2Scope)
//	if err != nil {
//		log.Fatalf("failed to convert scanConfig2Scope: %v", err)
//	}
//
//	scanConfig2ScopeB, err := scanConfig2ScopeType.MarshalJSON()
//	if err != nil {
//		log.Fatalf("failed to marshal scanConfig2ScopeType: %v", err)
//	}
//
//	single2 := models.SingleScheduleScanConfig{
//		OperationTime: time.Now(),
//	}
//	var scanConfig2Scheduled models.RuntimeScheduleScanConfigType
//	err = scanConfig2Scheduled.FromSingleScheduleScanConfig(single2)
//	if err != nil {
//		log.Fatalf("failed to create FromSingleScheduleScanConfig: %v", err)
//	}
//	scanConfig2ScheduledB, err := scanConfig2Scheduled.MarshalJSON()
//
//	scanConfigs := []ScanConfig{
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "1"),
//			},
//			Name:               utils.StringPtr("demo scan 1"),
//			ScanFamiliesConfig: scanConfig1FamiliesB,
//			Scheduled:          scanConfig1ScheduledB,
//			Scope:              scanConfig1ScopeB,
//		},
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "2"),
//			},
//			Name:               utils.StringPtr("demo scan 2"),
//			ScanFamiliesConfig: scanConfig2FamiliesB,
//			Scheduled:          scanConfig2ScheduledB,
//			Scope:              scanConfig2ScopeB,
//		},
//	}
//	if _, err := db.ScanConfigsTable().SaveScanConfig(&scanConfigs[0]); err != nil {
//		log.Fatalf("failed to save scan config 1: %v", err)
//	}
//	if _, err := db.ScanConfigsTable().SaveScanConfig(&scanConfigs[1]); err != nil {
//		log.Fatalf("failed to save scan config 2: %v", err)
//	}
//
//	// Create targets
//	targets := []Target{
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "1"),
//			},
//			Type:             "VMInfo",
//			Location:         utils.StringPtr("eu-central-1"),
//			InstanceID:       utils.StringPtr("i-instance-1-from-eu-central-1"),
//			InstanceProvider: utils.StringPtr("aws"),
//		},
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "2"),
//			},
//			Type:             "VMInfo",
//			Location:         utils.StringPtr("eu-central-1"),
//			InstanceID:       utils.StringPtr("i-instance-2-from-eu-central-1"),
//			InstanceProvider: utils.StringPtr("aws"),
//		},
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "3"),
//			},
//			Type:             "VMInfo",
//			Location:         utils.StringPtr("us-east-1"),
//			InstanceID:       utils.StringPtr("i-instance-1-from-us-east-1"),
//			InstanceProvider: utils.StringPtr("aws"),
//		},
//	}
//	if _, err := db.TargetsTable().SaveTarget(&targets[0]); err != nil {
//		log.Fatalf("failed to save target 1: %v", err)
//	}
//	if _, err := db.TargetsTable().SaveTarget(&targets[1]); err != nil {
//		log.Fatalf("failed to save target 2: %v", err)
//	}
//	if _, err := db.TargetsTable().SaveTarget(&targets[2]); err != nil {
//		log.Fatalf("failed to save target 3: %v", err)
//	}
//
//	// Create scans
//
//	// Create scan 1
//	scan1Start := time.Now()
//	scan1End := scan1Start.Add(10 * time.Hour)
//	scan1Targets := []string{targets[0].ID.String(), targets[1].ID.String()}
//	scan1TargetsB, err := json.Marshal(scan1Targets)
//	if err != nil {
//		log.Fatalf("failed to marshal scan1Targets: %v", err)
//	}
//	scan1Summary := &models.ScanSummary{
//		JobsCompleted:          utils.PointerTo[int](23),
//		JobsLeftToRun:          utils.PointerTo[int](107),
//		TotalExploits:          utils.PointerTo[int](14),
//		TotalMalware:           utils.PointerTo[int](44),
//		TotalMisconfigurations: utils.PointerTo[int](9),
//		TotalPackages:          utils.PointerTo[int](4221),
//		TotalRootkits:          utils.PointerTo[int](1),
//		TotalSecrets:           utils.PointerTo[int](0),
//		TotalVulnerabilities: &models.VulnerabilityScanSummary{
//			TotalCriticalVulnerabilities:   utils.PointerTo[int](9),
//			TotalHighVulnerabilities:       utils.PointerTo[int](12),
//			TotalLowVulnerabilities:        utils.PointerTo[int](424),
//			TotalMediumVulnerabilities:     utils.PointerTo[int](1551),
//			TotalNegligibleVulnerabilities: utils.PointerTo[int](132),
//		},
//	}
//	scan1SummaryB, err := json.Marshal(scan1Summary)
//	if err != nil {
//		log.Fatalf("failed to marshal scan1Summary: %v", err)
//	}
//
//	scan1ConfigSnapshot := &models.ScanConfigData{
//		Name:               utils.PointerTo[string]("Scan Config 1"),
//		ScanFamiliesConfig: scanConfig1Families,
//		Scheduled:          &scanConfig1Scheduled,
//		Scope:              &scanConfig1ScopeType,
//	}
//	scan1ConfigSnapshotB, err := json.Marshal(scan1ConfigSnapshot)
//	if err != nil {
//		log.Fatalf("failed to marshal scan1ConfigSnapshot: %v", err)
//	}
//
//	// Create scan 2
//	scan2Start := time.Now()
//	scan2Targets := []string{targets[2].ID.String()}
//	scan2TargetsB, err := json.Marshal(scan2Targets)
//	if err != nil {
//		log.Fatalf("failed to marshal scan2TargetsB: %v", err)
//	}
//
//	scan2Summary := &models.ScanSummary{
//		JobsCompleted:          utils.PointerTo[int](77),
//		JobsLeftToRun:          utils.PointerTo[int](98),
//		TotalExploits:          utils.PointerTo[int](6),
//		TotalMalware:           utils.PointerTo[int](0),
//		TotalMisconfigurations: utils.PointerTo[int](75),
//		TotalPackages:          utils.PointerTo[int](9778),
//		TotalRootkits:          utils.PointerTo[int](5),
//		TotalSecrets:           utils.PointerTo[int](557),
//		TotalVulnerabilities: &models.VulnerabilityScanSummary{
//			TotalCriticalVulnerabilities:   utils.PointerTo[int](11),
//			TotalHighVulnerabilities:       utils.PointerTo[int](52),
//			TotalLowVulnerabilities:        utils.PointerTo[int](241),
//			TotalMediumVulnerabilities:     utils.PointerTo[int](8543),
//			TotalNegligibleVulnerabilities: utils.PointerTo[int](73),
//		},
//	}
//	scan2SummaryB, err := json.Marshal(scan2Summary)
//	if err != nil {
//		log.Fatalf("failed to marshal scan2Summary: %v", err)
//	}
//
//	scan2ConfigSnapshot := &models.ScanConfigData{
//		Name:               utils.PointerTo[string]("Scan Config 2"),
//		ScanFamiliesConfig: scanConfig2Families,
//		Scheduled:          &scanConfig2Scheduled,
//		Scope:              &scanConfig2ScopeType,
//	}
//	scan2ConfigSnapshotB, err := json.Marshal(scan2ConfigSnapshot)
//	if err != nil {
//		log.Fatalf("failed to marshal scan2ConfigSnapshot: %v", err)
//	}
//
//	scans := []Scan{
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "1"),
//			},
//			ScanStartTime:      &scan1Start,
//			ScanEndTime:        &scan1End,
//			ScanConfigID:       utils.StringPtr(scanConfigs[0].ID.String()),
//			ScanConfigSnapshot: scan1ConfigSnapshotB,
//			State:              string(models.Done),
//			StateMessage:       "Scan was completed successfully",
//			StateReason:        string(models.ScanStateReasonSuccess),
//			Summary:            scan1SummaryB,
//			TargetIDs:          scan1TargetsB,
//		},
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "2"),
//			},
//			ScanStartTime:      &scan2Start,
//			ScanEndTime:        nil, // not ended
//			ScanConfigID:       utils.StringPtr(scanConfigs[1].ID.String()),
//			ScanConfigSnapshot: scan2ConfigSnapshotB,
//			State:              string(models.InProgress),
//			StateMessage:       "Scan is in progress",
//			StateReason:        string(models.ScanStateReasonSuccess),
//			Summary:            scan2SummaryB,
//			TargetIDs:          scan2TargetsB,
//		},
//	}
//	if _, err := db.ScansTable().SaveScan(&scans[0]); err != nil {
//		log.Fatalf("failed to save scan 1: %v", err)
//	}
//	if _, err := db.ScansTable().SaveScan(&scans[1]); err != nil {
//		log.Fatalf("failed to save scan 2: %v", err)
//	}
//
//	// Create scan results
//	scanFindingsSummary := &models.ScanFindingsSummary{
//		TotalExploits:          utils.PointerTo[int](6),
//		TotalMalware:           utils.PointerTo[int](0),
//		TotalMisconfigurations: utils.PointerTo[int](75),
//		TotalPackages:          utils.PointerTo[int](9778),
//		TotalRootkits:          utils.PointerTo[int](5),
//		TotalSecrets:           utils.PointerTo[int](557),
//		TotalVulnerabilities: &models.VulnerabilityScanSummary{
//			TotalCriticalVulnerabilities:   utils.PointerTo[int](11),
//			TotalHighVulnerabilities:       utils.PointerTo[int](52),
//			TotalLowVulnerabilities:        utils.PointerTo[int](241),
//			TotalMediumVulnerabilities:     utils.PointerTo[int](8543),
//			TotalNegligibleVulnerabilities: utils.PointerTo[int](73),
//		},
//	}
//	scanFindingsSummaryB, err := json.Marshal(scanFindingsSummary)
//	if err != nil {
//		log.Fatalf("failed to marshal scanFindingsSummary: %v", err)
//	}
//	ScanFindingsSummary2 := &models.ScanFindingsSummary{
//		TotalExploits:          utils.PointerTo[int](10),
//		TotalMalware:           utils.PointerTo[int](1),
//		TotalMisconfigurations: utils.PointerTo[int](0),
//		TotalPackages:          utils.PointerTo[int](10000001),
//		TotalRootkits:          utils.PointerTo[int](2),
//		TotalSecrets:           utils.PointerTo[int](666),
//		TotalVulnerabilities: &models.VulnerabilityScanSummary{
//			TotalCriticalVulnerabilities:   utils.PointerTo[int](1),
//			TotalHighVulnerabilities:       utils.PointerTo[int](3),
//			TotalLowVulnerabilities:        utils.PointerTo[int](5),
//			TotalMediumVulnerabilities:     utils.PointerTo[int](7),
//			TotalNegligibleVulnerabilities: utils.PointerTo[int](8),
//		},
//	}
//	ScanFindingsSummary2B, err := json.Marshal(ScanFindingsSummary2)
//	if err != nil {
//		log.Fatalf("failed to marshal ScanFindingsSummary2: %v", err)
//	}
//	scanResults := []ScanResult{
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "1"),
//			},
//			ScanID:   scans[0].ID.String(),
//			TargetID: targets[0].ID.String(),
//			Summary:  scanFindingsSummaryB,
//		},
//		{
//			Base: Base{
//				ID: uuid.NewV5(uuid.Nil, "2"),
//			},
//			ScanID:   scans[1].ID.String(),
//			TargetID: targets[1].ID.String(),
//			Summary:  ScanFindingsSummary2B,
//		},
//	}
//	if _, err := db.ScanResultsTable().SaveScanResult(&scanResults[0]); err != nil {
//		log.Fatalf("failed to save scanResults 1: %v", err)
//	}
//	if _, err := db.ScanResultsTable().SaveScanResult(&scanResults[1]); err != nil {
//		log.Fatalf("failed to save scanResults 2: %v", err)
//	}
//}

func InitializeDatabase(config types.DBConfig) (types.Database, error) {
	RegisterDrivers()

	if driver, ok := DBDrivers[config.DriverType]; ok {
		return driver(config)
	}
	return nil, fmt.Errorf("unknown DB driver %s", config.DriverType)
}
