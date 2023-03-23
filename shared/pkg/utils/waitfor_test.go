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

package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

type State struct {
	iterations int
}

func conditionFailing() ConditionFunc {
	return func(_ context.Context) (bool, error) {
		return false, errors.New("condition failed")
	}
}

func conditionNeverMet() ConditionFunc {
	return func(_ context.Context) (bool, error) {
		return false, nil
	}
}

func conditionMetAfter(s *State, iteration int) ConditionFunc {
	return func(_ context.Context) (bool, error) {
		if s.iterations == iteration {
			return true, nil
		}
		s.iterations++
		return false, nil
	}
}

func TestWaitFor(t *testing.T) {
	t.Run("returns an error when context deadline is exceeded", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := WaitFor(ctx, conditionNeverMet(), 10*time.Second, 2*time.Second)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("returns an error when parent context deadline is exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := WaitFor(ctx, conditionNeverMet(), 10*time.Second, 2*time.Second)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("returns an error when parent context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(5 * time.Second)
			defer cancel()
		}()

		err := WaitFor(ctx, conditionNeverMet(), 10*time.Second, 2*time.Second)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("succeeds with timeout set", func(t *testing.T) {
		s := &State{0}
		expectedIteration := 2
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := WaitFor(ctx, conditionMetAfter(s, expectedIteration), 10*time.Second, 2*time.Second)
		assert.NilError(t, err)
		assert.Assert(t, s.iterations == expectedIteration)
	})

	t.Run("succeeds without timeout set", func(t *testing.T) {
		s := &State{0}
		expectedIteration := 2
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := WaitFor(ctx, conditionMetAfter(s, expectedIteration), 0, 2*time.Second)
		assert.NilError(t, err)
		assert.Assert(t, s.iterations == expectedIteration)
	})

	t.Run("succeeds with timeout set in parent context", func(t *testing.T) {
		s := &State{0}
		expectedIteration := 2
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		err := WaitFor(ctx, conditionMetAfter(s, expectedIteration), 10*time.Second, 2*time.Second)
		assert.NilError(t, err)
		assert.Assert(t, s.iterations == expectedIteration)
	})

	t.Run("returns error if evaluating condition returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := WaitFor(ctx, conditionFailing(), 10*time.Second, 2*time.Second)
		assert.ErrorContains(t, err, "condition failed")
	})

	t.Run("returns validation error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := WaitFor(ctx, conditionNeverMet(), 2*time.Second, 10*time.Second)
		assert.ErrorIs(t, err, ErrWaitForInValidParameters)
	})
}
