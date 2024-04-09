package runner

import (
	"context"
	"fmt"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

const (
	proxyImage         = "traefik:v2.11"
	proxyContainerName = "plugin-scanner-proxy"
	proxyHostAddress   = "127.0.0.1:8080"
)

func (r *Runner) CreateProxyContainer(ctx context.Context) error {
	// Pull proxy image
	err := r.pullImage(ctx, proxyImage)
	if err != nil {
		return fmt.Errorf("failed to pull proxy image: %w", err)
	}

	// Get proxy port (HOST:CONTAINER)
	ports, bindings, _ := nat.ParsePortSpecs([]string{
		fmt.Sprintf("%s:80", proxyHostAddress),
	})

	// Create container
	container, err := r.dockerClient.ContainerCreate(
		ctx,
		&containertypes.Config{
			Image: proxyImage,
			Cmd: []string{
				"--api.insecure=true",
				"--providers.docker=true",
				"--entrypoints.web.address=:80",
			},
			ExposedPorts: ports,
		},
		&containertypes.HostConfig{
			PortBindings: bindings,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
			},
		},
		nil,
		nil,
		proxyContainerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create proxy container: %w", err)
	}

	// Start proxy container
	err = r.dockerClient.ContainerStart(context.Background(), container.ID, containertypes.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start proxy container: %w", err)
	}

	return nil
}
