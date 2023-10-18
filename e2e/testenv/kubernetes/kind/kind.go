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
	"path/filepath"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/client"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/fs"

	"github.com/openclarity/vmclarity/e2e/testenv/kubernetes/common"
	"github.com/openclarity/vmclarity/e2e/testenv/kubernetes/utils"
	envtypes "github.com/openclarity/vmclarity/e2e/testenv/types"
)

type KindEnv struct {
	name           string
	provider       *cluster.Provider
	kindConfigPath string
	kubeConfigPath string
	ChartHelper    *common.ChartHelper
	k8sClientSet   kubernetes.Interface
}

const (
	KindClusterName            = "vmclarity-e2e"
	KindConfigFilePath         = "testenv/kubernetes/kind/kind-config.yaml"
	KindClusterCreationTimeout = 2 * time.Minute
	KindAPIServerPort          = "30000"
)

// nolint:wrapcheck
func New(_ *envtypes.Config) (*KindEnv, error) {
	return &KindEnv{
		name:           KindClusterName,
		kindConfigPath: KindConfigFilePath,
		kubeConfigPath: path.Join(os.TempDir(), KindClusterName),
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

	k8sClientSet, err := utils.CreateK8sClient(e.kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create k8s clientset: %w", err)
	}
	e.k8sClientSet = k8sClientSet

	return common.CreateTestDeployment(e.Context(ctx))
}

// nolint:wrapcheck
func (e *KindEnv) Stop(ctx context.Context) error {
	if err := common.DeleteTestDeployment(e.Context(ctx)); err != nil {
		// TODO (pebalogh) maybe just log ???
		return fmt.Errorf("failed to delete test deployment: %w", err)
	}
	if err := e.ChartHelper.DeleteHelmChart(); err != nil {
		// TODO (pebalogh) maybe just log ???
		return fmt.Errorf("failed to delete VMclarity helm chart: %w", err)
	}

	return nil
}

// nolint:wrapcheck
func (e *KindEnv) SetUp(ctx context.Context) error {
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

	return e.loadContainerImagesToCluster(ctx)
}

// nolint:wrapcheck
func (e *KindEnv) TearDown(_ context.Context) error {
	return e.provider.Delete(e.name, e.kubeConfigPath)
}

// nolint:wrapcheck
func (e *KindEnv) ServicesReady(ctx context.Context) (bool, error) {
	// Note (pebalogh) checking services isn't necessary because of the `helm install --wait`
	return true, nil
}

// nolint:wrapcheck
func (e *KindEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, _ io.Writer) error {
	ctx = e.Context(ctx)
	for _, podName := range services {
		pod, err := utils.GetVMClarityPodByName(ctx, podName, common.VMClarityReleaseName)
		if err != nil {
			return fmt.Errorf("failed to get pod: %w", err)
		}
		logBytes, err := utils.GetPodLogs(ctx, pod, startTime)
		if err != nil {
			return fmt.Errorf("failed to get log for pod %s: %w", podName, err)
		}
		_, err = stdout.Write(logBytes)
		if err != nil {
			return fmt.Errorf("failed to write logs to stdout: %w", err)
		}
	}

	return nil
}

// nolint:wrapcheck
func (e *KindEnv) Services(ctx context.Context) ([]string, error) {
	podList, err := utils.ListVMClarityPods(e.Context(ctx), common.VMClarityNamespace)
	if err != nil {
		return nil, fmt.Errorf("failed to list VMClarity pods: %w", err)
	}

	podNames := make([]string, 0)
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.GetName())
	}

	return podNames, nil
}

// nolint:wrapcheck
func (e *KindEnv) VMClarityAPIURL() (*url.URL, error) {
	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s/api", KindAPIServerPort),
	}, nil
}

// nolint:wrapcheck
func (e *KindEnv) Context(ctx context.Context) context.Context {
	return context.WithValue(ctx, envtypes.KubernetesContextKey, e.k8sClientSet)
}

// nolint:wrapcheck
func (e *KindEnv) loadContainerImagesToCluster(ctx context.Context) error {
	nodeList, err := e.provider.ListNodes(e.name)
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	imagesMap := common.GetImageList()
	images := make([]string, 0)
	for _, image := range imagesMap {
		images = append(images, image)
	}

	// Setup the tar path where the images will be saved
	dir, err := fs.TempDir("", "images-tar")
	if err != nil {
		return fmt.Errorf("failed to create tempdir: %w", err)
	}
	defer os.RemoveAll(dir)
	imagesTarPath := filepath.Join(dir, "images.tar")
	// Save the images into a tar
	err = save(ctx, images, imagesTarPath)
	if err != nil {
		return fmt.Errorf("failed to save images to tar archive: %w", err)
	}

	// Load the images on the selected nodes
	fns := []func() error{}
	for _, selectedNode := range nodeList {
		selectedNode := selectedNode // capture loop variable
		fns = append(fns, func() error {
			return loadContainerImageToNode(selectedNode, imagesTarPath)
		})
	}

	return errors.UntilErrorConcurrent(fns)
}

// nolint:wrapcheck
func loadContainerImageToNode(node nodes.Node, imagesTarPath string) error {
	f, err := os.Open(imagesTarPath)
	if err != nil {
		return fmt.Errorf("failed to open image tar file: %w", err)
	}
	defer f.Close()

	return nodeutils.LoadImageArchive(node, f)
}

func save(ctx context.Context, images []string, tarName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}

	responseBody, err := cli.ImageSave(ctx, images)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	defer responseBody.Close()

	if err := command.CopyToFile(tarName, responseBody); err != nil {
		return fmt.Errorf("failed to copy image to tar file: %w", err)
	}

	return nil
}
