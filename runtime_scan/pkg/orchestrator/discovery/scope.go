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

package discovery

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
)

const (
	discoveryInterval = 2 * time.Minute
)

type ScopeDiscoverer struct {
	backendClient  *client.ClientWithResponses
	providerClient provider.Client
}

func CreateScopeDiscoverer(backendClient *client.ClientWithResponses, providerClient provider.Client) *ScopeDiscoverer {
	return &ScopeDiscoverer{
		backendClient:  backendClient,
		providerClient: providerClient,
	}
}

func (sd *ScopeDiscoverer) setScopes(ctx context.Context, scopes *models.Scopes) (*models.Scopes, error) {
	resp, err := sd.backendClient.PutDiscoveryScopesWithResponse(ctx, *scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to put scopes: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scopes: empty body")
		}
		return resp.JSON200, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to set scopes. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to set scopes. status code=%v", resp.StatusCode())
	}
}

func (sd *ScopeDiscoverer) Start(ctx context.Context) {
	go func() {
		for {
			log.Debug("Discovering available scopes")
			// nolint:contextcheck
			scopes, err := sd.providerClient.DiscoverScopes(ctx)
			if err != nil {
				log.Warnf("Failed to discover scopes: %v", err)
			} else {
				_, err := sd.setScopes(ctx, scopes)
				if err != nil {
					log.Warnf("Failed to set scopes: %v", err)
				}
			}
			select {
			case <-time.After(discoveryInterval):
				log.Debug("Discovery interval elapsed")
			case <-ctx.Done():
				log.Infof("Stop watching scan configs.")
				return
			}
		}
	}()
}
