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

package containerd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/errdefs"
	containerdImages "github.com/containerd/containerd/images"
	"github.com/containerd/containerd/images/archive"
	"github.com/containerd/containerd/leases"
	criConstants "github.com/containerd/containerd/pkg/cri/constants"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/nerdctl/pkg/imgutil"
	"github.com/containerd/nerdctl/pkg/imgutil/commit"
	"github.com/containerd/nerdctl/pkg/labels/k8slabels"
	"github.com/containers/image/v5/docker/reference"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/containerruntimediscovery/types"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"github.com/openclarity/vmclarity/utils/log"
)

type ContainerdDiscoverer struct {
	client *containerd.Client
}

var _ types.Discoverer = &ContainerdDiscoverer{}

func NewContainerdDiscoverer(ctx context.Context) (types.Discoverer, error) {
	// Containerd supports multiple namespaces so that a single daemon can
	// be used by multiple clients like Docker and Kubernetes and the
	// resources will not conflict etc. In order to discover all the
	// containers for kubernetes we need to set the kubernetes namespace as
	// the default for our client.
	client, err := containerd.New("/var/run/containerd/containerd.sock", containerd.WithDefaultNamespace(criConstants.K8sContainerdNamespace))
	if err != nil {
		return nil, fmt.Errorf("failed to create containerd client: %w", err)
	}

	_, err = client.ListImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	return &ContainerdDiscoverer{
		client: client,
	}, nil
}

type onFound func(image containerd.Image) (bool, error)

func (cd *ContainerdDiscoverer) imageIDWalk(ctx context.Context, imageID string, f onFound) (bool, error) {
	var found bool

	// ContainerD doesn't allow to filter images by config digest, so we
	// have to walk all the images to find all the images by ID and then
	// merge them together.
	images, err := cd.client.ListImages(ctx)
	if err != nil {
		return found, fmt.Errorf("failed to list images: %w", err)
	}

	for _, image := range images {
		configDescriptor, err := image.Config(ctx)
		if err != nil {
			return found, fmt.Errorf("failed to load image config descriptor: %w", err)
		}
		id := configDescriptor.Digest.String()

		if id != imageID {
			continue
		}

		found = true

		stop, err := f(image)
		if err != nil {
			return found, err
		}

		if stop {
			break
		}
	}

	return found, nil
}

func (cd *ContainerdDiscoverer) Images(ctx context.Context) ([]models.ContainerImageInfo, error) {
	logger := log.GetLoggerFromContextOrDefault(ctx)

	images, err := cd.client.ListImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	imageSet := map[string]models.ContainerImageInfo{}
	for _, image := range images {
		// Ignore our transient images used for snapshoting the
		// containers they will be cleaned up as soon as the export is
		// done.
		if strings.HasPrefix(image.Name(), "vmclarity.io/container-snapshot:") {
			continue
		}

		cii, err := cd.getContainerImageInfo(ctx, image)
		if err != nil {
			return nil, fmt.Errorf("unable to convert image %s to container image info: %w", image.Name(), err)
		}

		if cii.ImageID == "" {
			logger.Warnf("found image with empty ImageID: %s", cii.String())
			continue
		}

		existing, ok := imageSet[cii.ImageID]
		if ok {
			merged, err := existing.Merge(cii)
			if err != nil {
				return nil, fmt.Errorf("unable to merge image %v with %v: %w", existing, cii, err)
			}
			cii = merged
		}
		imageSet[cii.ImageID] = cii
	}

	result := []models.ContainerImageInfo{}
	for _, image := range imageSet {
		result = append(result, image)
	}
	return result, nil
}

func (cd *ContainerdDiscoverer) Image(ctx context.Context, imageID string) (models.ContainerImageInfo, error) {
	var result models.ContainerImageInfo
	found, err := cd.imageIDWalk(ctx, imageID, func(image containerd.Image) (bool, error) {
		cii, err := cd.getContainerImageInfo(ctx, image)
		if err != nil {
			return false, fmt.Errorf("unable to convert image %s to container image info: %w", image.Name(), err)
		}

		result, err = result.Merge(cii)
		if err != nil {
			return false, fmt.Errorf("unable to merge image %v with %v: %w", result, cii, err)
		}

		return false, nil
	})
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to walk all image: %w", err)
	}
	if !found {
		return models.ContainerImageInfo{}, types.ErrNotFound
	}

	return result, nil
}

