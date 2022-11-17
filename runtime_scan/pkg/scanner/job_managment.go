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
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
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

	// wait until scan of all instances is done - non blocking. once all done, notify on fullScanDone chan
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
				log.WithFields(s.logFields).Debugf("Scan process was canceled. instanceID=%v, scanUUID=%v", data.instance.GetID(), data.scanUUID)
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
			job, err := s.runJob(context.TODO(), data)
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

			s.deleteJobIfNeeded(context.TODO(), &job, data.success, data.completed)

			select {
			case done <- true:
			case <-ks:
				log.WithFields(s.logFields).Infof("Instance scan was canceled. instanceID=%v", data.instance.GetID())
			}
		case <-ks:
			log.WithFields(s.logFields).Debugf("worker #%v halted", workNumber)
			return
		}
	}
}

func (s *Scanner) waitForResult(data *scanData, ks chan bool) {
	log.WithFields(s.logFields).Infof("Waiting for result. instanceID=%+v", data.instance.GetID())
	ticker := time.NewTicker(s.scanConfig.JobResultTimeout)
	select {
	case <-data.resultChan:
		log.WithFields(s.logFields).Infof("Instance scanned result has arrived. instanceID=%v", data.instance.GetID())
	case <-ticker.C:
		errMsg := fmt.Errorf("job has timed out. instanceID=%v", data.instance.GetID())
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
		log.WithFields(s.logFields).Infof("Instance scan was canceled. instanceID=%v", data.instance.GetID())
	}
}

func (s *Scanner) runJob(ctx context.Context, data *scanData) (provider.Job, error) {
	instanceToScan := data.instance

	volume, err := instanceToScan.GetRootVolume(ctx)
	if err != nil {
		return provider.Job{}, fmt.Errorf("failed to get root volume of an instance %v: %v", instanceToScan.GetID(), err)
	}

	snapshot, err := volume.TakeSnapshot(ctx)
	if err != nil {
		return provider.Job{}, fmt.Errorf("failed to take snapshot of a volume: %v", err)

	}
	if err := snapshot.WaitForReady(ctx); err != nil {
		return provider.Job{}, fmt.Errorf("failed to wait for snapshot %v ready: %v", snapshot.GetID(), err)

	}

	cpySnapshot, err := snapshot.Copy(ctx, s.region)
	if err != nil {
		return provider.Job{}, fmt.Errorf("failed to copy snapshot %v: %v", snapshot.GetID(), err)
	}

	if err := cpySnapshot.WaitForReady(ctx); err != nil {
		return provider.Job{}, fmt.Errorf("failed wait for snapshot %v ready: %v", cpySnapshot.GetID(), err)
	}

	i, err := s.providerClient.LaunchInstance(ctx, cpySnapshot)
	if err != nil {
		return provider.Job{}, fmt.Errorf("failed to launch a new instance: %v", err)
	}

	return provider.Job{
		Instance:    i,
		SrcSnapshot: snapshot,
		DstSnapshot: cpySnapshot,
	}, nil
}

func (s *Scanner) deleteJobIfNeeded(ctx context.Context, job *provider.Job, isSuccessfulJob, isCompletedJob bool) {
	if job == nil {
		return
	}

	// delete uncompleted jobs - scan process was canceled
	if !isCompletedJob {
		s.deleteJob(ctx, job)
		return
	}

	switch s.scanConfig.DeleteJobPolicy {
	case config.DeleteJobPolicyNever:
		// do nothing
	case config.DeleteJobPolicyAll:
		s.deleteJob(ctx, job)
	case config.DeleteJobPolicySuccessful:
		if isSuccessfulJob {
			s.deleteJob(ctx, job)
		}
	}
}

func (s *Scanner) deleteJob(ctx context.Context, job *provider.Job) {
	if err := job.Instance.Delete(ctx); err != nil {
		log.Errorf("failed to delete instance: %v", err)
	}
	if err := job.SrcSnapshot.Delete(ctx); err != nil {
		log.Errorf("failed to delete source snapshot: %v", err)
	}
	if err := job.DstSnapshot.Delete(ctx); err != nil {
		log.Errorf("failed to delete dest snapshot: %v", err)
	}
}
