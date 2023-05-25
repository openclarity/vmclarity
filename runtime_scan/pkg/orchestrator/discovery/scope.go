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
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

const (
	discoveryInterval = 2 * time.Minute
)

type ScopeDiscoverer struct {
	backendClient  *backendclient.BackendClient
	providerClient provider.Client
}

func CreateScopeDiscoverer(backendClient *backendclient.BackendClient, providerClient provider.Client) *ScopeDiscoverer {
	return &ScopeDiscoverer{
		backendClient:  backendClient,
		providerClient: providerClient,
	}
}

func (sd *ScopeDiscoverer) Start(ctx context.Context) {
	go func() {
		for {
			log.Debug("Discovering available assets")
			err := sd.DiscoverAndCreateAssets(ctx)
			if err != nil {
				log.Warnf("Failed to discover assets: %v", err)
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

func (sd *ScopeDiscoverer) DiscoverAndCreateAssets(ctx context.Context) error {
	assets, err := sd.providerClient.DiscoverAssets(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover assets from provider: %w", err)
	}

	errs := []error{}
	for _, asset := range assets {
		_, err := sd.backendClient.PostTarget(ctx, asset)
		if err != nil && !errors.As(err, &backendclient.TargetConflictError{}) {
			// If there is an error, and its not a conflict telling
			// us that the target already exists, then we need to
			// keept track of it and log it as a failure to
			// complete discovery. We don't fail instantly here
			// because discovering the assets is a heavy operation
			// so we want to give the best chance to create all the
			// assets in the DB before failing.
			errs = append(errs, fmt.Errorf("failed to post target: %v", err))
		}
	}

	// TODO(sambetts) Compare the assets list to the assets in the DB and
	// mark missing assets as terminated.

	return errors.Join(errs...)
}
