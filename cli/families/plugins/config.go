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

import "github.com/openclarity/vmclarity/cli/families/types"

type Config struct {
	PluginList []Plugin      `yaml:"plugins_list" mapstructure:"plugins_list"`
	Inputs     []types.Input `yaml:"inputs" mapstructure:"inputs"`
}

type Plugin struct {
	// Enabled is a flag that determines if the plugin scanner is enabled
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
	// Name is the name of the plugin scanner
	Name string `yaml:"name" mapstructure:"name"`
	// ImageName is the name of the docker image that will be used to run the plugin scanner
	ImageName string `yaml:"image_name" mapstructure:"image_name"`
	// Output is a directory where the plugin scanner will store its results
	Output string `yaml:"output" mapstructure:"output"`
	// PluginConfig is a json string that will be passed to the plugin scanner
	PluginConfig string `yaml:"plugin_config" mapstructure:"plugin_config"`
}
