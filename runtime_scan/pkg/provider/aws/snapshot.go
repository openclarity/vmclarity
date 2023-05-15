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

package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type Snapshot struct {
	ID       string
	Region   string
	Metadata provider.ScanMetadata
	VolumeID string

	ec2Client *ec2.Client
}

func (s *Snapshot) Copy(ctx context.Context, region string) (*Snapshot, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(logrus.Fields{
		"SnapshotID":     s.ID,
		"Operation":      "Copy",
		"TargetVolumeID": s.VolumeID,
	})

	if s.Region == region {
		logger.Debugf("Copying snapshot is skipped. SourceRegion=%s TargetRegion=%s", s.Region, region)
		return s, nil
	}

	options := func(options *ec2.Options) {
		options.Region = region
	}

	ec2TagsForSnapshot := EC2TagsFromScanMetadata(s.Metadata)
	ec2TagsForSnapshot = append(ec2TagsForSnapshot, ec2types.Tag{
		Key:   utils.PointerTo(EC2TagKeyTargetVolumeID),
		Value: utils.PointerTo(s.VolumeID),
	})
	ec2Filters := EC2FiltersFromEC2Tags(ec2TagsForSnapshot)

	describeParams := &ec2.DescribeSnapshotsInput{
		Filters: ec2Filters,
	}
	describeOut, err := s.ec2Client.DescribeSnapshots(ctx, describeParams, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target snapshots. TargetRegion=%s SourceSnapshot=%s SourceRegion=%s: %w",
			region, s.ID, s.Region, err)
	}

	if len(describeOut.Snapshots) > 1 {
		logger.Warnf("Multiple snapshots found for target volume: %d", len(describeOut.Snapshots))
	}

	for _, snap := range describeOut.Snapshots {
		switch snap.State {
		case ec2types.SnapshotStateError, ec2types.SnapshotStateRecoverable:
			// We want to recreate the snapshot if it is in error or recoverable state. Cleanup will take care of
			// removing these as well.
		case ec2types.SnapshotStateRecovering, ec2types.SnapshotStatePending, ec2types.SnapshotStateCompleted:
			fallthrough
		default:
			return &Snapshot{
				ec2Client: s.ec2Client,
				ID:        *snap.SnapshotId,
				Region:    s.Region,
				Metadata:  s.Metadata,
				VolumeID:  *snap.VolumeId,
			}, nil
		}
	}

	createParams := &ec2.CopySnapshotInput{
		SourceRegion:     &s.Region,
		SourceSnapshotId: &s.ID,
		Description:      utils.PointerTo(EC2SnapshotDescription),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSnapshot,
				Tags:         ec2TagsForSnapshot,
			},
		},
	}

	snap, err := s.ec2Client.CopySnapshot(ctx, createParams, options)
	if err != nil {
		return nil, fmt.Errorf("failed to copy snapshot with %s id from region %s to region %s: %w",
			s.ID, s.Region, region, err)
	}

	return &Snapshot{
		ec2Client: s.ec2Client,
		ID:        *snap.SnapshotId,
		Region:    region,
		Metadata:  s.Metadata,
		VolumeID:  s.VolumeID,
	}, nil
}

func (s *Snapshot) Delete(ctx context.Context) error {
	if s == nil {
		return nil
	}

	_, err := s.ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
		SnapshotId: &s.ID,
	}, func(options *ec2.Options) {
		options.Region = s.Region
	})
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %v", err)
	}

	return nil
}

func (s *Snapshot) WaitForReady(ctx context.Context, timeout time.Duration, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			ready, err := s.IsReady(ctx)
			if err != nil {
				return fmt.Errorf("failed to get volume snapshot state. SnapshotID=%s: %w", s.ID, err)
			}
			if ready {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait until volume snapshot is in ready state. SnapshotID=%s: %w", s.ID, ctx.Err())
		}
	}
}

func (s *Snapshot) IsReady(ctx context.Context) (bool, error) {
	var ready bool

	out, err := s.ec2Client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: []string{s.ID},
	}, func(options *ec2.Options) {
		options.Region = s.Region
	})
	if err != nil {
		return ready, fmt.Errorf("failed to describe snapshot. SnapshotID=%s: %w", s.ID, err)
	}

	if len(out.Snapshots) != 1 {
		return ready, fmt.Errorf("got unexcpected number of snapshots (%d). Excpecting 1. SnapshotID=%s",
			len(out.Snapshots), s.ID)
	}

	if out.Snapshots[0].State == ec2types.SnapshotStateCompleted {
		ready = true
	}

	return ready, nil
}

func (s *Snapshot) CreateVolume(ctx context.Context, az string) (*Volume, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(logrus.Fields{
		"SnapshotID":     s.ID,
		"Operation":      "CreateVolume",
		"TargetVolumeID": s.VolumeID,
	})

	options := func(options *ec2.Options) {
		options.Region = s.Region
	}

	ec2TagsForVolume := EC2TagsFromScanMetadata(s.Metadata)
	ec2TagsForVolume = append(ec2TagsForVolume, ec2types.Tag{
		Key:   utils.PointerTo(EC2TagKeyTargetVolumeID),
		Value: utils.PointerTo(s.VolumeID),
	})

	ec2Filters := EC2FiltersFromEC2Tags(ec2TagsForVolume)
	ec2Filters = append(ec2Filters, ec2types.Filter{
		Name:   utils.PointerTo("snapshot-id"),
		Values: []string{s.ID},
	})

	descParams := &ec2.DescribeVolumesInput{
		Filters: ec2Filters,
	}

	describeOut, err := s.ec2Client.DescribeVolumes(ctx, descParams, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch volume for snapshot. SnapshotID=%s: %w", s.ID, err)
	}

	if len(describeOut.Volumes) > 1 {
		logger.Warnf("Multiple volumes found for snapshot: %d", len(describeOut.Volumes))
	}

	for _, vol := range describeOut.Volumes {
		switch vol.State {
		case ec2types.VolumeStateDeleted, ec2types.VolumeStateDeleting:
			// We want to re-create the volume if it is in deleting or deleted state
		case ec2types.VolumeStateAvailable, ec2types.VolumeStateCreating, ec2types.VolumeStateInUse, ec2types.VolumeStateError:
			// We want to return the volume even if it is in error state to avoid creating new volumes where
			// which might get into error state as well.
			fallthrough
		default:
			return &Volume{
				ec2Client: s.ec2Client,
				ID:        *vol.VolumeId,
				Region:    s.Region,
				Metadata:  s.Metadata,
			}, nil
		}
	}

	createParams := &ec2.CreateVolumeInput{
		AvailabilityZone: &az,
		SnapshotId:       &s.ID,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeVolume,
				Tags:         ec2TagsForVolume,
			},
		},
		VolumeType: ec2types.VolumeTypeGp2,
	}
	out, err := s.ec2Client.CreateVolume(ctx, createParams, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create volume from snapshot. SnapshotID=%s: %w", s.ID, err)
	}

	return &Volume{
		ec2Client: s.ec2Client,
		ID:        *out.VolumeId,
		Region:    s.Region,
		Metadata:  s.Metadata,
	}, nil
}
