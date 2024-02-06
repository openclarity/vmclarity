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

package scan

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/cli/pkg/utils"
	"github.com/openclarity/vmclarity/orchestrator/common"
	"github.com/openclarity/vmclarity/provider"
	"github.com/openclarity/vmclarity/utils/log"
)

type (
	ScanQueue      = common.Queue[ScanReconcileEvent]
	ScanPoller     = common.Poller[ScanReconcileEvent]
	ScanReconciler = common.Reconciler[ScanReconcileEvent]
)

func New(c Config) *Watcher {
	return &Watcher{
		backend:          c.Backend,
		provider:         c.Provider,
		pollPeriod:       c.PollPeriod,
		reconcileTimeout: c.ReconcileTimeout,
		scanTimeout:      c.ScanTimeout,
		queue:            common.NewQueue[ScanReconcileEvent](),
	}
}

type Watcher struct {
	backend          *client.BackendClient
	provider         provider.Provider
	pollPeriod       time.Duration
	reconcileTimeout time.Duration
	scanTimeout      time.Duration

	queue *ScanQueue
}

func (w *Watcher) Start(ctx context.Context) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("controller", "ScanWatcher")
	ctx = log.SetLoggerForContext(ctx, logger)

	poller := &ScanPoller{
		PollPeriod: w.pollPeriod,
		Queue:      w.queue,
		GetItems:   w.GetRunningScans,
	}
	poller.Start(ctx)

	reconciler := &ScanReconciler{
		ReconcileTimeout:  w.reconcileTimeout,
		Queue:             w.queue,
		ReconcileFunction: w.Reconcile,
	}
	reconciler.Start(ctx)
}

// nolint:cyclop
func (w *Watcher) GetRunningScans(ctx context.Context) ([]ScanReconcileEvent, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)
	logger.Debugf("Fetching running Scans")

	filter := fmt.Sprintf("status/state ne '%s' and status/state ne '%s'", types.ScanStatusStateDone, types.ScanStatusStateFailed)
	selector := "id"
	params := types.GetScansParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	}
	scans, err := w.backend.GetScans(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get running scans: %w", err)
	}

	switch {
	case scans.Items == nil && scans.Count == nil:
		return nil, fmt.Errorf("failed to fetch running Scans: invalid API response: %v", scans)
	case scans.Count != nil && *scans.Count <= 0:
		fallthrough
	case scans.Items != nil && len(*scans.Items) <= 0:
		return nil, nil
	}

	events := make([]ScanReconcileEvent, 0, *scans.Count)
	for _, scan := range *scans.Items {
		scanID, ok := scan.GetID()
		if !ok {
			logger.Warnf("Skipping to invalid Scan: ID is nil: %v", scan)
			continue
		}

		events = append(events, ScanReconcileEvent{
			ScanID: scanID,
		})
	}

	return events, nil
}

// nolint:cyclop
func (w *Watcher) Reconcile(ctx context.Context, event ScanReconcileEvent) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(event.ToFields())
	ctx = log.SetLoggerForContext(ctx, logger)

	params := types.GetScansScanIDParams{
		Expand: utils.PointerTo("scanConfig"),
	}
	scan, err := w.backend.GetScan(ctx, event.ScanID, params)
	if err != nil || scan == nil {
		return fmt.Errorf("failed to fetch Scan. ScanID=%s: %w", event.ScanID, err)
	}

	if scan.IsTimedOut(w.scanTimeout) {
		scan.Status = types.NewScanStatus(
			types.ScanStatusStateFailed,
			types.ScanStatusReasonTimeout,
			utils.PointerTo("Scan has timed out"),
		)

		err = w.backend.PatchScan(ctx, *scan.Id, &types.Scan{Status: scan.Status})
		if err != nil {
			return fmt.Errorf("failed to patch Scan. ScanID=%s: %w", event.ScanID, err)
		}
	}

	status, ok := scan.GetStatus()
	if !ok {
		return fmt.Errorf("failed to determine status of Scan. ScanID=%s", event.ScanID)
	}
	logger.Tracef("Reconciling Scan state: %s", status.State)

	switch status.State {
	case types.ScanStatusStatePending:
		if err = w.reconcilePending(ctx, scan); err != nil {
			return err
		}
	case types.ScanStatusStateDiscovered:
		if err = w.reconcileDiscovered(ctx, scan); err != nil {
			return err
		}
	case types.ScanStatusStateInProgress:
		if err = w.reconcileInProgress(ctx, scan); err != nil {
			return err
		}
	case types.ScanStatusStateAborted:
		if err = w.reconcileAborted(ctx, scan); err != nil {
			return err
		}
	case types.ScanStatusStateDone, types.ScanStatusStateFailed:
		logger.Debug("Reconciling Scan is skipped as it is already finished.")
		fallthrough
	default:
		return nil
	}

	return nil
}

