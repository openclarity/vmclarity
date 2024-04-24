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
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/openclarity/vmclarity/plugins/sdk/plugin"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/openclarity/vmclarity/core/to"
	runnerclient "github.com/openclarity/vmclarity/plugins/runner/internal/client"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
)

const (
	DefaultScannerInputDir = "/mnt/snapshot"

	DefaultScannerHostNetworkInterface = "127.0.0.1"
	DefaultScannerInternalServerPort   = nat.Port("8080/tcp")

	DefaultPollInterval = 2 * time.Second
	DefaultTimeout      = 60 * time.Second
)

var ErrScanNotDone = errors.New("scan has not finished yet")

type PluginRunner interface {
	Start(ctx context.Context) error
	WaitReady(ctx context.Context) error
	Run(ctx context.Context) error
	WaitDone(ctx context.Context) error

	// Logs returns log stream reader.
	// Logs() (io.ReadCloser, error)

	// Result
	// TODO: Can return the plugintypes.Result object loaded from JSON directly
	Result() (io.ReadCloser, error)

	// Remove cleans up all resources used by this PluginRunner. Should be called after New.
	// TODO: should first send HTTP request to container before removing resources.
	Remove(ctx context.Context) error
}

type PluginConfig struct {
	// Name is the name of the plugin scanner. This should be unique as only one PluginRunner with the same Name can exist.
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

func New(ctx context.Context, config PluginConfig) (PluginRunner, error) {
	// Load docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create plugin runner
	runner := &runner{
		config:       config,
		dockerClient: dockerClient,
	}
	if err = runner.create(ctx); err != nil {
		defer runner.Remove(ctx) //nolint:errcheck

		return nil, fmt.Errorf("failed to create plugin runner: %w", err)
	}

	return runner, nil
}

// create creates the plugin container (without starting it) and sets
// runner.containerID on success.
func (r *runner) create(ctx context.Context) error {
	// Pull scanner image if required
	err := pullImage(ctx, r.dockerClient, r.config.ImageName)
	if err != nil {
		return fmt.Errorf("failed to pull scanner image: %w", err)
	}

	// Get scanner container mounts
	mounts, err := r.getPluginContainerMounts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get scanner container mounts: %w", err)
	}

	// Get networking config
	networkingConfig, err := r.getNetworkingConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network ID: %w", err)
	}

	// Create scanner container
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
			Mounts: mounts,
		},
		networkingConfig,
		nil,
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to create scanner container: %w", err)
	}
	r.containerID = container.ID

	// Copy config file to container
	err = r.copyConfigToContainer(ctx)
	if err != nil {
		return fmt.Errorf("failed to copy scanner config to container: %w", err)
	}

	return nil
}

func (r *runner) Start(ctx context.Context) error {
	// Start container
	if err := r.dockerClient.ContainerStart(ctx, r.containerID, containertypes.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start scanner container: %w", err)
	}

	// Wait for the container to enter the running state
	if err := r.waitContainerReady(ctx); err != nil {
		return fmt.Errorf("failed to wait for scanner to start: %w", err)
	}

	// Load plugin client
	if err := r.loadPluginClient(ctx); err != nil {
		return fmt.Errorf("failed to create scanner container client: %w", err)
	}

	return nil
}

func (r *runner) WaitReady(ctx context.Context) error {
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
			OutputFile:     r.config.OutputFile,
			TimeoutSeconds: int(DefaultTimeout.Seconds()), // TODO: this should be configurable
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
	// Copy result file from container
	reader, _, err := r.dockerClient.CopyFromContainer(context.Background(), r.containerID, r.config.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to copy scanner result file: %w", err)
	}

	// Extract the tar file and read the content
	tr := tar.NewReader(reader)
	_, err = tr.Next()
	if err == io.EOF {
		return nil, ErrScanNotDone
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read tar file: %w", err)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, tr); err != nil {
		return nil, fmt.Errorf("failed to copy file contents: %w", err)
	}

	return io.NopCloser(buf), nil
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

// loadPluginClient loads http client into runner.client to interact with plugin
// server by trying to connect with the container either via internal container
// IP address or via host using exposed port. This method handles retry until
// DefaultTimeout is reached.
func (r *runner) loadPluginClient(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Get network data needed to reach the container
	inspect, err := r.dockerClient.ContainerInspect(ctx, r.containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect scanner container: %w", err)
	}
	hostPorts, ok := inspect.NetworkSettings.Ports[DefaultScannerInternalServerPort]
	if !ok {
		return fmt.Errorf("failed to get scanner ports: %w", err)
	}
	if len(hostPorts) != 1 {
		return errors.New("network port not attached to scanner container")
	}
	containerHostPort := hostPorts[0].HostPort
	containerIP := inspect.Config.Hostname

	// Try proper client for interacting with plugin server
	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	clientErrs := make([]error, 2) //nolint:gomnd
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking http clients of %s timed out: %w", r.config.Name, errors.Join(clientErrs...))

		case <-ticker.C:
			// Check if possible to connect directly via exposed port
			r.client, clientErrs[0] = newPluginClient(ctx,
				"http://localhost:"+containerHostPort,
			)
			if clientErrs[0] == nil {
				return nil
			}

			// Check if possible to connect via container IP and internal port
			r.client, clientErrs[1] = newPluginClient(ctx,
				"http://"+net.JoinHostPort(containerIP, DefaultScannerInternalServerPort.Port()),
			)
			if clientErrs[1] == nil {
				return nil
			}
		}
	}
}

func getScannerConfigSourcePath(name string) string {
	return filepath.Join(os.TempDir(), name+"-plugin.json")
}

func getScannerConfigDestinationPath() string {
	return filepath.Join("/plugin.json")
}

// newPluginClient creates a new client to interact with plugin server and pings
// the server to check if available or errors.
func newPluginClient(ctx context.Context, server string) (*runnerclient.ClientWithResponses, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second) //nolint:gomnd
	defer cancel()

	c, err := runnerclient.NewClientWithResponses(server)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin client: %w", err)
	}

	_, err = c.GetStatusWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping plugin scanner server: %w", err)
	}

	return c, nil
}

// waitContainerReady waits for the container to enter Running state to ensure
// that network binds are available and server is ready to receive requests.
func (r *runner) waitContainerReady(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	ticker := time.NewTicker(DefaultPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("checking ready of %s timed out", r.config.Name)

		case <-ticker.C:
			// Get state data needed to check the container
			inspect, err := r.dockerClient.ContainerInspect(ctx, r.containerID)
			if err != nil {
				return fmt.Errorf("failed to inspect scanner container: %w", err)
			}

			if inspect.State.Running {
				return nil
			}
		}
	}
}
