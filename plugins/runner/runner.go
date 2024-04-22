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
	"errors"
	"fmt"
	"github.com/openclarity/vmclarity/plugins/sdk/plugin"
	"io"
	"os"
	"path/filepath"
	"time"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/openclarity/vmclarity/core/to"
	runnerclient "github.com/openclarity/vmclarity/plugins/runner/internal/client"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
)

const (
	DefaultScannerInputDir  = "/asset"
	DefaultScannerOutputDir = "/export"

	DefaultScannerHostNetworkInterface = "127.0.0.1"
	DefaultScannerInternalServerPort   = nat.Port("8080/tcp")

	DefaultPollInterval = 2 * time.Second
	DefaultTimeout      = 60 * time.Second
)

var ErrScanNotDone = errors.New("scan has not finished yet")

type Runner interface {
	Start(ctx context.Context) error
	WaitReady(ctx context.Context) error
	Run(ctx context.Context) error
	WaitDone(ctx context.Context) error
	// Result
	// TODO: Can return the plugintypes.Result object loaded from JSON directly
	Result() (io.ReadCloser, error)
	// Remove removes all resources used by Runner. Should be called after New.
	Remove(ctx context.Context) error
	// TODO: implement stop to send the stop HTTP request to container
	// Stop(ctx context.Context) error
}

type PluginConfig struct {
	// Name is the name of the plugin scanner. This should be unique as only one Runner with the same Name can exist.
	Name string `yaml:"name" mapstructure:"name"`
	// ImageName is the name of the docker image that will be used to run the plugin scanner
	ImageName string `yaml:"image_name" mapstructure:"image_name"`
	// InputDir is a directory where the plugin scanner will read the asset filesystem
	InputDir string `yaml:"input_dir" mapstructure:"input_dir"`
	// OutputFile is a file where the plugin scanner will write the result
	OutputFile string `yaml:"output_file" mapstructure:"output_file"`
	// ScannerConfig is a json string that will be passed to the scanner in the plugin
	ScannerConfig string `yaml:"scanner_config" mapstructure:"scanner_config"`
}

type runner struct {
	config       PluginConfig
	client       runnerclient.ClientWithResponsesInterface
	dockerClient *client.Client
	containerID  string
}

func New(ctx context.Context, config PluginConfig) (Runner, error) {
	// Load docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create runner
	runner := &runner{
		config:       config,
		dockerClient: dockerClient,
	}

	if err = runner.create(ctx); err != nil {
		// Remove resources if needed
		var removeErr error
		if runner.containerID != "" {
			removeErr = runner.Remove(ctx)
		}

		// Return collected errors
		errs := errors.Join(err, removeErr)
		return nil, fmt.Errorf("failed to create runner: %w", errs)
	}

	return runner, nil
}

