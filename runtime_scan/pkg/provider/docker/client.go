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

package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Client struct {
	dockerClient *client.Client
	config       *Config
}

func New(_ context.Context) (*Client, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration. Provider=%s: %w", models.Docker, err)
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate provider configuration. Provider=%s: %w", models.Docker, err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to load provider configuration. Provider=%s: %w", models.Docker, err)
	}

	return &Client{
		dockerClient: dockerClient,
		config:       config,
	}, nil
}

func (c *Client) Kind() models.CloudProvider {
	return models.Docker
}

func (c *Client) DiscoverAssets(ctx context.Context) ([]models.AssetType, error) {
	// Get image assets
	imageAssets, err := c.getImages(ctx)
	if err != nil {
		return nil, provider.FatalErrorf("failed to get images. Provider=%s: %w", models.Docker, err)
	}

	// Get container assets
	containerAssets, err := c.getContainers(ctx)
	if err != nil {
		return nil, provider.FatalErrorf("failed to get containers. Provider=%s: %w", models.Docker, err)
	}

	// Combine assets
	assets := append(imageAssets, containerAssets...)

	return assets, nil
}

func (c *Client) RunAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {

	err := c.prepareScanVolume(ctx, config)
	if err != nil {
		return provider.FatalErrorf("failed to prepare scan volume. Provider=%s: %w", models.Docker, err)
	}

	err = c.createScanConfigFile(config)
	if err != nil {
		return provider.FatalErrorf("failed to create scanconfig.yaml file. Provider=%s: %w", models.Docker, err)
	}

	containerId, err := c.createScanContainer(ctx, config)
	if err != nil {
		return provider.FatalErrorf("failed to create scan container. Provider=%s: %w", models.Docker, err)
	}

	err = c.dockerClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		return provider.FatalErrorf("failed to start scan container. Provider=%s: %w", models.Docker, err)
	}

	return nil
}

func (c *Client) prepareScanVolume(ctx context.Context, config *provider.ScanJobConfig) error {
	scanName, err := getScanName(config)
	if err != nil {
		return err
	}

	// Create volume if not found
	resp, err := c.dockerClient.VolumeList(ctx, volume.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", scanName)),
	})
	if err != nil {
		return fmt.Errorf("failed to get volumes: %w", err)
	}
	if len(resp.Volumes) == 1 {
		fmt.Printf("scan volume already created")
		return nil
	}
	if len(resp.Volumes) == 0 {
		_, err = c.dockerClient.VolumeCreate(ctx, volume.CreateOptions{
			Name: scanName,
		})
		if err != nil {
			return fmt.Errorf("failed to create volume: %w", err)
		}
	}

	readCloser, err := c.export(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to export target: %w", err)
	}

	// Create an ephemeral container to populate volume with export output
	_, err = c.dockerClient.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull helper image: %w", err)
	}
	response, err := c.dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: "alpine",
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: scanName,
					Target: "/data",
				},
			},
		}, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create helper container: %w", err)
	}
	defer func() {
		err = c.dockerClient.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{})
		if err != nil {
			_ = fmt.Errorf("failed to remove helper container: %w", err)
		}
	}()
	err = c.dockerClient.CopyToContainer(ctx, response.ID, "/data", readCloser, types.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("failed to copy data to container: %w", err)
	}
	return nil
}

func (c *Client) export(ctx context.Context, config *provider.ScanJobConfig) (io.ReadCloser, error) {
	objectType, err := config.AssetInfo.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("failed to get asset object type: %w", err)
	}

	switch objectType {
	case "ContainerInfo":
		id, err := getAssetId(config)
		if err != nil {
			return nil, err
		}
		return c.dockerClient.ContainerExport(ctx, id)
	case "ContainerImageInfo":
		name, err := getAssetName(config)
		if err != nil {
			return nil, err
		}
		// Create an ephemeral container to export asset
		response, err := c.dockerClient.ContainerCreate(ctx,
			&container.Config{Image: name},
			nil,
			nil,
			nil,
			"",
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create helper container: %w", err)
		}
		defer func() {
			err = c.dockerClient.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{})
			if err != nil {
				_ = fmt.Errorf("failed to remove helper container: %w", err)
			}
		}()
		return c.dockerClient.ContainerExport(ctx, response.ID)
	default:
		return nil, fmt.Errorf("get raw contents not implemented for current object type (%s)", objectType)
	}
}

