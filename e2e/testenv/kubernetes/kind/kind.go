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

package kind

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"

	"github.com/openclarity/vmclarity/e2e/testenv/kubernetes/common"
	envtypes "github.com/openclarity/vmclarity/e2e/testenv/types"
)

type KindEnv struct {
	name           string
	provider       *cluster.Provider
	kindConfigPath string
	kubeConfigPath string
}

const (
	KindClusterPrefix          = "vmclarity-e2e"
	KindConfigFilePath         = "testenv/kubernetes/kind/kind-config.yaml"
	KindClusterCreationTimeout = 2 * time.Minute
)

// nolint:wrapcheck
func New(_ *envtypes.Config) (*KindEnv, error) {
	return &KindEnv{
		name:           common.RandomName(KindClusterPrefix, 8),
		kindConfigPath: KindConfigFilePath,
	}, nil
}

// nolint:wrapcheck
func (e *KindEnv) Start(ctx context.Context) error {
	return common.DeployHelmChart(e.kubeConfigPath)
}

// nolint:wrapcheck
func (e *KindEnv) Stop(ctx context.Context) error {
	return nil
}

func (e *KindEnv) SetUp(_ context.Context) error {
	providerOpts, err := cluster.DetectNodeProvider()
	if err != nil {
		return fmt.Errorf("failed to detect provider: %w", err)
	}
	provider := cluster.NewProvider(providerOpts)
	if err := provider.Create(e.name,
		cluster.CreateWithConfigFile(e.kindConfigPath),
		cluster.CreateWithWaitForReady(KindClusterCreationTimeout),
	); err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}
	e.provider = provider
	kubeConfigPath, err := provider.KubeConfig(e.name, false)
	if err != nil {
		return fmt.Errorf("failed to get kube config cluster: %w", err)
	}
	e.kindConfigPath = kubeConfigPath

	return nil
}

func (e *KindEnv) TearDown(_ context.Context) error {
	return e.provider.Delete(e.name, e.kubeConfigPath)
}

// nolint:wrapcheck
func (e *KindEnv) ServicesReady(ctx context.Context) (bool, error) {
	return true, nil
}

// nolint:wrapcheck
func (e *KindEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, stderr io.Writer) error {
	return nil
}

func (e *KindEnv) Services() []string {
	return nil
}

func (e *KindEnv) VMClarityAPIURL() (*url.URL, error) {
	return nil, nil
}

func (e *KindEnv) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, "", "")
}

func loadContainerImagesToCluster(cluster string) error {
	return nil
}

func loadContainerImageToCluster(cluster, image string) error {
	return nil
}