func (r *runner) create(ctx context.Context) error {
	// Write scanner config file to temp dir
	err := os.WriteFile(getScannerConfigSourcePath(r.config.Name), []byte(r.config.ScannerConfig), 0o600) // nolint:gomnd
	if err != nil {
		return fmt.Errorf("failed write scanner config file: %w", err)
	}

	// Pull scanner image if required
	err = pullImage(ctx, r.dockerClient, r.config.ImageName)
	if err != nil {
		return fmt.Errorf("failed to pull scanner image: %w", err)
	}

	// Create scanner container
	{
		container, err := r.dockerClient.ContainerCreate(
			ctx,
			&containertypes.Config{
				Image: r.config.ImageName,
				Env: []string{
					fmt.Sprintf("%s=0.0.0.0:%s", plugin.EnvListenAddress, DefaultScannerInternalServerPort.Port()),
				},
				ExposedPorts: nat.PortSet{DefaultScannerInternalServerPort: struct{}{}},
			},
			&containertypes.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{
					DefaultScannerInternalServerPort: {
						{
							HostIP:   DefaultScannerHostNetworkInterface,
							HostPort: "",
						},
					},
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: getScannerConfigSourcePath(r.config.Name),
						Target: getScannerConfigDestinationPath(),
					},
					{
						Type:   mount.TypeBind,
						Source: r.config.InputDir,
						Target: DefaultScannerInputDir,
					},
					{
						Type:   mount.TypeBind,
						Source: filepath.Dir(r.config.OutputFile),
						Target: DefaultScannerOutputDir,
					},
				},
			},
			nil,
			nil,
			r.config.Name,
		)
		if err != nil {
			return fmt.Errorf("failed to create scanner container: %w", err)
		}

		r.containerID = container.ID
	}

	// Get host port where the container is reachable from
	var hostPort string
	{
		inspect, err := r.dockerClient.ContainerInspect(ctx, r.containerID)
		if err != nil {
			return fmt.Errorf("failed to inspect scanner container: %w", err)
		}
		hostPorts, ok := inspect.NetworkSettings.Ports[DefaultScannerInternalServerPort]
		if !ok {
			return fmt.Errorf("failed to get scanner ports: %w", err)
		}
		if len(hostPorts) != 1 {
			return fmt.Errorf("network port not attached to scanner")
		}
		hostPort = hostPorts[0].HostPort
	}

	// Create client to interact with the plugin
	r.client, err = runnerclient.NewClientWithResponses(
		fmt.Sprintf("http://%s:%s", DefaultScannerHostNetworkInterface, hostPort),
	)
	if err != nil {
		return fmt.Errorf("failed to create plugin client: %w", err)
	}

	return nil
}

func (r *runner) Start(ctx context.Context) error {
	err := r.dockerClient.ContainerStart(ctx, r.containerID, containertypes.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start scanner container: %w", err)
	}

	return nil
}

func (r *runner) WaitReady(ctx context.Context) error {
	// TODO: give some time for the docker container to boot.
	// TODO: Can be done by adding retry logic (3 retries to WaitReady, 10s delay between retries)

	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking health of %s timed out", r.config.Name)

		case <-ticker.C:
			resp, err := r.client.GetHealthzWithResponse(ctx)
			if err != nil {
				return fmt.Errorf("failed to get scanner health: %w", err)
			}

			if resp.StatusCode() == 200 { //nolint:gomnd
				return nil
			}
		}
	}
}

func (r *runner) Run(ctx context.Context) error {
	_, err := r.client.PostConfigWithResponse(
		ctx,
		types.PostConfigJSONRequestBody{
			ScannerConfig:  to.Ptr(r.config.ScannerConfig),
			InputDir:       DefaultScannerInputDir,
			OutputFile:     filepath.Join(DefaultScannerOutputDir, filepath.Base(r.config.OutputFile)),
			TimeoutSeconds: int(DefaultTimeout),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to post scanner config: %w", err)
	}

	return nil
}

func (r *runner) WaitDone(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking status of %s timed out", r.config.Name)

		case <-ticker.C:
			resp, err := r.client.GetStatusWithResponse(ctx)
			if err != nil {
				return fmt.Errorf("failed to get scanner status: %w", err)
			}

			if resp.JSON200.State == types.Done {
				return nil
			}
		}
	}
}

func (r *runner) Result() (io.ReadCloser, error) {
	_, err := os.Stat(r.config.OutputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrScanNotDone
		}
		return nil, fmt.Errorf("failed to fetch scanner result file: %w", err)
	}

	file, err := os.Open(r.config.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open scanner result file: %w", err)
	}

	return file, nil
}

func (r *runner) Remove(ctx context.Context) error {
	err := r.dockerClient.ContainerStop(ctx, r.containerID, containertypes.StopOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("failed to stop scanner container: %w", err)
	}

	err = r.dockerClient.ContainerRemove(ctx, r.containerID, containertypes.RemoveOptions{})
	if err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("failed to remove scanner container: %w", err)
	}

	// Remove scanner config file
	err = os.RemoveAll(getScannerConfigSourcePath(r.config.Name))
	if err != nil {
		return fmt.Errorf("failed to remove scanner config file: %w", err)
	}

	return nil
}

func getScannerConfigSourcePath(name string) string {
	return filepath.Join(os.TempDir(), name+"-plugin.json")
}

func getScannerConfigDestinationPath() string {
	return filepath.Join("/plugin.json")
}
