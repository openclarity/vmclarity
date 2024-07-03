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

package scanner

import (
	"context"
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/exploits"
	"github.com/openclarity/vmclarity/scanner/families/infofinder"
	"github.com/openclarity/vmclarity/scanner/families/malware"
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration"
	"github.com/openclarity/vmclarity/scanner/families/plugins"
	"github.com/openclarity/vmclarity/scanner/families/rootkits"
	"github.com/openclarity/vmclarity/scanner/families/sbom"
	"github.com/openclarity/vmclarity/scanner/families/secrets"
	"github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/families/vulnerabilities"
	"github.com/openclarity/vmclarity/scanner/utils/containerrootfs"
	"github.com/openclarity/vmclarity/workflow"
	workflowtypes "github.com/openclarity/vmclarity/workflow/types"
)

type FamilyResult struct {
	FamilyType families.FamilyType
	Result     any
	Err        error
}

type FamilyNotifier interface {
	FamilyStarted(context.Context, families.FamilyType) error
	FamilyFinished(ctx context.Context, res FamilyResult) error
}

type Manager struct {
	config *Config
	tasks  []workflowtypes.Task[familyRunParams]
}

func New(config *Config) *Manager {
	manager := &Manager{
		config: config,
	}

	// Analyzers
	if config.SBOM.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			sbom.New(config.SBOM), "sbom", nil,
		))
	}

	// Scanners
	if config.Vulnerabilities.Enabled {
		// must run after SBOM to support the case when configured to use the output from sbom
		var deps []string
		if config.SBOM.Enabled {
			deps = append(deps, "sbom")
		}

		manager.tasks = append(manager.tasks, newWorkflowTask(
			vulnerabilities.New(config.Vulnerabilities), "vulnerabilities", deps,
		))
	}
	if config.Secrets.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			secrets.New(config.Secrets), "secrets", nil,
		))
	}
	if config.Rootkits.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			rootkits.New(config.Rootkits), "rootkits", nil,
		))
	}
	if config.Malware.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			malware.New(config.Malware), "malware", nil,
		))
	}
	if config.Misconfiguration.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			misconfiguration.New(config.Misconfiguration), "misconfiguration", nil,
		))
	}
	if config.InfoFinder.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			infofinder.New(config.InfoFinder), "infofinder", nil,
		))
	}
	if config.Plugins.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTask(
			plugins.New(config.Plugins), "plugins", nil,
		))
	}

	// Enrichers.
	if config.Exploits.Enabled {
		// must run after Vulnerabilities to support the case when configured to use the output from Vulnerabilities
		var deps []string
		if config.Vulnerabilities.Enabled {
			deps = append(deps, "vulnerabilities")
		}

		manager.tasks = append(manager.tasks, newWorkflowTask(
			exploits.New(config.Exploits), "exploits", deps,
		))
	}

	return manager
}

func (m *Manager) Run(ctx context.Context, notifier FamilyNotifier) []error {
	var oneOrMoreFamilyFailed bool
	var errs []error

	logger := log.GetLoggerFromContextOrDiscard(ctx)

	// Register container cache
	utils.ContainerRootfsCache = containerrootfs.NewCache()
	defer func() {
		err := utils.ContainerRootfsCache.CleanupAll()
		if err != nil {
			logger.WithError(err).Errorf("failed to cleanup all cached container rootfs files")
		}
	}()

	// Define channel to use to listen for processing errors
	errCh := make(chan error)

	// Run task processor in the background
	go func() {
		// Close error channel to allow listener to exit properly
		defer close(errCh)

		// Create families processor
		processor, err := workflow.New[familyRunParams, workflowtypes.Task[familyRunParams]](m.tasks)
		if err != nil {
			errCh <- fmt.Errorf("failed to create families processor: %w", err)
			return
		}

		params := familyRunParams{
			Notifier: notifier,
			Results:  families.NewResults(),
			ErrCh:    errCh,
		}

		// Run families processor
		if err := processor.Run(ctx, params); err != nil {
			errCh <- fmt.Errorf("failed to run families processor: %w", err)
			return
		}
	}()

	// Listen for processing errors
	for err := range errCh {
		if err == nil {
			continue
		}

		// Check if family run failed, otherwise add the error to slice
		var familyErr *runnerFamilyRunError
		if errors.As(err, &familyErr) {
			oneOrMoreFamilyFailed = true
		} else {
			errs = append(errs, err)
		}
	}

	if oneOrMoreFamilyFailed {
		errs = append(errs, errors.New("at least one family failed to run"))
	}

	return errs
}

type familyRunParams struct {
	Notifier FamilyNotifier
	Results  *families.Results
	ErrCh    chan<- error
}

type familyRunner[T any] struct {
	family families.Family[T]
}

func newFamilyRunner[T any](family families.Family[T]) *familyRunner[T] {
	return &familyRunner[T]{
		family: family,
	}
}

func (r *familyRunner[T]) Run(ctx context.Context, params familyRunParams) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	// Notify about start, return preemptively if it fails
	if err := params.Notifier.FamilyStarted(ctx, r.family.GetType()); err != nil {
		params.ErrCh <- fmt.Errorf("family started notification failed: %w", err)
		return
	}

	// Run family
	result, err := r.family.Run(ctx, params.Results)
	familyResult := FamilyResult{
		Result:     result,
		FamilyType: r.family.GetType(),
		Err:        err,
	}

	// Handle family result depending on returned data
	logger.Debugf("Received result from family %q: %v", r.family.GetType(), familyResult)
	if err != nil {
		logger.Errorf("Received error result from family %q: %v", r.family.GetType(), err)

		// Submit run error so that we can check if the errors on channel are from
		// notifiers or from the actual family run
		params.ErrCh <- &runnerFamilyRunError{
			Family: r.family.GetType(),
			Err:    err,
		}
	} else {
		params.Results.SetFamilyResult(result)
	}

	// Notify about finish
	if err := params.Notifier.FamilyFinished(ctx, familyResult); err != nil {
		params.ErrCh <- fmt.Errorf("family finished notification failed: %w", err)
	}
}

type runnerFamilyRunError struct {
	Family families.FamilyType
	Err    error
}

func (e *runnerFamilyRunError) Error() string {
	return fmt.Sprintf("family %s finished with error: %v", e.Family, e.Err)
}

func newWorkflowTask[T any](family families.Family[T], name string, deps []string) workflowtypes.Task[familyRunParams] {
	return workflowtypes.Task[familyRunParams]{
		Name: name,
		Deps: deps,
		Fn: func(ctx context.Context, params familyRunParams) error {
			newFamilyRunner(family).Run(ctx, params)
			return nil
		},
	}
}
