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

type SnapshotImpl struct {
	providerClient provider_service.ProviderClient
	id             string
	location       string
}

func (s *SnapshotImpl) GetID() string {
	return s.id
}

func (s *SnapshotImpl) GetLocation() string {
	return s.location
}

func (s *SnapshotImpl) Copy(ctx context.Context, dstRegion string) (types.Snapshot, error) {
	params := provider_service.CopySnapshotParams{
		SnapshotID:       s.GetID(),
		SnapshotLocation: s.GetLocation(),
		DestLocation:     dstRegion,
	}
	res, err := s.providerClient.CopySnapshot(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to copy snapshot. snapshotID=%v: %v", s.GetID(), err)
	}

	ret := SnapshotImpl{
		providerClient: s.providerClient,
		id:             res.GetSnapshot().GetId(),
		location:       res.GetSnapshot().GetLocation(),
	}

	return &ret, nil
}

func (s *SnapshotImpl) Delete(ctx context.Context) error {
	params := provider_service.DeleteSnapshotParams{
		SnapshotID:       s.GetID(),
		SnapshotLocation: s.GetLocation(),
	}
	_, err := s.providerClient.DeleteSnapshot(ctx, &params)

	return err
}

func (s *SnapshotImpl) WaitForReady(ctx context.Context) error {
	params := provider_service.WaitForSnapshotReadyParams{
		SnapshotID:       s.GetID(),
		SnapshotLocation: s.GetLocation(),
	}
	_, err := s.providerClient.WaitForSnapshotReady(ctx, &params)

	return err
}

// TODO availabilityZone is aws specific. we need it because in aws, volume should be created in
// the same availabilityZone as the instance that it is about to be attached to.
func (s *SnapshotImpl) CreateVolume(ctx context.Context, availabilityZone string) (types.Volume, error) {
	params := provider_service.CreateVolumeFromSnapshotParams{
		SnapshotID:       s.GetID(),
		SnapshotLocation: s.GetLocation(),
		AvailabilityZone: availabilityZone,
	}
	res, err := s.providerClient.CreateVolumeFromSnapshot(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume from snapshot. snapshotID=%v: %v", s.GetID(), err)
	}

	return &VolumeImpl{
		providerClient: s.providerClient,
		id:             res.GetVolume().GetId(),
		location:       res.GetVolume().GetLocation(),
	}, nil
}
