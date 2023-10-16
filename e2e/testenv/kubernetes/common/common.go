package common

import (
	"fmt"
	"math/rand"
	"os/exec"
)

const (
	VMClarityChartPath   = "../charts/vmclarity"
	HelmDriverEnvVar     = "HELM_DRIVER"
	VMClarityNamespace   = "vmclarity"
	VMClarityReleaseName = "vmaclarity-e2e"
	KubernetesProvider   = "kubernetes"
)

func RandomName(prefix string, length int) string {
	chars := "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return prefix + "-" + string(result)
}

func DeployHelmChart(kubeConfigPath string) error {
	// The commented out because of the https://github.com/helm/helm/issues/12357
	// before finding the proper solution we are using command to deploy helm chart

	//chart, err := loader.LoadDir(VMClarityChartPath)
	//if err != nil {
	//	return fmt.Errorf("failed to load VMClarity helm chart: %w", err)
	//}
	//
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
	//	return fmt.Errorf("failed to init action configuration: %w", err)
	//}
	//
	//client := action.NewInstall(actionConfig)
	//client.ReleaseName = VMClarityReleaseName
	//client.Namespace = namespace
	//
	//if _, err := client.Run(chart, createValues()); err != nil {
	//	return fmt.Errorf("failed to install VMClarity helm chart: %w", err)
	//}

	cmd := exec.Command("helm", "install", VMClarityReleaseName,
		"--namespace", VMClarityNamespace,
		"--create-namespace",
		VMClarityChartPath,
		"--kubeconfig", kubeConfigPath,
		"--set", "orchestrator.provider=kubernetes",
		"--wait",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {

		return fmt.Errorf("failed to install VMClarity helm chart: %w, %s", err, string(output))
	}

	return nil
}

func createValues() map[string]interface{} {
	return map[string]interface{}{
		"orchestrator": map[string]interface{}{
			"provider": KubernetesProvider,
		},
	}
}
