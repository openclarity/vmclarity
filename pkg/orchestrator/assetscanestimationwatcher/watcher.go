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

package assetscanestimationwatcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/orchestrator/common"
	"github.com/openclarity/vmclarity/pkg/orchestrator/provider"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"github.com/openclarity/vmclarity/pkg/shared/log"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

type (
	AssetScanEstimationQueue      = common.Queue[AssetScanEstimationReconcileEvent]
	AssetScanEstimationPoller     = common.Poller[AssetScanEstimationReconcileEvent]
	AssetScanEstimationReconciler = common.Reconciler[AssetScanEstimationReconcileEvent]
)

func New(c Config) *Watcher {
	return &Watcher{
		backend:          c.Backend,
		provider:         c.Provider,
		pollPeriod:       c.PollPeriod,
		reconcileTimeout: c.ReconcileTimeout,
		abortTimeout:     c.AbortTimeout,
		queue:            common.NewQueue[AssetScanEstimationReconcileEvent](),
	}
}

type Watcher struct {
	backend          *backendclient.BackendClient
	provider         provider.Provider
	pollPeriod       time.Duration
	reconcileTimeout time.Duration
	abortTimeout     time.Duration

	queue *AssetScanEstimationQueue
}

func (w *Watcher) Start(ctx context.Context) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("controller", "AssetScanEstimationWatcher")
	ctx = log.SetLoggerForContext(ctx, logger)

	poller := &AssetScanEstimationPoller{
		PollPeriod: w.pollPeriod,
		Queue:      w.queue,
		GetItems:   w.GetAssetScanEstimations,
	}
	poller.Start(ctx)

	reconciler := &AssetScanEstimationReconciler{
		ReconcileTimeout:  w.reconcileTimeout,
		Queue:             w.queue,
		ReconcileFunction: w.Reconcile,
	}
	reconciler.Start(ctx)
}

// nolint:cyclop
func (w *Watcher) GetAssetScanEstimations(ctx context.Context) ([]AssetScanEstimationReconcileEvent, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)
	logger.Debugf("Fetching AssetScanEstimations which need to be reconciled")

	filter := fmt.Sprintf("state/state ne '%s' and state/state ne '%s'",
		models.AssetScanEstimationStateStateDone, models.AssetScanEstimationStateStateFailed)
	selector := "id,scanEstimation/id,asset/id"
	params := models.GetAssetScanEstimationsParams{
		Filter: &filter,
		Select: &selector,
		Count:  utils.PointerTo(true),
	}
	assetScanEstimations, err := w.backend.GetAssetScanEstimations(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get AssetScanEstimations: %w", err)
	}

	switch {
	case assetScanEstimations.Items == nil && assetScanEstimations.Count == nil:
		return nil, fmt.Errorf("failed to fetch AssetScanEstimations: invalid API response: %v", assetScanEstimations)
	case assetScanEstimations.Count != nil && *assetScanEstimations.Count <= 0:
		fallthrough
	case assetScanEstimations.Items != nil && len(*assetScanEstimations.Items) <= 0:
		return nil, nil
	}

	events := make([]AssetScanEstimationReconcileEvent, 0, len(*assetScanEstimations.Items))
	for _, assetScanEstimation := range *assetScanEstimations.Items {
		assetScanEstimationID, ok := assetScanEstimation.GetID()
		if !ok {
			logger.Warnf("Skipping due to invalid AssetScanEstimation: ID is nil: %v", assetScanEstimation)
			continue
		}
		scanEstimationID, ok := assetScanEstimation.GetScanEstimationID()
		if !ok {
			logger.Warnf("Skipping due to invalid AssetScanEstimation: ScanEstimation.ID is nil: %v", assetScanEstimation)
			continue
		}
		assetID, ok := assetScanEstimation.GetAssetID()
		if !ok {
			logger.Warnf("Skipping due to invalid AssetScanEstimation: Asset.ID is nil: %v", assetScanEstimation)
			continue
		}

		events = append(events, AssetScanEstimationReconcileEvent{
			AssetScanEstimationID: assetScanEstimationID,
			ScanEstimationID:      scanEstimationID,
			AssetID:               assetID,
		})
	}

	return events, nil
}

