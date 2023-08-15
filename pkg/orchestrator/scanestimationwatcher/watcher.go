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
	"sync"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/orchestrator/common"
	"github.com/openclarity/vmclarity/pkg/orchestrator/provider"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"github.com/openclarity/vmclarity/pkg/shared/log"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
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
	backend               *backendclient.BackendClient
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

	filter := fmt.Sprintf("state ne '%s' and state ne '%s'", models.ScanEstimationStateDone, models.ScanEstimationStateFailed)
	selector := "id"
	params := models.GetScanEstimationsParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	}
	scanEstimations, err := w.backend.GetScanEstimations(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get running sScanEstimations: %v", err)
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

	params := models.GetScanEstimationsScanEstimationIDParams{}
	scanEstimation, err := w.backend.GetScanEstimation(ctx, event.ScanEstimationID, params)
	if err != nil || scanEstimation == nil {
		return fmt.Errorf("failed to fetch ScanEstimation. ScanEstimationID=%s: %w", event.ScanEstimationID, err)
	}

	if scanEstimation.IsTimedOut(w.scanEstimationTimeout) {
		scanEstimation.State = utils.PointerTo(models.ScanEstimationStateAborted)
		scanEstimation.StateMessage = utils.PointerTo("ScanEstimation has been timed out")
		scanEstimation.StateReason = utils.PointerTo(models.ScanEstimationStateReasonTimedOut)

		err = w.backend.PatchScanEstimation(ctx, *scanEstimation.Id, &models.ScanEstimation{
			State:        scanEstimation.State,
			StateMessage: scanEstimation.StateMessage,
			StateReason:  scanEstimation.StateReason,
		})
		if err != nil {
			return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", event.ScanEstimationID, err)
		}
	}

	state, ok := scanEstimation.GetState()
	if !ok {
		return fmt.Errorf("failed to determine state of ScanEstimation. ScanEstimationID=%s", event.ScanEstimationID)
	}
	logger.Tracef("Reconciling ScanEstimation state: %s", state)

	switch state {
	case models.ScanEstimationStatePending:
		if err = w.reconcilePending(ctx, scanEstimation); err != nil {
			return err
		}
	case models.ScanEstimationStateDiscovered:
		if err = w.reconcileDiscovered(ctx, scanEstimation); err != nil {
			return err
		}
	case models.ScanEstimationStateInProgress:
		if err = w.reconcileInProgress(ctx, scanEstimation); err != nil {
			return err
		}
	case models.ScanEstimationStateAborted:
		if err = w.reconcileAborted(ctx, scanEstimation); err != nil {
			return err
		}
	case models.ScanEstimationStateDone, models.ScanEstimationStateFailed:
		logger.Debug("Reconciling ScanEstimation is skipped as it is already finished.")
		fallthrough
	default:
		return nil
	}

	return nil
}

