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
	"os"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/openclarity/vmclarity/core/to"
	runnerclient "github.com/openclarity/vmclarity/plugins/runner/internal/client"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
)

const (
	DefaultScannerInputDir   = "/asset"
	DefaultScannerOutputDir  = "/export"
	DefaultScannerServerPort = "8080"
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
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// TBD API Client for the plugin scanner could be created here

	return &Runner{
		dockerClient: dockerClient,
		PluginConfig: config,
	}, nil
}

// Start scanner container:
// * get image for scanner that needs to be created
// * create bind mount for filesystem directories where asset filesystem is stored and where the findings should be stored
// * create bind mount for where the scanner config file is stored
// * open port for plugin container to listen on
// * create client for plugin container.
func (r *Runner) StartScanner() error {
	// Write scanner config file to temp dir
	err := os.WriteFile(getScannerConfigSourcePath(r.Name), []byte(r.ScannerConfig), 0o600) // nolint:gomnd
	if err != nil {
		return fmt.Errorf("failed write scanner config file: %w", err)
	}

	// Pull scanner image if required
	err = r.pullImage(context.Background(), r.ImageName)
	if err != nil {
		return fmt.Errorf("failed to pull scanner image: %w", err)
	}

	// Create scanner container
	//
	// Traefik redirects the requests to its API to our scanners so that requests to
	// localhost:TraefikContainerPort/{SCANNER_NAME} are redirected to the actual
	// {SCANNER_NAME} container. All traffic flow can be configured (e.g. auth,
	// encryption). This also enables having scanners on different hosts (although
	// not needed + requires additional configuration). The host network driver is
	// not limited by the ports anymore as we only use one (for the proxy), This way,
	// scanner containers wont overload the network driver.
	containerResp, err := r.dockerClient.ContainerCreate(
		context.Background(),
		&containertypes.Config{
			Image: r.ImageName,
			Env:   []string{"PLUGIN_SERVER_LISTEN_ADDRESS=0.0.0.0:" + DefaultScannerServerPort},
			Labels: map[string]string{
				"traefik.enable": "true",
				"traefik.http.routers." + r.Name + "-scanner.rule":                      "PathPrefix(`/" + r.Name + "/`)",
				"traefik.http.middlewares." + r.Name + "-scanner.stripprefix.prefixes":  "/" + r.Name,
				"traefik.http.routers." + r.Name + "-scanner.middlewares":               r.Name + "-scanner",
				"traefik.http.services." + r.Name + "-scanner.loadbalancer.server.port": DefaultScannerServerPort,
			},
		},
		&containertypes.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getScannerConfigSourcePath(r.Name),
					Target: getScannerConfigDestinationPath(),
				},
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

	r.client, err = runnerclient.NewClientWithResponses(
		fmt.Sprintf("http://%s/%s/", proxyHostAddress, r.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
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
			File:           to.Ptr(getScannerConfigDestinationPath()),
			InputDir:       DefaultScannerInputDir,
			OutputDir:      DefaultScannerOutputDir,
			OutputFormat:   "vmclarity-json",
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

	// Remove scanner config file
	err = os.RemoveAll(getScannerConfigSourcePath(r.Name))
	if err != nil {
		return fmt.Errorf("failed to remove scanner config file: %w", err)
	}

	return nil
}
