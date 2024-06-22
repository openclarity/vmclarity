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

package finding

import (
	"context"
	"fmt"
	"time"

	apiclient "github.com/openclarity/vmclarity/api/client"
	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/core/to"
	"github.com/openclarity/vmclarity/orchestrator/common"
	"github.com/openclarity/vmclarity/provider"
)

type (
	ScanQueue      = common.Queue[FindingReconcileEvent]
	ScanPoller     = common.Poller[FindingReconcileEvent]
	ScanReconciler = common.Reconciler[FindingReconcileEvent]
)

func New(c Config) *Watcher {
	return &Watcher{
		client:              c.Client,
		provider:            c.Provider,
		pollPeriod:          c.PollPeriod,
		reconcileTimeout:    c.ReconcileTimeout,
		summaryUpdatePeriod: c.SummaryUpdatePeriod,
		queue:               common.NewQueue[FindingReconcileEvent](),
	}
}

type Watcher struct {
	client              *apiclient.Client
	provider            provider.Provider
	pollPeriod          time.Duration
	reconcileTimeout    time.Duration
	summaryUpdatePeriod time.Duration

	queue *ScanQueue
}

func (w *Watcher) Start(ctx context.Context) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("controller", "ScanWatcher")
	ctx = log.SetLoggerForContext(ctx, logger)

	poller := &ScanPoller{
		PollPeriod: w.pollPeriod,
		Queue:      w.queue,
		GetItems:   w.GetFindingsWithOutdatedSummary,
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
func (w *Watcher) GetFindingsWithOutdatedSummary(ctx context.Context) ([]FindingReconcileEvent, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)
	logger.Debugf("Fetching Findings with outdated summary")

	// NOTE: we only care about package findings since other findings are not
	// connected to vulnerabilities
	findings, err := w.client.GetFindings(ctx, apitypes.GetFindingsParams{
		Filter: to.Ptr(fmt.Sprintf(
			"findingInfo/objectType eq 'Package' and summary/updatedAt lt %s",
			time.Now().Add(-w.summaryUpdatePeriod).Format(time.RFC3339)),
		),
		Select: to.Ptr("id,findingInfo,summary"),
		Count:  to.Ptr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get outdated findings: %w", err)
	}

	// Check returned results
	switch {
	case findings.Items == nil && findings.Count == nil:
		return nil, fmt.Errorf("failed to get outdated findings: invalid API response: %v", findings)
	case findings.Count != nil && *findings.Count <= 0:
		fallthrough
	case findings.Items != nil && len(*findings.Items) <= 0:
		return nil, nil
	}

	events := make([]FindingReconcileEvent, 0, *findings.Count)
	for _, finding := range *findings.Items {
		if finding.Id == nil || *finding.Id == "" {
			logger.Warnf("Skipping invalid Finding: ID is nil: %v", finding)
			continue
		}

		events = append(events, FindingReconcileEvent{
			FindingID: *finding.Id,
		})
	}

	return events, nil
}

// nolint:cyclop
func (w *Watcher) Reconcile(ctx context.Context, event FindingReconcileEvent) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(event.ToFields())
	ctx = log.SetLoggerForContext(ctx, logger)

	// Get finding
	finding, err := w.client.GetFinding(ctx, event.FindingID, apitypes.GetFindingsFindingIDParams{})
	if err != nil || finding == nil {
		return fmt.Errorf("failed to fetch Finding. FindingID=%s: %w", event.FindingID, err)
	}

	discriminator, err := finding.FindingInfo.ValueByDiscriminator()
	if err != nil {
		return fmt.Errorf("failed to extract Finding type. FindingID=%s: %w", event.FindingID, err)
	}

	// Construct filter based on discriminator type
	// Use info properties that make the finding unique, check package scanner/findingkey.
	switch info := discriminator.(type) {
	case apitypes.PackageFindingInfo:
		if err = w.reconcilePackageSummary(ctx, event.FindingID, info); err != nil {
			return err
		}
	default:
	}

	return nil
}

// nolint:cyclop
func (w *Watcher) reconcilePackageSummary(ctx context.Context, findingID string, pkg apitypes.PackageFindingInfo) error {
	// Get total vulnerabilities
	critialVuls, err := w.getPackageVulnerabilityFindingsCount(ctx, pkg, apitypes.CRITICAL)
	if err != nil {
		return fmt.Errorf("failed to list critial vulnerabilities: %w", err)
	}
	highVuls, err := w.getPackageVulnerabilityFindingsCount(ctx, pkg, apitypes.HIGH)
	if err != nil {
		return fmt.Errorf("failed to list high vulnerabilities: %w", err)
	}
	mediumVuls, err := w.getPackageVulnerabilityFindingsCount(ctx, pkg, apitypes.MEDIUM)
	if err != nil {
		return fmt.Errorf("failed to list medium vulnerabilities: %w", err)
	}
	lowVuls, err := w.getPackageVulnerabilityFindingsCount(ctx, pkg, apitypes.LOW)
	if err != nil {
		return fmt.Errorf("failed to list low vulnerabilities: %w", err)
	}
	negligibleVuls, err := w.getPackageVulnerabilityFindingsCount(ctx, pkg, apitypes.NEGLIGIBLE)
	if err != nil {
		return fmt.Errorf("failed to list negligible vulnerabilities: %w", err)
	}

	// Patch finding with updated summary
	err = w.client.PatchFinding(ctx, findingID, apitypes.Finding{
		Id: to.Ptr(findingID),
		Summary: &apitypes.FindingSummary{
			UpdatedAt: to.Ptr(time.Now().Format(time.RFC3339)),
			TotalVulnerabilities: &apitypes.VulnerabilitySeveritySummary{
				TotalCriticalVulnerabilities:   to.Ptr(critialVuls),
				TotalHighVulnerabilities:       to.Ptr(highVuls),
				TotalMediumVulnerabilities:     to.Ptr(mediumVuls),
				TotalLowVulnerabilities:        to.Ptr(lowVuls),
				TotalNegligibleVulnerabilities: to.Ptr(negligibleVuls),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to patch finding summary: %w", err)
	}

	return nil
}

func (w *Watcher) getPackageVulnerabilityFindingsCount(ctx context.Context, pkg apitypes.PackageFindingInfo, severity apitypes.VulnerabilitySeverity) (int, error) {
	activeFindings, err := w.client.GetFindings(ctx, apitypes.GetFindingsParams{
		Count: to.Ptr(true),
		Filter: to.Ptr(fmt.Sprintf(
			"findingInfo/objectType eq 'Vulnerability' and findingInfo/severity eq '%s' and findingInfo/package/Name eq '%s' and findingInfo/package/Name eq '%s'",
			string(severity), to.ValueOrZero(pkg.Name), to.ValueOrZero(pkg.Version)),
		),
		// select the smallest amount of data to return in items, we
		// only care about the count.
		Top:    to.Ptr(1),
		Select: to.Ptr("id"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list package vulnerability findings: %w", err)
	}

	return *activeFindings.Count, nil
}
