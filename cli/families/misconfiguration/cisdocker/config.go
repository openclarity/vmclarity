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

package cisdocker

import (
	"time"

	dockle_config "github.com/Portshift/dockle/config"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/cli/families/misconfiguration/types"
)

const (
	DefaultCISDockerTimeout = 2 * time.Minute
)

func createDockleConfig(logger *logrus.Entry, imageName string, config types.CISDockerConfig) *dockle_config.Config {
	var username, password string
	var insecure, nonSSL bool
	if config.Registry != nil {
		insecure = config.Registry.SkipVerifyTLS
		nonSSL = config.Registry.UseHTTP
		if len(config.Registry.Auths) > 0 {
			username = config.Registry.Auths[0].Username
			password = config.Registry.Auths[0].Password
		}
	}

	if config.Timeout == 0 {
		config.Timeout = DefaultCISDockerTimeout
	}

	return &dockle_config.Config{
		Debug:     logger.Logger.Level == logrus.DebugLevel,
		Timeout:   config.Timeout,
		Username:  username,
		Password:  password,
		Insecure:  insecure,
		NonSSL:    nonSSL,
		ImageName: imageName,
	}
}
