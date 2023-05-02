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

	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/grpc/proto"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

type InstanceImpl struct {
	providerClient provider_service.ProviderClient
	id             string
	location       string
	image          string
	instanceType   string
	platform       string
	tags           []types.Tag
	launchTime     time.Time
}

func (i *InstanceImpl) GetID() string {
	return i.id
}

func (i *InstanceImpl) GetImage() string {
	return i.image
}

func (i *InstanceImpl) GetType() string {
	return i.instanceType
}

func (i *InstanceImpl) GetLaunchTime() time.Time {
	return i.launchTime
}

func (i *InstanceImpl) GetPlatform() string {
	return i.platform
}

func (i *InstanceImpl) GetTags() []types.Tag {
	return i.tags
}

func (i *InstanceImpl) GetLocation() string {
	return i.location
}

func (i *InstanceImpl) GetSecurityGroups() []string {
	// relevant only for AWS
	return nil
}
func (i *InstanceImpl) GetAvailabilityZone() string {
	// relevant only for AWS
	return ""
}

func (i *InstanceImpl) GetRootVolume(ctx context.Context) (types.Volume, error) {
	params := provider_service.GetInstanceRootVolumeParams{
		InstanceID:       i.GetID(),
		InstanceLocation: i.GetLocation(),
	}

	res, err := i.providerClient.GetInstanceRootVolume(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance root volume: %v", err)
	}

	ret := VolumeImpl{
		providerClient: i.providerClient,
		id:             res.Volume.GetId(),
		location:       res.Volume.GetLocation(),
	}

	return &ret, nil
}

func (i *InstanceImpl) WaitForReady(ctx context.Context) error {
	params := provider_service.WaitForInstanceReadyParams{
		InstanceID:       i.GetID(),
		InstanceLocation: i.GetLocation(),
	}

	_, err := i.providerClient.WaitForInstanceReady(ctx, &params)

	return err
}

func (i *InstanceImpl) Delete(ctx context.Context) error {
	params := provider_service.DeleteInstanceParams{
		InstanceID:       i.GetID(),
		InstanceLocation: i.GetLocation(),
	}
	_, err := i.providerClient.DeleteInstance(ctx, &params)

	return err
}

func (i *InstanceImpl) AttachVolume(ctx context.Context, volume types.Volume, deviceName string) error {
	params := provider_service.AttachVolumeToInstanceParams{
		InstanceID:       i.GetID(),
		InstanceLocation: i.GetLocation(),
		VolumeID:         volume.GetID(),
		DeviceName:       deviceName,
	}
	_, err := i.providerClient.AttachVolumeToInstance(ctx, &params)

	return err
}
