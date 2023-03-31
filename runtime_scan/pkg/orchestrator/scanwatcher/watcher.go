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

package scanwatcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/orchestrator/common"
	runtimeScanUtils "github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultPollInterval     = time.Minute
	DefaultReconcileTimeout = time.Minute
)

type ScanQueue = common.Queue[ScanReconcileEvent]
type ScanPoller = common.Poller[ScanReconcileEvent]
type ScanReconciler = common.Reconciler[ScanReconcileEvent]

type ScanEventKind int8

const (
	Aborted ScanEventKind = iota
)

type ScanReconcileEvent struct {
	Kind ScanEventKind
	Id   string
}

type Config struct {
	Backend          *backendclient.BackendClient
	PollPeriod       time.Duration
	ReconcileTimeout time.Duration
}

func New(c Config) *Watcher {
	logger := log.WithFields(log.Fields{"controller": "ScanWatcher"})
	return &Watcher{
		logger,
		c.Backend,
		c.PollPeriod,
		c.ReconcileTimeout,
	}
}

type Watcher struct {
	logger           *log.Entry
	client           *backendclient.BackendClient
	pollPeriod       time.Duration
	reconcileTimeout time.Duration
}

func (w *Watcher) Start(ctx context.Context) {
	queue := common.NewQueue[ScanReconcileEvent]()

	poller := &ScanPoller{
		Logger:     w.logger,
		PollPeriod: w.pollPeriod,
		Queue:      queue,
		GetItems:   w.GetItems,
	}
	poller.Start(ctx)

	reconciler := &ScanReconciler{
		ReconcileTimeout:  w.reconcileTimeout,
		Queue:             queue,
		ReconcileFunction: w.Reconcile,
	}
	reconciler.Start(ctx)
}

func (w *Watcher) GetItems(ctx context.Context) ([]ScanReconcileEvent, error) {
	scans, err := w.getScansByState(ctx, models.ScanStateAborted)
	if err != nil || scans.Items == nil || len(*scans.Items) <= 0 {
		return nil, err
	}

	count := len(*scans.Items)
	r := make([]ScanReconcileEvent, count)
	for i, scan := range *scans.Items {
		r[i] = ScanReconcileEvent{
			Kind: Aborted,
			Id:   *scan.Id,
		}
	}

	return r, nil
}

func (w *Watcher) Reconcile(ctx context.Context, s ScanReconcileEvent) error {
	w.logger.Infof("reconciling scan event: %v", s)
	switch {
	case s.Kind == Aborted:
		return w.reconcileAborted(ctx, s)
	}

	return nil
}

func (w *Watcher) getScansByState(ctx context.Context, s models.ScanState) (models.Scans, error) {
	filter := fmt.Sprintf("state eq '%s'", s)
	selector := "id,state,stateReason,targetIDs"
	params := models.GetScansParams{
		Filter: &filter,
		Select: &selector,
	}
	scans, err := w.client.GetScans(ctx, params)
	if err != nil {
		err = fmt.Errorf("getting Scan(s) by their state failed: %v", err)
	}
	if scans == nil {
		scans = &models.Scans{}
	}

	return *scans, err
}

func (w *Watcher) getScanResultsByScanID(ctx context.Context, id string) (models.TargetScanResults, error) {
	filter := fmt.Sprintf("scan/id eq '%s'", id)
	selector := "id,scan,status,target"
	params := models.GetScanResultsParams{
		Filter: &filter,
		Select: &selector,
	}

	scanResutls, err := w.client.GetScanResults(ctx, params)
	if err != nil {
		err = fmt.Errorf("getting ScanResult(s) by Scan ID failed: %v", err)
	}

	return scanResutls, err
}

func (w *Watcher) reconcileAborted(ctx context.Context, s ScanReconcileEvent) error {
	scanResults, err := w.getScanResultsByScanID(ctx, s.Id)
	if err != nil {
		return err
	}

	if scanResults.Items == nil || len(*scanResults.Items) <= 0 {
		w.logger.Debug("nothing to reconcile")
		return nil
	}

	var wg sync.WaitGroup
	for _, scanResult := range *scanResults.Items {
		if scanResult.Id == nil {
			continue
		}
		id := *scanResult.Id

		wg.Add(1)
		go func() {
			defer wg.Done()
			sr := models.TargetScanResult{
				Status: &models.TargetScanStatus{
					General: &models.TargetScanState{
						State: runtimeScanUtils.PointerTo(models.ABORTED),
					},
				},
			}

			err := w.client.PatchScanResult(ctx, sr, id)
			if err != nil {
				w.logger.Errorf("failed to patch ScanResult with id: %s", id)
				return
			}
		}()
	}
	wg.Wait()

	return nil
}
