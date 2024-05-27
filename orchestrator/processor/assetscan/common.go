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

package assetscan

import (
	"context"
	"errors"
	"fmt"
	"time"

	apiclient "github.com/openclarity/vmclarity/api/client"
	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
)

func (asp *AssetScanProcessor) createOrUpdateDBFinding(ctx context.Context, info *apitypes.FindingInfo, assetScanID string, completedTime time.Time) (string, error) {
	// Create new finding
	finding := apitypes.Finding{
		FirstSeen: &completedTime,
		LastSeen:  &completedTime,
		LastSeenBy: &apitypes.AssetScanRelationship{
			Id: assetScanID,
		},
		FindingInfo: info,
	}

	fd, err := asp.client.PostFinding(ctx, finding)
	if err == nil {
		return *fd.Id, nil
	}

	var conflictError apiclient.FindingConflictError
	if !errors.As(err, &conflictError) {
		return "", fmt.Errorf("failed to create finding: %w", err)
	}

	var id string
	// Update existing finding if newer
	if conflictError.ConflictingFinding.LastSeen.Before(completedTime) {
		id = *conflictError.ConflictingFinding.Id
		finding := apitypes.Finding{
			LastSeen: &completedTime,
			LastSeenBy: &apitypes.AssetScanRelationship{
				Id: assetScanID,
			},
			FindingInfo: info,
		}

		err = asp.client.PatchFinding(ctx, id, finding)
		if err != nil {
			return id, fmt.Errorf("failed to patch finding: %w", err)
		}
	}

	return id, nil
}

func (asp *AssetScanProcessor) invalidateOlderFindingsByType(ctx context.Context, findingType string, assetID string, completedTime time.Time) error {
	// Invalidate any findings of this type for this asset where foundOn is
	// older than this asset scan, and has not already been invalidated by
	// an asset scan older than this asset scan.
	findingsToInvalidate, err := asp.client.GetFindings(ctx, apitypes.GetFindingsParams{
		Filter: to.Ptr(fmt.Sprintf(
			"findingInfo/objectType eq '%s' and asset/id eq '%s' and foundOn lt %s and (invalidatedOn gt %s or invalidatedOn eq null)",
			findingType, assetID, completedTime.Format(time.RFC3339), completedTime.Format(time.RFC3339))),
	})
	if err != nil {
		return fmt.Errorf("failed to query findings to invalidate: %w", err)
	}

	for _, finding := range *findingsToInvalidate.Items {
		finding.InvalidatedOn = &completedTime

		err := asp.client.PatchFinding(ctx, *finding.Id, finding)
		if err != nil {
			return fmt.Errorf("failed to update existing finding %s: %w", *finding.Id, err)
		}
	}

	return nil
}

func (asp *AssetScanProcessor) getActiveFindingsByType(ctx context.Context, findingType string, assetID string) (int, error) {
	filter := fmt.Sprintf("findingInfo/objectType eq '%s' and asset/id eq '%s' and invalidatedOn eq null",
		findingType, assetID)
	activeFindings, err := asp.client.GetFindings(ctx, apitypes.GetFindingsParams{
		Count:  to.Ptr(true),
		Filter: &filter,

		// select the smallest amount of data to return in items, we
		// only care about the count.
		Top:    to.Ptr(1),
		Select: to.Ptr("id"),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list all active findings: %w", err)
	}
	return *activeFindings.Count, nil
}
