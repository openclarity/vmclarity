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
	"os"
)

const (
	EnvLogLevel      = "PLUGIN_SERVER_LOG_LEVEL"
	EnvListenAddress = "PLUGIN_SERVER_LISTEN_ADDRESS"

	DefaultLogLevel      = "info"
	DefaultListenAddress = "0.0.0.0:8080"
)

type Config struct {
	LogLevel      string `json:"log-level,omitempty" mapstructure:"log_level"`
	ListenAddress string `json:"listen-address,omitempty" mapstructure:"listen_address"`
}

func NewConfig() *Config {
	config := &Config{
		LogLevel:      DefaultLogLevel,
		ListenAddress: DefaultListenAddress,
	}

	if logLevel, ok := os.LookupEnv(EnvLogLevel); ok {
		config.LogLevel = logLevel
	}

	if listenAddress, ok := os.LookupEnv(EnvListenAddress); ok {
		config.ListenAddress = listenAddress
	}

	return config
}
