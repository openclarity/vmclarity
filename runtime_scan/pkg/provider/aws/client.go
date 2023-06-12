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
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	awstype "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider/cloudinit"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type Client struct {
	ec2Client *ec2.Client
	awsConfig *Config
}

func New(ctx context.Context, config *Config) (*Client, error) {
	awsClient := Client{
		awsConfig: config,
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	// nolint:contextcheck
	awsClient.ec2Client = ec2.NewFromConfig(cfg)

	return &awsClient, nil
}

func (c Client) Kind() models.CloudProvider {
	return models.AWS
}

func (c *Client) DiscoverScopes(ctx context.Context) (*models.Scopes, error) {
	regions, err := c.ListAllRegions(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to list all regions: %w", err)
	}

	scopes := models.ScopeType{}
	err = scopes.FromAwsAccountScope(models.AwsAccountScope{
		Regions: convertToAPIRegions(regions),
	})
	if err != nil {
		return nil, FatalError{
			Err: fmt.Errorf("failed to cast AwsAccountScope to ScopeType: %w", err),
		}
	}

	return &models.Scopes{
		ScopeInfo: &scopes,
	}, nil
}

// nolint:cyclop
func (c *Client) DiscoverTargets(ctx context.Context, scanScope *models.ScanScopeType) ([]models.TargetType, error) {
	var filters []ec2types.Filter

	awsScanScope, err := scanScope.AsAwsScanScope()
	if err != nil {
		return nil, FatalError{
			Err: fmt.Errorf("failed to cast ScanScopeType to AwsAccountScope: %w", err),
		}
	}

	scope := convertFromAPIScanScope(&awsScanScope)

	regions, err := c.getRegionsToScan(ctx, scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions to scan: %w", err)
	}
	if len(regions) == 0 {
		return nil, FatalError{
			Err: errors.New("no regions to scan"),
		}
	}

	filters = append(filters, EC2FiltersFromTags(scope.TagSelector)...)

	instanceStates := []ec2types.InstanceStateName{ec2types.InstanceStateNameRunning}
	if scope.ScanStopped {
		instanceStates = append(instanceStates, ec2types.InstanceStateNameStopped)
	}
	filters = append(filters, EC2FiltersFromInstanceState(instanceStates...)...)

	targets := make([]models.TargetType, 0)
	for _, region := range regions {
		// if no vpcs, that mean that we don't need any vpc filters
		if len(region.VPCs) == 0 {
			instances, err := c.GetInstances(ctx, filters, scope.ExcludeTags, region.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %w", err)
			}

			for _, instance := range instances {
				target, err := getVMInfoFromInstance(instance)
				if err != nil {
					return nil, FatalError{
						Err: fmt.Errorf("failed convert EC2 Instance to TargetType: %w", err),
					}
				}
				targets = append(targets, target)
			}

			continue
		}

		// need to do a per vpc call for DescribeInstances
		for _, vpc := range region.VPCs {
			vpcFilters := append(filters, createVPCFilters(vpc)...)

			instances, err := c.GetInstances(ctx, vpcFilters, scope.ExcludeTags, region.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %w", err)
			}
			for _, instance := range instances {
				target, err := getVMInfoFromInstance(instance)
				if err != nil {
					return nil, FatalError{
						Err: fmt.Errorf("failed convert EC2 Instance to TargetType: %w", err),
					}
				}
				targets = append(targets, target)
			}
		}
	}

	return targets, nil
}

// nolint:nilnil
func (c *Client) getInstanceWithID(ctx context.Context, id string, region string) (*ec2types.Instance, error) {
	out, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	}, func(options *ec2.Options) {
		options.Region = region
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch VM isntance. InstanceID=%s: %w", id, err)
	}

	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			ec2Instance := i

			if ec2Instance.InstanceId == nil || *ec2Instance.InstanceId != id {
				continue
			}

			return &ec2Instance, nil
		}
	}

	return nil, nil
}