func (w *Watcher) reconcilePending(ctx context.Context, scan *types.Scan) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scan == nil {
		return errors.New("invalid Scan: object is nil")
	}

	scanID, ok := scan.GetID()
	if !ok {
		return errors.New("invalid Scan: Id is nil")
	}

	// We don't want to scan terminated assets. These are historic assets
	// and we'll just get a not found error during the asset scan.
	assetFilter := "terminatedOn eq null"

	scope, ok := scan.GetScope()
	if !ok {
		return fmt.Errorf("invalid Scan: Scope is nil. ScanID=%s", scanID)
	}

	// If the scan has a scope configured, 'and' it with the check for
	// not terminated to make sure that we take both into account.
	if scope != "" {
		assetFilter = fmt.Sprintf("(%s) and (%s)", assetFilter, scope)
	}

	assets, err := w.backend.GetAssets(ctx, types.GetAssetsParams{
		Filter: &assetFilter,
		Select: utils.PointerTo("id"),
	})
	if err != nil {
		return fmt.Errorf("failed to discover Assets for Scan. ScanID=%s: %w", scanID, err)
	}

	numOfAssets := len(*assets.Items)

	if numOfAssets > 0 {
		assetIds := []string{}
		for _, asset := range *assets.Items {
			assetIds = append(assetIds, *asset.Id)
		}
		scan.AssetIDs = &assetIds
		scan.Status = types.NewScanStatus(
			types.ScanStatusStateDiscovered,
			types.ScanStatusReasonAssetsDiscovered,
			utils.PointerTo("Assets for Scan are successfully discovered"),
		)
	} else {
		scan.Status = types.NewScanStatus(
			types.ScanStatusStateDone,
			types.ScanStatusReasonNothingToScan,
			utils.PointerTo("No instances found in scope for Scan"),
		)
	}
	logger.Debugf("%d Asset(s) have been created for Scan", numOfAssets)

	scanPatch := &types.Scan{
		AssetIDs: scan.AssetIDs,
		Status:   scan.Status,
	}

	if err = w.backend.PatchScan(ctx, scanID, scanPatch); err != nil {
		return fmt.Errorf("failed to patch Scan. ScanID=%s: %w", scanID, err)
	}

	return nil
}

func (w *Watcher) reconcileDiscovered(ctx context.Context, scan *types.Scan) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scan == nil {
		return errors.New("invalid Scan: object is nil")
	}

	scanID, ok := scan.GetID()
	if !ok {
		return errors.New("invalid Scan: Id is nil")
	}

	if err := w.createAssetScansForScan(ctx, scan); err != nil {
		return fmt.Errorf("failed to creates AssetScan(s) for Scan. ScanID=%s: %w", scanID, err)
	}

	scanPatch := &types.Scan{
		Status: types.NewScanStatus(
			types.ScanStatusStateInProgress,
			types.ScanStatusReasonAssetScansRunning,
			nil,
		),
		Summary:  scan.Summary,
		AssetIDs: scan.AssetIDs,
	}
	err := w.backend.PatchScan(ctx, scanID, scanPatch)
	if err != nil {
		return fmt.Errorf("failed to update Scan. ScanID=%s: %w", scanID, err)
	}

	logger.Infof("Total %d unique assets for Scan", len(*scan.AssetIDs))

	return nil
}

