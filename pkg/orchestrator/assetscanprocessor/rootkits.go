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

package assetscanprocessor

import (
	"context"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/shared/findingkey"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
	"github.com/openclarity/vmclarity/utils/log"
)

func (asp *AssetScanProcessor) getExistingRootkitFindingsForScan(ctx context.Context, assetScan models.AssetScan) (map[findingkey.RootkitKey]string, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	existingMap := map[findingkey.RootkitKey]string{}

	existingFilter := fmt.Sprintf("findingInfo/objectType eq 'Rootkit' and foundBy/id eq '%s'", *assetScan.Id)
	existingFindings, err := asp.client.GetFindings(ctx, models.GetFindingsParams{
		Filter: &existingFilter,
		Select: utils.PointerTo("id,findingInfo/rootkitName,findingInfo/rootkitType,findingInfo/path"),
	})
	if err != nil {
		return existingMap, fmt.Errorf("failed to query for findings: %w", err)
	}

	for _, finding := range *existingFindings.Items {
		info, err := (*finding.FindingInfo).AsRootkitFindingInfo()
		if err != nil {
			return existingMap, fmt.Errorf("unable to get rootkit finding info: %w", err)
		}

		key := findingkey.GenerateRootkitKey(info)
		if _, ok := existingMap[key]; ok {
			return existingMap, fmt.Errorf("found multiple matching existing findings for rootkit %v", key)
		}
		existingMap[key] = *finding.Id
	}

	logger.Infof("Found %d existing rootkit findings for this scan", len(existingMap))
	logger.Debugf("Existing rootkit map: %v", existingMap)

	return existingMap, nil
}

// nolint:cyclop
func (asp *AssetScanProcessor) reconcileResultRootkitsToFindings(ctx context.Context, assetScan models.AssetScan) error {
	completedTime := assetScan.Status.LastTransitionTime

	newerFound, newerTime, err := asp.newerExistingFindingTime(ctx, assetScan.Asset.Id, "Rootkit", completedTime)
	if err != nil {
		return fmt.Errorf("failed to check for newer existing rootkit findings: %w", err)
	}

	// Build a map of existing findings for this scan to prevent us
	// recreating existings ones as we might be re-reconciling the same
	// asset scan because of downtime or a previous failure.
	existingMap, err := asp.getExistingRootkitFindingsForScan(ctx, assetScan)
	if err != nil {
		return fmt.Errorf("failed to check existing rootkit findings: %w", err)
	}

	if assetScan.Rootkits != nil && assetScan.Rootkits.Rootkits != nil {
		// Create new or update existing findings all the rootkits found by the
		// scan.
		for _, item := range *assetScan.Rootkits.Rootkits {
			itemFindingInfo := models.RootkitFindingInfo{
				Message:     item.Message,
				RootkitName: item.RootkitName,
				RootkitType: item.RootkitType,
			}

			findingInfo := models.Finding_FindingInfo{}
			err = findingInfo.FromRootkitFindingInfo(itemFindingInfo)
			if err != nil {
				return fmt.Errorf("unable to convert RootkitFindingInfo into FindingInfo: %w", err)
			}

			finding := models.Finding{
				Asset: &models.AssetRelationship{
					Id: assetScan.Asset.Id,
				},
				FoundBy: &models.AssetScanRelationship{
					Id: *assetScan.Id,
				},
				FoundOn:     &assetScan.Status.LastTransitionTime,
				FindingInfo: &findingInfo,
			}

			// Set InvalidatedOn time to the FoundOn time of the oldest
			// finding, found after this asset scan.
			if newerFound {
				finding.InvalidatedOn = &newerTime
			}

			key := findingkey.GenerateRootkitKey(itemFindingInfo)
			if id, ok := existingMap[key]; ok {
				err = asp.client.PatchFinding(ctx, id, finding)
				if err != nil {
					return fmt.Errorf("failed to create finding: %w", err)
				}
			} else {
				_, err = asp.client.PostFinding(ctx, finding)
				if err != nil {
					return fmt.Errorf("failed to create finding: %w", err)
				}
			}
		}
	}

	// Invalidate any findings of this type for this asset where foundOn is
	// older than this asset scan, and has not already been invalidated by
	// an asset scan older than this asset scan.
	err = asp.invalidateOlderFindingsByType(ctx, "Rootkit", assetScan.Asset.Id, completedTime)
	if err != nil {
		return fmt.Errorf("failed to invalidate older rootkit finding: %w", err)
	}

	// Get all findings which aren't invalidated, and then update the asset's summary
	asset, err := asp.client.GetAsset(ctx, assetScan.Asset.Id, models.GetAssetsAssetIDParams{})
	if err != nil {
		return fmt.Errorf("failed to get asset %s: %w", assetScan.Asset.Id, err)
	}
	if asset.Summary == nil {
		asset.Summary = &models.ScanFindingsSummary{}
	}

	totalRootkits, err := asp.getActiveFindingsByType(ctx, "Rootkit", assetScan.Asset.Id)
	if err != nil {
		return fmt.Errorf("failed to get active rootkit findings: %w", err)
	}
	asset.Summary.TotalRootkits = &totalRootkits

	err = asp.client.PatchAsset(ctx, asset, assetScan.Asset.Id)
	if err != nil {
		return fmt.Errorf("failed to patch asset %s: %w", assetScan.Asset.Id, err)
	}

	return nil
}