// nolint:cyclop
func (c *Client) createInstance(ctx context.Context, config *provider.ScanJobConfig) (*Instance, error) {
	options := func(options *ec2.Options) {
		options.Region = config.Region
	}

	ec2TagsForInstance := EC2TagsFromScanMetadata(config.ScanMetadata)
	ec2Filters := EC2FiltersFromEC2Tags(ec2TagsForInstance)

	describeParams := &ec2.DescribeInstancesInput{
		Filters: ec2Filters,
	}
	describeOut, err := c.ec2Client.DescribeInstances(ctx, describeParams, options)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scanner VM instances: %w", err)
	}

	for _, r := range describeOut.Reservations {
		for _, i := range r.Instances {
			ec2Instance := i

			if ec2Instance.InstanceId == nil || ec2Instance.State == nil {
				continue
			}

			ec2State := ec2Instance.State.Name
			if ec2State == ec2types.InstanceStateNameRunning || ec2State == ec2types.InstanceStateNamePending {
				return instanceFromEC2Instance(&ec2Instance, c.ec2Client, config), nil
			}
		}
	}

	userData, err := cloudinit.New(config)
	if err != nil {
		return nil, FatalError{
			Err: fmt.Errorf("failed to generate cloud-init: %w", err),
		}
	}
	userDataBase64 := base64.StdEncoding.EncodeToString([]byte(userData))

	runParams := &ec2.RunInstancesInput{
		MaxCount:     utils.PointerTo[int32](1),
		MinCount:     utils.PointerTo[int32](1),
		ImageId:      &c.awsConfig.AmiID,
		InstanceType: ec2types.InstanceType(c.awsConfig.InstanceType),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags:         ec2TagsForInstance,
			},
			{
				ResourceType: ec2types.ResourceTypeVolume,
				Tags:         ec2TagsForInstance,
			},
		},
		UserData: &userDataBase64,
		MetadataOptions: &ec2types.InstanceMetadataOptionsRequest{
			HttpEndpoint:         ec2types.InstanceMetadataEndpointStateEnabled,
			InstanceMetadataTags: ec2types.InstanceMetadataTagsStateEnabled,
		},
	}

	// Create network interface in the scanner subnet with the scanner security group.
	runParams.NetworkInterfaces = []ec2types.InstanceNetworkInterfaceSpecification{
		{
			AssociatePublicIpAddress: utils.PointerTo(false),
			DeleteOnTermination:      utils.PointerTo(true),
			DeviceIndex:              utils.PointerTo[int32](0),
			Groups:                   []string{c.awsConfig.SecurityGroupID},
			SubnetId:                 &c.awsConfig.SubnetID,
		},
	}

	var retryMaxAttempts int
	// Use spot instances if there is a configuration for it.
	if config.ScannerInstanceCreationConfig.UseSpotInstances {
		runParams.InstanceMarketOptions = &ec2types.InstanceMarketOptionsRequest{
			MarketType: ec2types.MarketTypeSpot,
			SpotOptions: &ec2types.SpotMarketOptions{
				InstanceInterruptionBehavior: ec2types.InstanceInterruptionBehaviorTerminate,
				SpotInstanceType:             ec2types.SpotInstanceTypeOneTime,
				MaxPrice:                     config.ScannerInstanceCreationConfig.MaxPrice,
			},
		}
	}
	// In the case of spot instances, we have higher probability to start an instance
	// by increasing RetryMaxAttempts
	if config.ScannerInstanceCreationConfig.RetryMaxAttempts != nil {
		retryMaxAttempts = *config.ScannerInstanceCreationConfig.RetryMaxAttempts
	}

	if config.KeyPairName != "" {
		// Set a key-pair to the instance.
		runParams.KeyName = &config.KeyPairName
	}

	// if retryMaxAttempts value is 0 it will be ignored
	out, err := c.ec2Client.RunInstances(ctx, runParams, options, func(options *ec2.Options) {
		options.RetryMaxAttempts = retryMaxAttempts
		options.RetryMode = awstype.RetryModeStandard
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}
	if len(out.Instances) < 1 {
		return nil, errors.New("failed to create instance: 0 instance in response")
	}

	return instanceFromEC2Instance(&out.Instances[0], c.ec2Client, config), nil
}

