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

	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/grpc/proto"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

type VolumeImpl struct {
	providerClient provider_service.ProviderClient
	id             string
	location       string
}

func (v *VolumeImpl) GetID() string {
	return v.id
}

func (v *VolumeImpl) TakeSnapshot(ctx context.Context) (types.Snapshot, error) {
	params := provider_service.TakeVolumeSnapshotParams{
		VolumeID:       v.GetID(),
		VolumeLocation: v.location,
	}
	res, err := v.providerClient.TakeVolumeSnapshot(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to take volume snapshot. volumeID=%v: %v", v.GetID(), err)
	}

	return &SnapshotImpl{
		providerClient: v.providerClient,
		id:             res.GetSnapshot().GetId(),
		location:       res.GetSnapshot().GetLocation(),
	}, nil
}

func (v *VolumeImpl) WaitForReady(ctx context.Context) error {
	params := provider_service.WaitForVolumeReadyParams{
		VolumeID:       v.GetID(),
		VolumeLocation: v.location,
	}
	_, err := v.providerClient.WaitForVolumeReady(ctx, &params)

	return err
}

func (v *VolumeImpl) WaitForAttached(ctx context.Context) error {
	params := provider_service.WaitForVolumeAttachedParams{
		VolumeID:       v.GetID(),
		VolumeLocation: v.location,
	}
	_, err := v.providerClient.WaitForVolumeAttached(ctx, &params)

	return err
}

func (v *VolumeImpl) Delete(ctx context.Context) error {
	params := provider_service.DeleteVolumeParams{
		VolumeID:       v.GetID(),
		VolumeLocation: v.location,
	}

	_, err := v.providerClient.DeleteVolume(ctx, &params)

	return err
}
