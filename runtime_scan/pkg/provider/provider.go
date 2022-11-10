package provider

import types2 "github.com/openclarity/vmclarity/runtime_scan/pkg/types"

type Client interface {
	// list VM instance ids in the account according to the filters.
	Discover(filters *types2.ScanScope) ([]types2.Instance, error)
	// create a snapshot of a volume, and return the snapshot id.
	CreateSnapshot(types2.Volume) (types2.Snapshot, error)
	// copy the snapshot from src region to dest region and return the snapshot id.
	CopySnapshot(snapshot types2.Snapshot, dstRegion string) (types2.Snapshot, error)
	//
	WaitForSnapshotReady(snapshot types2.Snapshot) error
	// get the instance root volume
	GetInstanceRootVolume(instance types2.Instance) (types2.Volume, error)
	// attach a volume to an instance
	//AttachVolume(volumeID, instanceID, region string) (string, error)
	//
	LaunchInstance(ami, deviceName string, snapshot types2.Snapshot) (types2.Instance, error)
	//
	DeleteInstance(types2.Instance) error
	//
	DeleteSnapshot(types2.Snapshot) error
}
