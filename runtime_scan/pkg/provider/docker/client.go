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
	types "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Client struct {
	dockerClient *client.Client
	config       *Config
}

func New(ctx context.Context) (*Client, error) {
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
	var err error
	wg := &sync.WaitGroup{}
	errs := make(chan error, 2)
	assets := make(chan models.AssetType)

	wg.Add(1)
	go c.getImages(ctx, assets, errs, wg)

	wg.Add(1)
	go c.getContainers(ctx, assets, errs, wg)

	go func() {
		wg.Wait()
		close(errs)
		close(assets)
	}()

	var ret []models.AssetType
	for t := range assets {
		ret = append(ret, t)
	}

	for e := range errs {
		if e != nil {
			// nolint:typecheck
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Client) getImages(ctx context.Context, assets chan models.AssetType, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	images, err := c.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		errs <- provider.FatalError{
			Err: fmt.Errorf("failed to get images. Provider=%s: %w", models.Docker, err),
		}
		return
	}

	for _, i := range images {
		asset, err := c.getAssetFromImage(ctx, i.ID)
		if err != nil {
			errs <- provider.FatalError{
				Err: fmt.Errorf("failed to create AssetType from ContainerImageInfo. Provider=%s: %w", models.Docker, err),
			}
			return
		}
		assets <- asset
	}
}

func (c *Client) getAssetFromImage(ctx context.Context, id string) (models.AssetType, error) {
	asset := models.AssetType{}

	info, err := c.getContainerImageInfoFromImage(ctx, id)
	if err != nil {
		return asset, err
	}

	err = asset.FromContainerImageInfo(info)
	if err != nil {
		return asset, fmt.Errorf("failed to create AssetType from ContainerImageInfo. Provider=%s: %w", models.Docker, err)
	}

	return asset, err
}

func (c *Client) getContainers(ctx context.Context, assets chan models.AssetType, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	containers, err := c.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		errs <- provider.FatalError{
			Err: fmt.Errorf("failed to get containers. Provider=%s: %w", models.Docker, err),
		}
		return
	}

	for _, container := range containers {
		inspect, err := c.dockerClient.ContainerInspect(ctx, container.ID)
		if err != nil {
			errs <- provider.FatalError{
				Err: fmt.Errorf("failed to get container. Provider=%s: %w", models.Docker, err),
			}
			return
		}

		created, err := time.Parse(time.RFC3339, inspect.Created)
		if err != nil {
			errs <- provider.FatalError{
				Err: fmt.Errorf("failed to parse time. Provider=%s: %w", models.Docker, err),
			}
			return
		}

		imageInfo, err := c.getContainerImageInfoFromImage(ctx, container.ImageID)
		if err != nil {
			// TODO (paralta) If image not required this should not be fatal
			errs <- provider.FatalError{
				Err: err,
			}
			return
		}

		asset := models.AssetType{}
		err = asset.FromContainerInfo(models.ContainerInfo{
			ContainerName: utils.PointerTo(inspect.Name),
			CreatedAt:     utils.PointerTo(created),
			Id:            utils.PointerTo(container.ID),
			Image:         utils.PointerTo(imageInfo),
			Labels:        convertTags(container.Labels),
			Location:      nil, // TODO (paralta) Clarify what is location
			ObjectType:    "ContainerInfo",
		})
		if err != nil {
			errs <- provider.FatalError{
				Err: fmt.Errorf("failed to create AssetType from ContainerInfo. Provider=%s: %w", models.Docker, err),
			}
			return
		}
		assets <- asset
	}
}

func (c *Client) getContainerImageInfoFromImage(ctx context.Context, id string) (models.ContainerImageInfo, error) {
	i, _, err := c.dockerClient.ImageInspectWithRaw(ctx, id)
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to get image. Provider=%s: %w", models.Docker, err)
	}

	return models.ContainerImageInfo{
		Architecture: utils.PointerTo(i.Architecture),
		Id:           utils.PointerTo(i.ID),
		Labels:       convertTags(i.Config.Labels),
		Name:         utils.PointerTo(i.Config.Image),
		ObjectType:   "ContainerImageInfo",
		Os:           utils.PointerTo(i.Os),
		Size:         utils.PointerTo(int(i.Size)),
	}, err
}

