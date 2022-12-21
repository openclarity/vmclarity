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

package scanner

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

type Scanner struct {
	instanceIDToScanData map[string]*scanData
	progress             types.ScanProgress
	scanConfig           *_config.ScanConfig
	killSignal           chan bool
	providerClient       provider.Client
	logFields            log.Fields
	backendClient        *client.ClientWithResponses

	region string

	sync.Mutex
}

type scanData struct {
	instance  types.Instance
	scanUUID  string
	success   bool // Needed for deletion policy in case we want to access the logs
	timeout   bool
	completed bool
}

func CreateScanner(config *_config.Config, providerClient provider.Client, backendClient *client.ClientWithResponses) *Scanner {
	s := &Scanner{
		progress: types.ScanProgress{
			State: types.Idle,
		},
		killSignal:     make(chan bool),
		providerClient: providerClient,
		logFields:      log.Fields{"scanner id": uuid.NewV4().String()},
		region:         config.Region,
		Mutex:          sync.Mutex{},
		backendClient:  backendClient,
	}

	return s
}

// initScan Calculate properties of scan targets
// nolint:cyclop,unparam
func (s *Scanner) initScan(ctx context.Context, instances []types.Instance) error {

	scanID, err := s.createScan(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a scan: %v", err)
	}
	instanceIDToScanData := make(map[string]*scanData)

	// Populate the instance to scanData map
	for _, instance := range instances {
		targetID := instance.GetID()
		if err := s.createTargetIfNotExist(ctx, targetID, instance.GetLocation()); err != nil {
			return fmt.Errorf("failed to create target. targetId=%v: %v", targetID, err)
		}
		if err := s.createInitScanStatus(ctx, scanID, targetID); err != nil {
			log.Errorf("Failed to create an init scan result. instance id=%v, scan id=%v: %v", instance.GetID(), scanID, err)
			continue
		}
		instanceIDToScanData[instance.GetID()] = &scanData{
			instance:  instance,
			scanUUID:  scanID,
			success:   false,
			completed: false,
			timeout:   false,
		}
	}

	s.instanceIDToScanData = instanceIDToScanData
	s.progress.InstancesToScan = uint32(len(instanceIDToScanData))

	log.WithFields(s.logFields).Infof("Total %d unique instances to scan", s.progress.InstancesToScan)

	return nil
}

func (s *Scanner) Scan(ctx context.Context, scanConfig *_config.ScanConfig, instances []types.Instance, scanDone chan struct{}) error {
	s.Lock()
	defer s.Unlock()

	s.scanConfig = scanConfig

	log.WithFields(s.logFields).Infof("Start scanning...")

	s.progress.State = types.ScanInit

	err := s.initScan(ctx, instances)
	if err != nil {
		return fmt.Errorf("failed to init scan: %v", err)
	}

	if s.progress.InstancesToScan == 0 {
		log.WithFields(s.logFields).Info("Nothing to scan")
		s.progress.SetStatus(types.NothingToScan)
		nonBlockingNotification(scanDone)
		return nil
	}

	s.progress.SetStatus(types.Scanning)
	go func() {
		s.jobBatchManagement(ctx, scanDone)

		s.Lock()
		s.progress.SetStatus(types.DoneScanning)
		s.Unlock()
	}()

	return nil
}

func (s *Scanner) GetScanStatus(ctx context.Context, data *scanData) (*models.TargetScanStatus, error) {
	resp, err := s.backendClient.GetScansScanIDTargetsTargetIDScanStatusWithResponse(ctx, data.scanUUID, data.instance.GetID())
	if err != nil {
		return nil, fmt.Errorf("failed to get scan status: %v", err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get scan status: status code=%v", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to get scan status: empty body")
	}

	return resp.JSON200, nil
}

func (s *Scanner) SetScanStatusCompletionError(ctx context.Context, data *scanData, errMsg string) error {
	// Get the status and set the completion error
	status, err := s.GetScanStatus(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to get scan status: %v", err)
	}
	var errors []string
	if status.General.Errors != nil {
		errors = *status.General.Errors
	}
	errors = append(errors, errMsg)
	status.General.Errors = &errors
	done := models.DONE
	status.General.State = &done

	// Update the status
	resp, err := s.backendClient.PutScansScanIDTargetsTargetIDScanStatusWithResponse(ctx, data.scanUUID, data.instance.GetID(), *status)
	if err != nil {
		return fmt.Errorf("failed to set scan status: %v", err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set scan status: status code=%v", resp.HTTPResponse.StatusCode)
	}

	return nil
}

func (s *Scanner) ScanProgress() types.ScanProgress {
	return types.ScanProgress{
		InstancesToScan:          s.progress.InstancesToScan,
		InstancesStartedToScan:   atomic.LoadUint32(&s.progress.InstancesStartedToScan),
		InstancesCompletedToScan: atomic.LoadUint32(&s.progress.InstancesCompletedToScan),
		State:                    s.progress.State,
	}
}

func (s *Scanner) Clear() {
	s.Lock()
	defer s.Unlock()

	log.WithFields(s.logFields).Infof("Clearing...")
	close(s.killSignal)
}

func (s *Scanner) createScan(ctx context.Context) (string, error) {
	startTime := time.Now()
	scan := models.Scan{
		EndTime:   nil,
		Id:        nil,
		StartTime: &startTime,
	}
	resp, err := s.backendClient.PostScansWithResponse(ctx, scan)
	if err != nil {
		return "", fmt.Errorf("failed to post scan: %v", err)
	}
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create a scan: status code=%v", resp.HTTPResponse.StatusCode)
	}
	if resp.JSON409 != nil {
		log.Warnf("Scan already exists. scanID=%v", resp.JSON409.Id)
		if resp.JSON409.Id == nil {
			return "", fmt.Errorf("scan already exists but has no ID")
		}
		return *resp.JSON409.Id, nil
	}
	if resp.JSON201 == nil {
		return "", fmt.Errorf("failed to create a scan: empty body")
	}
	if resp.JSON201.Id == nil {
		return "", fmt.Errorf("scan has no ID")
	}

	return *resp.JSON409.Id, nil
}
