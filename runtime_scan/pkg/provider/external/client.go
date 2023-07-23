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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/external/proto"
)

type Client struct {
	providerClient provider_service.ProviderClient
	config         *Config
	conn           *grpc.ClientConn
}

func New(_ context.Context) (*Client, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	client := Client{
		config: config,
	}

	var opts []grpc.DialOption
	// TODO secure connections
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(config.ProviderPluginAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial grpc. address=%v: %v", config.ProviderPluginAddress, err)
	}
	client.conn = conn
	client.providerClient = provider_service.NewProviderClient(conn)

	return &client, nil
}

func (c Client) Kind() models.CloudProvider {
	return models.External
}

func (c *Client) DiscoverAssets(ctx context.Context) ([]models.AssetType, error) {
	var ret []models.AssetType

	res, err := c.providerClient.DiscoverAssets(ctx, &provider_service.DiscoverAssetsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to discover assets: %v", err)
	}

	assets := res.GetAssets()
	for _, asset := range assets {
		modelsAsset, err := convertAssetToModels(asset)
		if err != nil {
			return nil, fmt.Errorf("failed to convert asset to models asset: %v", err)
		}

		ret = append(ret, *modelsAsset.AssetInfo)
	}

	return ret, nil
}

func (c *Client) RunAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	scanJobConfig, err := convertScanJobConfig(config)
	if err != nil {
		return fmt.Errorf("failed to convert scan job config: %v", err)
	}

	res, err := c.providerClient.RunAssetScan(ctx, &provider_service.RunAssetScanParams{
		ScanJobConfig: scanJobConfig,
	})
	if err != nil {
		return provider.FatalErrorf("failed to run asset scan: %v", err)
	}

	if res.Err == nil {
		return provider.FatalErrorf("failed to run asset scan: an error type must be set")
	}

	switch res.Err.ErrorType.(type) {
	case *provider_service.Error_ErrNone:
		return nil
	case *provider_service.Error_ErrRetry:
		retryableErr := res.GetErr().GetErrRetry()
		return provider.RetryableErrorf(time.Second*time.Duration(retryableErr.After), retryableErr.Err)
	case *provider_service.Error_ErrFatal:
		fatalErr := res.GetErr().GetErrFatal()
		return provider.FatalErrorf("failed to run asset scan: %v", fatalErr.Err)
	default:
		return provider.FatalErrorf("failed to run asset scan: error type is not supported: %t", res.Err.GetErrorType())
	}
}

func (c *Client) RemoveAssetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	scanJobConfig, err := convertScanJobConfig(config)
	if err != nil {
		return fmt.Errorf("failed to convert scan job config: %v", err)
	}

	_, err = c.providerClient.RemoveAssetScan(ctx, &provider_service.RemoveAssetScanParams{
		ScanJobConfig: scanJobConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to remove asset scan: %v", err)
	}
	return nil
}
