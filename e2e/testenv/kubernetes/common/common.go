package common

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const VMClarityChartPath = "../charts/vmclarity"

func RandomName(prefix string, length int) string {
	chars := "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return prefix + "-" + string(result)
}

func DeployHelmChart(kubeConfigPath string) error {
	chart, err := loader.LoadDir(VMClarityChartPath)
	if err != nil {
		return fmt.Errorf("failed to load VMClarity helm chart: %w", err)
	}

	actionConfig := new(action.Configuration)
	namespace := "default"
	restClientGetter := genericclioptions.NewConfigFlags(true)
	restClientGetter.Namespace = &namespace
	restClientGetter.KubeConfig = &kubeConfigPath
	if err := actionConfig.Init(
		restClientGetter,
		namespace,
		os.Getenv("HELM_DRIVER"),
		logrus.Printf,
	); err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = "test-vmclarity"
	client.Namespace = namespace

	// define values
	vals := map[string]interface{}{
		"orchestrator": map[string]interface{}{
			"provider": "kubernetes",
		},
	}
	if _, err := client.Run(chart, vals); err != nil {
		return err
	}

	return nil
}
