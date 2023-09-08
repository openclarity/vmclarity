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

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/orchestrator/provider"
	"github.com/openclarity/vmclarity/pkg/shared/backendclient"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

const (
	discoveryInterval = 2 * time.Minute
)

type Discoverer struct {
	backendClient  *backendclient.BackendClient
	providerClient provider.Provider
}

func New(config Config) *Discoverer {
	return &Discoverer{
		backendClient:  config.Backend,
		providerClient: config.Provider,
	}
}

func (d *Discoverer) Start(ctx context.Context) {
	go func() {
		for {
			log.Debug("Discovering available assets")
			err := d.DiscoverAndCreateAssets(ctx)
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

func (d *Discoverer) DiscoverAndCreateAssets(ctx context.Context) error {
	discoveryTime := time.Now()

	discoverer := d.providerClient.DiscoverAssets(ctx)

	errs := []error{}
	failedPatchAssets := make(map[string]struct{})
	for assetType := range discoverer.Chan() {
		assetData := models.Asset{
			AssetInfo: utils.PointerTo(assetType),
			LastSeen:  &discoveryTime,
			FirstSeen: &discoveryTime,
		}
		_, err := d.backendClient.PostAsset(ctx, assetData)
		if err == nil {
			continue
		}

		var conflictError backendclient.AssetConflictError
		if !errors.As(err, &conflictError) {
			// If there is an error, and it's not a conflict telling
			// us that the asset already exists, then we need to
			// keep track of it and log it as a failure to
			// complete discovery. We don't fail instantly here
			// because discovering the assets is a heavy operation,
			// so we want to give the best chance to create all the
			// assets in the DB before failing.
			errs = append(errs, fmt.Errorf("failed to post asset: %w", err))
			continue
		}

		// As we got a conflict it means there is an existing asset
		// which matches the unique properties of this asset, in this
		// case we'll patch the just AssetInfo and FirstSeen instead.
		assetData.FirstSeen = nil
		err = d.backendClient.PatchAsset(ctx, assetData, *conflictError.ConflictingAsset.Id)
		if err != nil {
			failedPatchAssets[*conflictError.ConflictingAsset.Id] = struct{}{}
			errs = append(errs, fmt.Errorf("failed to patch asset: %w", err))
		}
	}

	if err := discoverer.Err(); err != nil {
		return fmt.Errorf("failed to discover assets: %w", err)
	}

	// Find all assets which are not already terminatedOn and were not
	// updated or created by this discovery run by comparing their
	// lastSeen time to this discovery's time stamp
	//
	// TODO(sambetts) when we add multiple providers/standalone support we
	// need to filter these assets by provider so that we don't find assets
	// which don't belong to us. We need to give the provider some kind of
	// identity in this case.
	assetResp, err := d.backendClient.GetAssets(ctx, models.GetAssetsParams{
		Filter: utils.PointerTo(fmt.Sprintf("terminatedOn eq null and (lastSeen eq null or lastSeen lt %s)", discoveryTime.Format(time.RFC3339))),
		Select: utils.PointerTo("id"),
	})
	if err != nil {
		return fmt.Errorf("failed to get existing Assets: %w", err)
	}

	// Patch all assets which were not found by this discovery as
	// terminated by setting terminatedOn.
	for _, asset := range *assetResp.Items {
		// Skip mark terminated if asset found but patch failed.
		if _, ok := failedPatchAssets[*asset.Id]; ok {
			continue
		}

		assetData := models.Asset{
			TerminatedOn: &discoveryTime,
		}

		err := d.backendClient.PatchAsset(ctx, assetData, *asset.Id)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to patch asset: %w", err))
		}
	}

	return errors.Join(errs...)
}
