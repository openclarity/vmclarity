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
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

// run jobs.
func (s *Scanner) jobBatchManagement(scanDone chan struct{}) {
	s.Lock()
	instanceIDToScanData := s.instanceIDToScanData
	numberOfWorkers := s.scanConfig.MaxScanParallelism
	instancesStartedToScan := &s.progress.InstancesStartedToScan
	instancesCompletedToScan := &s.progress.InstancesCompletedToScan
	s.Unlock()

	// queue of scan data
	q := make(chan *scanData)
	// done channel takes the result of the job
	done := make(chan bool)

	fullScanDone := make(chan bool)

	// spawn workers
	for i := 0; i < numberOfWorkers; i++ {
		go s.worker(q, i, done, s.killSignal)
	}

	// wait until scan of all images is done - non blocking. once all done, notify on fullScanDone chan
	go func() {
		for c := 0; c < len(instanceIDToScanData); c++ {
			select {
			case <-done:
				atomic.AddUint32(instancesCompletedToScan, 1)
			case <-s.killSignal:
				log.WithFields(s.logFields).Debugf("Scan process was canceled - stop waiting for finished jobs")
				return
			}
		}

		fullScanDone <- true
	}()

	// send all scan data on scan data queue, for workers to pick it up.
	for _, data := range instanceIDToScanData {
		go func(data *scanData, ks chan bool) {
			select {
			case q <- data:
				atomic.AddUint32(instancesStartedToScan, 1)
			case <-ks:
				log.WithFields(s.logFields).Debugf("Scan process was canceled. instanceID=%v, scanUUID=%v", data.instance.ID, data.scanUUID)
				return
			}
		}(data, s.killSignal)
	}

	// wait for killSignal or fullScanDone
	select {
	case <-s.killSignal:
		log.WithFields(s.logFields).Info("Scan process was canceled")
	case <-fullScanDone:
		log.WithFields(s.logFields).Infof("All jobs has finished")
		// Nonblocking notification of a finished scan
		nonBlockingNotification(scanDone)
	}
}

// worker waits for data on the queue, runs a scan job and waits for results from that scan job. Upon completion, done is notified to the caller.
func (s *Scanner) worker(queue chan *scanData, workNumber int, done, ks chan bool) {
	for {
		select {
		case data := <-queue:
			job, err := s.runJob(data)
			if err != nil {
				errMsg := fmt.Errorf("failed to run job: %v", err)
				log.WithFields(s.logFields).Error(errMsg)
				s.Lock()
				data.success = false
				data.scanErr = &types.ScanError{
					ErrMsg:    err.Error(),
					ErrType:   string(types.JobRun),
					ErrSource: types.ScanErrSourceJob,
				}
				data.completed = true
				s.Unlock()
			} else {
				s.waitForResult(data, ks)
			}

			s.deleteJobIfNeeded(&job, data.success, data.completed)

			select {
			case done <- true:
			case <-ks:
				log.WithFields(s.logFields).Infof("Image scan was canceled. imageID=%v", data.instance.ID)
			}
		case <-ks:
			log.WithFields(s.logFields).Debugf("worker #%v halted", workNumber)
			return
		}
	}
}

func (s *Scanner) waitForResult(data *scanData, ks chan bool) {
	log.WithFields(s.logFields).Infof("Waiting for result. instanceID=%+v", data.instance.ID)
	ticker := time.NewTicker(s.scanConfig.JobResultTimeout)
	select {
	case <-data.resultChan:
		log.WithFields(s.logFields).Infof("Instance scanned result has arrived. instanceID=%v", data.instance.ID)
	case <-ticker.C:
		errMsg := fmt.Errorf("job has timed out. instanceID=%v", data.instance.ID)
		log.WithFields(s.logFields).Warn(errMsg)
		s.Lock()
		data.success = false
		data.scanErr = &types.ScanError{
			ErrMsg:    errMsg.Error(),
			ErrType:   string(types.JobTimeout),
			ErrSource: types.ScanErrSourceJob,
		}
		data.timeout = true
		data.completed = true
		s.Unlock()
	case <-ks:
		log.WithFields(s.logFields).Infof("Instance scan was canceled. instanceID=%v", data.instance.ID)
	}
}

func (s *Scanner) runJob(data *scanData) (types.Job, error) {
	rootVolume, err := s.providerClient.GetInstanceRootVolume(data.instance)
	if err != nil {
		return types.Job{}, fmt.Errorf("failed to get instance root volume. instance id=%v: %v", data.instance.ID, err)
	}

	// create a snapshot of the root volume
	srcSnapshot, err := s.providerClient.CreateSnapshot(rootVolume)
	if err != nil {
		return types.Job{}, fmt.Errorf("failed to create snapshot: %v", err)
	}
	if err := s.providerClient.WaitForSnapshotReady(srcSnapshot); err != nil {
		return types.Job{}, fmt.Errorf("failed to wait for snapshot to be ready: %v", err)
	}

	//copy the snapshot to the scanner region
	// TODO check if scanner region is same as snapshot region?
	cpySnapshot, err := s.providerClient.CopySnapshot(srcSnapshot, s.region)
	if err != nil {
		return types.Job{}, fmt.Errorf("failed to copy snapshot: %v", err)
	}
	if err := s.providerClient.WaitForSnapshotReady(cpySnapshot); err != nil {
		return types.Job{}, fmt.Errorf("failed to wait for snapshot to be ready: %v", err)
	}

	// create the scanner job (vm) with a boot script
	launchedInstance, err := s.providerClient.LaunchInstance(s.jobAMI, s.deviceName, cpySnapshot)
	if err != nil {
		return types.Job{}, fmt.Errorf("failed to launch instance: %v", err)
	}
	if err := s.providerClient.WaitForInstanceReady(launchedInstance); err != nil {
		return types.Job{}, fmt.Errorf("failed to wait for instance to be ready: %v", err)
	}

	return types.Job{
		Instance:    launchedInstance,
		SrcSnapshot: srcSnapshot,
		DstSnapshot: cpySnapshot,
	}, nil
}

func (s *Scanner) deleteJobIfNeeded(job *types.Job, isSuccessfulJob, isCompletedJob bool) {
	if job == nil {
		return
	}

	// delete uncompleted jobs - scan process was canceled
	if !isCompletedJob {
		s.deleteJob(job)
		return
	}

	switch s.scanConfig.DeleteJobPolicy {
	case config.DeleteJobPolicyNever:
		// do nothing
	case config.DeleteJobPolicyAll:
		s.deleteJob(job)
	case config.DeleteJobPolicySuccessful:
		if isSuccessfulJob {
			s.deleteJob(job)
		}
	}
}

func (s *Scanner) deleteJob(job *types.Job) {
	if err := s.providerClient.DeleteInstance(job.Instance); err != nil {
		log.Errorf("failed to delete instance: %v", err)
	}
	if err := s.providerClient.DeleteSnapshot(job.SrcSnapshot); err != nil {
		log.Errorf("failed to delete source snapshot: %v", err)
	}
	if err := s.providerClient.DeleteSnapshot(job.DstSnapshot); err != nil {
		log.Errorf("failed to delete dest snapshot: %v", err)
	}
}
