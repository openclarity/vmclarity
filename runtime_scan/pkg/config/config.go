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
	"github.com/spf13/viper"
)

const (
	ScannerJobResultListenPort        = "SCANNER_JOB_RESULT_LISTEN_PORT"
	ScannerRegion                     = "SCANNER_REGION"
	defaultScannerRegion              = "us-east-1"
	ScannerJobImageID                 = "SCANNER_JOB_IMAGE_ID"
	defaultScannerJobImageID          = "ami-0568773882d492fc8" // ubuntu server 22.04 LTS (HVM), SSD volume type
	ScannerAttachedVolumeDeviceName   = "SCANNER_ATTACHED_VOLUME_DEVICE_NAME"
	defaultAttachedVolumeDeviceName   = "xvdh"
	defaultScannerJobResultListenPort = 8888
)

type Config struct {
	ScannerJobResultListenPort int
	Region                     string // scanner region
	AmiID                      string // image id of a scanner job
	DeviceName                 string // the name of the block device to attach to the scanner instance (mounted snapshot)
}

func setConfigDefaults() {
	// TODO defaults for region and ami ID
	viper.SetDefault(ScannerJobResultListenPort, defaultScannerJobResultListenPort)
	viper.SetDefault(ScannerRegion, defaultScannerRegion)
	viper.SetDefault(ScannerJobImageID, defaultScannerJobImageID)
	viper.SetDefault(ScannerAttachedVolumeDeviceName, defaultAttachedVolumeDeviceName)

	viper.AutomaticEnv()
}

func LoadConfig() (*Config, error) {
	setConfigDefaults()

	config := &Config{
		ScannerJobResultListenPort: viper.GetInt(ScannerJobResultListenPort),
		Region:                     viper.GetString(ScannerRegion),
		AmiID:                      viper.GetString(ScannerJobImageID),
		DeviceName:                 viper.GetString(ScannerAttachedVolumeDeviceName),
	}

	return config, nil
}
