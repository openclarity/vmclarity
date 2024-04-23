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
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	imagetypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func pullImage(ctx context.Context, client *client.Client, imageName string) error {
	images, err := client.ImageList(ctx, imagetypes.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return fmt.Errorf("failed to get images: %w", err)
	}

	if len(images) == 0 {
		pullResp, err := client.ImagePull(ctx, imageName, imagetypes.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}

		// consume output
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}

	return nil
}

func (r *runner) copyConfigToContainer(ctx context.Context) error {
	// Write scanner config file to temp dir
	src := getScannerConfigSourcePath(r.config.Name)
	err := os.WriteFile(src, []byte(r.config.ScannerConfig), 0o400) // nolint:gomnd
	if err != nil {
		return fmt.Errorf("failed write scanner config file: %w", err)
	}

	// Create tar archive from scan config file
	srcInfo, err := archive.CopyInfoSourcePath(src, false)
	if err != nil {
		return fmt.Errorf("failed to get copy info: %w", err)
	}
	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return fmt.Errorf("failed to create tar archive: %w", err)
	}
	defer srcArchive.Close()

	// Prepare archive for copy
	dstInfo := archive.CopyInfo{Path: getScannerConfigDestinationPath()}
	dst, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return fmt.Errorf("failed to prepare archive: %w", err)
	}
	defer preparedArchive.Close()

	// Copy scan config file to container
	err = r.dockerClient.CopyToContainer(ctx, r.containerID, dst, preparedArchive, types.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("failed to copy config file to container: %w", err)
	}

	return nil
}

func (r *runner) getPluginContainerMounts(ctx context.Context) ([]mount.Mount, error) {
	// Get container ID
	containerID, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get container ID: %w", err)
	}

	// Get container mounts
	container, err := r.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			// Not running in a container
			return []mount.Mount{
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
			}, nil
		}
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	// Running in a container
	// Convert MountPoint to mount.Mount
	var mounts []mount.Mount
	for _, p := range container.Mounts {
		if p.Destination == r.config.InputDir {
			mounts = append(mounts, mount.Mount{
				Type:   p.Type,
				Source: p.Source,
				Target: p.Destination,
			})
		}
	}

	return mounts, nil
}

func (r *runner) getNetworkingConfig(ctx context.Context) (*network.NetworkingConfig, error) {
	// Get container ID
	containerID, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get container ID: %w", err)
	}

	// Get container mounts
	container, err := r.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			// Not running in a container
			return nil, nil
		}
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	// Running in a container
	for _, net := range container.NetworkSettings.Networks {
		return &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				r.containerID: {
					NetworkID: net.NetworkID,
				},
			},
		}, nil
	}

	return nil, nil
}
