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
	"github.com/openclarity/vmclarity/cli/pkg/cicd/vminfoprovider"
	cliutils "github.com/openclarity/vmclarity/cli/pkg/utils"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/families"
)

type VMClarityInitiator struct {
	client             *backendclient.BackendClient
	input              string
	targetInfoType     string
	scanConfigName     string
	scanConfigID       string
	scanConfigFamilies *models.ScanFamiliesConfig
}

func NewInitiator(
	client *backendclient.BackendClient,
	config *families.Config,
	scanConfigID,
	scanConfigName,
	input,
	inputType string,
) (*VMClarityInitiator, error) {
	if client == nil {
		return nil, errors.New("backend client must not be nil")
	}
	targetInfoType, err := getTargetInfoType(inputType)
	if err != nil {
		return nil, fmt.Errorf("failed to get target type by inputType=%s: %v", inputType, err)
	}
	return &VMClarityInitiator{
		client:             client,
		input:              input,
		targetInfoType:     targetInfoType,
		scanConfigName:     scanConfigName,
		scanConfigID:       scanConfigID,
		scanConfigFamilies: cliutils.ConvertScanFamiliesConfigToAPIModel(config),
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

func (i *VMClarityInitiator) InitResults(ctx context.Context) (string, string, error) {
	targetID, err := i.createTarget(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to init results: %v", err)
	}
	scanID, err := i.createScan(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to init results: %v", err)
	}
	scanResultID, err := i.createScanResult(ctx, targetID, scanID)
	if err != nil {
		return scanID, "", fmt.Errorf("failed to init results: %v", err)
	}

	return scanID, scanResultID, nil
}

// nolint:cyclop
func (i *VMClarityInitiator) createTarget(ctx context.Context) (string, error) {
	// Now we are support only directory and vm input in the CI/CD mode
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
		instanceID, location, err := vmInfoProvider.GetVMInfo()
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
		// Scan config relationship is not set in CI/CD mode
		// to avoid uniqueness check of a scan
		ScanConfigSnapshot: &models.ScanConfigData{
			Name:               utils.PointerTo(i.scanConfigName),
			ScanFamiliesConfig: i.scanConfigFamilies,
		},
		StartTime: &now,
		Summary:   createInitScanSummary(),
	}

	createdScan, err := i.client.PostScan(ctx, *scan)
	if err != nil {
		var conErr backendclient.ScanConflictError
		if errors.As(err, &conErr) {
			logrus.Infof("Scan already exist. scan id=%v.", *conErr.ConflictingScan.Id)
			return *conErr.ConflictingScan.Id, nil
		}
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
		var conErr backendclient.ScanResultConflictError
		if errors.As(err, &conErr) {
			logrus.Infof("Scan results already exist. scan result id=%v.", *conErr.ConflictingScanResult.Id)
			return *conErr.ConflictingScanResult.Id, nil
		}
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
