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

package config

import (
	"fmt"
	"net"

	"github.com/spf13/viper"

	"github.com/openclarity/vmclarity/shared/pkg/families"
)

type Config struct {
	// Asset
	Asset     *Asset              `json:"asset,omitempty" yaml:"asset,omitempty" mapstructure:"asset,omitempty"`
	Addresses *families.Addresses `json:"addresses,omitempty" yaml:"addresses,omitempty" mapstructure:"addresses,omitempty"`
	Paths     *families.Paths     `json:"paths,omitempty" yaml:"paths,omitempty" mapstructure:"paths,omitempty"`
	*families.Config
}

type Asset struct {
	Type       string `json:"type" yaml:"type" mapstructure:"type"`
	Location   string `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`
	InstanceID string `json:"instanceID,omitempty" yaml:"instanceID,omitempty" mapstructure:"instanceID,omitempty"`
}

func setDefaultPaths() {
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile#L33
	viper.SetDefault("paths.gitleaksBinaryPath", "/artifacts/gitleaks")
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile#L35
	viper.SetDefault("paths.lynisInstallPath", "/artifacts/lynis")
	// https://github.com/openclarity/vmclarity-tools-base/blob/main/Dockerfile
	viper.SetDefault("paths.chkrootkitBinaryPath", "/artifacts/chkrootkit")
	viper.SetDefault("paths.clamBinaryPath", "clamscan")
	viper.SetDefault("paths.freshclamBinaryPath", "freshclam")
}

func setDefaultAddresses() {
	viper.SetDefault("addresses.exploitsDBAddress", fmt.Sprintf("http://%s", net.JoinHostPort("localhost", "1326")))
}

func GetPaths() families.Paths {
	setDefaultPaths()
	return families.Paths{
		GitleaksBinaryPath:            viper.GetString("paths.gitleaksBinaryPath"),
		ClamBinaryPath:                viper.GetString("paths.clamBinaryPath"),
		FreshclamBinaryPath:           viper.GetString("paths.freshclamBinaryPath"),
		AlternativeFreshclamMirrorURL: viper.GetString("paths.alternativeFreshclamMirrorURL"),
		LynisInstallPath:              viper.GetString("paths.lynisInstallPath"),
		ChkrootkitBinaryPath:          viper.GetString("paths.chkrootkitBinaryPath"),
	}
}

func GetAddresses() families.Addresses {
	setDefaultAddresses()
	return families.Addresses{
		ExploitsDBAddress:  viper.GetString("addresses.exploitsDBAddress"),
		GrypeServerAddress: viper.GetString("addresses.grypeServerAddress"),
		TrivyServerAddress: viper.GetString("addresses.trivyServerAddress"),
	}
}