func (w *Watcher) reconcilePending(ctx context.Context, scanEstimation *models.ScanEstimation) error {
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

	assets, err := w.backend.GetAssets(ctx, models.GetAssetsParams{
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
		scanEstimation.State = utils.PointerTo(models.ScanEstimationStateDiscovered)
		scanEstimation.StateMessage = utils.PointerTo("Assets for Scan estimation are successfully discovered")
	} else {
		scanEstimation.State = utils.PointerTo(models.ScanEstimationStateDone)
		scanEstimation.StateReason = utils.PointerTo(models.ScanEstimationStateReasonNothingToEstimate)
		scanEstimation.StateMessage = utils.PointerTo("No instances found in scope for Scan estimation")
	}
	logger.Debugf("%d Asset(s) have been created for Scan estimation", numOfAssets)

	scanEstimationPatch := &models.ScanEstimation{
		AssetIDs:     scanEstimation.AssetIDs,
		State:        scanEstimation.State,
		StateReason:  scanEstimation.StateReason,
		StateMessage: scanEstimation.StateMessage,
		Summary: &models.ScanEstimationSummary{
			JobsCompleted: utils.PointerTo(0),
			JobsLeftToRun: utils.PointerTo(numOfAssets),
		},
	}

	if err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch); err != nil {
		return fmt.Errorf("failed to patch Scan estimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}

func (w *Watcher) reconcileDiscovered(ctx context.Context, scanEstimation *models.ScanEstimation) error {
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
	scanEstimation.State = utils.PointerTo(models.ScanEstimationStateInProgress)

	scanEstimationPatch := &models.ScanEstimation{
		State:    scanEstimation.State,
		Summary:  scanEstimation.Summary,
		AssetIDs: scanEstimation.AssetIDs,
	}
	err := w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to update Scan estimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	logger.Infof("Total %d unique assets for ScanEstimation", len(*scanEstimation.AssetIDs))

	return nil
}

func (w *Watcher) createAssetScanEstimationsForScanEstimation(ctx context.Context, scanEstimation *models.ScanEstimation) error {
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
		scanEstimation.Summary = &models.ScanEstimationSummary{}
	}
	scanEstimation.Summary.JobsLeftToRun = utils.PointerTo(numOfAssets)

	return nil
}

func (w *Watcher) newAssetScanEstimationFromScanEstimation(ctx context.Context, scanEstimation *models.ScanEstimation, assetID string) (*models.AssetScanEstimation, error) {

	if scanEstimation.ScanTemplate == nil {
		return nil, fmt.Errorf("empty scan template")
	}

	// get asset from backend
	asset, err := w.backend.GetAsset(ctx, assetID, models.GetAssetsAssetIDParams{})
	if err != nil {
		return nil, err
	}

	assetWriteable := &models.Asset{
		Id:         asset.Id,
		Revision:   asset.Revision,
		ScansCount: asset.ScansCount,
		Summary:    asset.Summary,
		AssetInfo:  asset.AssetInfo,
	}

	return &models.AssetScanEstimation{
		Asset:             assetWriteable,
		AssetScanTemplate: scanEstimation.ScanTemplate.AssetScanTemplate,
		Estimation:        nil,
		Id:                nil,
		ScanEstimation: &models.ScanEstimationFakeRelationship{
			Id: scanEstimation.Id,
		},
		State: &models.AssetScanEstimationState{
			Errors:             nil,
			LastTransitionTime: utils.PointerTo(time.Now()),
			State:              utils.PointerTo(models.AssetScanEstimationStateStatePending),
		},
	}, nil
}

func (w *Watcher) createAssetScanEstimationForAsset(ctx context.Context, scanEstimation *models.ScanEstimation, assetID string) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	assetScanEstimationData, err := w.newAssetScanEstimationFromScanEstimation(ctx, scanEstimation, assetID)
	if err != nil {
		return fmt.Errorf("failed to generate new AssetScanEstimation for ScanEstimation. ScanEstimationID=%s, AssetID=%s: %w", *scanEstimation.Id, assetID, err)
	}

	_, err = w.backend.PostAssetScanEstimation(ctx, *assetScanEstimationData)
	if err != nil {
		var conErr backendclient.AssetScanEstimationConflictError
		if errors.As(err, &conErr) {
			assetScanEstimationID := *conErr.ConflictingAssetScanEstimation.Id
			logger.WithField("AssetScanEstimationID", assetScanEstimationID).Debug("AssetScanEstimation already exist.")
			return nil
		}
		return fmt.Errorf("failed to post AssetScanEstimation to backend API: %w", err)
	}
	return nil
}

func updateScanEstimationSummaryFromAssetScanEstimation(scanEstimation *models.ScanEstimation, result models.AssetScanEstimation) error {
	state, ok := result.GetState()
	if !ok {
		return fmt.Errorf("state must not be nil for AssetScan. AssetScanID=%s", *result.Id)
	}

	s := scanEstimation.Summary

	switch state {
	case models.AssetScanEstimationStateStatePending:
		s.JobsLeftToRun = utils.PointerTo(*s.JobsLeftToRun + 1)
	case models.AssetScanEstimationStateStateDone, models.AssetScanEstimationStateStateAborted, models.AssetScanEstimationStateStateFailed:
		s.JobsCompleted = utils.PointerTo(*s.JobsCompleted + 1)
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileInProgress(ctx context.Context, scanEstimation *models.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid ScanEstimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	filter := fmt.Sprintf("scanEstimation/id eq '%s'", scanEstimationID)
	selector := "id,state"
	assetScanEstimations, err := w.backend.GetAssetScanEstimations(ctx, models.GetAssetScanEstimationsParams{
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
	scanEstimation.Summary = &models.ScanEstimationSummary{
		JobsCompleted: utils.PointerTo(0),
		JobsLeftToRun: utils.PointerTo(0),
	}

	var assetScanEstimationsWithErr int
	for _, assetScanEstimation := range *assetScanEstimations.Items {
		assetScanEstimationID, ok := assetScanEstimation.GetID()
		if !ok {
			return errors.New("invalid AssetScanEstimation: ID is nil")
		}

		if err := updateScanEstimationSummaryFromAssetScanEstimation(scanEstimation, assetScanEstimation); err != nil {
			return fmt.Errorf("failed to update ScanEstimation Summary from AssetScanEstimation. ScanEstimationID=%s AssetScanEstimationID=%s: %w",
				scanEstimationID, assetScanEstimationID, err)
		}

		errs := assetScanEstimation.GetGeneralErrors()
		if len(errs) > 0 {
			assetScanEstimationsWithErr++
		}
	}
	logger.Tracef("ScanEstimation Summary updated. JobCompleted=%d JobLeftToRun=%d", *scanEstimation.Summary.JobsCompleted,
		*scanEstimation.Summary.JobsLeftToRun)

	if *scanEstimation.Summary.JobsLeftToRun <= 0 {
		if assetScanEstimationsWithErr > 0 {
			scanEstimation.State = utils.PointerTo(models.ScanEstimationStateFailed)
			scanEstimation.StateReason = utils.PointerTo(models.ScanEstimationStateReasonOneOrMoreAssetFailedToEstimate)
		} else {
			scanEstimation.State = utils.PointerTo(models.ScanEstimationStateDone)
			scanEstimation.StateReason = utils.PointerTo(models.ScanEstimationStateReasonSuccess)
		}
		scanEstimation.StateMessage = utils.PointerTo(fmt.Sprintf("%d succeeded, %d failed out of %d total asset scan estimations",
			*assetScanEstimations.Count-assetScanEstimationsWithErr, assetScanEstimationsWithErr, *assetScanEstimations.Count))

		scanEstimation.EndTime = utils.PointerTo(time.Now())
	}

	scanEstimationPatch := &models.ScanEstimation{
		State:        scanEstimation.State,
		Summary:      scanEstimation.Summary,
		StateMessage: scanEstimation.StateMessage,
		EndTime:      scanEstimation.EndTime,
		AssetIDs:     scanEstimation.AssetIDs,
	}
	err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileAborted(ctx context.Context, scanEstimation *models.ScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	if scanEstimation == nil {
		return errors.New("invalid ScanEstimation: object is nil")
	}

	scanEstimationID, ok := scanEstimation.GetID()
	if !ok {
		return errors.New("invalid ScanEstimation: ID is nil")
	}

	filter := fmt.Sprintf("scanEstimation/id eq '%s' and state/state ne '%s' and state/state ne '%s'",
		scanEstimationID, models.AssetScanEstimationStateStateAborted, models.AssetScanEstimationStateStateDone)
	selector := "id,state"
	params := models.GetAssetScanEstimationsParams{
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
				ase := models.AssetScanEstimation{
					State: &models.AssetScanEstimationState{
						State: utils.PointerTo(models.AssetScanEstimationStateStateAborted),
					},
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
	scanEstimation.State = utils.PointerTo(models.ScanEstimationStateFailed)
	scanEstimation.StateReason = utils.PointerTo(models.ScanEstimationStateReasonAborted)
	scanEstimation.StateMessage = utils.PointerTo("ScanEstimation has been aborted")

	scanEstimationPatch := &models.ScanEstimation{
		State:        scanEstimation.State,
		EndTime:      scanEstimation.EndTime,
		StateReason:  scanEstimation.StateReason,
		StateMessage: scanEstimation.StateMessage,
		AssetIDs:     scanEstimation.AssetIDs,
	}
	err = w.backend.PatchScanEstimation(ctx, scanEstimationID, scanEstimationPatch)
	if err != nil {
		return fmt.Errorf("failed to patch ScanEstimation. ScanEstimationID=%s: %w", scanEstimationID, err)
	}

	return nil
}
