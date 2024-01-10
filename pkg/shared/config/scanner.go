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
)

type Scanner struct {
	ScannersList []string           `yaml:"scanners_list" mapstructure:"scanners_list"`
	GrypeConfig  GrypeConfig        `yaml:"grype_config" mapstructure:"grype_config"`
	TrivyConfig  ScannerTrivyConfig `yaml:"trivy_config" mapstructure:"trivy_config"`
}

const (
	ScannersList = "SCANNERS_LIST"
)

func setScannerConfigDefaults() {
	viper.SetDefault(ScannersList, []string{"grype"})
}

func LoadScannerConfig() *Scanner {
	setScannerConfigDefaults()
	return &Scanner{
		ScannersList: viper.GetStringSlice(ScannersList),
		GrypeConfig:  LoadGrypeConfig(),
		TrivyConfig:  LoadScannerTrivyConfig(),
	}
}
