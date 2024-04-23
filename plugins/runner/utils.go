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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	imagetypes "github.com/docker/docker/api/types/image"
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
