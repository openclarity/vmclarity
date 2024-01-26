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

package scanestimationwatcher

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/cli/pkg/utils"
	"github.com/openclarity/vmclarity/pkg/orchestrator/common"
	"github.com/openclarity/vmclarity/pkg/orchestrator/provider"
	"github.com/openclarity/vmclarity/utils/log"
)

type (
	ScanEstimationQueue      = common.Queue[ScanEstimationReconcileEvent]
	ScanEstimationPoller     = common.Poller[ScanEstimationReconcileEvent]
	ScanEstimationReconciler = common.Reconciler[ScanEstimationReconcileEvent]
)

func New(c Config) *Watcher {
	return &Watcher{
		backend:               c.Backend,
		provider:              c.Provider,
		pollPeriod:            c.PollPeriod,
		reconcileTimeout:      c.ReconcileTimeout,
		scanEstimationTimeout: c.ScanEstimationTimeout,
		queue:                 common.NewQueue[ScanEstimationReconcileEvent](),
	}
}

type Watcher struct {
	backend               *client.BackendClient
	provider              provider.Provider
	pollPeriod            time.Duration
	reconcileTimeout      time.Duration
	scanEstimationTimeout time.Duration

	queue *ScanEstimationQueue
}

func (w *Watcher) Start(ctx context.Context) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("controller", "ScanEstimationWatcher")
	ctx = log.SetLoggerForContext(ctx, logger)

	poller := &ScanEstimationPoller{
		PollPeriod: w.pollPeriod,
		Queue:      w.queue,
		GetItems:   w.GetScanEstimations,
	}
	poller.Start(ctx)

	reconciler := &ScanEstimationReconciler{
		ReconcileTimeout:  w.reconcileTimeout,
		Queue:             w.queue,
		ReconcileFunction: w.Reconcile,
	}
	reconciler.Start(ctx)
}

func (w *Watcher) GetScanEstimations(ctx context.Context) ([]ScanEstimationReconcileEvent, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)
	logger.Debugf("Fetching running ScanEstimations")

	filter := fmt.Sprintf("(status/state ne '%s' and status/state ne '%s') or (deleteAfter eq null or deleteAfter lt %s)",
		types.ScanEstimationStatusStateDone, types.ScanEstimationStatusStateFailed, time.Now().Format(time.RFC3339))
	selector := "id"
	params := types.GetScanEstimationsParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	}
	scanEstimations, err := w.backend.GetScanEstimations(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get running ScanEstimations: %w", err)
	}

	switch {
	case scanEstimations.Items == nil && scanEstimations.Count == nil:
		return nil, fmt.Errorf("failed to fetch running ScanEstimations: invalid API response: %v", scanEstimations)
	case scanEstimations.Count != nil && *scanEstimations.Count <= 0:
		fallthrough
	case scanEstimations.Items != nil && len(*scanEstimations.Items) <= 0:
		return nil, nil
	}

	events := make([]ScanEstimationReconcileEvent, 0, *scanEstimations.Count)
	for _, scanEstimation := range *scanEstimations.Items {
		scanEstimationID, ok := scanEstimation.GetID()
		if !ok {
			logger.Warnf("Skipping due to invalid ScanEstimation: ID is nil: %v", scanEstimation)
			continue
		}

		events = append(events, ScanEstimationReconcileEvent{
			ScanEstimationID: scanEstimationID,
		})
	}

	return events, nil
}

