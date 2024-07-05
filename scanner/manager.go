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
	tasks  []workflowtypes.Task[workflowParams]
}

func New(config *Config) *Manager {
	manager := &Manager{
		config: config,
	}

	// Analyzers
	if config.SBOM.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("sbom", sbom.New(config.SBOM)))
	}

	// Scanners
	if config.Vulnerabilities.Enabled {
		// must run after SBOM to support the case when configured to use the output from sbom
		var deps []string
		if config.SBOM.Enabled {
			deps = append(deps, "sbom")
		}

		manager.tasks = append(manager.tasks, newWorkflowTaskFor("vulnerabilities", vulnerabilities.New(config.Vulnerabilities), deps...))
	}
	if config.Secrets.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("secrets", secrets.New(config.Secrets)))
	}
	if config.Rootkits.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("rootkits", rootkits.New(config.Rootkits)))
	}
	if config.Malware.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("malware", malware.New(config.Malware)))
	}
	if config.Misconfiguration.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("misconfiguration", misconfiguration.New(config.Misconfiguration)))
	}
	if config.InfoFinder.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("infofinder", infofinder.New(config.InfoFinder)))
	}
	if config.Plugins.Enabled {
		manager.tasks = append(manager.tasks, newWorkflowTaskFor("plugins", plugins.New(config.Plugins)))
	}

	// Enrichers.
	if config.Exploits.Enabled {
		// must run after Vulnerabilities to support the case when configured to use the output from Vulnerabilities
		var deps []string
		if config.Vulnerabilities.Enabled {
			deps = append(deps, "vulnerabilities")
		}

		manager.tasks = append(manager.tasks, newWorkflowTaskFor("exploits", exploits.New(config.Exploits), deps...))
	}

	return manager
}

func (m *Manager) Run(ctx context.Context, notifier FamilyNotifier) []error {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	// Register container cache
	utils.ContainerRootfsCache = containerrootfs.NewCache()
	defer func() {
		err := utils.ContainerRootfsCache.CleanupAll()
		if err != nil {
			logger.WithError(err).Errorf("failed to cleanup all cached container rootfs files")
		}
	}()

	// Define channel to use to listen for all processing errors
	errCh := make(chan error)

	// Run task processor in the background
	go func() {
		// Close error channel to allow listener to exit properly
		defer close(errCh)

		// Create families processor
		processor, err := workflow.New[workflowParams, workflowtypes.Task[workflowParams]](m.tasks)
		if err != nil {
			errCh <- fmt.Errorf("failed to create families processor: %w", err)
			return
		}

		// Run families processor and wait until all family workflow tasks have completed
		err = processor.Run(ctx, workflowParams{
			Notifier: notifier,
			Results:  families.NewResults(),
			ErrCh:    errCh,
		})
		if err != nil {
			errCh <- fmt.Errorf("failed to run families processor: %w", err)
		}
	}()

	var oneOrMoreFamilyFailed bool
	var errs []error

	// Listen for processing errors
	for err := range errCh {
		if err == nil {
			continue
		}

		// Check if family run failed, otherwise add the error to the slice
		var familyErr *familyFailedError
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

// workflowParams defines parameters for familyRunner workflow tasks
type workflowParams struct {
	Notifier FamilyNotifier
	Results  *families.Results
	ErrCh    chan<- error
}

// newWorkflowTaskFor returns a wrapped familyRunner as a workflow task
func newWorkflowTaskFor[T any](name string, family families.Family[T], deps ...string) workflowtypes.Task[workflowParams] {
	return workflowtypes.Task[workflowParams]{
		Name: name,
		Deps: deps,
		Fn: func(ctx context.Context, params workflowParams) error {
			// Execute family using family runner and forward collected errors
			errs := newFamilyRunner(family).Run(ctx, params.Notifier, params.Results)
			for _, err := range errs {
				params.ErrCh <- err
			}

			return nil
		},
	}
}
