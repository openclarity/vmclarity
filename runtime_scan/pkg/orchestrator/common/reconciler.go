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

package common

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type Reconciler[T comparable] struct {
	Logger *log.Entry

	// Reconcile function which will be called whenever there is an event on EventChan
	ReconcileFunction func(context.Context, T) error

	// Maximum amount of time to spend trying to reconcile one item before
	// moving onto the next item.
	ReconcileTimeout time.Duration

	// The queue which the reconciler will receive events to reconcile on.
	Queue Dequeuer[T]
}

func (r *Reconciler[T]) Start(ctx context.Context) {
	go func() {
		for {
			// queue.Get will block until an item is available to
			// return.
			item, err := r.Queue.Dequeue(ctx)
			if err != nil {
				r.Logger.Errorf("Failed to get item from queue: %v", err)
			} else {
				timeoutCtx, cancel := context.WithTimeout(ctx, r.ReconcileTimeout)
				err := r.ReconcileFunction(timeoutCtx, item)
				if err != nil {
					r.Logger.Errorf("Failed to reconcile item: %v", err)
				}
				cancel()
			}

			// Check if the parent context done if so we also need
			// to exit.
			select {
			case <-ctx.Done():
				r.Logger.Infof("Shutting down: %v", ctx.Err())
				return
			default:
			}
		}
	}()
}
