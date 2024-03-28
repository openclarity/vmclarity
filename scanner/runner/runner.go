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

package runner

import (
	"fmt"
	"time"

	"github.com/docker/docker/client"
	runnerclient "github.com/openclarity/vmclarity/scanner/runner/internal/runner"
)

const (
	DefaultScannerInputDir  = "/asset"
	DefaultScannerOutputDir = "/export"
	DefaultScannerConfig    = "plugin.json"
)

type PluginConfig struct {
	// Name is the name of the plugin scanner
	Name string `yaml:"name" mapstructure:"name"`
	// ImageName is the name of the docker image that will be used to run the plugin scanner
	ImageName string `yaml:"image_name" mapstructure:"image_name"`
	// InputDir is a directory where the plugin scanner will read the asset filesystem
	InputDir string `yaml:"input_dir" mapstructure:"input_dir"`
	// Output is a directory where the plugin scanner will store its results
	OutputDir string `yaml:"output_dir" mapstructure:"output_dir"`
	// ScannerConfig is a json string that will be passed to the scanner in the plugin
	ScannerConfig string `yaml:"scanner_config" mapstructure:"scanner_config"`
}

type Runner struct {
	client       runnerclient.ClientWithResponsesInterface
	dockerClient *client.Client
	containerID  string

	PluginConfig
}

func New(config PluginConfig) (*Runner, error) {
	return &Runner{
		dockerClient: dockerClient,
		PluginConfig: config,
	}, nil
}

// Start scanner container:
// * DONE get image for scanner that needs to be created
// * DONE create bind mount for filesystem directories where asset filesystem is stored and where the findings should be stored
// * TDB create bind mount for where the scanner config file is stored or copy config file to container
// * configure socket for communication with Plugin API client
// * create client for plugin container
func (r *Runner) StartScanner() error {
	return fmt.Errorf("not implemented")
}

// Wait for scanner to be ready:
// * poll the plugin container's /healthz endpoint until its healthy
func (r *Runner) WaitScannerReady(pollInterval, timeout time.Duration) error {
	return fmt.Errorf("not implemented")
}

// Post scanner configuration:
// * send scanner configuration file parsed from the AssetScan configuration received
// * send directories where the asset filesystem is stored and where the scanner findings should be saved
func (r *Runner) RunScanner() error {
	return fmt.Errorf("not implemented")
}

// Wait for scanner to be done:
// * poll plugin container's /status endpoint
func (r *Runner) WaitScannerDone() error {
	return fmt.Errorf("not implemented")
}

// Stop scanner
// * once runner receives a scanner status Done, kill scanner container
func (r *Runner) StopScanner() error {
	return fmt.Errorf("not implemented")
}
