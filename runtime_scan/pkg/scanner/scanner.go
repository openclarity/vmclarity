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
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

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
	backendClient        *client.Client

	region string

	sync.Mutex
}

type scanData struct {
	instance    types.Instance
	scanUUID    string
	scanResults []string
	success     bool
	completed   bool
	timeout     bool
	scanErr     *types.ScanError
}

func CreateScanner(config *_config.Config, providerClient provider.Client, backendClient *client.Client) *Scanner {
	s := &Scanner{
		progress: types.ScanProgress{
			Status: types.Idle,
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
func (s *Scanner) initScan() error {
	instanceIDToScanData := make(map[string]*scanData)

	// Populate the instance to scanData map
	for _, instance := range s.scanConfig.Instances {
		instanceIDToScanData[instance.GetID()] = &scanData{
			instance:    instance,
			scanUUID:    uuid.NewV4().String(),
			scanResults: nil, // list of expected scan results get from scanner job config implement later
			success:     false,
			completed:   false,
			timeout:     false,
			scanErr:     nil,
		}
	}

	s.instanceIDToScanData = instanceIDToScanData
	s.progress.InstancesToScan = uint32(len(instanceIDToScanData))

	log.WithFields(s.logFields).Infof("Total %d unique instances to scan", s.progress.InstancesToScan)

	return nil
}

func (s *Scanner) Scan(scanConfig *_config.ScanConfig, scanDone chan struct{}) error {
	s.Lock()
	defer s.Unlock()

	s.scanConfig = scanConfig

	log.WithFields(s.logFields).Infof("Start scanning...")

	s.progress.Status = types.ScanInit

	if err := s.initScan(); err != nil {
		s.progress.SetStatus(types.ScanInitFailure)
		return fmt.Errorf("failed to initiate scan: %v", err)
	}

	if s.progress.InstancesToScan == 0 {
		log.WithFields(s.logFields).Info("Nothing to scan")
		s.progress.SetStatus(types.NothingToScan)
		nonBlockingNotification(scanDone)
		return nil
	}

	s.progress.SetStatus(types.Scanning)
	go func() {
		s.jobBatchManagement(scanDone)

		s.Lock()
		s.progress.SetStatus(types.DoneScanning)
		s.Unlock()
	}()

	return nil
}

func (s *Scanner) GetScanStatus(data *scanData) *types.InstanceScanResult {
	scanResult := &types.InstanceScanResult{
		Instance: data.instance,
		Success:  false,
		Status:   types.Scanning,
	}
	resp, err := s.backendClient.GetTargetsTargetIDScanResultsScanID(context.TODO(), data.instance.GetID(), data.scanUUID)
	if err != nil {
		log.WithFields(s.logFields).Errorf("Failed to get scan results %s for target: %s", data.scanUUID, data.instance.GetID())
		scanResult.Status = types.NothingToScan
		// TODO use map for scan errors in the case of scan types later
		scanResult.ScanError = &types.ScanError{
			ErrMsg:    err.Error(),
			ErrType:   string(types.JobRun),
			ErrSource: types.ScanErrSourceJob,
		}
		return scanResult
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.WithFields(s.logFields).Errorf("Failed to get scan results %s for target: %s", data.scanUUID, data.instance.GetID())
		scanResult.Status = types.NothingToScan
		// TODO use map for scan errors in the case of scan types later
		scanResult.ScanError = &types.ScanError{
			ErrMsg:    err.Error(),
			ErrType:   string(types.JobRun),
			ErrSource: types.ScanErrSourceJob,
		}
		return scanResult
	}
	var results models.ScanResults
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.WithFields(s.logFields).Errorf("Failed to decode scan results %s for target: %s", data.scanUUID, data.instance.GetID())
		scanResult.Status = types.NothingToScan
		// TODO use map for scan errors in the case of scan types later
		scanResult.ScanError = &types.ScanError{
			ErrMsg:    err.Error(),
			ErrType:   string(types.JobRun),
			ErrSource: types.ScanErrSourceJob,
		}
		return scanResult
	}
	if results.Sboms != nil {
		log.WithFields(s.logFields).Infof("Scan results %s exist for target %s.", data.scanUUID, data.instance.GetID())
		scanResult.Success = true
		scanResult.Status = types.DoneScanning
		return scanResult
	}

	log.WithFields(s.logFields).Infof("Scan results %s not exist for target %s. waiting for results...", data.scanUUID, data.instance.GetID())
	return scanResult
}

func (s *Scanner) ScanProgress() types.ScanProgress {
	return types.ScanProgress{
		InstancesToScan:          s.progress.InstancesToScan,
		InstancesStartedToScan:   atomic.LoadUint32(&s.progress.InstancesStartedToScan),
		InstancesCompletedToScan: atomic.LoadUint32(&s.progress.InstancesCompletedToScan),
		Status:                   s.progress.Status,
	}
}

func (s *Scanner) Clear() {
	s.Lock()
	defer s.Unlock()

	log.WithFields(s.logFields).Infof("Clearing...")
	close(s.killSignal)
}
