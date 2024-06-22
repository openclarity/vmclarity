// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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
	"time"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/provider"
)

const (
	DefaultPollInterval        = 30 * time.Second
	DefaultReconcileTimeout    = 5 * time.Minute
	DefaultSummaryUpdatePeriod = 5 * time.Minute
)

type Config struct {
	Client              *client.Client
	Provider            provider.Provider
	PollPeriod          time.Duration `mapstructure:"poll_period"`
	ReconcileTimeout    time.Duration `mapstructure:"reconcile_timeout"`
	SummaryUpdatePeriod time.Duration `mapstructure:"summary_update_period"`
}

func (c Config) WithBackendClient(client *client.Client) Config {
	c.Client = client
	return c
}

func (c Config) WithProviderClient(p provider.Provider) Config {
	c.Provider = p
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

func (c Config) WithSummaryUpdatePeriod(t time.Duration) Config {
	c.SummaryUpdatePeriod = t
	return c
}
