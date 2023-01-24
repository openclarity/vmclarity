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

package vulnerabilities

import (
	"github.com/openclarity/kubeclarity/shared/pkg/config"

	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

type Config struct {
	Enabled        bool           `yaml:"enabled" mapstructure:"enabled"`
	ScannersList   []string       `yaml:"scanners_list" mapstructure:"scanners_list"`
	Inputs         []Input        `yaml:"inputs" mapstructure:"inputs"`
	InputFromSbom  bool           `yaml:"input_from_sbom" mapstructure:"input_from_sbom"`
	ScannersConfig *config.Config `yaml:"scanners_config" mapstructure:"scanners_config"`
}

type Input struct {
	Input     string `yaml:"input" mapstructure:"input"`
	InputType string `yaml:"input_type" mapstructure:"input_type"`
}

type InputFromFamily struct {
	FamilyType types.FamilyType `yaml:"family_type" mapstructure:"family_type"`
}