// nolint:cyclop,gocognit,maintidx
func (c *Client) RunTargetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	vmInfo, err := config.TargetInfo.AsVMInfo()
	if err != nil {
		return FatalError{Err: err}
	}

	logger := log.GetLoggerFromContextOrDefault(ctx).WithFields(logrus.Fields{
		"TargetInstanceID": vmInfo.InstanceID,
		"TargetLocation":   vmInfo.Location,
		"ScannerLocation":  config.Region,
		"Provider":         string(c.Kind()),
	})

	// Note(chrisgacsal): In order to speed up the initialization process the scanner instance and the volume are created
	//                    in parallel.

	// Create scanner instance
	numOfGoroutines := 2
	errs := make(chan error, numOfGoroutines)
	var scannnerInstance *Instance
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Trace("Creating scanner VM instance")

		var err error
		scannnerInstance, err = c.createInstance(ctx, config)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to create scanner VM instance: %w", err))
			return
		}

		ready, err := scannnerInstance.IsReady(ctx)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to get scanner VM instance state: %w", err))
			return
		}
		logger.WithFields(logrus.Fields{
			"ScannerInstanceID": scannnerInstance.ID,
		}).Debugf("Scanner instance is ready: %t", ready)
		if !ready {
			errs <- RetryableError{
				Err:   errors.New("scanner instance is not ready"),
				After: InstanceReadynessAfter,
			}
		}
	}()

	// Create volume snapshot from the root volume of the Target Instance used for scanning by:
	// * fetching the Target Instance from provider
	// * creating a volume snapshot from the root volume of the Target Instance
	// * copying the volume snapshot to the region/location of the scanner instance if they are deployed in separate locations
	var destVolSnapshot *Snapshot
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Debug("Getting target VM instance")

		targetVMLocation, err := NewLocation(vmInfo.Location)
		if err != nil {
			errs <- FatalError{
				Err: fmt.Errorf("failed to parse Location for target VM instance: %w", err),
			}
			return
		}

		var SrcEC2Instance *ec2types.Instance
		SrcEC2Instance, err = c.getInstanceWithID(ctx, vmInfo.InstanceID, targetVMLocation.Region)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to fetch target VM instance: %w", err))
			return
		}
		if SrcEC2Instance == nil {
			errs <- FatalError{
				Err: fmt.Errorf("failed to find target VM instance. InstanceID=%s", vmInfo.InstanceID),
			}
			return
		}
		srcInstance := instanceFromEC2Instance(SrcEC2Instance, c.ec2Client, config)

		logger.WithField("TargetInstanceID", srcInstance.ID).Trace("Found target VM instance")

		srcVol := srcInstance.RootVolume()
		if srcVol == nil {
			errs <- FatalError{
				Err: errors.New("failed to get root block device for target VM instance"),
			}
			return
		}

		logger.WithField("TargetVolumeID", srcVol.ID).Debug("Creating target volume snapshot for target VM instance")
		srcVolSnapshot, err := srcVol.CreateSnapshot(ctx)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to create volume snapshot from target volume. TargetVolumeID=%s: %w",
				srcVol.ID, err))
			return
		}

		ready, err := srcVolSnapshot.IsReady(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get volume snapshot state. TargetVolumeSnapshotID=%s: %w",
				srcVolSnapshot.ID, err)
			errs <- WrapError(err)
			return
		}

		logger.WithFields(logrus.Fields{
			"TargetVolumeID":         srcVol.ID,
			"TargetVolumeSnapshotID": srcVolSnapshot.ID,
		}).Debugf("Target volume snapshot is ready: %t", ready)
		if !ready {
			errs <- RetryableError{
				Err:   errors.New("target volume snapshot is not ready"),
				After: SnapshotReadynessAfter,
			}
			return
		}

		logger.WithFields(logrus.Fields{
			"TargetVolumeID":         srcVol.ID,
			"TargetVolumeSnapshotID": srcVolSnapshot.ID,
		}).Debug("Copying target volume snapshot to scanner location")
		destVolSnapshot, err = srcVolSnapshot.Copy(ctx, config.Region)
		if err != nil {
			err = fmt.Errorf("failed to copy target volume snapshot to location. TargetVolumeSnapshotID=%s Location=%s: %w",
				srcVolSnapshot.ID, config.Region, err)
			errs <- WrapError(err)
			return
		}

		ready, err = destVolSnapshot.IsReady(ctx)
		if err != nil {
			err = fmt.Errorf("failed to get volume snapshot state. ScannerVolumeSnapshotID=%s: %w",
				srcVolSnapshot.ID, err)
			errs <- WrapError(err)
			return
		}

		logger.WithFields(logrus.Fields{
			"TargetVolumeID":          srcVol.ID,
			"TargetVolumeSnapshotID":  srcVolSnapshot.ID,
			"ScannerVolumeSnapshotID": destVolSnapshot.ID,
		}).Debugf("Scanner volume snapshot is ready: %t", ready)

		if !ready {
			errs <- RetryableError{
				Err:   errors.New("scanner volume snapshot is not ready"),
				After: SnapshotReadynessAfter,
			}
			return
		}
	}()
	wg.Wait()
	close(errs)

	// NOTE: make sure to drain results channel
	for e := range errs {
		if e != nil {
			// nolint:typecheck
			err = errors.Join(err, e)
		}
	}
	if err != nil {
		return err
	}

	// Create volume to be scanned from snapshot
	scannerInstanceAZ := scannnerInstance.AvailabilityZone
	logger.WithFields(logrus.Fields{
		"ScannerAvailabilityZone": scannerInstanceAZ,
		"ScannerVolumeSnapshotID": destVolSnapshot.ID,
	}).Debug("Creating scanner volume from volume snapshot for scanner VM instance")
	scannerVol, err := destVolSnapshot.CreateVolume(ctx, scannerInstanceAZ)
	if err != nil {
		err = fmt.Errorf("failed to create volume from snapshot. SnapshotID=%s: %w", destVolSnapshot.ID, err)
		return WrapError(err)
	}

	ready, err := scannerVol.IsReady(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check if scanner volume is ready. ScannerVolumeID=%s: %w", scannerVol.ID, err)
		return WrapError(err)
	}
	logger.WithFields(logrus.Fields{
		"ScannerAvailabilityZone": scannerInstanceAZ,
		"ScannerVolumeSnapshotID": destVolSnapshot.ID,
		"ScannerVolumeID":         scannerVol.ID,
	}).Debugf("Scanner volume is ready: %t", ready)
	if !ready {
		return RetryableError{
			Err:   fmt.Errorf("scanner volume is not ready. ScannerVolumeID=%s", scannerVol.ID),
			After: VolumeReadynessAfter,
		}
	}

	// Attach volume to be scanned to scanner instance
	logger.WithFields(logrus.Fields{
		"ScannerAvailabilityZone": scannerInstanceAZ,
		"ScannerVolumeSnapshotID": destVolSnapshot.ID,
		"ScannerVolumeID":         scannerVol.ID,
		"ScannerIntanceID":        scannnerInstance.ID,
	}).Debug("Attaching scanner volume to scanner VM instance")
	err = scannnerInstance.AttachVolume(ctx, scannerVol, config.BlockDeviceName)
	if err != nil {
		err = fmt.Errorf("failed to attach volume to scanner instance. ScannerVolumeID=%s ScannerInstanceID=%s: %w",
			scannerVol.ID, scannnerInstance.ID, err)
		return WrapError(err)
	}

	// Wait until the volume is attached to the scanner instance
	logger.WithFields(logrus.Fields{
		"ScannerAvailabilityZone": scannerInstanceAZ,
		"ScannerVolumeSnapshotID": destVolSnapshot.ID,
		"ScannerVolumeID":         scannerVol.ID,
		"ScannerIntanceID":        scannnerInstance.ID,
	}).Debug("Checking if scanner volume is attached to scanner VM instance")
	ready, err = scannerVol.IsAttached(ctx)
	if err != nil {
		err = fmt.Errorf("failed to check if volume is attached to scanner instance. ScannerVolumeID=%s ScannerInstanceID=%s: %w",
			scannerVol.ID, scannnerInstance.ID, err)
		return WrapError(err)
	}
	if !ready {
		return RetryableError{
			Err:   fmt.Errorf("scanner volume is not attached yet. ScannerVolumeID=%s", scannerVol.ID),
			After: VolumeAttachmentReadynessAfter,
		}
	}

	return nil
}

