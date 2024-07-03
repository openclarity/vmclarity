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
	"github.com/openclarity/vmclarity/scanner/types"
	"time"
)

type Config struct {
	Timeout     int             `yaml:"timeout" mapstructure:"timeout"`
	ServerAddr  string          `yaml:"server_addr" mapstructure:"server_addr"`
	ServerToken string          `yaml:"server_token" mapstructure:"server_token"`
	CacheDir    string          `yaml:"cache_dir" mapstructure:"cache_dir"`
	TempDir     string          `yaml:"temp_dir" mapstructure:"temp_dir"`
	Registry    *types.Registry `yaml:"registry" mapstructure:"registry"`
}

func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