func (w *Watcher) createAssetScansForScan(ctx context.Context, scan *types.Scan) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scan.AssetIDs == nil || *scan.AssetIDs == nil {
		return nil
	}
	numOfAssets := len(*scan.AssetIDs)

	errs := make(chan error, numOfAssets)
	var wg sync.WaitGroup
	for _, id := range *scan.AssetIDs {
		wg.Add(1)
		assetID := id
		go func() {
			defer wg.Done()

			err := w.createAssetScanForAsset(ctx, scan, assetID)
			if err != nil {
				logger.WithField("AssetID", assetID).Errorf("Failed to create AssetScan: %v", err)
				errs <- err

				return
			}
		}()
	}
	wg.Wait()
	close(errs)

	assetErrs := make([]error, 0, numOfAssets)
	for err := range errs {
		assetErrs = append(assetErrs, err)
	}
	numOfErrs := len(assetErrs)

	if numOfErrs > 0 {
		return fmt.Errorf("failed to create %d AssetScan(s) for Scan. ScanID=%s: %w", numOfErrs, *scan.Id, assetErrs[0])
	}

	scan.Summary.JobsLeftToRun = utils.PointerTo(numOfAssets)

	return nil
}

func (w *Watcher) createAssetScanForAsset(ctx context.Context, scan *types.Scan, assetID string) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	assetScanData, err := newAssetScanFromScan(scan, assetID)
	if err != nil {
		return fmt.Errorf("failed to generate new AssetScan for Scan. ScanID=%s, AssetID=%s: %w", *scan.Id, assetID, err)
	}

	_, err = w.backend.PostAssetScan(ctx, *assetScanData)
	if err != nil {
		var conErr client.AssetScanConflictError
		if errors.As(err, &conErr) {
			assetScanID := *conErr.ConflictingAssetScan.Id
			logger.WithField("AssetScanID", assetScanID).Debug("AssetScan already exist.")
			return nil
		}
		return fmt.Errorf("failed to post AssetScan to backend API: %w", err)
	}
	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileInProgress(ctx context.Context, scan *types.Scan) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scan == nil {
		return errors.New("invalid Scan: object is nil")
	}

	scanID, ok := scan.GetID()
	if !ok {
		return errors.New("invalid Scan: ID is nil")
	}

	// FIXME(chrisgacsal):a add pagination to API queries in poller/reconciler logic by using Top/Skip
	filter := fmt.Sprintf("scan/id eq '%s'", scanID)
	selector := "id,status,summary"
	assetScans, err := w.backend.GetAssetScans(ctx, types.GetAssetScansParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve AssetScans for Scan. ScanID=%s: %w", scanID, err)
	}

	if assetScans.Count == nil || assetScans.Items == nil {
		return fmt.Errorf("invalid response for getting AssetScans for Scan. ScanID=%s: Count and/or Items parameters are nil", scanID)
	}

	// Reset Scan Summary as it is going to be recalculated
	scan.Summary = newScanSummary()

	var failedAssetScans int
	for _, assetScan := range *assetScans.Items {
		assetScanID, ok := assetScan.GetID()
		if !ok {
			return errors.New("invalid AssetScan: ID is nil")
		}

		if err := updateScanSummaryFromAssetScan(scan, assetScan); err != nil {
			return fmt.Errorf("failed to update Scan Summary from AssetScan. ScanID=%s AssetScanID=%s: %w",
				scanID, assetScanID, err)
		}

		status, ok := assetScan.GetStatus()
		if !ok {
			return fmt.Errorf("status must not be nil for AssetScan. AssetScanID=%s", *assetScan.Id)
		}

		if status.State == types.AssetScanStatusStateFailed {
			failedAssetScans++
		}
	}
	logger.Tracef("Scan Summary updated. JobCompleted=%d JobLeftToRun=%d", *scan.Summary.JobsCompleted,
		*scan.Summary.JobsLeftToRun)

	if *scan.Summary.JobsLeftToRun <= 0 {
		message := utils.PointerTo(
			fmt.Sprintf(
				"%d succeeded, %d failed out of %d total asset scans",
				*assetScans.Count-failedAssetScans,
				failedAssetScans,
				*assetScans.Count,
			),
		)

		if failedAssetScans > 0 {
			scan.Status = types.NewScanStatus(
				types.ScanStatusStateFailed,
				types.ScanStatusReasonError,
				message,
			)
		} else {
			scan.Status = types.NewScanStatus(
				types.ScanStatusStateDone,
				types.ScanStatusReasonSuccess,
				message,
			)
		}
		scan.EndTime = utils.PointerTo(time.Now())
	}

	scanPatch := &types.Scan{
		Status:   scan.Status,
		Summary:  scan.Summary,
		EndTime:  scan.EndTime,
		AssetIDs: scan.AssetIDs,
	}
	err = w.backend.PatchScan(ctx, scanID, scanPatch)
	if err != nil {
		return fmt.Errorf("failed to patch Scan. ScanID=%s: %w", scanID, err)
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileAborted(ctx context.Context, scan *types.Scan) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scan == nil {
		return errors.New("invalid Scan: object is nil")
	}

	scanID, ok := scan.GetID()
	if !ok {
		return errors.New("invalid Scan: ID is nil")
	}

	filter := fmt.Sprintf("scan/id eq '%s' and status/state ne '%s' and status/state ne '%s' and status/state ne '%s'",
		scanID, types.AssetScanStatusStateAborted, types.AssetScanStatusStateDone, types.AssetScanStatusStateFailed)
	selector := "id,status"
	params := types.GetAssetScansParams{
		Filter: &filter,
		Select: &selector,
	}

	assetScans, err := w.backend.GetAssetScans(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to fetch AssetScan(s) for Scan. ScanID=%s: %w", scanID, err)
	}

	if assetScans.Items != nil && len(*assetScans.Items) > 0 {
		var reconciliationFailed bool
		var wg sync.WaitGroup

		for _, assetScan := range *assetScans.Items {
			if assetScan.Id == nil {
				continue
			}
			assetScanID := *assetScan.Id

			wg.Add(1)
			go func() {
				defer wg.Done()
				as := types.AssetScan{
					Status: types.NewAssetScanStatus(
						types.AssetScanStatusStateAborted,
						types.AssetScanStatusReasonCancellation,
						nil,
					),
				}

				err = w.backend.PatchAssetScan(ctx, as, assetScanID)
				if err != nil {
					logger.WithField("AssetScanID", assetScanID).Error("Failed to patch AssetScan")
					reconciliationFailed = true
					return
				}
			}()
		}
		wg.Wait()

		// NOTE: reconciliationFailed is used to track errors returned by patching AssetScans
		//       as setting the state of Scan to types.ScanStateFailed must be skipped in case
		//       even a single error occurred to allow reconciling re-running for this Scan.
		if reconciliationFailed {
			return errors.New("updating one or more AssetScans failed")
		}
	}

	scan.EndTime = utils.PointerTo(time.Now())
	scan.Status = types.NewScanStatus(
		types.ScanStatusStateFailed,
		types.ScanStatusReasonCancellation,
		utils.PointerTo("Scan has been aborted"),
	)

	scanPatch := &types.Scan{
		Status:   scan.Status,
		EndTime:  scan.EndTime,
		AssetIDs: scan.AssetIDs,
	}
	err = w.backend.PatchScan(ctx, scanID, scanPatch)
	if err != nil {
		return fmt.Errorf("failed to patch Scan. ScanID=%s: %w", scanID, err)
	}

	return nil
}