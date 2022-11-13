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

import types2 "github.com/openclarity/vmclarity/runtime_scan/pkg/types"

type Client interface {
	// Discover - list VM instances in the account according to the filters.
	Discover(filters *types2.ScanScope) ([]types2.Instance, error)
	// Create a snapshot of a volume.
	// return the newly created snapshot.
	CreateSnapshot(volume types2.Volume) (types2.Snapshot, error)
	// CopySnapshot - copy the snapshot from src region to dest region.
	// return the newly created (copy) snapshot.
	CopySnapshot(snapshot types2.Snapshot, dstRegion string) (types2.Snapshot, error)
	// WaitForSnapshotReady - wait until snapshot state is 'completed'.
	WaitForSnapshotReady(snapshot types2.Snapshot) error
	// WaitForSnapshotReady - wait until instance state is 'running'.
	WaitForInstanceReady(instance types2.Instance) error
	// GetInstanceRootVolume - get the instance's root volume.
	GetInstanceRootVolume(instance types2.Instance) (types2.Volume, error)
	// LaunchInstance - launch an instance. the snapshot will be attached to the instance at launch.
	// return the launched instance
	LaunchInstance(ami, deviceName string, snapshot types2.Snapshot) (types2.Instance, error)
	// DeleteInstance - delete an instance.
	DeleteInstance(instance types2.Instance) error
	// DeleteSnapshot - delete a snapshot.
	DeleteSnapshot(snapshot types2.Snapshot) error
}
