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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

type VMClarityInitiator struct {
	client         *backendclient.BackendClient
	input          string
	inputType      string
	scanConfigName string
	scanConfigID   string
}

func NewInitiator(
	client *backendclient.BackendClient,
	scanConfigID,
	scanConfigName,
	input,
	inputType string) (*VMClarityInitiator, error) {
	if client == nil {
		return nil, errors.New("backend client must not be nil")
	}
	return &VMClarityInitiator{
		client:         client,
		input:          input,
		inputType:      inputType,
		scanConfigName: scanConfigName,
		scanConfigID:   scanConfigID,
	}, nil
}

func (i *VMClarityInitiator) InitResults(ctx context.Context) (string, error) {
	targetID, err := i.createTarget(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to init results: %v", err)
	}
	scanID, err := i.createScan(ctx, targetID)
	if err != nil {
		return "", fmt.Errorf("failed to init results: %v", err)
	}
	scanResultID, err := i.createScanResult(ctx, targetID, scanID)
	if err != nil {
		return "", fmt.Errorf("failed to init results: %v", err)
	}

	return scanResultID, nil
}

func (i *VMClarityInitiator) createTarget(ctx context.Context) (string, error) {
	info := models.TargetType{}
	err := info.FromDirInfo(models.DirInfo{
		DirName: utils.PointerTo(i.input),
		// TODO what should be the location ???
		Location: utils.PointerTo(i.input),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create DirInfo: %v", err)
	}

	createdTarget, err := i.client.PostTarget(ctx, models.Target{
		TargetInfo: &info,
	})
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

func (i *VMClarityInitiator) createScan(ctx context.Context, targetID string) (string, error) {
	now := time.Now().UTC()
	scan := &models.Scan{
		ScanConfig: &models.ScanConfigRelationship{
			Id: i.scanConfigID,
		},
		ScanConfigSnapshot: &models.ScanConfigData{
			Name: utils.PointerTo(i.scanConfigName),
		},
		StartTime: &now,
		TargetIDs: utils.PointerTo([]string{targetID}),
	}
	var scanID string
	createdScan, err := i.client.PostScan(ctx, *scan)
	if err != nil {
		var conErr backendclient.ScanConflictError
		if errors.As(err, &conErr) {
			logrus.Infof("Scan already exist. scan id=%v.", *conErr.ConflictingScan.Id)
			scanID = *conErr.ConflictingScan.Id
		} else {
			return "", fmt.Errorf("failed to post scan: %v", err)
		}
	} else {
		scanID = *createdScan.Id
	}
	return scanID, nil
}

func (i *VMClarityInitiator) createScanResult(ctx context.Context, targetID, scanID string) (string, error) {
	scanResult := models.TargetScanResult{
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
