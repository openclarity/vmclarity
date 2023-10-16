package common

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"github.com/docker/distribution/reference"
	"helm.sh/helm/v3/pkg/action"
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
	ActionConfig   *action.Configuration
	Namespace      string
	KubeConfigPath string
	ReleaseName    string
}

func RandomName(prefix string, length int) string {
	chars := "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return prefix + "-" + string(result)
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
	cmd := exec.Command("helm", "install", c.ReleaseName,
		VMClarityChartPath,
		"--namespace", VMClarityNamespace,
		"--create-namespace",
		"--kubeconfig", c.KubeConfigPath,
		"--set", "orchestrator.provider=kubernetes",
		"--set", "orchestrator.serviceAccount.automountServiceAccountToken=true",
		"--wait",
	)

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
	cmd := exec.Command("helm", "delete", c.ReleaseName,
		"--namespace", VMClarityNamespace,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete VMClarity helm chart: %w, %s", err, string(output))
	}

	return nil
}

func createValues(imageList map[string]string) (map[string]interface{}, error) {
	var parsedImageList map[string]map[string]string
	var err error
	for k, v := range imageList {
		parsedImageList[k], err = getImageRegistryRepositoryTag(v)
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

func getImageRegistryRepositoryTag(image string) (map[string]string, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image: %s", image)
	}

	registry := reference.Domain(named)
	repository := reference.Path(named)
	tagged, ok := named.(reference.Tagged)
	if !ok {
		return nil, fmt.Errorf("failed to get image tag from image name: %s", image)
	}
	tag := tagged.Tag()

	return map[string]string{
		"registry":   registry,
		"repository": repository,
		"tag":        tag,
	}, nil
}
