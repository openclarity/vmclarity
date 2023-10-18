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
	"fmt"
	"os"
	"os/exec"

	"github.com/openclarity/vmclarity/e2e/testenv/kubernetes/utils"
)

const (
	VMClarityChartPath         = "../charts/vmclarity"
	HelmDriverEnvVar           = "HELM_DRIVER"
	VMClarityNamespace         = "vmclarity"
	VMClarityReleaseName       = "vmclarity-e2e"
	KubernetesProvider         = "kubernetes"
	APIServerContainerImage    = "APIServerContainerImage"
	OrchestratorContainerImage = "OrchestratorContainerImage"
	ScannerContainerImage      = "ScannerContainerImage"
	UIContainerImage           = "UIContainerImage"
	UIBackendContainerImage    = "UIBackendContainerImage"
)

type ChartHelper struct {
	//	ActionConfig   *action.Configuration
	Namespace      string
	KubeConfigPath string
	ReleaseName    string
}

func NewChartHelper(kubeConfigPath string) (*ChartHelper, error) {
	// Commented out because of the https://github.com/helm/helm/issues/12357
	// before finding the proper solution we are using command to deploy helm chart

	//actionConfig := new(action.Configuration)
	//namespace := VMClarityNamespace
	//restClientGetter := genericclioptions.NewConfigFlags(true)
	//restClientGetter.Namespace = &namespace
	//restClientGetter.KubeConfig = &kubeConfigPath
	//if err := actionConfig.Init(
	//	restClientGetter,
	//	namespace,
	//	os.Getenv(HelmDriverEnvVar),
	//	logrus.Printf,
	//); err != nil {
	//	return nil, fmt.Errorf("failed to init action configuration: %w", err)
	//}
	//
	//return &ChartHelper{
	//	ActionConfig:   actionConfig,
	//	Namespace:      namespace,
	//	KubeConfigPath: kubeConfigPath,
	//	ReleaseName:    VMClarityReleaseName,
	//}, nil

	// TODO (pebalogh) remove after issue above is solved
	return &ChartHelper{
		Namespace:      VMClarityNamespace,
		KubeConfigPath: kubeConfigPath,
		ReleaseName:    VMClarityReleaseName,
	}, nil
}

func (c *ChartHelper) DeployHelmChart() error {
	// Commented out because of the https://github.com/helm/helm/issues/12357
	// before finding the proper solution we are using command to deploy helm chart

	//chart, err := loader.LoadDir(VMClarityChartPath)
	//if err != nil {
	//	return fmt.Errorf("failed to load VMClarity helm chart: %w", err)
	//}
	//
	//client := action.NewInstall(c.ActionConfig)
	//client.ReleaseName = c.ReleaseName
	//client.Namespace = c.Namespace
	//client.Wait = true
	//
	//values, err := createValues(GetImageList())
	//if err != nil {
	//	return fmt.Errorf("failed to create values: %w", err)
	//}
	//
	//if _, err := client.Run(chart, values); err != nil {
	//	return fmt.Errorf("failed to install VMClarity helm chart: %w", err)
	//}

	// TODO (pebalogh) remove this after the issue above is solved
	parsedImageList := make(map[string]map[string]string)
	var err error
	for k, v := range GetImageList() {
		parsedImageList[k], err = utils.GetImageRegistryRepositoryTag(v)
		if err != nil {
			return fmt.Errorf("failed to parse %s image: %s", k, v)
		}
	}
	var cmdArgs []string
	cmdArgs = append(cmdArgs,
		"--set", fmt.Sprintf("apiserver.image.registry=%s", parsedImageList[APIServerContainerImage]["registry"]),
		"--set", fmt.Sprintf("apiserver.image.repository=%s", parsedImageList[APIServerContainerImage]["repository"]),
		"--set", fmt.Sprintf("apiserver.image.tag=%s", parsedImageList[APIServerContainerImage]["tag"]),
	)
	cmdArgs = append(cmdArgs,
		"--set", fmt.Sprintf("orchestrator.image.registry=%s", parsedImageList[OrchestratorContainerImage]["registry"]),
		"--set", fmt.Sprintf("orchestrator.image.repository=%s", parsedImageList[OrchestratorContainerImage]["repository"]),
		"--set", fmt.Sprintf("orchestrator.image.tag=%s", parsedImageList[OrchestratorContainerImage]["tag"]),
	)
	cmdArgs = append(cmdArgs,
		"--set", fmt.Sprintf("ui.image.registry=%s", parsedImageList[UIContainerImage]["registry"]),
		"--set", fmt.Sprintf("ui.image.repository=%s", parsedImageList[UIContainerImage]["repository"]),
		"--set", fmt.Sprintf("ui.image.tag=%s", parsedImageList[UIContainerImage]["tag"]),
	)
	cmdArgs = append(cmdArgs,
		"--set", fmt.Sprintf("uibackend.image.registry=%s", parsedImageList[UIBackendContainerImage]["registry"]),
		"--set", fmt.Sprintf("uibackend.image.repository=%s", parsedImageList[UIBackendContainerImage]["repository"]),
		"--set", fmt.Sprintf("uibackend.image.tag=%s", parsedImageList[UIBackendContainerImage]["tag"]),
	)
	cmdArgs = append(cmdArgs,
		"--set", fmt.Sprintf("orchestrator.scannerImage.registry=%s", parsedImageList[ScannerContainerImage]["registry"]),
		"--set", fmt.Sprintf("orchestrator.scannerImage.repository=%s", parsedImageList[ScannerContainerImage]["repository"]),
		"--set", fmt.Sprintf("orchestrator.scannerImage.tag=%s", parsedImageList[ScannerContainerImage]["tag"]),
	)

	args := []string{
		"install",
		c.ReleaseName,
		VMClarityChartPath,
		"--namespace", VMClarityNamespace,
		"--create-namespace",
		"--kubeconfig", c.KubeConfigPath,
		"--wait",
		"--set", "orchestrator.provider=kubernetes",
		"--set", "orchestrator.serviceAccount.automountServiceAccountToken=true",
		"--set", "gateway.service.type=NodePort",
		"--set", "gateway.service.nodePort=30000",
	}
	args = append(args, cmdArgs...)

	cmd := exec.Command("helm", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install VMClarity helm chart: %w, %s", err, string(output))
	}

	return nil
}

