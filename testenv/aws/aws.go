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
	"fmt"
	"io"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	"github.com/openclarity/vmclarity/testenv/types"
)

type ContextKeyType string

const AWSClientContextKey ContextKeyType = "AWSClient"

// AWS Environment.
type AWSEnv struct {
	client *cloudformation.Client
	meta   map[string]interface{}
}

// Setup AWS test environment from cloud formation template.
func (e *AWSEnv) SetUp(ctx context.Context) error {
	// Create a new CloudFormation stack from template
	// _, err := e.client.CreateStack()

	panic("not implemented")
}

func (e *AWSEnv) TearDown(ctx context.Context) error {
	// Delete the CloudFormation stack
	// _, err := e.client.DeleteStack()

	panic("not implemented")
}

func (e *AWSEnv) ServicesReady(ctx context.Context) (bool, error) {
	// Check if all services are ready
	// _, err := e.client.DescribeStacks()

	panic("not implemented")
}

func (e *AWSEnv) ServiceLogs(ctx context.Context, services []string, startTime time.Time, stdout, stderr io.Writer) error {
	panic("not implemented")
}

func (e *AWSEnv) Services(ctx context.Context) (types.Services, error) {
	// Get list of services from CloudFormation stack
	// _, err := e.client.ListStackResources()

	panic("not implemented")
}

func (e *AWSEnv) Endpoints(_ context.Context) (*types.Endpoints, error) {
	// Get IP of VMClarity Server from CloudFormation stack
	host := "localhost"
	// Get port of VMClarity Server
	port := "8080"

	endpoints := new(types.Endpoints)
	endpoints.SetAPI("http", host, port, "/api")
	endpoints.SetUIBackend("http", host, port, "/ui/api")

	return endpoints, nil
}

func (e *AWSEnv) Context(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, AWSClientContextKey, e.client), nil
}

func New(config *Config, opts ...ConfigOptFn) (*AWSEnv, error) {
	if err := applyConfigWithOpts(config, opts...); err != nil {
		return nil, fmt.Errorf("failed to apply config options: %w", err)
	}

	// Load AWS configuration
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	// Create AWS CloudFormation client
	client := cloudformation.NewFromConfig(cfg)

	return &AWSEnv{
		client: client,
		meta: map[string]interface{}{
			"environment": "aws",
			"name":        config.EnvName,
		},
	}, nil
}
