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

package external

import (
	"context"
	"io"
	"net/url"
	"time"

	envtypes "github.com/openclarity/vmclarity/e2e/testenv/types"
)

type KubernetesEnv struct {
	name string
}

// nolint:wrapcheck
func New(_ *envtypes.Config) (*KubernetesEnv, error) {
	return nil, nil
}

// nolint:wrapcheck
func (e *KubernetesEnv) Start(ctx context.Context) error {
	return nil
}

// nolint:wrapcheck
func (e *KubernetesEnv) Stop(ctx context.Context) error {
	return nil
}

func (e *KubernetesEnv) SetUp(_ context.Context) error {
	// NOTE(chrisgacsal): nothing to do
	return nil
}

func (e *KubernetesEnv) TearDown(_ context.Context) error {
	// NOTE(chrisgacsal): nothing to do
	return nil
}

// nolint:wrapcheck
func (e *KubernetesEnv) ServicesReady(ctx context.Context) (bool, error) {
	return true, nil
}

// nolint:wrapcheck
func (e *KubernetesEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, stderr io.Writer) error {
	return nil
}

func (e *KubernetesEnv) Services() []string {
	return nil
}

func (e *KubernetesEnv) VMClarityAPIURL() (*url.URL, error) {
	return nil, nil
}

func (e *KubernetesEnv) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, "", "")
}
