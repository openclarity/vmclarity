// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package config

import (
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
)

const (
	ScannerJobResultListenPort           = "SCANNER_JOB_RESULT_LISTEN_PORT"
	CredsSecretNamespace                 = "CREDS_SECRET_NAMESPACE" // nolint: gosec
	ScannerJobTemplateConfigMapName      = "SCANNER_JOB_TEMPLATE_CONFIG_MAP_NAME"
	ScannerJobTemplateConfigMapNamespace = "SCANNER_JOB_TEMPLATE_CONFIG_MAP_NAMESPACE"
	defaultScannerJobResultListenPort    = 8888
)

type Config struct {
	ScannerJobResultListenPort int
	CredsSecretNamespace       string
	Region string
	VpcID string
	SubnetID string
	AmiID string
}

func setConfigDefaults() {
	viper.SetDefault(CredsSecretNamespace, "kubeclarity")
	viper.SetDefault(ScannerJobTemplateConfigMapName, "")
	viper.SetDefault(ScannerJobTemplateConfigMapNamespace, "kubeclarity")
	viper.SetDefault(ScannerJobResultListenPort, defaultScannerJobResultListenPort)

	viper.AutomaticEnv()
}

func LoadConfig(clientset kubernetes.Interface) (*Config, error) {
	setConfigDefaults()



	config := &Config{
		ScannerJobResultListenPort: viper.GetInt(ScannerJobResultListenPort),
		CredsSecretNamespace:       viper.GetString(CredsSecretNamespace),
	}

	return config, nil
}
