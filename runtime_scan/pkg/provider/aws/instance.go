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

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type Instance struct {
	ID               string
	Region           string
	VpcID            string
	SecurityGroups   []models.SecurityGroup
	AvailabilityZone string
	Image            string
	InstanceType     string
	Platform         string
	Tags             []models.Tag
	LaunchTime       time.Time
	RootDeviceName   string
	Volumes          []Volume

	Metadata provider.ScanMetadata

	ec2Client *ec2.Client
}

func (i *Instance) GetLocation() string {
	return i.Region + "/" + i.VpcID
}

func (i *Instance) GetRootVolume() *Volume {
	var root Volume
	for idx, vol := range i.Volumes {
		if idx == 0 {
			root = vol
		}
		if vol.BlockDeviceName == i.RootDeviceName {
			root = vol
			break
		}
	}

	return &root
}

func (i *Instance) WaitForReady(ctx context.Context, timeout time.Duration, interval time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			ready, err := i.IsReady(ctx)
			if err != nil {
				return fmt.Errorf("failed to get VM instance state. InstanceID=%s: %w", i.ID, err)
			}
			if ready {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait until VM Instance is in ready state. InstanceID=%s: %w", i.ID, ctx.Err())
		}
	}
}

func (i *Instance) IsReady(ctx context.Context) (bool, error) {
	var ready bool

	out, err := i.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{i.ID},
	}, func(options *ec2.Options) {
		options.Region = i.Region
	})
	if err != nil {
		return ready, fmt.Errorf("failed to get VM instance. InstanceID=%s: %w", i.ID, err)
	}

	state := getInstanceState(out, i.ID)
	if state == ec2types.InstanceStateNameRunning {
		ready = true
	}

	return ready, nil
}

func (i *Instance) Delete(ctx context.Context) error {
	if i == nil {
		return nil
	}

	_, err := i.ec2Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{i.ID},
	}, func(options *ec2.Options) {
		options.Region = i.Region
	})
	if err != nil {
		return fmt.Errorf("failed to terminate instances: %v", err)
	}

	return nil
}

func (i *Instance) AttachVolume(ctx context.Context, volume *Volume, deviceName string) error {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithFields(logrus.Fields{
		"InstanceID": i.ID,
		"Operation":  "AttachVolume",
		"VolumeID":   volume.ID,
	})

	options := func(options *ec2.Options) {
		options.Region = volume.Region
	}

	descParams := &ec2.DescribeVolumesInput{
		VolumeIds: []string{volume.ID},
	}

	describeOut, err := i.ec2Client.DescribeVolumes(ctx, descParams, options)
	if err != nil {
		return fmt.Errorf("failed to fetch volume. VolumeID=%s: %w", volume.ID, err)
	}

	logger.Tracef("Found %d volumes", len(describeOut.Volumes))

	var volumeAttached bool
	for _, vol := range describeOut.Volumes {
		logger.WithFields(logrus.Fields{
			"VolumeState": vol.State,
		}).Trace("Found volume")

		switch vol.State {
		case ec2types.VolumeStateInUse:
			for _, attachment := range vol.Attachments {
				if *attachment.VolumeId == volume.ID && *attachment.InstanceId == i.ID {
					logger.Trace("Volume is already attached to the instance")
					volumeAttached = true
					break
				}
			}
		case ec2types.VolumeStateAvailable:
			logger.Trace("Attaching volume to instance")
			_, err := i.ec2Client.AttachVolume(ctx, &ec2.AttachVolumeInput{
				Device:     utils.PointerTo(deviceName),
				InstanceId: utils.PointerTo(i.ID),
				VolumeId:   utils.PointerTo(volume.ID),
			}, func(options *ec2.Options) {
				options.Region = i.Region
			})
			if err != nil {
				return fmt.Errorf("failed to attach volume: %v", err)
			}
			volumeAttached = true
			break
		default:
			continue
		}
	}

	if !volumeAttached {
		return fmt.Errorf("failed to attach volume due to its state. VolumeID=%s", volume.ID)
	}

	return nil
}
