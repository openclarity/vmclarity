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
	"fmt"
	"time"
)

var ErrWaitForInValidParameters = errors.New("condition will never get evaluated as timeout < interval")

// ConditionFunc is a function which returns ok indicating whether the condition is met or not.
// It also returns a non-nil error if the error occurs during evaluating the condition.
type ConditionFunc func(ctx context.Context) (bool, error)

// WaitFor takes f function and periodically runs with interval until:
//   - f returns true or error
//   - parent ctx is cancelled
//   - timeout is reached
//
// Setting timeout to <= 0 means that f condition will be continuously evaluated until f returns true or an error.
// Providing timeout and interval parameters returns ErrWaitForInValidParameters error if timeout > 0 and interval > timeout.
//
// Wait for condition with timeout:
//
//	f := func(ctx context.Context) (bool, error) {
//		expected := 5
//		actual := rand.Intn(10)
//		if actual == expected {
//			return true, nil
//		}
//		return false, nil
//	}
//	err := utils.WaitFor(ctx, f, time.Minute, 10 * time.Second)
//	if err != nil {
//		panic(err)
//	}
//
// Wait for condition without timeout:
//
//	err := utils.WaitFor(ctx, f, 0, 10 * time.Second)
//	if err != nil {
//		panic(err)
//	}
func WaitFor(ctx context.Context, f ConditionFunc, timeout time.Duration, interval time.Duration) error {
	var cancelFn context.CancelFunc

	if timeout > 0 && timeout < interval {
		return ErrWaitForInValidParameters
	}

	if timeout <= 0 {
		ctx, cancelFn = context.WithCancel(ctx)
	} else {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	defer cancelFn()

	timer := time.NewTimer(interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			{
				cond, err := f(ctx)
				if err != nil {
					return fmt.Errorf("condition failed: %w", err)
				}
				if cond {
					return nil
				}
				timer.Reset(interval)
			}
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return fmt.Errorf("waiting for condition was cancelled: %w", ctx.Err())
		}
	}
}