// nolint:cyclop
func (w *Watcher) Reconcile(ctx context.Context, event ScanEstimationReconcileEvent) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(event.ToFields())
	ctx = log.SetLoggerForContext(ctx, logger)

	params := types.GetScanEstimationsScanEstimationIDParams{}
	scanEstimation, err := w.backend.GetScanEstimation(ctx, event.ScanEstimationID, params)
	if err != nil || scanEstimation == nil {
		return fmt.Errorf("failed to fetch ScanEstimation. ScanEstimationID=%s: %w", event.ScanEstimationID, err)
	}

	if scanEstimation.IsTimedOut(w.scanEstimationTimeout) {
		err = w.backend.PatchScanEstimation(ctx, *scanEstimation.Id, &types.ScanEstimation{
			Status: types.NewScanEstimationStatus(
				types.ScanEstimationStatusStateFailed,
				types.ScanEstimationStatusReasonTimeout,
				utils.PointerTo("ScanEstimation has been timed out"),
			),
		})
		if err != nil {
			return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", event.ScanEstimationID, err)
		}
	}

	status, ok := scanEstimation.GetStatus()
	if !ok {
		if err = w.reconcileNoState(ctx, scanEstimation); err != nil {
			return err
		}
		return nil
	}

	logger.Tracef("Reconciling ScanEstimation state: %s", status.State)

	switch status.State {
	case types.ScanEstimationStatusStatePending:
		if err = w.reconcilePending(ctx, scanEstimation); err != nil {
			return err
		}
	case types.ScanEstimationStatusStateDiscovered:
		if err = w.reconcileDiscovered(ctx, scanEstimation); err != nil {
			return err
		}
	case types.ScanEstimationStatusStateInProgress:
		if err = w.reconcileInProgress(ctx, scanEstimation); err != nil {
			return err
		}
	case types.ScanEstimationStatusStateAborted:
		if err = w.reconcileAborted(ctx, scanEstimation); err != nil {
			return err
		}
	case types.ScanEstimationStatusStateDone, types.ScanEstimationStatusStateFailed:
		if err = w.reconcileDone(ctx, scanEstimation); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func (w *Watcher) reconcileDone(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	if scanEstimation.EndTime == nil {
		scanEstimation.EndTime = utils.PointerTo(time.Now())
	}
	if scanEstimation.TTLSecondsAfterFinished == nil {
		scanEstimation.TTLSecondsAfterFinished = utils.PointerTo(DefaultScanEstimationTTLSeconds)
	}

	endTime := *scanEstimation.EndTime
	ttl := *scanEstimation.TTLSecondsAfterFinished

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	timeNow := time.Now()

	if scanEstimation.DeleteAfter == nil {
		scanEstimation.DeleteAfter = utils.PointerTo(endTime.Add(time.Duration(ttl) * time.Second))
		// if delete time has already pass, no need to patch the object, just delete it.
		if !timeNow.After(*scanEstimation.DeleteAfter) {
			scanEstimationPatch := types.ScanEstimation{
				DeleteAfter:             scanEstimation.DeleteAfter,
				EndTime:                 scanEstimation.EndTime,
				TTLSecondsAfterFinished: scanEstimation.TTLSecondsAfterFinished,
			}
			err := w.backend.PatchScanEstimation(ctx, scanEstimationID, &scanEstimationPatch)
			if err != nil {
				return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%v: %w", scanEstimationID, err)
			}
			return nil
		}
	}

	if timeNow.After(*scanEstimation.DeleteAfter) {
		err := w.backend.DeleteScanEstimation(ctx, scanEstimationID)
		if err != nil {
			return fmt.Errorf("failed to delete ScanEstimation. ScanEstimationID=%v: %w", scanEstimationID, err)
		}
	}

	return nil
}

func (w *Watcher) reconcileNoState(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	scanEstimationPatch := types.ScanEstimation{
		Status: types.NewScanEstimationStatus(
			types.ScanEstimationStatusStatePending,
			types.ScanEstimationStatusReasonCreated,
			nil,
		),
	}
	err := w.backend.PatchScanEstimation(ctx, scanEstimationID, &scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to update ScanEstimation. ScanEstimationID=%v: %w", scanEstimationID, err)
	}
	return nil
}

func (w *Watcher) reconcilePending(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid ScanEstimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: Id is nil")
	}

	// We don't want to scan terminated assets. These are historic assets
	// and we'll just get a not found error during the asset scan.
	assetFilter := "terminatedOn eq null"

	scope, ok := scanEstimation.GetScope()
	if !ok {
		return fmt.Errorf("invalid ScanEstimation: Scope is nil. ScanEstimationID=%s", scanEstimationID)
	}

	// If the scan estimation has a scope configured, 'and' it with the check for
	// not terminated to make sure that we take both into account.
	if scope != "" {
		assetFilter = fmt.Sprintf("(%s) and (%s)", assetFilter, scope)
	}

	assets, err := w.backend.GetAssets(ctx, types.GetAssetsParams{
		Filter: &assetFilter,
		Select: utils.PointerTo("id"),
	})
	if err != nil {
		return fmt.Errorf("failed to discover Assets for Scan estimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	numOfAssets := len(*assets.Items)

	if numOfAssets > 0 {
		assetIds := []string{}
		for _, asset := range *assets.Items {
			assetIds = append(assetIds, *asset.Id)
		}
		scanEstimation.AssetIDs = &assetIds
		scanEstimation.Status = types.NewScanEstimationStatus(
			types.ScanEstimationStatusStateDiscovered,
			types.ScanEstimationStatusReasonSuccessfulDiscovery,
			utils.PointerTo("Assets for Scan estimation are successfully discovered"),
		)
	} else {
		scanEstimation.Status = types.NewScanEstimationStatus(
			types.ScanEstimationStatusStateDone,
			types.ScanEstimationStatusReasonNothingToEstimate,
			utils.PointerTo("No instances found in scope for Scan estimation"),
		)
	}
	logger.Debugf("%d Asset(s) have been created for Scan estimation", numOfAssets)

	// Set default ttl if not set.
	if scanEstimation.TTLSecondsAfterFinished == nil {
		scanEstimation.TTLSecondsAfterFinished = utils.PointerTo(DefaultScanEstimationTTLSeconds)
	}

	scanEstimationPatch := &types.ScanEstimation{
		StartTime:               utils.PointerTo(time.Now()),
		TTLSecondsAfterFinished: scanEstimation.TTLSecondsAfterFinished,
		AssetIDs:                scanEstimation.AssetIDs,
		Status:                  scanEstimation.Status,
		Summary: &types.ScanEstimationSummary{
			JobsCompleted: utils.PointerTo(0),
			JobsLeftToRun: utils.PointerTo(numOfAssets),
		},
	}

	if err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch); err != nil {
		return fmt.Errorf("failed to patch Scan estimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}

func (w *Watcher) reconcileDiscovered(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid Scan estimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid Scan estimation: Id is nil")
	}

	if err := w.createAssetScanEstimationsForScanEstimation(ctx, scanEstimation); err != nil {
		return fmt.Errorf("failed to creates AssetScanEstimation(s) for ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	scanEstimationPatch := &types.ScanEstimation{
		Status: types.NewScanEstimationStatus(
			types.ScanEstimationStatusStateInProgress,
			types.ScanEstimationStatusReasonRunning,
			nil,
		),
		Summary:              scanEstimation.Summary,
		AssetIDs:             scanEstimation.AssetIDs,
		AssetScanEstimations: scanEstimation.AssetScanEstimations,
	}
	err := w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to update Scan estimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	logger.Infof("Total %d unique assets for ScanEstimation", len(*scanEstimation.AssetIDs))

	return nil
}

func (w *Watcher) createAssetScanEstimationsForScanEstimation(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation.AssetIDs == nil || *scanEstimation.AssetIDs == nil {
		return nil
	}
	numOfAssets := len(*scanEstimation.AssetIDs)

	errs := make(chan error, numOfAssets)
	var wg sync.WaitGroup
	for _, id := range *scanEstimation.AssetIDs {
		wg.Add(1)
		assetID := id
		go func() {
			defer wg.Done()

			err := w.createAssetScanEstimationForAsset(ctx, scanEstimation, assetID)
			if err != nil {
				logger.WithField("AssetID", assetID).Errorf("Failed to create AssetScanEstimation: %v", err)
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
		return fmt.Errorf("failed to create %d AssetScanEstimation(s) for ScanEstimation. ScanEstimationID=%s: %w", numOfErrs, *scanEstimation.Id, assetErrs[0])
	}

	if scanEstimation.Summary == nil {
		scanEstimation.Summary = &types.ScanEstimationSummary{}
	}
	scanEstimation.Summary.JobsLeftToRun = utils.PointerTo(numOfAssets)

	return nil
}

func (w *Watcher) newAssetScanEstimationFromScanEstimation(scanEstimation *types.ScanEstimation, assetID string) (types.AssetScanEstimation, error) {
	if scanEstimation == nil {
		return types.AssetScanEstimation{}, fmt.Errorf("empty scan estimation")
	}

	if scanEstimation.ScanTemplate == nil {
		return types.AssetScanEstimation{}, fmt.Errorf("empty scan template")
	}

	return types.AssetScanEstimation{
		TTLSecondsAfterFinished: utils.PointerTo(DefaultScanEstimationTTLSeconds),
		Asset: &types.AssetRelationship{
			Id: assetID,
		},
		AssetScanTemplate: scanEstimation.ScanTemplate.AssetScanTemplate,
		ScanEstimation: &types.ScanEstimationRelationship{
			Id: scanEstimation.Id,
		},
		Status: types.NewAssetScanEstimationStatus(
			types.AssetScanEstimationStatusStatePending,
			types.AssetScanEstimationStatusReasonCreated,
			nil,
		),
	}, nil
}

func (w *Watcher) createAssetScanEstimationForAsset(ctx context.Context, scanEstimation *types.ScanEstimation, assetID string) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	assetScanEstimationData, err := w.newAssetScanEstimationFromScanEstimation(scanEstimation, assetID)
	if err != nil {
		return fmt.Errorf("failed to generate new AssetScanEstimation for ScanEstimation. ScanEstimationID=%s, AssetID=%s: %w", *scanEstimation.Id, assetID, err)
	}

	ret, err := w.backend.PostAssetScanEstimation(ctx, assetScanEstimationData)
	if err != nil {
		var conErr client.AssetScanEstimationConflictError
		if errors.As(err, &conErr) {
			assetScanEstimationID := *conErr.ConflictingAssetScanEstimation.Id
			logger.WithField("AssetScanEstimationID", assetScanEstimationID).Debug("AssetScanEstimation already exist.")
			return nil
		}
		return fmt.Errorf("failed to post AssetScanEstimation to backend API: %w", err)
	}

	if scanEstimation.AssetScanEstimations == nil {
		scanEstimation.AssetScanEstimations = &[]types.AssetScanEstimationRelationship{}
	}
	*scanEstimation.AssetScanEstimations = append(*scanEstimation.AssetScanEstimations, types.AssetScanEstimationRelationship{Id: ret.Id})

	return nil
}

func updateScanEstimationSummaryFromAssetScanEstimation(scanEstimation *types.ScanEstimation, assetScanEstimation types.AssetScanEstimation) error {
	status, ok := assetScanEstimation.GetStatus()
	if !ok {
		return fmt.Errorf("state must not be nil for AssetScan. AssetScanID=%s", *assetScanEstimation.Id)
	}

	s := scanEstimation.Summary

	switch status.State {
	case types.AssetScanEstimationStatusStatePending:
		s.JobsLeftToRun = utils.PointerTo(*s.JobsLeftToRun + 1)
	case types.AssetScanEstimationStatusStateDone:
		if s.TotalScanTime == nil {
			s.TotalScanTime = utils.PointerTo(0)
		}
		if s.TotalScanSize == nil {
			s.TotalScanSize = utils.PointerTo(0)
		}
		if s.TotalScanCost == nil {
			s.TotalScanCost = utils.PointerTo(float32(0))
		}
		*(s.TotalScanTime) += utils.ValueOrZero(assetScanEstimation.Estimation.Duration)
		*(s.TotalScanSize) += utils.ValueOrZero(assetScanEstimation.Estimation.Size)
		*(s.TotalScanCost) += utils.ValueOrZero(assetScanEstimation.Estimation.Cost)
		fallthrough
	case types.AssetScanEstimationStatusStateAborted, types.AssetScanEstimationStatusStateFailed:
		s.JobsCompleted = utils.PointerTo(*s.JobsCompleted + 1)
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileInProgress(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid ScanEstimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	filter := fmt.Sprintf("scanEstimation/id eq '%s'", scanEstimationID)
	selector := "id,status,estimation"
	assetScanEstimations, err := w.backend.GetAssetScanEstimations(ctx, types.GetAssetScanEstimationsParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve AssetScanEstimations for ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	if assetScanEstimations.Count == nil || assetScanEstimations.Items == nil {
		return fmt.Errorf("invalid response for getting AssetScanEstimations for ScanEstimation. ScanEstimationID=%s: Count and/or Items parameters are nil", scanEstimationID)
	}

	// Reset Scan Summary as it is going to be recalculated
	scanEstimation.Summary = &types.ScanEstimationSummary{
		JobsCompleted: utils.PointerTo(0),
		JobsLeftToRun: utils.PointerTo(0),
	}

	var failedAssetScanEstimations int
	for _, assetScanEstimation := range *assetScanEstimations.Items {
		assetScanEstimationID, ok := assetScanEstimation.GetID()
		if !ok {
			return errors.New("invalid AssetScanEstimation: ID is nil")
		}

		if err := updateScanEstimationSummaryFromAssetScanEstimation(scanEstimation, assetScanEstimation); err != nil {
			return fmt.Errorf("failed to update ScanEstimation Summary from AssetScanEstimation. ScanEstimationID=%s AssetScanEstimationID=%s: %w",
				scanEstimationID, assetScanEstimationID, err)
		}

		status, ok := assetScanEstimation.GetStatus()
		if !ok {
			logger.Warnf("Failed to get assetScanEstimation %v state", assetScanEstimationID)
		} else if status.State == types.AssetScanEstimationStatusStateFailed {
			failedAssetScanEstimations++
		}
	}
	logger.Tracef("ScanEstimation Summary updated. JobCompleted=%d JobLeftToRun=%d", *scanEstimation.Summary.JobsCompleted,
		*scanEstimation.Summary.JobsLeftToRun)

	message := utils.PointerTo(
		fmt.Sprintf(
			"%d succeeded, %d failed out of %d total asset scan estimations",
			*assetScanEstimations.Count-failedAssetScanEstimations,
			failedAssetScanEstimations,
			*assetScanEstimations.Count,
		),
	)

	if *scanEstimation.Summary.JobsLeftToRun <= 0 {
		if failedAssetScanEstimations > 0 {
			scanEstimation.Status = types.NewScanEstimationStatus(
				types.ScanEstimationStatusStateFailed,
				types.ScanEstimationStatusReasonError,
				message,
			)
		} else {
			scanEstimation.Status = types.NewScanEstimationStatus(
				types.ScanEstimationStatusStateDone,
				types.ScanEstimationStatusReasonSuccess,
				message,
			)
		}
		scanEstimation.EndTime = utils.PointerTo(time.Now())

		if err := updateTotalScanTimeWithParallelScans(scanEstimation); err != nil {
			return fmt.Errorf("failed to update scan time from paraller scans: %w", err)
		}
		scanEstimation.DeleteAfter = utils.PointerTo(scanEstimation.EndTime.Add(time.Duration(*scanEstimation.TTLSecondsAfterFinished) * time.Second))
	}

	scanEstimationPatch := &types.ScanEstimation{
		DeleteAfter: scanEstimation.DeleteAfter,
		Status:      scanEstimation.Status,
		Summary:     scanEstimation.Summary,
		EndTime:     scanEstimation.EndTime,
		AssetIDs:    scanEstimation.AssetIDs,
	}
	err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}

func updateTotalScanTimeWithParallelScans(scanEstimation *types.ScanEstimation) error {
	if scanEstimation == nil {
		return fmt.Errorf("empty scan estimation")
	}

	if scanEstimation.ScanTemplate == nil {
		return fmt.Errorf("empty scan template")
	}

	if scanEstimation.Summary == nil {
		return fmt.Errorf("empty summary")
	}

	if scanEstimation.Summary.JobsCompleted == nil {
		return fmt.Errorf("jobsCompleted is not set")
	}

	if *scanEstimation.Summary.JobsCompleted == 0 {
		return fmt.Errorf("0 completed jobs in summary")
	}

	maxParallelScanners := utils.ValueOrZero(scanEstimation.ScanTemplate.MaxParallelScanners)

	if maxParallelScanners > 1 {
		numberOfJobs := *scanEstimation.Summary.JobsCompleted

		actualParallelScanners := int(math.Min(float64(maxParallelScanners), float64(numberOfJobs)))

		// Note: This is a rough estimation, as we don't know which jobs will be running in parallel.
		*scanEstimation.Summary.TotalScanTime = *scanEstimation.Summary.TotalScanTime / actualParallelScanners
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileAborted(ctx context.Context, scanEstimation *types.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid ScanEstimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	filter := fmt.Sprintf("scanEstimation/id eq '%s' and status/state ne '%s' and status/state ne '%s'",
		scanEstimationID, types.AssetScanEstimationStatusStateAborted, types.AssetScanEstimationStatusStateDone)
	selector := "id,status"
	params := types.GetAssetScanEstimationsParams{
		Filter: &filter,
		Select: &selector,
	}

	assetScanEstimations, err := w.backend.GetAssetScanEstimations(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to fetch AssetScanEstimation(s) for ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	if assetScanEstimations.Items != nil && len(*assetScanEstimations.Items) > 0 {
		var reconciliationFailed bool
		var wg sync.WaitGroup

		for _, assetScanEstimation := range *assetScanEstimations.Items {
			if assetScanEstimation.Id == nil {
				continue
			}
			assetScanEstimationID := *assetScanEstimation.Id

			wg.Add(1)
			go func() {
				defer wg.Done()
				ase := types.AssetScanEstimation{
					Status: types.NewAssetScanEstimationStatus(
						types.AssetScanEstimationStatusStateAborted,
						types.AssetScanEstimationStatusReasonCancellation,
						nil,
					),
				}

				err = w.backend.PatchAssetScanEstimation(ctx, ase, assetScanEstimationID)
				if err != nil {
					logger.WithField("AssetScanEstimationID", assetScanEstimationID).Error("Failed to patch AssetScanEstimation")
					reconciliationFailed = true
					return
				}
			}()
		}
		wg.Wait()

		// NOTE: reconciliationFailed is used to track errors returned by patching AssetScanEstimations
		//       as setting the state of ScanEstimation to Failed must be skipped in case
		//       even a single error occurred to allow reconciling re-running for this ScanEstimation.
		if reconciliationFailed {
			return errors.New("updating one or more AssetScanEstimations failed")
		}
	}

	scanEstimation.EndTime = utils.PointerTo(time.Now())

	scanEstimationPatch := &types.ScanEstimation{
		EndTime: scanEstimation.EndTime,
		Status: types.NewScanEstimationStatus(
			types.ScanEstimationStatusStateFailed,
			types.ScanEstimationStatusReasonAborted,
			utils.PointerTo("ScanEstimation has been aborted"),
		),
		AssetIDs: scanEstimation.AssetIDs,
	}
	err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}
