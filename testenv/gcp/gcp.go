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

package gcp

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/installation"
	"github.com/openclarity/vmclarity/testenv/utils/docker"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/openclarity/vmclarity/testenv/types"
	envtypes "github.com/openclarity/vmclarity/testenv/types"
	"github.com/openclarity/vmclarity/testenv/utils"
	"google.golang.org/api/deploymentmanager/v2"
	"google.golang.org/api/iam/v1"

	"github.com/openclarity/vmclarity/testenv/gcp/asset"
)

type ContextKeyType string

const (
	GCPClientContextKey ContextKeyType = "GCPClient"
	DefaultRemoteUser                  = "vmclarity"
	ProjectID                          = "gcp-osedev-nprd-52462"
	Zone                               = "us-central1-f"
	TestAssetName                      = "vmclarity-test-asset"

	VMClarityInstallScriptSchema = "components/vmclarity_install_script.py.schema"
	VMClarityServerSchema        = "components/vmclarity-server.py.schema"
	VMClaritySchema              = "vmclarity.py.schema"
	VMClarity                    = "vmclarity.py"
	VMClarityInstallScript       = "components/vmclarity_install_script.py"
	VMClarityInstall             = "components/vmclarity-install.sh"
	VMClarityServer              = "components/vmclarity-server.py"
	Network                      = "components/network.py"
	FirewallRules                = "components/firewall-rules.py"
	StaticIP                     = "components/static-ip.py"
	ServiceAccount               = "components/service-account.py"
	Roles                        = "components/roles.py"
	CloudRouter                  = "components/cloud-router.py"
)

type GCPEnv struct {
	workDir         string
	dm              *deploymentmanager.Service
	instancesClient *compute.InstancesClient
	iamService      *iam.Service
	serverIP        *string
	envName         string

	sshKeyPair          *utils.SSHKeyPair
	sshPortForwardInput *utils.SSHForwardInput
	sshPortForward      *utils.SSHPortForward

	*docker.DockerHelper
}

