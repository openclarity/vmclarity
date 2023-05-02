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

package external

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/grpc/proto"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/config/external_provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

type Client struct {
	providerClient provider_service.ProviderClient
	config         *external_provider.Config
	conn           *grpc.ClientConn
}

func Create(_ context.Context, config *external_provider.Config) (*Client, error) {
	var opts []grpc.DialOption
	// TODO secure connections
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := Client{}

	conn, err := grpc.Dial(config.ProviderPluginAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial grpc. address=%v: %v", config.ProviderPluginAddress, err)
	}
	client.conn = conn
	client.providerClient = provider_service.NewProviderClient(conn)

	return &client, nil
}

func (c *Client) DiscoverScopes(ctx context.Context) (*models.Scopes, error) {
	res, err := c.providerClient.DiscoverScopes(ctx, &provider_service.DiscoverScopesParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to discover scopes: %v", err)
	}

	scopes := res.GetScopes()

	resScopes := models.Scopes{
		ScopeInfo: &models.ScopeType{},
	}

	err = resScopes.ScopeInfo.UnmarshalJSON([]byte(scopes))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal scopes: %v", err)
	}

	return &resScopes, nil
}

func (c *Client) DiscoverInstances(ctx context.Context, scanScope *models.ScanScopeType) ([]types.Instance, error) {
	var ret []types.Instance

	scopesB, err := scanScope.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshak scan scope: %v", err)
	}

	res, err := c.providerClient.DiscoverInstances(ctx, &provider_service.DiscoverInstancesParams{ScanScopes: string(scopesB)})
	if err != nil {
		return nil, fmt.Errorf("failed to discover instances: %v", err)
	}

	instances := res.GetInstances()

	for _, instance := range instances {
		ret = append(ret, &InstanceImpl{
			providerClient: c.providerClient,
			id:             instance.GetId(),
			location:       instance.GetLocation(),
			image:          instance.GetImage(),
			instanceType:   instance.GetInstanceType(),
			platform:       instance.GetPlatform(),
			tags:           toTypesTags(instance.GetTags()),
			launchTime:     instance.LaunchTime.AsTime(),
		})
	}

	return ret, nil
}

// TODO I don't think we should tell the user where to run the instance, we should remove region from the interafce args
func (c *Client) RunScanningJob(ctx context.Context, region, id string, config provider.ScanningJobConfig) (types.Instance, error) {
	params := provider_service.RunScanningJobParams{
		Location: region,
		Id:       id,
		Config: &provider_service.ScanningJobConfig{
			ScannerImage:     config.ScannerImage,
			ScannerCLIConfig: config.ScannerCLIConfig,
			VmClarityAddress: config.VMClarityAddress,
			ScanResultID:     config.ScanResultID,
			KeyPairName:      config.KeyPairName,
		},
	}
	if config.ScannerInstanceCreationConfig != nil {
		params.Config.ScannerInstanceCreationConfig = &provider_service.ScannerInstanceCreationConfig{
			UseSpotInstances: config.ScannerInstanceCreationConfig.UseSpotInstances}
		if config.ScannerInstanceCreationConfig.MaxPrice != nil {
			params.Config.ScannerInstanceCreationConfig.MaxPrice = *config.ScannerInstanceCreationConfig.MaxPrice
		}
		if config.ScannerInstanceCreationConfig.RetryMaxAttempts != nil {
			params.Config.ScannerInstanceCreationConfig.RetryMaxAttempts = int32(*config.ScannerInstanceCreationConfig.RetryMaxAttempts)
		}
	}

	res, err := c.providerClient.RunScanningJob(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to run scanning job: %v", err)
	}

	instance := res.GetInstance()

	return &InstanceImpl{
		providerClient: c.providerClient,
		id:             instance.GetId(),
		location:       instance.GetLocation(),
		image:          instance.GetImage(),
		instanceType:   instance.GetInstanceType(),
		platform:       instance.GetPlatform(),
		tags:           toTypesTags(instance.GetTags()),
		launchTime:     instance.GetLaunchTime().AsTime(),
	}, nil
}

func toTypesTags(tags []*provider_service.Tag) []types.Tag {
	var ret []types.Tag
	for _, tag := range tags {
		ret = append(ret, types.Tag{
			Key: tag.Key,
			Val: tag.Val,
		})
	}
	return ret
}