// nolint:cyclop
func (w *Watcher) Reconcile(ctx context.Context, event AssetScanEstimationReconcileEvent) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(event.ToFields())
	ctx = log.SetLoggerForContext(ctx, logger)

	assetScanEstimation, err := w.backend.GetAssetScanEstimation(ctx, event.AssetScanEstimationID, models.GetAssetScanEstimationsAssetScanEstimationIDParams{
		Expand: utils.PointerTo("asset"),
	})
	if err != nil {
		return fmt.Errorf("failed to get AssetScanEstimation with %s id: %w", event.AssetScanEstimationID, err)
	}

	state, ok := assetScanEstimation.GetState()
	if !ok {
		return fmt.Errorf("cannot determine state of AssetScanEstimation with %s id", event.AssetScanEstimationID)
	}

	logger.Tracef("Reconciling AssetScanEstimation state: %s", state)

	switch state {
	case models.AssetScanEstimationStateStatePending:
		if err = w.reconcilePending(ctx, &assetScanEstimation); err != nil {
			return err
		}
	case models.AssetScanEstimationStateStateAborted:
		if err = w.reconcileAborted(ctx, &assetScanEstimation); err != nil {
			return err
		}
	default:
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcilePending(ctx context.Context, assetScanEstimation *models.AssetScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	logger.Debugf("reconciling pending asset scan estimations")

	assetScanEstimationID, ok := assetScanEstimation.GetID()
	if !ok {
		return errors.New("invalid AssetScanEstimation: ID is nil")
	}

	if assetScanEstimation.Asset == nil || assetScanEstimation.Asset.AssetInfo == nil {
		return errors.New("invalid AssetScanEstimation: Asset or AssetInfo is nil")
	}

	asset := assetScanEstimation.Asset

	stats, _ := w.getLatestAssetScanStats(ctx, asset)

	estimation, err := w.provider.Estimate(ctx, stats, asset, assetScanEstimation.AssetScanTemplate)
	if err != nil {
		assetScanEstimation.State.State = utils.PointerTo(models.AssetScanEstimationStateStateFailed)
		assetScanEstimation.State.Errors = utils.PointerTo(utils.UnwrapErrorStrings(err))
	} else {
		assetScanEstimation.State.State = utils.PointerTo(models.AssetScanEstimationStateStateDone)
	}

	assetScanEstimation.State.LastTransitionTime = utils.PointerTo(time.Now())

	assetScanEstimation.Estimation = &models.Estimation{
		Cost: estimation.Cost,
		Size: estimation.Size,
		Time: estimation.Time,
	}

	assetScanEstimationPatch := models.AssetScanEstimation{
		State:      assetScanEstimation.State,
		Estimation: assetScanEstimation.Estimation,
	}
	err = w.backend.PatchAssetScanEstimation(ctx, assetScanEstimationPatch, assetScanEstimationID)
	if err != nil {
		return fmt.Errorf("failed to update AssetScanEstimation. AssetScanEstimation=%s: %w", assetScanEstimationID, err)
	}
	return nil
}

// nolint:cyclop
func (w *Watcher) reconcileAborted(ctx context.Context, assetScanEstimation *models.AssetScanEstimation) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	assetScanEstimationID, ok := assetScanEstimation.GetID()
	if !ok {
		return errors.New("invalid AssetScanEstimation: ID is nil")
	}

	// Check if AssetScanEstimation is in aborted state for more time than the timeout allows
	if assetScanEstimation.State == nil || assetScanEstimation.State.State == nil {
		return errors.New("invalid AssetScanEstimation: State is nil")
	}

	var transitionTimeToAbort time.Time
	if assetScanEstimation.State.LastTransitionTime != nil {
		transitionTimeToAbort = *assetScanEstimation.State.LastTransitionTime
		logger.Debugf("AssetScanEstimation moved to aborted state: %s", transitionTimeToAbort)
	}

	now := time.Now()
	abortTimedOut := now.After(transitionTimeToAbort.Add(w.abortTimeout))
	if !abortTimedOut {
		logger.Tracef("AssetScanEstimation in aborted state is not timed out yet. TransitionTime=%s Timeout=%s",
			transitionTimeToAbort, w.abortTimeout)
		return nil
	}
	logger.Tracef("AssetScanEstimation in aborted state is timed out. TransitionTime=%s Timeout=%s",
		transitionTimeToAbort, w.abortTimeout)

	assetScanEstimation.State.State = utils.PointerTo(models.AssetScanEstimationStateStateDone)
	assetScanEstimation.State.LastTransitionTime = utils.PointerTo(now)
	assetScanEstimation.State.Errors = utils.PointerTo([]string{
		fmt.Sprintf("failed to wait for scanner to finish graceful shutdown on abort after: %s", w.abortTimeout),
	})

	assetScanEstimationPatch := models.AssetScanEstimation{
		State: assetScanEstimation.State,
	}
	err := w.backend.PatchAssetScanEstimation(ctx, assetScanEstimationPatch, assetScanEstimationID)
	if err != nil {
		return fmt.Errorf("failed to update AssetScanEstimation. AssetScanEstimation=%s: %w", assetScanEstimationID, err)
	}

	return nil
}