func (e *GCPEnv) SetUp(ctx context.Context) error {
	vmclarityConfigExampleYaml, err := installation.GCPManifestBundle.ReadFile("vmclarity-config.example.yaml")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	example := strings.Replace(string(vmclarityConfigExampleYaml), "<SSH Public Key>", string(e.sshKeyPair.PublicKey), -1)

	vmclarityPy, err := installation.GCPManifestBundle.ReadFile(VMClarity)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityPySchema, err := installation.GCPManifestBundle.ReadFile(VMClaritySchema)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityServerPy, err := installation.GCPManifestBundle.ReadFile(VMClarityServer)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	networkPy, err := installation.GCPManifestBundle.ReadFile(Network)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	firewallRulesPy, err := installation.GCPManifestBundle.ReadFile(FirewallRules)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	staticIpPy, err := installation.GCPManifestBundle.ReadFile(StaticIP)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	serviceAccountPy, err := installation.GCPManifestBundle.ReadFile(ServiceAccount)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	rolesPy, err := installation.GCPManifestBundle.ReadFile(Roles)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	cloudRouterPy, err := installation.GCPManifestBundle.ReadFile(CloudRouter)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityInstallScriptPy, err := installation.GCPManifestBundle.ReadFile(VMClarityInstallScript)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityInstallSh, err := installation.GCPManifestBundle.ReadFile(VMClarityInstall)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityServerPySchema, err := installation.GCPManifestBundle.ReadFile(VMClarityServerSchema)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	vmclarityInstallScriptPySchema, err := installation.GCPManifestBundle.ReadFile(VMClarityInstallScriptSchema)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	op, err := e.dm.Deployments.Insert(
		ProjectID,
		&deploymentmanager.Deployment{
			Name: e.envName,
			Target: &deploymentmanager.TargetConfiguration{
				Config: &deploymentmanager.ConfigFile{
					Content: example,
				},
				Imports: []*deploymentmanager.ImportFile{
					{
						Content: string(vmclarityPy),
						Name:    VMClarity,
					},
					{
						Content: string(vmclarityPySchema),
						Name:    VMClaritySchema,
					},
					{
						Content: string(vmclarityServerPy),
						Name:    VMClarityServer,
					},
					{
						Content: string(networkPy),
						Name:    Network,
					},
					{
						Content: string(firewallRulesPy),
						Name:    FirewallRules,
					},
					{
						Content: string(staticIpPy),
						Name:    StaticIP,
					},
					{
						Content: string(serviceAccountPy),
						Name:    ServiceAccount,
					},
					{
						Content: string(rolesPy),
						Name:    Roles,
					},
					{
						Content: string(cloudRouterPy),
						Name:    CloudRouter,
					},
					{
						Content: string(vmclarityInstallScriptPy),
						Name:    filepath.Base(VMClarityInstallScript),
					},
					{
						Content: string(vmclarityInstallSh),
						Name:    filepath.Base(VMClarityInstall),
					},
					{
						Content: string(vmclarityServerPySchema),
						Name:    VMClarityServerSchema,
					},
					{
						Content: string(vmclarityInstallScriptPySchema),
						Name:    VMClarityInstallScriptSchema,
					},
				},
			},
		},
	).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to set up the deployment: %w", err)
	}

	for {
		op, err = e.dm.Operations.Get(ProjectID, op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
		if op.Status == "DONE" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	err = asset.Create(ctx, e.instancesClient, ProjectID, Zone, TestAssetName)
	if err != nil {
		return fmt.Errorf("failed to create test asset: %w", err)
	}

	if err = e.afterSetUp(ctx); err != nil {
		return fmt.Errorf("failed to run after setup: %w", err)
	}

	return nil
}

func (e *GCPEnv) TearDown(ctx context.Context) error {
	e.sshPortForward.Stop()

	op, err := e.dm.Deployments.Delete(ProjectID, e.envName).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("unable to delete deployment: %w", err)
	}
	for {
		op, err = e.dm.Operations.Get(ProjectID, op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
		if op.Status == "DONE" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	err = asset.Delete(ctx, e.instancesClient, ProjectID, Zone, TestAssetName)
	if err != nil {
		return fmt.Errorf("failed to delete test asset: %w", err)
	}

	err = os.RemoveAll(e.workDir)
	if err != nil {
		return fmt.Errorf("failed to remove work directory: %w", err)
	}

	_, err = e.iamService.Projects.Roles.Undelete("projects/"+ProjectID+"/roles/vmclarity_"+e.envName+"_discoverer_snapshotter", &iam.UndeleteRoleRequest{}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to undelete discoverer snapshotter role: %w", err)
	}
	_, err = e.iamService.Projects.Roles.Undelete("projects/"+ProjectID+"/roles/vmclarity_"+e.envName+"_scanner", &iam.UndeleteRoleRequest{}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to undelete scanner role: %w", err)
	}

	return nil
}

func (e *GCPEnv) ServiceLogs(_ context.Context, _ []string, startTime time.Time, stdout, stderr io.Writer) error {
	input := &utils.SSHJournalctlInput{
		PrivateKey: e.sshKeyPair.PrivateKey,
		PublicKey:  e.sshKeyPair.PublicKey,
		User:       DefaultRemoteUser,
		Host:       *e.serverIP,
		WorkDir:    e.workDir,
		Service:    "docker",
	}

	err := utils.GetServiceLogs(input, startTime, stdout, stderr)
	if err != nil {
		return fmt.Errorf("failed to get service logs: %w", err)
	}

	return nil
}

func (e *GCPEnv) Endpoints(_ context.Context) (*envtypes.Endpoints, error) {
	apiURL, err := url.Parse("http://" + e.sshPortForwardInput.LocalAddressPort() + "/api")
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	uiBackendURL, err := url.Parse("http://" + e.sshPortForwardInput.LocalAddressPort() + "/ui/api")
	if err != nil {
		return nil, fmt.Errorf("failed to parse Backend API URL: %w", err)
	}

	return &types.Endpoints{
		API:       apiURL,
		UIBackend: uiBackendURL,
	}, nil
}

func (e *GCPEnv) Context(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, GCPClientContextKey, e.dm), nil
}

func New(config *Config, opts ...ConfigOptFn) (*GCPEnv, error) {
	if err := applyConfigWithOpts(config, opts...); err != nil {
		return nil, fmt.Errorf("failed to apply config options: %w", err)
	}

	dm, err := deploymentmanager.NewService(config.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create deploymentmanager: %w", err)
	}

	instancesClient, err := compute.NewInstancesRESTClient(config.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance client: %w", err)
	}

	iamService, err := iam.NewService(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM service: %w", err)
	}

	sshKeyPair := &utils.SSHKeyPair{}
	if config.PublicKeyFile != "" && config.PrivateKeyFile != "" {
		err = sshKeyPair.Load(config.PrivateKeyFile, config.PublicKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load ssh key pair: %w", err)
		}
	} else {
		sshKeyPair, err = utils.GenerateSSHKeyPair()
		if err != nil {
			return nil, fmt.Errorf("failed to generate ssh key pair: %w", err)
		}
	}

	return &GCPEnv{
		workDir:         config.WorkDir,
		dm:              dm,
		instancesClient: instancesClient,
		iamService:      iamService,
		sshKeyPair:      sshKeyPair,
		envName:         config.EnvName,
	}, nil
}
