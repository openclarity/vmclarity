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

package aws

import (
	"context"
	"errors"

	envtypes "github.com/openclarity/vmclarity/testenv/types"
)

const DefaultAWSRegion = "eu-central-1"

// Config defines configuration for AWS environment.
//
//nolint:containedctx
type Config struct {
	// WorkDir absolute path to the directory where the deployment files prior performing actions
	WorkDir string `mapstructure:"work_dir"`
	// EnvName the name of the stack to be created
	EnvName string `mapstructure:"env_name"`
	// Region the AWS region to be used
	Region string `mapstructure:"region"`
	// PublicKey the public key to be used for the key pair
	PublicKey string `mapstructure:"public_key"`

	// ctx used during project initialization
	ctx context.Context
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Region == "" {
		return errors.New("parameter region must be provided")
	}

	if c.PublicKey == "" {
		return errors.New("parameter public_key must be provided")
	}

	return nil
}

// ConfigOptFn defines transformer function for Config.
type ConfigOptFn func(*Config) error

var applyConfigWithOpts = envtypes.WithOpts[Config, ConfigOptFn]

func WithContext(ctx context.Context) ConfigOptFn {
	return func(config *Config) error {
		config.ctx = ctx

		return nil
	}
}

// WithWorkDir set workDir for Config.
func WithWorkDir(dir string) ConfigOptFn {
	return func(config *Config) error {
		config.WorkDir = dir

		return nil
	}
}
