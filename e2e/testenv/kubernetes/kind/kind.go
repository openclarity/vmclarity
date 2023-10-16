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
	"os"
	"path"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/client"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"

	"github.com/openclarity/vmclarity/e2e/testenv/kubernetes/common"
	envtypes "github.com/openclarity/vmclarity/e2e/testenv/types"
)

type KindEnv struct {
	name           string
	provider       *cluster.Provider
	kindConfigPath string
	kubeConfigPath string
	ChartHelper    *common.ChartHelper
}

const (
	KindClusterPrefix          = "vmclarity-e2e"
	KindConfigFilePath         = "testenv/kubernetes/kind/kind-config.yaml"
	KindClusterCreationTimeout = 2 * time.Minute
)

// nolint:wrapcheck
func New(_ *envtypes.Config) (*KindEnv, error) {
	kindClusterName := common.RandomName(KindClusterPrefix, 8)
	return &KindEnv{
		name:           kindClusterName,
		kindConfigPath: KindConfigFilePath,
		kubeConfigPath: path.Join(os.TempDir(), kindClusterName),
	}, nil
}

// nolint:wrapcheck
func (e *KindEnv) Start(ctx context.Context) error {
	chartHelper, err := common.NewChartHelper(e.kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create chart helper: %w", err)
	}
	e.ChartHelper = chartHelper
	if err := e.ChartHelper.DeployHelmChart(); err != nil {
		return fmt.Errorf("failed to deploy VMClarity helm chart: %w", err)
	}
	// TODO (pebalogh) deploy a test pod/deployment/etc
	return nil
}

// nolint:wrapcheck
func (e *KindEnv) Stop(ctx context.Context) error {
	// TODO (pebalogh) remove test pod/deployment/etc
	if err := e.ChartHelper.DeleteHelmChart(); err != nil {
		// TODO (pebalogh) just log
		return fmt.Errorf("failed to delete VMclarity helm chart: %w", err)
	}
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
	if err := provider.ExportKubeConfig(e.name, e.kubeConfigPath, false); err != nil {
		return fmt.Errorf("failed to get kubeconfig for kind cluster: %w", err)
	}

	return e.loadContainerImagesToCluster()
}

func (e *KindEnv) TearDown(_ context.Context) error {
	return e.provider.Delete(e.name, e.kubeConfigPath)
}

// nolint:wrapcheck
func (e *KindEnv) ServicesReady(ctx context.Context) (bool, error) {
	// TODO check if services are ready
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

	return common.GetVMClarityAPIURL(), nil
	return nil, nil
}

func (e *KindEnv) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, "", "")
}

func (e *KindEnv) loadContainerImagesToCluster() error {
	nodeList, err := e.provider.ListNodes(e.name)
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	images := common.GetImageList()

	for k, image := range images {
		imageTarName, err := save(k, image)
		if err != nil {
			return fmt.Errorf("failed to save docker image: %w", err)
		}
		if err := loadContainerImageToCluster(nodeList, imageTarName); err != nil {
			return err
		}
	}

	return nil
}

func loadContainerImageToCluster(nodeList []nodes.Node, imageTarName string) error {
	f, err := os.Open(imageTarName)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, node := range nodeList {
		return nodeutils.LoadImageArchive(node, f)
	}

	return nil
}

func save(pattern, image string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	responseBody, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer responseBody.Close()

	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for image archive: %w", err)
	}
	defer file.Close()

	if err := command.CopyToFile(file.Name(), responseBody); err != nil {
		return "", err
	}

	return file.Name(), nil
}
