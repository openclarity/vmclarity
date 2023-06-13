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
	"fmt"

	"github.com/spf13/viper"
)

const (
	EnvScannerRegion       = "VMCLARITY_AWS_SCANNER_REGION"
	EnvSubnetID            = "VMCLARITY_AWS_SUBNET_ID"
	EnvSecurityGroupID     = "VMCLARITY_AWS_SECURITY_GROUP_ID"
	EnvKeyPairName         = "VMCLARITY_AWS_KEYPAIR_NAME"
	EnvScannerImage        = "VMCLARITY_AWS_SCANNER_AMI_ID"
	EnvScannerInstanceType = "VMCLARITY_AWS_SCANNER_INSTANCE_TYPE"
	EnvBlockDeviceName     = "VMCLARITY_AWS_BLOCK_DEVICE_NAME"

	DefaultScannerInstanceType = "t2.large"
	DefaultBlockDeviceName     = "xvdh"
)

type Config struct {
	// Region where the Scanner instance needs to be created
	ScannerRegion string
	// SubnetID where the Scanner instance needs to be created
	SubnetID string
	// SecurityGroupID which needs to be attached to the Scanner instance
	SecurityGroupID string
	// KeyPairName is the name of the SSH KeyPair to use for Scanner instance launch
	KeyPairName string
	// ScannerImage is the AMI image used for creating Scanner instance
	ScannerImage string
	// ScannerInstanceType is the instance type used for Scanner instance
	ScannerInstanceType string
	// BlockDeviceName contains the block device name used for attaching Scanner volume to the Scanner instance
	BlockDeviceName string
}

func (c *Config) Validate() error {
	if c.ScannerRegion == "" {
		return fmt.Errorf("parameter Region must not be nil")
	}

	if c.SubnetID == "" {
		return fmt.Errorf("parameter SubnetID must not be nil")
	}

	if c.SecurityGroupID == "" {
		return fmt.Errorf("parameter SecurityGroupID must not be nil")
	}

	if c.ScannerImage == "" {
		return fmt.Errorf("parameter ScannerImage must not be nil")
	}

	if c.ScannerInstanceType == "" {
		return fmt.Errorf("parameter ScannerInstanceType must not be nil")
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Avoid modifying the global instance
	v := viper.New()

	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	_ = v.BindEnv("ScannerRegion", EnvScannerRegion)
	_ = v.BindEnv("SubnetID", EnvSubnetID)
	_ = v.BindEnv("SecurityGroupID", EnvSecurityGroupID)
	_ = v.BindEnv("KeyPairName", EnvKeyPairName)
	_ = v.BindEnv("ScannerImage", EnvScannerImage)

	_ = v.BindEnv("ScannerInstanceType", EnvScannerInstanceType)
	v.SetDefault("ScannerInstanceType", DefaultScannerInstanceType)

	_ = v.BindEnv("BlockDeviceName", EnvBlockDeviceName)
	v.SetDefault("BlockDeviceName", DefaultBlockDeviceName)

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to parse provider configuration. Provider=AWS: %w", err)
	}

	return config, nil
}