func (cd *ContainerdDiscoverer) getContainerImageInfo(ctx context.Context, image containerd.Image) (models.ContainerImageInfo, error) {
	configDescriptor, err := image.Config(ctx)
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to load image config descriptor: %w", err)
	}
	id := configDescriptor.Digest.String()

	imageSpec, err := image.Spec(ctx)
	if err != nil {
		return models.ContainerImageInfo{}, fmt.Errorf("failed to load image spec: %w", err)
	}

	// NOTE(sambetts) We can not use image.Size as it gives us the size of
	// the compressed layers and not the real size of the content.
	snapshotter := cd.client.SnapshotService(containerd.DefaultSnapshotter)
	// NOTE(chrisgacsal): ignore error as determining size of the image is not critical
	size, _ := imgutil.UnpackedImageSize(ctx, snapshotter, image)

	repoTags, repoDigests := ParseImageReferences([]string{image.Name()})

	return models.ContainerImageInfo{
		ImageID:      id,
		Architecture: utils.PointerTo(imageSpec.Architecture),
		Labels:       models.MapToTags(imageSpec.Config.Labels),
		RepoTags:     &repoTags,
		RepoDigests:  &repoDigests,
		ObjectType:   "ContainerImageInfo",
		Os:           utils.PointerTo(imageSpec.OS),
		Size:         utils.PointerTo(int(size)),
	}, nil
}

// TODO(sambetts) Support auth config for fetching private images if they are missing.
func (cd *ContainerdDiscoverer) ExportImage(ctx context.Context, imageID string) (io.ReadCloser, error) {
	var img containerd.Image
	found, err := cd.imageIDWalk(ctx, imageID, func(image containerd.Image) (bool, error) {
		img = image
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk all images: %w", err)
	}
	if !found {
		return nil, types.ErrNotFound
	}

	// NOTE(sambetts) When running in Kubernetes containerd can be
	// configured to garbage collect the un-expanded blobs from the content
	// store after they are converted to a rootfs snapshot that is used to
	// boot containers. For this reason we need to re-fetch the image to
	// ensure that all the required blobs for export are in the content
	// store.
	// nolint: dogsled
	_, _, _, missing, err := containerdImages.Check(ctx, cd.client.ContentStore(), img.Target(), platforms.Default())
	if err != nil {
		return nil, fmt.Errorf("unable to check image in content store: %w", err)
	}
	if len(missing) > 0 {
		imageInfo, err := cd.Image(ctx, imageID)
		if err != nil {
			return nil, fmt.Errorf("failed to get image info to export: %w", err)
		}
		if imageInfo.RepoDigests == nil || len(*imageInfo.RepoDigests) == 0 {
			return nil, fmt.Errorf("image has no known repo digests can not safely fetch it")
		}

		// TODO(sambetts) Maybe try all the digests in case one has gone missing?
		ref := (*imageInfo.RepoDigests)[0]
		img, err = cd.client.Pull(ctx, ref)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch image %s: %w", ref, err)
		}
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := cd.client.Export(
			ctx,
			pw,
			archive.WithImage(cd.client.ImageService(), img.Name()),
			archive.WithPlatform(platforms.Default()),
		)
		if err != nil {
			log.GetLoggerFromContextOrDefault(ctx).Errorf("failed to export image: %v", err)
		}
	}()
	return pr, nil
}

// ParseImageReferences parses a list of arbitrary image references and returns
// the repotags and repodigests.
func ParseImageReferences(refs []string) ([]string, []string) {
	var tags, digests []string
	for _, ref := range refs {
		parsed, err := reference.ParseAnyReference(ref)
		if err != nil {
			continue
		}
		if _, ok := parsed.(reference.Canonical); ok {
			digests = append(digests, parsed.String())
		} else if _, ok := parsed.(reference.Tagged); ok {
			tags = append(tags, parsed.String())
		}
	}
	return tags, digests
}

func (cd *ContainerdDiscoverer) Containers(ctx context.Context) ([]models.ContainerInfo, error) {
	containers, err := cd.client.Containers(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list containers: %w", err)
	}

	result := make([]models.ContainerInfo, len(containers))
	for i, container := range containers {
		// Get container info
		info, err := cd.getContainerInfo(ctx, container)
		if err != nil {
			return nil, fmt.Errorf("failed to convert container to ContainerInfo: %w", err)
		}
		result[i] = info
	}
	return result, nil
}

func (cd *ContainerdDiscoverer) Container(ctx context.Context, containerID string) (models.ContainerInfo, error) {
	container, err := cd.client.LoadContainer(ctx, containerID)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return models.ContainerInfo{}, types.ErrNotFound
		}
		return models.ContainerInfo{}, fmt.Errorf("failed to get container from store: %w", err)
	}

	return cd.getContainerInfo(ctx, container)
}