// deleteInstances terminates all instances which meet the conditions defined by the filters argument.
// It returns:
//   - nil, error: if error happened during the operation
//   - false, nil: if operation is in-progress
//   - true, nil:  if the operation is done and safe to assume that related resources are released (netif, vol, etc)
func (c *Client) deleteInstances(ctx context.Context, filters []ec2types.Filter, region string) (bool, error) {
	options := func(options *ec2.Options) {
		options.Region = region
	}

	describeParams := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	describeOut, err := c.ec2Client.DescribeInstances(ctx, describeParams, options)
	if err != nil {
		return false, fmt.Errorf("failed to fetch instances: %w", err)
	}

	instances := make([]string, 0)
	for _, r := range describeOut.Reservations {
		for _, i := range r.Instances {
			if i.InstanceId == nil || i.State == nil {
				continue
			}
			if i.State.Name != ec2types.InstanceStateNameTerminated {
				instances = append(instances, *i.InstanceId)
			}
		}
	}

	if len(instances) > 0 {
		terminateParams := &ec2.TerminateInstancesInput{
			InstanceIds: instances,
		}
		_, err = c.ec2Client.TerminateInstances(ctx, terminateParams, options)
		if err != nil {
			return false, fmt.Errorf("failed to terminate instances %v: %w", instances, err)
		}
		return false, nil
	}

	return true, nil
}

