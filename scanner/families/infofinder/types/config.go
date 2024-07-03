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
	sshtopologyconfig "github.com/openclarity/vmclarity/scanner/families/infofinder/sshtopology/config"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
)

type Config struct {
	Enabled         bool                     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	ScannersList    []string                 `yaml:"scanners_list" mapstructure:"scanners_list"`
	StripInputPaths bool                     `yaml:"strip_input_paths" mapstructure:"strip_input_paths"`
	Inputs          []scannertypes.ScanInput `yaml:"inputs" mapstructure:"inputs"`
	ScannersConfig  ScannersConfig           `yaml:"scanners_config" mapstructure:"scanners_config"`
}

type ScannersConfig struct {
	SSHTopology sshtopologyconfig.Config `yaml:"ssh_topology" mapstructure:"ssh_topology"`
}