// nolint: cyclop
func (cd *ContainerdDiscoverer) ExportContainer(ctx context.Context, containerID string) (io.ReadCloser, func(), error) {
	clean := &types.Cleanup{}
	defer clean.Clean()

	container, err := cd.client.LoadContainer(ctx, containerID)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, func() {}, types.ErrNotFound
		}
		return nil, func() {}, fmt.Errorf("failed to get container from store: %w", err)
	}

	img, err := container.Image(ctx)
	if err != nil {
		return nil, func() {}, fmt.Errorf("unable to get image from container %s: %w", containerID, err)
	}

	// NOTE(sambetts) When running in Kubernetes containerd can be
	// configured to garbage collect the un-expanded blobs from the content
	// store after they are converted to a rootfs snapshot that is used to
	// boot containers. For this reason we need to re-fetch the image to
	// ensure that all the required blobs for export are in the content
	// store.
	// nolint: dogsled
	_, _, _, missing, err := containerdImages.Check(ctx, cd.client.ContentStore(), img.Target(), platforms.Default())
	if err != nil {
		return nil, func() {}, fmt.Errorf("unable to check image in content store: %w", err)
	}
	if len(missing) > 0 {
		configDescriptor, err := img.Config(ctx)
		if err != nil {
			return nil, func() {}, fmt.Errorf("failed to load image config descriptor: %w", err)
		}
		imageID := configDescriptor.Digest.String()

		imageInfo, err := cd.Image(ctx, imageID)
		if err != nil {
			return nil, func() {}, fmt.Errorf("failed to get image info to export: %w", err)
		}
		if imageInfo.RepoDigests == nil || len(*imageInfo.RepoDigests) == 0 {
			return nil, func() {}, fmt.Errorf("image has no known repo digests can not safely fetch it")
		}

		// TODO(sambetts) Maybe try all the digests in case one has gone missing?
		ref := (*imageInfo.RepoDigests)[0]
		_, err = cd.client.Pull(ctx, ref)
		if err != nil {
			return nil, func() {}, fmt.Errorf("failed to fetch image %s: %w", ref, err)
		}
	}

	ctx, done, err := cd.client.WithLease(ctx, leases.WithRandomID(), leases.WithExpiration(1*time.Hour))
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to get lease from containerd: %w", err)
	}
	clean.Add(func() {
		err := done(ctx)
		if err != nil {
			log.GetLoggerFromContextOrDefault(ctx).Errorf("failed to release lease: %v", err)
		}
	})

	imageName := fmt.Sprintf("vmclarity.io/container-snapshot:%s", containerID)
	_, err = commit.Commit(ctx, cd.client, container, &commit.Opts{
		Author:  "VMClarity",
		Message: fmt.Sprintf("Snapshot of container %s for security scanning", containerID),
		Ref:     imageName,
		Pause:   false,
	})
	if err != nil {
		return nil, func() {}, fmt.Errorf("unable to commit container to image: %w", err)
	}
	clean.Add(func() {
		err := cd.client.ImageService().Delete(ctx, imageName)
		if err != nil {
			log.GetLoggerFromContextOrDefault(ctx).Errorf("failed to clean up snapshot %s for container %s: %v", imageName, containerID, err)
		}
	})

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := cd.client.Export(
			ctx,
			pw,
			archive.WithImage(cd.client.ImageService(), imageName),
			archive.WithPlatform(platforms.Default()),
		)
		if err != nil {
			log.GetLoggerFromContextOrDefault(ctx).Errorf("failed to export container snapshot: %v", err)
		}
	}()

	return pr, clean.Release(), nil
}

func (cd *ContainerdDiscoverer) getContainerInfo(ctx context.Context, container containerd.Container) (models.ContainerInfo, error) {
	id := container.ID()

	labels, err := container.Labels(ctx)
	if err != nil {
		return models.ContainerInfo{}, fmt.Errorf("unable to get labels for container %s: %w", id, err)
	}
	// If this doesn't exist then use empty string as the name. Containerd
	// doesn't have the concept of a Name natively.
	name := labels[k8slabels.ContainerName]

	info, err := container.Info(ctx)
	if err != nil {
		return models.ContainerInfo{}, fmt.Errorf("unable to get info for container %s: %w", id, err)
	}
	createdAt := info.CreatedAt

	image, err := container.Image(ctx)
	if err != nil {
		return models.ContainerInfo{}, fmt.Errorf("unable to get image from container %s: %w", id, err)
	}

	configDescriptor, err := image.Config(ctx)
	if err != nil {
		return models.ContainerInfo{}, fmt.Errorf("failed to load image config descriptor: %w", err)
	}
	imageID := configDescriptor.Digest.String()

	imageInfo, err := cd.Image(ctx, imageID)
	if err != nil {
		return models.ContainerInfo{}, fmt.Errorf("unable to convert image %s to container image info: %w", image.Name(), err)
	}

	return models.ContainerInfo{
		ContainerID:   container.ID(),
		ContainerName: utils.PointerTo(name),
		CreatedAt:     utils.PointerTo(createdAt),
		Image:         utils.PointerTo(imageInfo),
		Labels:        models.MapToTags(labels),
		ObjectType:    "ContainerInfo",
	}, nil
}
