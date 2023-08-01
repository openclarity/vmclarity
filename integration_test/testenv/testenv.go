package testenv

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"time"
)

const (
	vmClarityBackendContainerName = "vmclarity-backend"
	stateRunning                  = "running"
	healthStateHealthy            = "healthy"
)

type Environment struct {
	composer api.Service
	project  *types.Project
	reuse    bool
}

func New(o *cli.ProjectOptions, reuse bool) (*Environment, error) {
	project, err := cli.ProjectFromOptions(o)
	if err != nil {
		return nil, err
	}

	for i, service := range project.Services {
		service.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     service.Name,
			api.WorkingDirLabel:  project.WorkingDir,
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False",
		}
		project.Services[i] = service
	}

	cmd, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}

	opts := cliflags.NewClientOptions()

	if err = cmd.Initialize(opts); err != nil {
		return nil, err
	}

	return &Environment{
		composer: compose.NewComposeService(cmd),
		project:  project,
		reuse:    reuse,
	}, nil
}

func (e *Environment) Start(ctx context.Context) error {
	if e.reuse {
		return nil
	}

	timeout := 1 * time.Minute
	opts := api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: true,
			QuietPull:     true,
			Timeout:       &timeout,
			Services:      e.Services(),
			Inherit:       false,
		},
		Start: api.StartOptions{
			Project:     e.project,
			Wait:        true,
			WaitTimeout: 2 * time.Minute,
			Services:    e.Services(),
		},
	}
	return e.composer.Up(ctx, e.project, opts)
}

func (e *Environment) Stop(ctx context.Context) error {
	if e.reuse {
		return nil
	}

	timeout := 1 * time.Minute
	opts := api.DownOptions{
		RemoveOrphans: true,
		Project:       e.project,
		Volumes:       true,
		Timeout:       &timeout,
	}
	return e.composer.Down(ctx, e.project.Name, opts)
}

func (e *Environment) ServicesReady(ctx context.Context) (bool, error) {
	services := e.Services()

	ps, err := e.composer.Ps(
		ctx,
		e.project.Name,
		api.PsOptions{
			Services: services,
			Project:  e.project,
		},
	)
	if err != nil {
		return false, err
	}

	if len(services) != len(ps) {
		return false, nil
	}

	for _, c := range ps {
		if c.State != stateRunning && c.Health != healthStateHealthy {
			return false, nil
		}
	}

	return true, nil
}

func (e *Environment) Services() []string {
	services := make([]string, len(e.project.Services))
	for i, srv := range e.project.Services {
		services[i] = srv.Name
	}
	return services
}

func (e *Environment) VMClarityURL() (*url.URL, error) {
	var vmClarityBackend types.ServiceConfig
	var ok bool

	for _, srv := range e.project.Services {
		if srv.Name == vmClarityBackendContainerName {
			vmClarityBackend = srv
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.Errorf("container with name %s is not available", vmClarityBackendContainerName)
	}

	if len(vmClarityBackend.Ports) < 1 {
		return nil, errors.Errorf("container with name %s has no ports published", vmClarityBackendContainerName)
	}

	port := vmClarityBackend.Ports[0].Published
	hostIP := vmClarityBackend.Ports[0].HostIP
	if hostIP == "" {
		hostIP = "127.0.0.1"
	}

	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", hostIP, port),
		Path:   "api",
	}, nil
}
