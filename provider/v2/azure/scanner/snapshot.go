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

package scanner

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"

	"github.com/openclarity/vmclarity/core/to"
	"github.com/openclarity/vmclarity/provider"
	"github.com/openclarity/vmclarity/provider/v2/azure/common"
)

var (
	SnapshotCreateEstimateProvisionTime = 2 * time.Minute
)

func snapshotNameFromJobConfig(config *provider.ScanJobConfig) string {
	return "snapshot-" + config.AssetScanID
}

func (s *Scanner) ensureSnapshotForVMRootVolume(ctx context.Context, config *provider.ScanJobConfig, vm armcompute.VirtualMachine) (armcompute.Snapshot, error) {
	snapshotName := snapshotNameFromJobConfig(config)

	snapshotRes, err := s.SnapshotsClient.Get(ctx, s.Config.ScannerResourceGroup, snapshotName, nil)
	if err == nil {
		if *snapshotRes.Properties.ProvisioningState != provisioningStateSucceeded {
			return snapshotRes.Snapshot, provider.RetryableErrorf(SnapshotCreateEstimateProvisionTime, "snapshot is not ready yet")
		}

		// Everything is good, the snapshot exists and is provisioned successfully
		return snapshotRes.Snapshot, nil
	}

	notFound, err := common.HandleAzureRequestError(err, "getting snapshot %s", snapshotName)
	if !notFound {
		return armcompute.Snapshot{}, err
	}

	_, err = s.SnapshotsClient.BeginCreateOrUpdate(ctx, s.Config.ScannerResourceGroup, snapshotName, armcompute.Snapshot{
		Location: vm.Location,
		Properties: &armcompute.SnapshotProperties{
			CreationData: &armcompute.CreationData{
				CreateOption:     to.Ptr(armcompute.DiskCreateOptionCopy),
				SourceResourceID: vm.Properties.StorageProfile.OSDisk.ManagedDisk.ID,
			},
		},
	}, nil)
	if err != nil {
		_, err := common.HandleAzureRequestError(err, "creating snapshot %s", snapshotName)
		return armcompute.Snapshot{}, err
	}

	return armcompute.Snapshot{}, provider.RetryableErrorf(SnapshotCreateEstimateProvisionTime, "snapshot creating")
}
