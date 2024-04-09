package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/filters"
	imagetypes "github.com/docker/docker/api/types/image"
)

func (r *Runner) pullImage(ctx context.Context, imageName string) error {
	images, err := r.dockerClient.ImageList(ctx, imagetypes.ListOptions{
		Filters: filters.NewArgs(filters.Arg("reference", imageName)),
	})
	if err != nil {
		return fmt.Errorf("failed to get images: %w", err)
	}

	if len(images) == 0 {
		pullResp, err := r.dockerClient.ImagePull(ctx, imageName, imagetypes.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}

		// consume output
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}

	return nil
}

func getScannerConfigSourcePath(name string) string {
	return filepath.Join(os.TempDir(), name+"-plugin.json")
}

func getScannerConfigDestinationPath() string {
	return filepath.Join("/plugin.json")
}