func (c *Client) RunAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {

	err := c.prepareScanVolume(ctx, config)
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to prepare scan volume. Provider=%s: %w", models.Docker, err),
		}
	}

	err = c.createScanConfigFile(config)
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to create scanconfig.yaml file. Provider=%s: %w", models.Docker, err),
		}
	}

	containerId, err := c.createScanContainer(ctx, config)
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to create scan container. Provider=%s: %w", models.Docker, err),
		}
	}

	err = c.dockerClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to start scan container. Provider=%s: %w", models.Docker, err),
		}
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
		return fmt.Errorf("failed to get volumes. Provider=%s: %w", models.Docker, err)
	}
	if len(resp.Volumes) == 1 {
		fmt.Printf("scan volume already created. Provider=%s", models.Docker)
		return nil
	}
	if len(resp.Volumes) == 0 {
		_, err = c.dockerClient.VolumeCreate(ctx, volume.CreateOptions{
			Name: scanName,
		})
		if err != nil {
			return fmt.Errorf("failed to create volume. Provider=%s: %w", models.Docker, err)
		}
	}

	readCloser, err := c.export(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to export target. Provider=%s: %w", models.Docker, err)
	}

	// Create an ephemeral container to populate volume with export output
	_, err = c.dockerClient.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull helper image. Provider=%s: %w", models.Docker, err)
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
		return fmt.Errorf("failed to create helper container. Provider=%s: %w", models.Docker, err)
	}
	defer func() {
		err = c.dockerClient.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{})
		if err != nil {
			_ = fmt.Errorf("failed to remove helper container. Provider=%s: %w", models.Docker, err)
		}
	}()
	err = c.dockerClient.CopyToContainer(ctx, response.ID, "/data", readCloser, types.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("failed to copy data to container. Provider=%s: %w", models.Docker, err)
	}
	return nil
}

func (c *Client) export(ctx context.Context, config *provider.ScanJobConfig) (io.ReadCloser, error) {
	objectType, err := config.AssetInfo.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("failed to get asset object type. Provider=%s: %w", models.Docker, err)
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
			&container.Config{
				Image: name,
			}, nil, nil, nil, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create helper container. Provider=%s: %w", models.Docker, err)
		}
		defer func() {
			err = c.dockerClient.ContainerRemove(ctx, response.ID, types.ContainerRemoveOptions{})
			if err != nil {
				_ = fmt.Errorf("failed to remove helper container. Provider=%s: %w", models.Docker, err)
			}
		}()
		return c.dockerClient.ContainerExport(ctx, response.ID)
	default:
		return nil, fmt.Errorf("get raw contents not implemented for current object type (%s). Provider=%s", models.Docker, objectType)
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

	resp, err := c.dockerClient.ContainerCreate(
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
		nil,
		nil,
		scanName,
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Client) getContainerIdFromContainerName(ctx context.Context, scanName string) (string, error) {

	containers, err := c.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", scanName)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get containers. Provider=%s: %w", models.Docker, err)
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("scan container not found. Provider=%s: %w", models.Docker, err)
	}
	if len(containers) > 1 {
		return "", fmt.Errorf("found more than one scan container. Provider=%s: %w", models.Docker, err)
	}
	return containers[0].ID, nil
}

func (c *Client) RemoveAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	scanName, err := getScanName(config)
	if err != nil {
		return provider.FatalError{Err: err}
	}
	err = c.dockerClient.VolumeRemove(ctx, scanName, false)
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to remove volume. Provider=%s: %w", models.Docker, err),
		}
	}

	scanConfigFileName, err := getScanConfigFileName(config)
	if err != nil {
		return provider.FatalError{Err: err}
	}
	err = os.Remove(path.Dir(scanConfigFileName))
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to remove scan config file. Provider=%s: %w", models.Docker, err),
		}
	}

	containerId, err := c.getContainerIdFromContainerName(ctx, scanName)
	if err != nil {
		return provider.FatalError{Err: err}
	}
	err = c.dockerClient.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return provider.FatalError{
			Err: fmt.Errorf("failed to remove scan container. Provider=%s: %w", models.Docker, err),
		}
	}
	return nil
}
