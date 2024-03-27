// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package plugins

import "github.com/openclarity/vmclarity/scanner/families/types"

type Config struct {
	Enabled        bool                     `yaml:"enabled" mapstructure:"enabled"`
	ScannersList   []string                 `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs         []types.Input            `yaml:"inputs" mapstructure:"inputs"`
	ScannersConfig map[string]ScannerConfig `yaml:"scanners_config" mapstructure:"scanners_config"`
}

type ScannerConfig struct {
	// ImageName is the name of the docker image that will be used to run the plugin scanner
	ImageName string `yaml:"image_name" mapstructure:"image_name"`
	// Output is a directory where the plugin scanner will store its results
	OutputDir string `yaml:"output_dir" mapstructure:"output_dir"`
	// Config is a json string that will be passed to the plugin scanner
	Config string `yaml:"config" mapstructure:"config"`
}