func (c *ChartHelper) DeleteHelmChart() error {
	// Commented out because of the https://github.com/helm/helm/issues/12357
	// before finding the proper solution we are using command to deploy helm chart

	//uninstall := action.NewUninstall(c.ActionConfig)
	//if _, err := uninstall.Run(c.ReleaseName); err != nil {
	//	return fmt.Errorf("failed to delete VMClarity helm chart: %w", err)
	//}

	// TODO (pebalogh) remove this after the issue above is solved
	cmd := exec.Command("helm", "delete", c.ReleaseName, "--namespace", VMClarityNamespace) // nolint: gosec

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete VMClarity helm chart: %w, %s", err, string(output))
	}

	return nil
}

// nolint: unused
func createValues(imageList map[string]string) (map[string]interface{}, error) {
	parsedImageList := make(map[string]map[string]string)
	var err error
	for k, v := range imageList {
		parsedImageList[k], err = utils.GetImageRegistryRepositoryTag(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s image: %s", k, v)
		}
	}

	return map[string]interface{}{
		"apiserver": map[string]interface{}{
			"image": parsedImageList[APIServerContainerImage],
		},
		"orchestrator": map[string]interface{}{
			"provider": KubernetesProvider,
			"serviceAccount": map[string]interface{}{
				"automountServiceAccountToken": true,
			},
			"image":        parsedImageList[OrchestratorContainerImage],
			"scannerImage": parsedImageList[ScannerContainerImage],
		},
		"ui": map[string]interface{}{
			"image": parsedImageList[UIContainerImage],
		},
		"uibackend": map[string]interface{}{
			"image": parsedImageList[UIBackendContainerImage],
		},
		// Note (pebalogh) this is for only kind at the moment.
		// If we implement external K8S we will need to set different values.
		"gateway": map[string]interface{}{
			"service": map[string]interface{}{
				"type":     "NodePort",
				"nodePort": "30000",
			},
		},
	}, nil
}

func GetImageList() map[string]string {
	return map[string]string{
		APIServerContainerImage:    os.Getenv(APIServerContainerImage),
		OrchestratorContainerImage: os.Getenv(OrchestratorContainerImage),
		ScannerContainerImage:      os.Getenv(ScannerContainerImage),
		UIContainerImage:           os.Getenv(UIContainerImage),
		UIBackendContainerImage:    os.Getenv(UIBackendContainerImage),
	}
}