func (c *Client) createScanConfigFile(config *provider.ScanJobConfig) error {
	scanConfigFilePath, err := getScanConfigFileName(config)
	if err != nil {
		return err
	}

	_, err = os.Stat(scanConfigFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(scanConfigFilePath, []byte(config.ScannerCLIConfig), 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (c *Client) createScanContainer(ctx context.Context, config *provider.ScanJobConfig) (string, error) {
	scanName, err := getScanName(config)
	if err != nil {
		return "", err
	}

	containerId, err := c.getContainerIdFromContainerName(ctx, scanName)
	if containerId != "" {
		return containerId, nil
	}

	pl, err := c.dockerClient.ImagePull(ctx, config.ScannerImage, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	_, _ = io.Copy(io.Discard, pl)
	_ = pl.Close()

	scanConfigFilePath, err := getScanConfigFileName(config)
	if err != nil {
		return "", err
	}

	networkResp, err := c.dockerClient.NetworkCreate(
		ctx,
		scanName,
		types.NetworkCreate{
			CheckDuplicate: true,
			Driver:         "bridge",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to scan network: %w", err)
	}

	containerResp, err := c.dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: config.ScannerImage,
			Entrypoint: []string{"sh", "-c",
				fmt.Sprintf(
					"/app/vmclarity-cli --config /tmp/%s --server %s --asset-scan-id %s",
					filepath.Base(scanConfigFilePath),
					config.VMClarityAddress,
					config.AssetScanID,
				),
			},
		},
		&container.HostConfig{
			Binds: []string{fmt.Sprintf("%s:/tmp", path.Dir(scanConfigFilePath))},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: scanName,
					Target: "/mnt/snapshot",
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				scanName: {
					NetworkID: networkResp.ID,
				},
			},
		},
		nil,
		scanName,
	)
	if err != nil {
		return "", err
	}

	return containerResp.ID, nil
}

func (c *Client) getContainerIdFromContainerName(ctx context.Context, scanName string) (string, error) {

	containers, err := c.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", scanName)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("scan container not found: %w", err)
	}
	if len(containers) > 1 {
		return "", fmt.Errorf("found more than one scan container: %w", err)
	}
	return containers[0].ID, nil
}

func (c *Client) RemoveAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	scanName, err := getScanName(config)
	if err != nil {
		return provider.FatalErrorf("failed to get scan name. Provider=%s: %w", models.Docker, err)
	}

	err = c.dockerClient.VolumeRemove(ctx, scanName, false)
	if err != nil {
		return provider.FatalErrorf("failed to remove volume. Provider=%s: %w", models.Docker, err)
	}

	scanConfigFileName, err := getScanConfigFileName(config)
	if err != nil {
		return provider.FatalErrorf("failed to get scan config file name. Provider=%s: %w", models.Docker, err)
	}
	err = os.Remove(path.Dir(scanConfigFileName))
	if err != nil {
		return provider.FatalErrorf("failed to remove scan config file. Provider=%s: %w", models.Docker, err)
	}

	containerId, err := c.getContainerIdFromContainerName(ctx, scanName)
	if err != nil {
		return provider.FatalErrorf("failed to get scan container id. Provider=%s: %w", models.Docker, err)
	}
	err = c.dockerClient.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return provider.FatalErrorf("failed to remove scan container. Provider=%s: %w", models.Docker, err)
	}

	networkId, err := c.getNetworkIdFromNetworkName(ctx, scanName)
	if err != nil {
		return provider.FatalErrorf("failed to get scan network id. Provider=%s: %w", models.Docker, err)
	}
	err = c.dockerClient.NetworkRemove(ctx, networkId)
	if err != nil {
		return provider.FatalErrorf("failed to remove scan network. Provider=%s: %w", models.Docker, err)
	}

	return nil
}

func (c *Client) getNetworkIdFromNetworkName(ctx context.Context, scanName string) (string, error) {

	networks, err := c.dockerClient.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.Arg("name", scanName)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list networks: %w", err)
	}
	if len(networks) == 0 {
		return "", fmt.Errorf("scan network not found: %w", err)
	}
	if len(networks) > 1 {
		return "", fmt.Errorf("found more than one scan network: %w", err)
	}
	return networks[0].ID, nil
}
