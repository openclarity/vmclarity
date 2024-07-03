// // Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// // All rights reserved.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.
//
// // TODO(ramizpolic): improve runner usage and workflow clarity
package scanner

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/openclarity/vmclarity/scanner/families"
//	"github.com/openclarity/vmclarity/workflow"
//	workflowtypes "github.com/openclarity/vmclarity/workflow/types"
//
//	"github.com/openclarity/vmclarity/core/log"
//	"github.com/openclarity/vmclarity/scanner/families/exploits"
//	"github.com/openclarity/vmclarity/scanner/families/infofinder"
//	"github.com/openclarity/vmclarity/scanner/families/malware"
//	"github.com/openclarity/vmclarity/scanner/families/misconfiguration"
//	"github.com/openclarity/vmclarity/scanner/families/plugins"
//	"github.com/openclarity/vmclarity/scanner/families/rootkits"
//	"github.com/openclarity/vmclarity/scanner/families/sbom"
//	"github.com/openclarity/vmclarity/scanner/families/secrets"
//	"github.com/openclarity/vmclarity/scanner/families/types"
//	"github.com/openclarity/vmclarity/scanner/families/utils"
//	"github.com/openclarity/vmclarity/scanner/families/vulnerabilities"
//	"github.com/openclarity/vmclarity/scanner/utils/containerrootfs"
//)
//
//type FamilyResult struct {
//	FamilyType families.FamilyType
//	Result     any
//	Err        error
//}
//
//type FamilyNotifier interface {
//	FamilyStarted(context.Context, families.FamilyType) error
//	FamilyFinished(ctx context.Context, res FamilyResult) error
//}
//
//type Manager struct {
//	config *Config
//	tasks  []workflowtypes.Task[runner]
//}
//
//func New(config *Config) *Manager {
//	manager := &Manager{
//		config: config,
//	}
//
//	// Analyzers
//	if config.SBOM.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "sbom",
//			Fn:   withFamilyRunner(sbom.New(config.SBOM)),
//		})
//	}
//
//	// Scanners
//	if config.Vulnerabilities.Enabled {
//		// must run after SBOM to support the case when configured to use the output from sbom
//		var deps []string
//		if config.SBOM.Enabled {
//			deps = append(deps, "sbom")
//		}
//
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "vulnerabilities",
//			Deps: deps,
//			Fn:   withFamilyRunner(vulnerabilities.New(config.Vulnerabilities)),
//		})
//	}
//	if config.Secrets.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "secrets",
//			Fn:   withFamilyRunner(secrets.New(config.Secrets)),
//		})
//	}
//	if config.Rootkits.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "rootkits",
//			Fn:   withFamilyRunner(rootkits.New(config.Rootkits)),
//		})
//	}
//	if config.Malware.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "malware",
//			Fn:   withFamilyRunner(malware.New(config.Malware)),
//		})
//	}
//	if config.Misconfiguration.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "misconfiguration",
//			Fn:   withFamilyRunner(misconfiguration.New(config.Misconfiguration)),
//		})
//	}
//	if config.InfoFinder.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "infofinder",
//			Fn:   withFamilyRunner(infofinder.New(config.InfoFinder)),
//		})
//	}
//	if config.Plugins.Enabled {
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "plugins",
//			Fn:   withFamilyRunner(plugins.New(config.Plugins)),
//		})
//	}
//
//	// Enrichers.
//	if config.Exploits.Enabled {
//		// must run after Vulnerabilities to support the case when configured to use the output from Vulnerabilities
//		var deps []string
//		if config.Vulnerabilities.Enabled {
//			deps = append(deps, "vulnerabilities")
//		}
//
//		manager.tasks = append(manager.tasks, workflowtypes.Task[runner]{
//			Name: "exploits",
//			Deps: deps,
//			Fn:   withFamilyRunner(exploits.New(config.Exploits)),
//		})
//	}
//
//	return manager
//}
//
//func (m *Manager) Run(ctx context.Context, notifier FamilyNotifier) []error {
//	var oneOrMoreFamilyFailed bool
//	var errs []error
//
//	logger := log.GetLoggerFromContextOrDiscard(ctx)
//
//	// Register container cache
//	utils.ContainerRootfsCache = containerrootfs.NewCache()
//	defer func() {
//		err := utils.ContainerRootfsCache.CleanupAll()
//		if err != nil {
//			logger.WithError(err).Errorf("failed to cleanup all cached container rootfs files")
//		}
//	}()
//
//	// Define channel to use to listen for processing errors
//	errCh := make(chan error)
//
//	// Run task processor in the background
//	go func() {
//		// Close error channel to allow listener to exit properly
//		defer close(errCh)
//
//		// Create families processor
//		processor, err := workflow.New[runner, workflowtypes.Task[runner]](m.tasks)
//		if err != nil {
//			errCh <- fmt.Errorf("failed to create families processor: %w", err)
//			return
//		}
//
//		// Run families processor
//		if err := processor.Run(ctx, runner{
//			Notifier: notifier,
//			Results:  types.NewFamiliesResults(),
//			ErrCh:    errCh,
//		}); err != nil {
//			errCh <- fmt.Errorf("failed to run families processor: %w", err)
//			return
//		}
//	}()
//
//	// Listen for processing errors
//	for err := range errCh {
//		if err == nil {
//			continue
//		}
//
//		// Check if family run failed, otherwise add the error to slice
//		var familyErr *runnerFamilyRunError
//		if errors.As(err, &familyErr) {
//			oneOrMoreFamilyFailed = true
//		} else {
//			errs = append(errs, err)
//		}
//	}
//
//	if oneOrMoreFamilyFailed {
//		errs = append(errs, errors.New("at least one family failed to run"))
//	}
//
//	return errs
//}
//
//// withFamilyRunner returns a function that will handle workflow execution for
//// the given family using provided runner.
//func withFamilyRunner(family families.Family) func(context.Context, runner) error {
//	return func(ctx context.Context, runner runner) error {
//		// NOTE(ramizpolic): We do not return errors at all as returning an error in
//		// workflow function will cancel the whole execution. This is problematic as
//		// other families could still be able to run. Instead, we write all execution
//		// errors to a channel.
//		runner.Run(ctx, family)
//		return nil
//	}
//}
//
//type runner struct {
//	Notifier FamilyNotifier
//	Results  *types.FamiliesResults
//	ErrCh    chan<- error
//}
//
//func (r *runner) Run(ctx context.Context, family families.Family) {
//	logger := log.GetLoggerFromContextOrDiscard(ctx)
//
//	// Notify about start, return preemptively if it fails
//	if err := r.Notifier.FamilyStarted(ctx, family.GetType()); err != nil {
//		r.ErrCh <- fmt.Errorf("family started notification failed: %w", err)
//		return
//	}
//
//	// Run family
//	result, err := family.Run(ctx, r.Results)
//	familyResult := FamilyResult{
//		Result:     result,
//		FamilyType: family.GetType(),
//		Err:        err,
//	}
//
//	// Handle family result depending on returned data
//	logger.Debugf("Received result from family %q: %v", family.GetType(), familyResult)
//	if err != nil {
//		logger.Errorf("Received error result from family %q: %v", family.GetType(), err)
//
//		// Submit run error so that we can check if the errors on channel are from
//		// notifiers or from the actual family run
//		r.ErrCh <- &runnerFamilyRunError{
//			Family: family.GetType(),
//			Err:    err,
//		}
//	} else {
//		r.Results.SetFamilyResult(result)
//	}
//
//	// Notify about finish
//	if err := r.Notifier.FamilyFinished(ctx, familyResult); err != nil {
//		r.ErrCh <- fmt.Errorf("family finished notification failed: %w", err)
//	}
//}
//
//type runnerFamilyRunError struct {
//	Family families.FamilyType
//	Err    error
//}
//
//func (e *runnerFamilyRunError) Error() string {
//	return fmt.Sprintf("family %s finished with error: %v", e.Family, e.Err)
//}
