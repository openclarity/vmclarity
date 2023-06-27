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
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
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

	dockerClient := &client.Client{}

	err = client.FromEnv(dockerClient)
	if err != nil {
		return nil, fmt.Errorf("failed to load provider configuration. Provider=%s: %w", models.Docker, err)
	}

	return &Client{
		dockerClient: dockerClient,
		config:       config,
	}, nil
}

func (c Client) Kind() models.CloudProvider {
	return models.Docker
}

func (c Client) DiscoverAssets(ctx context.Context) ([]models.AssetType, error) {
	var ret []models.AssetType

	// TODO: collect container images
	// TODO (paralta) split collect images and collect containers to other funcs
	// TODO (paralta) add go routines
	images, err := c.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, provider.FatalError{
			Err: fmt.Errorf("failed to get images. Provider=%s: %w", models.Docker, err),
		}
	}

	for _, i := range images {
		asset, err := c.getAssetFromImage(ctx, i.ID)
		if err != nil {
			return nil, provider.FatalError{
				Err: fmt.Errorf("failed to create AssetType from ContainerImageInfo. Provider=%s: %w", models.Docker, err),
			}
		}
		ret = append(ret, asset)
	}

	// TODO: collect containers
	containers, err := c.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, provider.FatalError{
			Err: fmt.Errorf("failed to get containers. Provider=%s: %w", models.Docker, err),
		}
	}

	for _, container := range containers {
		inspect, err := c.dockerClient.ContainerInspect(ctx, container.ID)
		if err != nil {
			return nil, provider.FatalError{
				Err: fmt.Errorf("failed to get container. Provider=%s: %w", models.Docker, err),
			}
		}

		created, err := time.Parse(time.RFC3339, inspect.Created)
		if err != nil {
			return nil, provider.FatalError{
				Err: fmt.Errorf("failed to parse time. Provider=%s: %w", models.Docker, err),
			}
		}

		imageInfo, err := c.getContainerImageInfoFromImage(ctx, container.ImageID)
		if err != nil {
			// TODO (paralta) If image not required this should not be fatal
			return nil, provider.FatalError{
				Err: err,
			}
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
			return nil, provider.FatalError{
				Err: fmt.Errorf("failed to create AssetType from ContainerInfo. Provider=%s: %w", models.Docker, err),
			}
		}
		ret = append(ret, asset)
	}

	return ret, nil
}

func (c Client) RunAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	//TODO implement me
	panic("implement me")
}

func (c Client) RemoveAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	//TODO implement me
	panic("implement me")
}

func (c Client) getContainerImageInfoFromImage(ctx context.Context, id string) (models.ContainerImageInfo, error) {
	i, _, err := c.dockerClient.ImageInspectWithRaw(ctx, id)
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to get image. Provider=%s: %w", models.Docker, err)
	}

	return models.ContainerImageInfo{
		Architecture: utils.PointerTo(i.Architecture),
		Id:           utils.PointerTo(i.ID),
		Labels:       convertTags(i.Config.Labels),
		ObjectType:   "ContainerImageInfo",
		Os:           utils.PointerTo(i.Os),
		Size:         utils.PointerTo(int(i.Size)),
	}, err
}

func (c Client) getAssetFromImage(ctx context.Context, id string) (models.AssetType, error) {
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

func convertTags(tags map[string]string) *[]models.Tag {
	ret := make([]models.Tag, 0, len(tags))
	for key, val := range tags {
		ret = append(ret, models.Tag{
			Key:   key,
			Value: val,
		})
	}
	return &ret
}
