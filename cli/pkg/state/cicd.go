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

package state

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	cicdinitiator "github.com/openclarity/vmclarity/cli/pkg/cicd/initiator"
	cicdupdater "github.com/openclarity/vmclarity/cli/pkg/cicd/updater"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type CICDState struct {
	client          *backendclient.BackendClient
	scanResultID    models.ScanResultID
	initiatorConfig cicdinitiator.Config
	updater         cicdupdater.Updater
}

func (c *CICDState) WaitForVolumeAttachment(context.Context) error {
	return nil
}

func (c *CICDState) MarkInProgress(ctx context.Context) error {
	log.Info("Scanning is in progress CICD")
	var err error
	var scanID string
	if c.scanResultID == "" {
		scanID, c.scanResultID, err = cicdinitiator.InitResults(ctx, c.initiatorConfig)
		if err != nil {
			return fmt.Errorf("failed to init scan result: %w", err)
		}
	}
	u, err := cicdupdater.NewVMClarityUpdater(c.client, scanID, c.scanResultID)
	if err != nil {
		return fmt.Errorf("failed to create VMClarity updater: %w", err)
	}
	// If the scanResultID is defined by the user we get the scan ID from it.
	if err := u.SetScanIDIfNeeded(ctx); err != nil {
		return fmt.Errorf("failed to set scan ID: %w", err)
	}
	c.updater = u

	return nil
}

// nolint:cyclop
func (c *CICDState) MarkDone(ctx context.Context, errors []error) error {
	log.Info("Scanning is done CICD")
	scanResult, err := c.client.GetScanResult(ctx, c.scanResultID, models.GetScanResultsScanResultIDParams{})
	if err != nil {
		return fmt.Errorf("failed to get scan result: %w", err)
	}

	if scanResult.Status == nil {
		scanResult.Status = &models.TargetScanStatus{}
	}
	if scanResult.Status.General == nil {
		scanResult.Status.General = &models.TargetScanState{}
	}

	state := models.DONE
	scanResult.Status.General.State = &state
	scanResult.Status.General.LastTransitionTime = utils.PointerTo(time.Now())

	// If we had any errors running the family or exporting results add it
	// to the general errors
	if len(errors) > 0 {
		var errorStrs []string
		// Pull the errors list out so that we can append to it (if there are
		// any errors at this point I would have hoped the orcestrator wouldn't
		// have spawned the VM) but we never know.
		if scanResult.Status.General.Errors != nil {
			errorStrs = *scanResult.Status.General.Errors
		}
		for _, err := range errors {
			if err != nil {
				errorStrs = append(errorStrs, err.Error())
			}
		}
		if len(errorStrs) > 0 {
			scanResult.Status.General.Errors = &errorStrs
		}
	}

	err = c.client.PatchScanResult(ctx, scanResult, c.scanResultID)
	if err != nil {
		return fmt.Errorf("failed to patch scan result: %w", err)
	}

	// In CI/CD mode, the scan objects needs to be updated in order to calculate ScanSummary,
	// update scan state end endTime of the scan.
	if err := c.updater.UpdateScanStateAndSummary(ctx); err != nil {
		return fmt.Errorf("failed to udate scan: %v", err)
	}

	return nil
}

func (c *CICDState) IsAborted(context.Context) (bool, error) {
	return false, nil
}

func (c *CICDState) GetScanResultID() string {
	return c.scanResultID
}

func NewCICDState(
	client *backendclient.BackendClient,
	scanResultID models.ScanResultID,
	cicdinitiatorConfig cicdinitiator.Config,
) (*CICDState, error) {
	return &CICDState{
		client:          client,
		scanResultID:    scanResultID,
		initiatorConfig: cicdinitiatorConfig,
	}, nil
}
