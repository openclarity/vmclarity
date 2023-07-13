package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

func (c *Client) getContainers(ctx context.Context) ([]models.AssetType, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	// List all docker containers
	containers, err := c.dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Results will be written to assets concurrently
	assetMu := sync.Mutex{}
	assets := make([]models.AssetType, 0, len(containers))

	// Process each container in an independent processor goroutine
	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	for _, container := range containers {
		errGroup.Go(
			// errGroup expects a function with empty signature, so we use a function
			// generator to enable adding arguments. This avoids issues when using loop
			// variables in goroutines.
			//
			// If any processor returns an error, it will stop all processors.
			// TODO: Decide what the acceptance criteria should be (e.g. >= 50% container processed)
			func(container types.Container) func() error {
				return func() error {
					// Get container info
					info, err := c.dockerClient.ContainerInspect(errGroupCtx, container.ID)
					if err != nil {
						return fmt.Errorf("failed to inspect container: %w", err)
					}

					containerCreatedAt, err := time.Parse(time.RFC3339, info.Created)
					if err != nil {
						return fmt.Errorf("failed to parse time: %w", err)
					}

					// Get container image info
					imageInfo, err := c.getContainerImageInfo(errGroupCtx, container.ImageID)
					if err != nil {
						// TODO (paralta) If image not required this should not be fatal -- resolved
						logger.Warnf("Failed to get container. id=%v: %v", container.ID, err)
						return nil
					}

					asset := models.AssetType{}
					err = asset.FromContainerInfo(models.ContainerInfo{
						ContainerName: utils.PointerTo(info.Name),
						CreatedAt:     utils.PointerTo(containerCreatedAt),
						Id:            utils.PointerTo(container.ID),
						Image:         utils.PointerTo(imageInfo),
						Labels:        convertTags(container.Labels),
						Location:      nil, // TODO (paralta) Clarify what is location
						ObjectType:    "ContainerInfo",
					})
					if err != nil {
						return fmt.Errorf("failed to create AssetType from ContainerInfo: %w", err)
					}

					// Write to assets
					assetMu.Lock()
					assets = append(assets, asset)
					assetMu.Unlock()

					return nil
				}
			}(container),
		)
	}

	// This will block until all the processors have executed successfully or until
	// first error. If an error is returned by any processors, errGroup will cancel
	// execution via errGroupCtx and return that error.
	err = errGroup.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to process containers: %w", err)
	}

	return assets, nil
}
