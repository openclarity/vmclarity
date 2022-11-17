// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package provider

import "context"

type Tag struct {
	Key string
	Val string
}

type SecurityGroup struct {
	ID string
}

type Job struct {
	Instance    Instance
	SrcSnapshot Snapshot
	DstSnapshot Snapshot
}

type JobConfig struct {
	InstanceToScan Instance
	Region         string
	ImageID        string
	DeviceName     string
	SubnetID       string
}

type VPC struct {
	ID             string
	SecurityGroups []SecurityGroup
}

type Region struct {
	ID string
	VPCs []VPC
}

type Instance interface {
	GetID() string
	GetRootVolume(ctx context.Context) (Volume, error)
	WaitForReady(ctx context.Context) error
	Delete(ctx context.Context) error
}

type Volume interface {
	TakeSnapshot(ctx context.Context) (Snapshot, error)
}

type Snapshot interface {
	GetID() string
	GetRegion() string
	Copy(ctx context.Context, dstRegion string) (Snapshot, error)
	Delete(ctx context.Context) error
	WaitForReady(ctx context.Context) error
}
