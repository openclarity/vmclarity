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
)

func (c *Client) getImages(ctx context.Context) ([]models.AssetType, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	// List all docker images
	images, err := c.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Results will be written to assets concurrently
	assetMu := sync.Mutex{}
	assets := make([]models.AssetType, 0, len(images))

	// Process each image in an independent processor goroutine
	processGroup, processCtx := errgroup.WithContext(ctx)
	for _, image := range images {
		processGroup.Go(
			// processGroup expects a function with empty signature, so we use a function
			// generator to enable adding arguments. This avoids issues when using loop
			// variables in goroutines via shared memory space.
			//
			// If any processor returns an error, it will stop all processors.
			// IDEA: Decide what the acceptance criteria should be (e.g. >= 50% images processed)
			func(image types.ImageSummary) func() error {
				return func() error {
					// Get container image info
					info, err := c.getContainerImageInfo(processCtx, image.ID)
					if err != nil {
						logger.Warnf("Failed to get image. id=%v: %v", image.ID, err)
						return nil // skip fail
					}

					// Convert to asset
					asset := models.AssetType{}
					err = asset.FromContainerImageInfo(info)
					if err != nil {
						return fmt.Errorf("failed to create AssetType from ContainerImageInfo: %w", err)
					}

					// Write to assets
					assetMu.Lock()
					assets = append(assets, asset)
					assetMu.Unlock()

					return nil
				}
			}(image),
		)
	}

	// This will block until all the processors have executed successfully or until
	// the first error. If an error is returned by any processors, processGroup will
	// cancel execution via processCtx and return that error.
	err = processGroup.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to process images: %w", err)
	}

	return assets, nil
}

func (c *Client) getContainerImageInfo(ctx context.Context, imageID string) (models.ContainerImageInfo, error) {
	image, _, err := c.dockerClient.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to inspect image: %w", err)
	}

	return models.ContainerImageInfo{
		Architecture: utils.PointerTo(image.Architecture),
		Id:           utils.PointerTo(image.ID),
		Labels:       convertTags(image.Config.Labels),
		Name:         utils.PointerTo(image.Config.Image),
		ObjectType:   "ContainerImageInfo",
		Os:           utils.PointerTo(image.Os),
		Size:         utils.PointerTo(int(image.Size)),
	}, nil
}
