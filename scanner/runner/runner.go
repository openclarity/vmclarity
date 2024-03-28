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
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	imagetypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	runnerclient "github.com/openclarity/vmclarity/scanner/runner/internal/runner"
	"github.com/openclarity/vmclarity/scanner/types"
)

const (
	DefaultScannerInputDir   = "/asset"
	DefaultScannerOutputDir  = "/export"
	DefaultScannerSocketFile = "/var/run/plugin.sock"
	DefaultScannerConfig     = "plugin.json"
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
	socketFile   string

	PluginConfig
}

func New(config PluginConfig) (*Runner, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// TBD API Client for the plugin scanner could be created here

	return &Runner{
		dockerClient: dockerClient,
		PluginConfig: config,
		socketFile:   "/tmp/socket/plugin.sock", // TODO(paralta): make this configurable with plugin name
	}, nil
}

// Start scanner container:
// * get image for scanner that needs to be created
// * create bind mount for filesystem directories where asset filesystem is stored and where the findings should be stored
// * TDB create bind mount for where the scanner config file is stored or copy config file to container
// * configure unix domain socket for communication with Plugin API client
// * create client for plugin container.
func (r *Runner) StartScanner() error {
	err := os.Mkdir(filepath.Dir(r.socketFile), 0o777) //nolint:gomnd
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create socket dir: %w", err)
	}

	// Pull scanner image if required
	images, err := r.dockerClient.ImageList(context.Background(), imagetypes.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", r.ImageName)),
	})
	if err != nil {
		return fmt.Errorf("failed to get images: %w", err)
	}
	if len(images) == 0 {
		_, err = r.dockerClient.ImagePull(context.Background(), r.ImageName, imagetypes.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull scanner image: %w", err)
		}
	}

	containerResp, err := r.dockerClient.ContainerCreate(
		context.Background(),
		&containertypes.Config{
			Image: r.ImageName,
			Cmd:   []string{"sleep", "infinity"},
		},
		&containertypes.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: r.InputDir,
					Target: DefaultScannerInputDir,
				},
				{
					Type:   mount.TypeBind,
					Source: r.OutputDir,
					Target: DefaultScannerOutputDir,
				},
				{
					Type:   mount.TypeBind,
					Source: filepath.Dir(r.socketFile),
					Target: filepath.Dir(DefaultScannerSocketFile),
				},
			},
		},
		nil,
		nil,
		r.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to create scanner container: %w", err)
	}
	r.containerID = containerResp.ID

	err = r.dockerClient.ContainerStart(context.Background(), r.containerID, containertypes.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start scanner container: %w", err)
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", r.socketFile)
			},
		},
	}
	r.client, err = runnerclient.NewClientWithResponses(
		"http://unix/",
		runnerclient.WithHTTPClient(&httpClient),
	)
	if err != nil {
		return fmt.Errorf("failed to create plugin client: %w", err)
	}

	return nil
}

// Wait for scanner to be ready:
// * poll the plugin container's /healthz endpoint until its healthy
//
//nolint:gomnd
func (r *Runner) WaitScannerReady(pollInterval, timeout time.Duration) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking health of %s timed out", r.PluginConfig.Name)

		case <-ticker.C:
			resp, err := r.client.GetHealthzWithResponse(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get scanner health: %w", err)
			}

			if resp.StatusCode() == 200 {
				return nil
			}
		}
	}
}

// Post scanner configuration:
// * send scanner configuration file parsed from the AssetScan configuration received
// * send directories where the asset filesystem is stored and where the scanner findings should be saved.
func (r *Runner) RunScanner() error {
	_, err := r.client.PostConfigWithResponse(
		context.Background(),
		types.PostConfigJSONRequestBody{
			File:           r.ScannerConfig,
			InputDir:       DefaultScannerInputDir,
			OutputDir:      DefaultScannerOutputDir,
			TimeoutSeconds: 60, //nolint:gomnd
		},
	)
	if err != nil {
		return fmt.Errorf("failed to post scanner config: %w", err)
	}

	return nil
}

// Wait for scanner to be done:
// * poll plugin container's /status endpoint.
func (r *Runner) WaitScannerDone(pollInterval, timeout time.Duration) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking status of %s timed out", r.PluginConfig.Name)

		case <-ticker.C:
			resp, err := r.client.GetStatusWithResponse(context.Background())
			if err != nil {
				return fmt.Errorf("failed to get scanner status: %w", err)
			}

			if resp.JSON200.State == types.Done {
				return nil
			}
		}
	}
}

// Stop and remove scanner
// * once runner receives a scanner status Done, kill scanner container.
func (r *Runner) StopScanner() error {
	err := r.dockerClient.ContainerStop(context.Background(), r.containerID, containertypes.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop scanner container: %w", err)
	}

	err = r.dockerClient.ContainerRemove(context.Background(), r.containerID, containertypes.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove scanner container: %w", err)
	}

	// Remove socket dir
	err = os.RemoveAll(filepath.Dir(r.socketFile))
	if err != nil {
		return fmt.Errorf("failed to remove socket file: %w", err)
	}

	return nil
}
