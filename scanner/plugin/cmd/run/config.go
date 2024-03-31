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
	DefaultEnvPrefix  = "VMCLARITY_SCANNER_SERVER"
	DefaultLogLevel   = "info"
	DefaultSocketFile = "/var/run/plugin.sock"
)

type Config struct {
	LogLevel   string `json:"log-level,omitempty" mapstructure:"log_level"`
	SocketFile string `json:"socket-file,omitempty" mapstructure:"socket_file"`
}

func NewConfig() (*Config, error) {
	v := viper.NewWithOptions(
		viper.KeyDelimiter("."),
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
	)

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("log_level")
	v.SetDefault("log_level", DefaultLogLevel)

	_ = v.BindEnv("socket_file")
	v.SetDefault("socket_file", DefaultSocketFile)

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to load API Server configuration: %w", err)
	}

	return config, nil
}
