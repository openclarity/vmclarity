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

package docker

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultEnvPrefix = "VMCLARITY_DOCKER"
)

type Config struct {
	// ScannerImage is a container image used for creating Scanner container
	ScannerImage string `mapstructure:"scanner_image"`
}

func (c *Config) Validate() error {
	if c.ScannerImage == "" {
		return fmt.Errorf("parameter ScannerImage must be provided")
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Avoid modifying the global instance
	v := viper.New()

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("scanner_image")

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to parse provider configuration. Provider=Docker: %w", err)
	}

	return config, nil
}