// deleteVolumes deletes all volumes which meet the conditions defined by the filters argument.
// It returns:
//   - nil, error: if error happened during the operation
//   - false, nil: if operation is in-progress
//   - true, nil:  if the operation is done
func (c *Client) deleteVolumes(ctx context.Context, filters []ec2types.Filter, region string) (bool, error) {
	options := func(options *ec2.Options) {
		options.Region = region
	}

	describeParams := &ec2.DescribeVolumesInput{
		Filters: filters,
	}
	describeOut, err := c.ec2Client.DescribeVolumes(ctx, describeParams, options)
	if err != nil {
		return false, fmt.Errorf("failed to fetch volumes: %w", err)
	}

	volumes := make([]string, 0)
	for _, vol := range describeOut.Volumes {
		if vol.State == ec2types.VolumeStateDeleted || vol.State == ec2types.VolumeStateDeleting {
			continue
		}
		volumes = append(volumes, *vol.VolumeId)
	}

	if len(volumes) > 0 {
		for _, vol := range volumes {
			terminateParams := &ec2.DeleteVolumeInput{
				VolumeId: utils.PointerTo(vol),
			}
			_, err = c.ec2Client.DeleteVolume(ctx, terminateParams, options)
			if err != nil {
				return false, fmt.Errorf("failed to delete volume with %s id: %w", vol, err)
			}
		}
		return false, nil
	}

	return true, nil
}

