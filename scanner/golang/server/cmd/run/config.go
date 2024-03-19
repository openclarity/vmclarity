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

package run

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	DefaultEnvPrefix     = "VMCLARITY_SCANNER_SERVER"
	DefaultListenAddress = "0.0.0.0:8765"
	DefaultLogLevel      = "info"
)

type Config struct {
	ListenAddress string `json:"listen-address,omitempty" mapstructure:"listen_address"`
	LogLevel      string `json:"log-level,omitempty" mapstructure:"log_level"`
}

func NewConfig() (*Config, error) {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("."),
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
	)

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("listen_address")
	v.SetDefault("listen_address", DefaultListenAddress)

	_ = v.BindEnv("log_level")
	v.SetDefault("log_level", DefaultLogLevel)

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to load API Server configuration: %w", err)
	}

	return config, nil
}
