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

package initiator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/cli/pkg/standalone/vminfoprovider"
	cliutils "github.com/openclarity/vmclarity/cli/pkg/utils"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

type VMClarityInitiator struct {
	targetInfoType     string
	scanConfigFamilies *models.ScanFamiliesConfig
	Config
}

func newVMClarityInitiator(config Config) (*VMClarityInitiator, error) {
	if config.client == nil {
		return nil, errors.New("backend client must not be nil")
	}
	targetInfoType, err := getTargetInfoType(config.asset.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to get target type by inputType=%s: %v", config.asset.Type, err)
	}
	return &VMClarityInitiator{
		targetInfoType:     targetInfoType,
		scanConfigFamilies: cliutils.ConvertScanFamiliesConfigToAPIModel(config.fmConfig),
		Config:             config,
	}, nil
}

func getTargetInfoType(inputType string) (string, error) {
	switch inputType {
	case "dir", "DIR":
		return "DIRInfo", nil
	case "vm", "VM":
		return "VMInfo", nil
	default:
		return "", errors.New("unsupported target type")
	}
}

// initResults creates the necessary objects for exporting the results into the VMClarity backend,
// these objects are `target`, `scan`, `targetScanResult`.
// The function is returns the scanID and scanResultID that required for the export.
func (i *VMClarityInitiator) initResults(ctx context.Context) (string, string, error) {
	targetID, err := i.createTarget(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to init target: %v", err)
	}
	scanID, err := i.createScan(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to init scan: %v", err)
	}
	scanResultID, err := i.createScanResult(ctx, targetID, scanID)
	if err != nil {
		return scanID, "", fmt.Errorf("failed to init scan result: %v", err)
	}

	return scanID, scanResultID, nil
}

// nolint:cyclop
func (i *VMClarityInitiator) createTarget(ctx context.Context) (string, error) {
	// Now we are support only directory and vm input in the standalone mode
	info := models.TargetType{}
	switch i.targetInfoType {
	case "DIRInfo":
		hostName, err := os.Hostname()
		if err != nil {
			return "", fmt.Errorf("failed to get hostname: %v", err)
		}
		absPath, err := filepath.Abs(i.input)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path of %s: %v", i.input, err)
		}
		if i.asset.Location != "" {
			hostName = i.asset.Location
		}
		err = info.FromDirInfo(models.DirInfo{
			DirName:  utils.PointerTo(absPath),
			Location: utils.PointerTo(hostName),
		})
		if err != nil {
			return "", fmt.Errorf("failed to create DirInfo: %v", err)
		}
	case "VMInfo":
		// TODO(pebalogh) now we are supporting AWS cloud provider only
		vmInfoProvider := vminfoprovider.CreateNewAWSInfoProvider()
		// Get VM info from vm.
		instanceID, location, err := vmInfoProvider.GetVMInfo()
		// If the asset location and instanceID are set in config use them instead.
		if i.asset.Location != "" {
			location = i.asset.Location
		}
		if i.asset.InstanceID != "" {
			instanceID = i.asset.InstanceID
		}
		if err != nil {
			return "", fmt.Errorf("failed to get VMInfo: %v", err)
		}
		err = info.FromVMInfo(models.VMInfo{
			InstanceID: instanceID,
			Location:   location,
		})
		if err != nil {
			return "", fmt.Errorf("failed to create VMInfo: %v", err)
		}
	default:
		return "", errors.New("unsupported target type")
	}

	createdTarget, err := i.client.PostTarget(ctx, models.Target{TargetInfo: &info})
	if err != nil {
		var conErr backendclient.TargetConflictError
		if errors.As(err, &conErr) {
			logrus.Infof("Target already exist. target id=%v.", *conErr.ConflictingTarget.Id)
			return *conErr.ConflictingTarget.Id, nil
		}
		return "", fmt.Errorf("failed to post target: %v", err)
	}
	return *createdTarget.Id, nil
}

func (i *VMClarityInitiator) createScan(ctx context.Context) (string, error) {
	now := time.Now()
	scan := &models.Scan{
		// Scan config relationship is not set in standalone mode
		// to avoid uniqueness check of a scan
		ScanConfigSnapshot: &models.ScanConfigSnapshot{
			Name:               utils.PointerTo(i.scanConfigName),
			ScanFamiliesConfig: i.scanConfigFamilies,
		},
		StartTime: &now,
		Summary:   createInitScanSummary(),
	}

	createdScan, err := i.client.PostScan(ctx, *scan)
	if err != nil {
		return "", fmt.Errorf("failed to post scan: %v", err)
	}

	return *createdScan.Id, nil
}

func (i *VMClarityInitiator) createScanResult(ctx context.Context, targetID, scanID string) (string, error) {
	scanResult := models.TargetScanResult{
		Summary: createInitScanResultSummary(),
		Scan: &models.ScanRelationship{
			Id: scanID,
		},
		Target: &models.TargetRelationship{
			Id: targetID,
		},
	}
	createdScanResult, err := i.client.PostScanResult(ctx, scanResult)
	if err != nil {
		return "", fmt.Errorf("failed to post scan result: %v", err)
	}
	return *createdScanResult.Id, nil
}

func createInitScanResultSummary() *models.ScanFindingsSummary {
	return &models.ScanFindingsSummary{
		TotalExploits:          utils.PointerTo(0),
		TotalMalware:           utils.PointerTo(0),
		TotalMisconfigurations: utils.PointerTo(0),
		TotalPackages:          utils.PointerTo(0),
		TotalRootkits:          utils.PointerTo(0),
		TotalSecrets:           utils.PointerTo(0),
		TotalVulnerabilities: &models.VulnerabilityScanSummary{
			TotalCriticalVulnerabilities:   utils.PointerTo(0),
			TotalHighVulnerabilities:       utils.PointerTo(0),
			TotalMediumVulnerabilities:     utils.PointerTo(0),
			TotalLowVulnerabilities:        utils.PointerTo(0),
			TotalNegligibleVulnerabilities: utils.PointerTo(0),
		},
	}
}

func createInitScanSummary() *models.ScanSummary {
	return &models.ScanSummary{
		JobsCompleted:          utils.PointerTo(0),
		JobsLeftToRun:          utils.PointerTo(1),
		TotalExploits:          utils.PointerTo(0),
		TotalMalware:           utils.PointerTo(0),
		TotalMisconfigurations: utils.PointerTo(0),
		TotalPackages:          utils.PointerTo(0),
		TotalRootkits:          utils.PointerTo(0),
		TotalSecrets:           utils.PointerTo(0),
		TotalVulnerabilities: &models.VulnerabilityScanSummary{
			TotalCriticalVulnerabilities:   utils.PointerTo(0),
			TotalHighVulnerabilities:       utils.PointerTo(0),
			TotalMediumVulnerabilities:     utils.PointerTo(0),
			TotalLowVulnerabilities:        utils.PointerTo(0),
			TotalNegligibleVulnerabilities: utils.PointerTo(0),
		},
	}
}
