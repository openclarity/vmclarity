package scanner

import (
	"context"
	"time"

	"cloud.google.com/go/compute/apiv1/computepb"

	"github.com/openclarity/vmclarity/provider"
	"github.com/openclarity/vmclarity/provider/v2/gcp/common"
)

var (
	SnapshotCreateEstimateProvisionTime = 2 * time.Minute
	SnapshotDeleteEstimateTime          = 2 * time.Minute
)

func snapshotNameFromJobConfig(config *provider.ScanJobConfig) string {
	return "snapshot-" + config.AssetScanID
}

func (s *Scanner) ensureSnapshotFromAttachedDisk(ctx context.Context, config *provider.ScanJobConfig, disk *computepb.AttachedDisk) (*computepb.Snapshot, error) {
	snapshotName := snapshotNameFromJobConfig(config)

	snapshotRes, err := s.SnapshotsClient.Get(ctx, &computepb.GetSnapshotRequest{
		Project:  s.ProjectID,
		Snapshot: snapshotName,
	})
	if err == nil {
		if *snapshotRes.Status != ProvisioningStateReady {
			return snapshotRes, provider.RetryableErrorf(SnapshotCreateEstimateProvisionTime, "snapshot is not ready yet. status: %v", *snapshotRes.Status)
		}

		// Everything is good, the snapshot exists and is provisioned successfully
		return snapshotRes, nil
	}

	notFound, err := common.HandleGcpRequestError(err, "getting snapshot %s", snapshotName)
	if !notFound {
		return nil, err
	}

	// Snapshot not found, Create the snapshot
	req := &computepb.InsertSnapshotRequest{
		Project: s.ProjectID,
		SnapshotResource: &computepb.Snapshot{
			Name:       &snapshotName,
			SourceDisk: disk.Source,
		},
	}

	_, err = s.SnapshotsClient.Insert(ctx, req)
	if err != nil {
		_, err := common.HandleGcpRequestError(err, "create snapshot %s", snapshotName)
		return nil, err
	}

	return &computepb.Snapshot{}, provider.RetryableErrorf(SnapshotCreateEstimateProvisionTime, "snapshot creating")
}

func (s *Scanner) ensureSnapshotDeleted(ctx context.Context, config *provider.ScanJobConfig) error {
	snapshotName := snapshotNameFromJobConfig(config)

	return common.EnsureDeleted(
		"snapshot",
		func() error {
			_, err := s.SnapshotsClient.Get(ctx, &computepb.GetSnapshotRequest{
				Project:  s.ProjectID,
				Snapshot: snapshotName,
			})
			return err // nolint: wrapcheck
		},
		func() error {
			_, err := s.SnapshotsClient.Delete(ctx, &computepb.DeleteSnapshotRequest{
				Project:  s.ProjectID,
				Snapshot: snapshotName,
			})
			return err // nolint: wrapcheck
		},
		SnapshotDeleteEstimateTime,
	)
}
