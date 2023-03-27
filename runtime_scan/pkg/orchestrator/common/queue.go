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
	"fmt"
	"sync"
)

type Queue[T comparable] struct {
	itemAdded chan struct{}
	queue     []T
	inqueue   map[T]struct{}
	l         sync.Mutex
}

func NewQueue[T comparable]() *Queue[T] {
	return &Queue[T]{
		itemAdded: make(chan struct{}),
		queue:     make([]T, 0),
		inqueue:   map[T]struct{}{},
	}
}

func (q *Queue[T]) Get(ctx context.Context) (T, error) {
	if len(q.queue) == 0 {
		// If the queue is empty, block waiting for the itemAdded
		// notification or context timeout.
		select {
		case <-q.itemAdded:
			// continue
		case <-ctx.Done():
			var empty T
			return empty, fmt.Errorf("failed to get item: %w", ctx.Err())
		}
	} else {
		// If the queue isn't empty, consume any item added notification
		// so that its reset for the empty case
		select {
		case <-q.itemAdded:
		default:
		}
	}

	q.l.Lock()
	defer q.l.Unlock()

	item := q.queue[0]
	q.queue = q.queue[1:]
	delete(q.inqueue, item)

	return item, nil
}

func (q *Queue[T]) Add(item T) {
	q.l.Lock()
	defer q.l.Unlock()

	if _, ok := q.inqueue[item]; !ok {
		q.queue = append(q.queue, item)
		q.inqueue[item] = struct{}{}
	}

	select {
	case q.itemAdded <- struct{}{}:
	default:
	}
}

func (q *Queue[T]) Length() int {
	q.l.Lock()
	defer q.l.Unlock()

	return len(q.queue)
}

func (q *Queue[T]) Has(item T) bool {
	q.l.Lock()
	defer q.l.Unlock()

	_, ok := q.inqueue[item]
	return ok
}
