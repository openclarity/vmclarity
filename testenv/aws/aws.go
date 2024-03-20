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

package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/openclarity/vmclarity/testenv/aws/asset"
	"github.com/openclarity/vmclarity/testenv/types"
	"github.com/openclarity/vmclarity/testenv/utils"
)

type ContextKeyType string

const (
	AWSClientContextKey ContextKeyType = "AWSClient"
	AWSKeyName          string         = "vmclarity-testenv-key"
)

// AWS Environment.
type AWSEnv struct {
	client      *cloudformation.Client
	ec2Client   *ec2.Client
	s3Client    *s3.Client
	testAsset   *asset.Asset
	server      *Server
	stackName   string
	templateURL string
	workDir     string
	region      string
	sshKeyPair  *utils.SSHKeyPair
	meta        map[string]interface{}
}

type Server struct {
	InstanceID string
	PublicIP   string
}

// Setup AWS test environment from cloud formation template.
// * Create a new CloudFormation stack from template
// (upload template file to S3 is required since the template is larger than 51,200 bytes).
// * Create test asset.
func (e *AWSEnv) SetUp(ctx context.Context) error {
	// Prepare stack
	err := e.prepareStack(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare stack: %w", err)
	}

	// Create a new CloudFormation stack from template
	_, err = e.client.CreateStack(
		ctx,
		&cloudformation.CreateStackInput{
			StackName:    &e.stackName,
			Capabilities: []cloudformationtypes.Capability{cloudformationtypes.CapabilityCapabilityIam},
			TemplateURL:  &e.templateURL,
			Parameters: []cloudformationtypes.Parameter{
				{ParameterKey: aws.String("KeyName"), ParameterValue: aws.String(AWSKeyName)},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	// Create a new test asset
	err = e.testAsset.Create(ctx, e.ec2Client)
	if err != nil {
		return fmt.Errorf("failed to create test asset: %w", err)
	}

	return nil
}

func (e *AWSEnv) TearDown(ctx context.Context) error {
	// Delete the CloudFormation stack
	_, err := e.client.DeleteStack(
		ctx,
		&cloudformation.DeleteStackInput{
			StackName: &e.stackName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	// Cleanup stack
	err = e.cleanupStack(ctx)
	if err != nil {
		return fmt.Errorf("failed to cleanup stack: %w", err)
	}

	// Delete the test asset
	err = e.testAsset.Delete(ctx, e.ec2Client)
	if err != nil {
		return fmt.Errorf("failed to delete test asset: %w", err)
	}

	return nil
}

func (e *AWSEnv) ServicesReady(ctx context.Context) (bool, error) {
	// Check infrastructure status before checking services
	ready, err := e.infrastructureReady(ctx)
	if err != nil || !ready {
		return false, err
	}

	// Get list of services
	services, err := e.Services(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get services: %w", err)
	}

	// Assume that the services are ready if all server containers are in ready state
	ready = true
	for _, service := range services {
		if service.GetState() != types.ServiceStateReady {
			ready = false
			break
		}
	}

	return ready, nil
}

func (e *AWSEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, stderr io.Writer) error {
	panic("not implemented")
}

func (e *AWSEnv) Services(ctx context.Context) (types.Services, error) {
	// Get Docker Container list from VMClarity Server
	containerList, err := utils.GetRemoteDockerContainerList(ctx, e.sshKeyPair.PrivateKeyFile, e.server.PublicIP)
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %w", err)
	}

	var services types.Services
	for _, container := range *containerList {
		services = append(services, &Service{
			ID:    container.ID,
			State: convertStateFromDocker(container.State),
		})
	}

	return services, nil
}

func (e *AWSEnv) Endpoints(ctx context.Context) (*types.Endpoints, error) {
	remoteHost := e.server.PublicIP
	remotePort := "80"
	localHost := "localhost"
	localPort := "8080"

	// Run SSH tunnel to remote VMClarity server
	go utils.RunSSHTunnel(ctx, e.sshKeyPair.PrivateKeyFile, remoteHost, remotePort, localPort)

	// Wait for SSH tunnel to be ready
	time.Sleep(10 * time.Second) // nolint:gomnd

	endpoints := new(types.Endpoints)
	endpoints.SetAPI("http", localHost, localPort, "/api")
	endpoints.SetUIBackend("http", localHost, localPort, "/ui/api")

	return endpoints, nil
}

func (e *AWSEnv) Context(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, AWSClientContextKey, e.client), nil
}

func New(config *Config, opts ...ConfigOptFn) (*AWSEnv, error) {
	if err := applyConfigWithOpts(config, opts...); err != nil {
		return nil, fmt.Errorf("failed to apply config options: %w", err)
	}

	// Load default AWS configuration and set region
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	cfg.Region = config.Region

	// Create AWS CloudFormation client
	client := cloudformation.NewFromConfig(cfg)

	// Create AWS EC2 client
	ec2Client := ec2.NewFromConfig(cfg)

	// Create AWS S3 client
	s3Client := s3.NewFromConfig(cfg)

	return &AWSEnv{
		client:    client,
		ec2Client: ec2Client,
		s3Client:  s3Client,
		stackName: config.EnvName,
		workDir:   config.WorkDir,
		region:    config.Region,
		sshKeyPair: &utils.SSHKeyPair{
			PublicKeyFile:  config.PublicKeyFile,
			PrivateKeyFile: config.PrivateKeyFile,
			Temporary:      false,
		},
		testAsset: &asset.Asset{},
		meta: map[string]interface{}{
			"environment": "aws",
			"name":        config.EnvName,
		},
	}, nil
}
