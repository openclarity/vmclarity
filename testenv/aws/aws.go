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

package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/openclarity/vmclarity/testenv/types"
)

type ContextKeyType string

const AWSClientContextKey ContextKeyType = "AWSClient"

// AWS Environment.
type AWSEnv struct {
	client    *cloudformation.Client
	ec2Client *ec2.Client
	StackName string
	meta      map[string]interface{}
}

// Setup AWS test environment from cloud formation template.
func (e *AWSEnv) SetUp(ctx context.Context) error {
	// Read template file
	template, err := os.ReadFile("../installation/aws/VmClarity.cfn")
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Get template body
	templateBody := string(template)

	// Create a new CloudFormation stack from template
	_, err = e.client.CreateStack(
		ctx,
		&cloudformation.CreateStackInput{
			StackName:    &e.StackName,
			TemplateBody: &templateBody,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	return nil
}

func (e *AWSEnv) TearDown(ctx context.Context) error {
	// Delete the CloudFormation stack
	_, err := e.client.DeleteStack(
		ctx,
		&cloudformation.DeleteStackInput{
			StackName: &e.StackName,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete stack: %w", err)
	}

	return nil
}

func (e *AWSEnv) ServicesReady(ctx context.Context) (bool, error) {
	// Get the description of the CloudFormation stack
	stacks, err := e.client.DescribeStacks(
		ctx,
		&cloudformation.DescribeStacksInput{
			StackName: &e.StackName,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to describe stack: %w", err)
	}

	// Check if one stack found
	if len(stacks.Stacks) != 1 {
		return false, errors.New("failed to find stack")
	}

	// If the stack status is CREATE_COMPLETE, then the services are ready
	if stacks.Stacks[0].StackStatus == "CREATE_COMPLETE" {
		return true, nil
	}

	return false, nil
}

func (e *AWSEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, stderr io.Writer) error {
	panic("not implemented")
}

func (e *AWSEnv) Services(ctx context.Context) (types.Services, error) {
	// List all services in the stack
	resources, err := e.client.ListStackResources(
		ctx,
		&cloudformation.ListStackResourcesInput{
			StackName: &e.StackName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list stack resources: %w", err)
	}

	// Get VMClarity Server service
	for _, resource := range resources.StackResourceSummaries {
		if *resource.ResourceType == "AWS::EC2::Instance" {
			return types.Services{
				&Service{
					ID:          *resource.PhysicalResourceId,
					Namespace:   e.StackName,
					Application: "vmclarity",
					Component:   "server",
					State:       convertStateFromAWS(resource.ResourceStatus),
				},
			}, nil
		}
	}

	return types.Services{}, errors.New("failed to get services")
}

func (e *AWSEnv) Endpoints(ctx context.Context) (*types.Endpoints, error) {
	// Get VMClarity Server EC2 instance
	instances, err := e.Services(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	if len(instances) != 1 {
		return nil, errors.New("failed to get VMClarity Server instance")
	}

	// Describe VMClarity Server EC2 instance
	output, err := e.ec2Client.DescribeInstances(
		ctx,
		&ec2.DescribeInstancesInput{
			InstanceIds: []string{instances[0].GetID()},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	// Get Public IP of VMClarity Server
	host := output.Reservations[0].Instances[0].PublicIpAddress

	// Get port of VMClarity Server
	port := "8080"

	endpoints := new(types.Endpoints)
	endpoints.SetAPI("http", *host, port, "/api")
	endpoints.SetUIBackend("http", *host, port, "/ui/api")

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

	return &AWSEnv{
		client:    client,
		ec2Client: ec2Client,
		StackName: config.EnvName,
		meta: map[string]interface{}{
			"environment": "aws",
			"name":        config.EnvName,
		},
	}, nil
}