// deleteVolumeSnapshots deletes all volume snapshots which meet the conditions defined by the filters argument.
// It returns:
//   - nil, error: if error happened during the operation
//   - false, nil: if operation is in-progress
//   - true, nil:  if the operation is done
func (c *Client) deleteVolumeSnapshots(ctx context.Context, filters []ec2types.Filter, region string) (bool, error) {
	options := func(options *ec2.Options) {
		options.Region = region
	}

	describeParams := &ec2.DescribeSnapshotsInput{
		Filters: filters,
	}
	describeOut, err := c.ec2Client.DescribeSnapshots(ctx, describeParams, options)
	if err != nil {
		return false, fmt.Errorf("failed to fetch volume snapshots: %w", err)
	}

	snapshots := make([]string, 0)
	for _, snap := range describeOut.Snapshots {
		snapshots = append(snapshots, *snap.SnapshotId)
	}

	for _, snap := range snapshots {
		deleteParams := &ec2.DeleteSnapshotInput{
			SnapshotId: utils.PointerTo(snap),
		}
		_, err = c.ec2Client.DeleteSnapshot(ctx, deleteParams, options)
		if err != nil {
			return false, fmt.Errorf("failed to delete volume snapshot with %s id: %w", snap, err)
		}
	}

	return true, nil
}

// RemoveTargetScan removes all the cloud resources associated with a Scan defined by config parameter.
// The operation is idempotent, therefore it is safe to call it multiple times.
// nolint:cyclop,gocognit
func (c *Client) RemoveTargetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	vmInfo, err := config.TargetInfo.AsVMInfo()
	if err != nil {
		return FatalError{Err: err}
	}

	logger := log.GetLoggerFromContextOrDefault(ctx).WithFields(logrus.Fields{
		"ScannerLocation": config.Region,
		"Provider":        string(c.Kind()),
	})

	ec2Tags := EC2TagsFromScanMetadata(config.ScanMetadata)
	ec2Filters := EC2FiltersFromEC2Tags(ec2Tags)

	numOfGoroutines := 3
	errs := make(chan error, numOfGoroutines)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Delete scanner instance
		logger.Debug("Deleting scanner VM Instance.")
		done, err := c.deleteInstances(ctx, ec2Filters, config.Region)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to delete scanner VM instance: %w", err))
			return
		}
		// Deleting scanner VM instance is in-progress, thus cannot proceed with deleting the scanner volume.
		if !done {
			errs <- RetryableError{
				Err:   errors.New("deleting Scanner VM instance is in-progress"),
				After: InstanceReadynessAfter,
			}
			return
		}

		// Delete scanner volume
		logger.Debug("Deleting scanner volume.")
		done, err = c.deleteVolumes(ctx, ec2Filters, config.Region)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to delete scanner volume: %w", err))
			return
		}

		if !done {
			errs <- RetryableError{
				Err:   errors.New("deleting Scanner volume is in-progress"),
				After: VolumeReadynessAfter,
			}
			return
		}
	}()

	// Delete volume snapshot created for the scanner volume
	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Debug("Deleting scanner volume snapshot.")
		done, err := c.deleteVolumeSnapshots(ctx, ec2Filters, config.Region)
		if err != nil {
			errs <- WrapError(fmt.Errorf("failed to delete scanner volume snapshot: %w", err))
			return
		}
		if !done {
			errs <- RetryableError{
				Err:   errors.New("deleting Scanner volume snapshot is in-progress"),
				After: SnapshotReadynessAfter,
			}
			return
		}
	}()

	// Delete volume snapshot created from the Target instance volume
	wg.Add(1)
	go func() {
		defer wg.Done()

		location, err := NewLocation(vmInfo.Location)
		if err != nil {
			errs <- FatalError{
				Err: fmt.Errorf("failed to parse Location string. Location=%s: %w", vmInfo.Location, err),
			}
			return
		}

		if location.Region == config.Region {
			return
		}

		logger.WithField("TargetLocation", vmInfo.Location).Debug("Deleting target volume snapshot.")
		done, err := c.deleteVolumeSnapshots(ctx, ec2Filters, location.Region)
		if err != nil {
			errs <- fmt.Errorf("failed to delete target volume snapshot: %w", err)
			return
		}

		if !done {
			errs <- RetryableError{
				Err:   errors.New("deleting Target volume snapshot is in-progress"),
				After: SnapshotReadynessAfter,
			}
			return
		}
	}()
	wg.Wait()
	close(errs)

	// NOTE: make sure to drain results channel
	for e := range errs {
		if e != nil {
			// nolint:typecheck
			err = errors.Join(err, e)
		}
	}

	return err
}

