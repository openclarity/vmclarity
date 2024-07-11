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

package types

import (
	grypeconfig "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/grype/config"
	trivyconfig "github.com/openclarity/vmclarity/scanner/families/vulnerabilities/trivy/config"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
)

type Config struct {
	Enabled        bool                     `yaml:"enabled" mapstructure:"enabled"`
	ScannersList   []string                 `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs         []scannertypes.ScanInput `yaml:"inputs" mapstructure:"inputs"`
	InputFromSbom  bool                     `yaml:"input_from_sbom" mapstructure:"input_from_sbom"`
	ScannersConfig ScannersConfig           `yaml:"scanners_config" mapstructure:"scanners_config"`
}

type ScannersConfig struct {
	Grype grypeconfig.Config `yaml:"grype" mapstructure:"grype"`
	Trivy trivyconfig.Config `yaml:"trivy" mapstructure:"trivy"`
}
