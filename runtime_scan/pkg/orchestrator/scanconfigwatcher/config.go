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

package scanconfigwatcher

import (
	"time"

	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
)

const (
	DefaultPollInterval     = 15 * time.Second
	DefaultReconcileTimeout = 5 * time.Minute
)

type Config struct {
	Backend          *backendclient.BackendClient
	PollPeriod       time.Duration
	ReconcileTimeout time.Duration
}

func (c Config) WithBackendClient(b *backendclient.BackendClient) Config {
	c.Backend = b
	return c
}

func (c Config) WithReconcileTimeout(t time.Duration) Config {
	c.ReconcileTimeout = t
	return c
}

func (c Config) WithPollPeriod(t time.Duration) Config {
	c.PollPeriod = t
	return c
}