func (c *Client) GetInstances(ctx context.Context, filters []ec2types.Filter, excludeTags []models.Tag, regionID string) ([]Instance, error) {
	ret := make([]Instance, 0)

	out, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: utils.PointerTo[int32](maxResults), // TODO what will be a good number?
	}, func(options *ec2.Options) {
		options.Region = regionID
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}
	ret = append(ret, c.getInstancesFromDescribeInstancesOutput(ctx, out, excludeTags, regionID)...)

	// use pagination
	// TODO we can make it better by not saving all results in memory. See https://github.com/openclarity/vmclarity/pull/3#discussion_r1021656861
	for out.NextToken != nil {
		out, err = c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
			Filters:    filters,
			MaxResults: utils.PointerTo[int32](maxResults),
			NextToken:  out.NextToken,
		}, func(options *ec2.Options) {
			options.Region = regionID
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		ret = append(ret, c.getInstancesFromDescribeInstancesOutput(ctx, out, excludeTags, regionID)...)
	}

	return ret, nil
}

func (c *Client) getInstancesFromDescribeInstancesOutput(ctx context.Context, result *ec2.DescribeInstancesOutput, excludeTags []models.Tag, regionID string) []Instance {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	var ret []Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if hasExcludeTags(excludeTags, instance.Tags) {
				continue
			}
			if err := validateInstanceFields(instance); err != nil {
				logger.Errorf("Instance validation failed. instance id=%v: %v", getPointerValOrEmpty(instance.InstanceId), err)
				continue
			}
			ret = append(ret, Instance{
				ID:               *instance.InstanceId,
				Region:           regionID,
				AvailabilityZone: *instance.Placement.AvailabilityZone,
				Image:            *instance.ImageId,
				InstanceType:     string(instance.InstanceType),
				Platform:         *instance.PlatformDetails,
				Tags:             getTagsFromECTags(instance.Tags),
				LaunchTime:       *instance.LaunchTime,
				VpcID:            *instance.VpcId,
				SecurityGroups:   getSecurityGroupsIDs(instance.SecurityGroups),

				ec2Client: c.ec2Client,
			})
		}
	}
	return ret
}

func (c *Client) getRegionsToScan(ctx context.Context, scope *ScanScope) ([]Region, error) {
	if scope.AllRegions {
		return c.ListAllRegions(ctx, false)
	}

	return scope.Regions, nil
}

func (c *Client) ListAllRegions(ctx context.Context, isRecursive bool) ([]Region, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	ret := make([]Region, 0)
	out, err := c.ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: nil, // display also disabled regions?
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %w", err)
	}

	for _, region := range out.Regions {
		ret = append(ret, Region{
			Name: *region.RegionName,
		})
	}

	if isRecursive {
		for i, region := range ret {
			// List region VPCs
			vpcs, err := c.ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
				MaxResults: utils.PointerTo[int32](maxResults),
			}, func(options *ec2.Options) {
				options.Region = region.Name
			})
			if err != nil {
				logger.Errorf("Failed to describe vpcs. Region=%s: %v", region.Name, err)
				continue
			}
			ret[i].VPCs = convertAwsVPCs(vpcs.Vpcs)
			for i2, vpc := range ret[i].VPCs {
				// List VPC's security groups
				securityGroups, err := c.ec2Client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
					Filters: []ec2types.Filter{
						{
							Name:   utils.PointerTo(VpcIDFilterName),
							Values: []string{vpc.ID},
						},
					},
					MaxResults: utils.PointerTo[int32](maxResults),
				}, func(options *ec2.Options) {
					options.Region = region.Name
				})
				if err != nil {
					logger.Errorf("Failed to describe security groups. Region=%s, VpcID=%s: %v", region.Name, vpc.ID, err)
					continue
				}
				ret[i].VPCs[i2].SecurityGroups = getSecurityGroupsFromEC2SecurityGroups(securityGroups.SecurityGroups)
			}
		}
	}

	return ret, nil
}
