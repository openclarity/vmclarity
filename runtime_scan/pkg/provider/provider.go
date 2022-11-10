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
