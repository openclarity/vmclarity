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
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
)

// familyRunner handles a specific family execution operations
type familyRunner[T any] struct {
	family families.Family[T]
}

func (r *familyRunner[T]) Run(ctx context.Context, notifier families.FamilyNotifier, results *families.Results) []error {
	var errs []error

	// Inject family data into logger
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", r.family.GetType())
	ctx = log.SetLoggerForContext(ctx, logger)

	// Notify about start, return preemptively if it fails since we won't be able to
	// collect scan results anyway.
	if err := notifier.FamilyStarted(ctx, r.family.GetType()); err != nil {
		errs = append(errs, fmt.Errorf("family started notification failed: %w", err))
		return errs
	}

	// Run family
	result, err := r.family.Run(ctx, results)
	familyResult := families.FamilyResult{
		Result:     result,
		FamilyType: r.family.GetType(),
		Err:        err,
	}

	// Handle family result depending on returned data
	logger.Debugf("Received result from family: %v", familyResult)
	if err != nil {
		logger.WithError(err).Errorf("Family finished with error")

		// Submit run error so that we can check if the error are from the notifier or
		// from the actual family run
		errs = append(errs, &familyFailedError{
			FamilyType: r.family.GetType(),
			Err:        err,
		})
	} else {
		logger.Info("Family finished with success")

		// Set result in shared object for the family
		results.SetFamilyResult(result)
	}

	// Notify about finish
	if err := notifier.FamilyFinished(ctx, familyResult); err != nil {
		errs = append(errs, fmt.Errorf("family finished notification failed: %w", err))
	}

	return errs
}

func newFamilyRunner[T any](family families.Family[T]) *familyRunner[T] {
	return &familyRunner[T]{family: family}
}

// familyFailedError defines families.Family run fail error
type familyFailedError struct {
	FamilyType families.FamilyType
	Err        error
}

func (e *familyFailedError) Error() string {
	return fmt.Sprintf("family %s failed with %v", e.FamilyType, e.Err)
}
