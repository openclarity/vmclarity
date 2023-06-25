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

package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/shared/pkg/log"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type Client struct {
	snapshotsClient *compute.SnapshotsClient
	disksClient     *compute.DisksClient
	instancesClient *compute.InstancesClient
	regionsClient   *compute.RegionsClient

	gcpConfig Config
}

func New(ctx context.Context) (*Client, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	client := Client{
		gcpConfig: config,
	}

	regionsClient, err := compute.NewRegionsRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create regions client: %w", err)
	}
	client.regionsClient = regionsClient

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance client: %w", err)
	}
	client.instancesClient = instancesClient

	snapshotsClient, err := compute.NewSnapshotsRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot client: %w", err)
	}
	client.snapshotsClient = snapshotsClient

	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create disks client: %w", err)
	}
	client.disksClient = disksClient

	return &client, nil
}

func (c Client) Kind() models.CloudProvider {
	return models.GCP
}

// nolint:cyclop
func (c *Client) RunTargetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	// convert TargetInfo to vmInfo
	vminfo, err := config.TargetInfo.AsVMInfo()
	if err != nil {
		return provider.FatalErrorf("unable to get vminfo from target: %w", err)
	}

	logger := log.GetLoggerFromContextOrDefault(ctx).WithFields(logrus.Fields{
		"ScanResultID":   config.ScanResultID,
		"TargetLocation": vminfo.Location,
		"InstanceID":     vminfo.InstanceID,
		"ScannerZone":    c.gcpConfig.ScannerZone,
		"Provider":       string(c.Kind()),
	})
	logger.Debugf("Running target scan")

	targetName := vminfo.InstanceID
	targetZone := vminfo.Location

	// get the target instance to scan from gcp.
	targetVM, err := c.instancesClient.Get(ctx, &computepb.GetInstanceRequest{
		Instance: targetName,
		Project:  c.gcpConfig.ProjectID,
		Zone:     targetZone,
	})
	if err != nil {
		_, err := handleGcpRequestError(err, "getting target virtual machine %v", targetName)
		return err
	}
	logger.Debugf("Got target VM: %v", targetVM.Name)

	// get target instance boot disk
	bootDisk, err := getInstanceBootDisk(targetVM)
	if err != nil {
		return provider.FatalErrorf("unable to get instance boot disk: %w", err)
	}
	logger.Debugf("Got target boot disk: %v", bootDisk.GetSource())

	// ensure that a snapshot was created from the target instance root disk. (create if not)
	snapshot, err := c.ensureSnapshotFromAttachedDisk(ctx, config, bootDisk)
	if err != nil {
		return fmt.Errorf("failed to ensure snapshot for vm root volume: %w", err)
	}
	logger.Debugf("Created snapshot: %v", snapshot.Name)

	// create a disk from the snapshot.
	// Snapshots are global resources, so any snapshot is accessible by any resource within the same project.
	var diskFromSnapshot *computepb.Disk
	diskFromSnapshot, err = c.ensureDiskFromSnapshot(ctx, config, snapshot)
	if err != nil {
		return fmt.Errorf("failed to ensure disk created from snapshot: %w", err)
	}
	logger.Debugf("Created disk from snapshot: %v", diskFromSnapshot.Name)

	// create the scanner instance
	scannerVM, err := c.ensureScannerVirtualMachine(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to ensure scanner virtual machine: %w", err)
	}
	logger.Debugf("Created scanner virtual machine: %v", scannerVM.Name)

	// attach the disk from snapshot to the scanner instance
	err = c.ensureDiskAttachedToScannerVM(ctx, scannerVM, diskFromSnapshot)
	if err != nil {
		return fmt.Errorf("failed to ensure target disk is attached to virtual machine: %w", err)
	}
	logger.Debugf("Attached disk to scanner virtual machine")

	return nil
}

func (c *Client) RemoveTargetScan(ctx context.Context, config *provider.ScanJobConfig) error {
	logger := log.GetLoggerFromContextOrDefault(ctx).WithFields(logrus.Fields{
		"scanResultID": config.ScanResultID,
		"ScannerZone":  c.gcpConfig.ScannerZone,
		"Provider":     string(c.Kind()),
	})

	err := c.ensureScannerVirtualMachineDeleted(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to ensure scanner virtual machine deleted: %w", err)
	}
	logger.Debugf("Deleted scanner virtual machine")

	err = c.ensureTargetDiskDeleted(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to ensure target disk deleted: %w", err)
	}
	logger.Debugf("Deleted disk")

	err = c.ensureSnapshotDeleted(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to ensure snapshot deleted: %w", err)
	}
	logger.Debugf("Deleted snapshot")

	return nil
}

func (c *Client) DiscoverScopes(ctx context.Context) (*models.Scopes, error) {
	var ret models.Scopes
	ret.ScopeInfo = &models.ScopeType{}
	var regions []models.GcpRegion

	gcpRegions, err := c.listAllRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list regions: %w", err)
	}

	for _, region := range gcpRegions {
		regions = append(regions, models.GcpRegion{
			Name:  *region.Name,
			Zones: utils.PointerTo(getZonesLastPart(region.Zones)),
		})
	}

	err = ret.ScopeInfo.FromGcpProjectScope(models.GcpProjectScope{
		ProjectID: c.gcpConfig.ProjectID,
		Regions:   &regions,
	})
	if err != nil {
		return nil, provider.FatalErrorf("failed to convert from gcp project scope: %v", err)
	}

	return &ret, nil
}

// nolint: cyclop
func (c *Client) DiscoverTargets(ctx context.Context, scanScope *models.ScanScopeType) ([]models.TargetType, error) {
	var ret []models.TargetType

	gcpScanScope, err := scanScope.AsGcpScanScope()
	if err != nil {
		return nil, fmt.Errorf("failed to convert as gcp scan scope: %v", err)
	}

	// get list of zones to scan
	zones, err := c.getZonesToScanFromScanScope(ctx, gcpScanScope)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones to scan from scan scope: %w", err)
	}

	// prepare include and exclude tags filter
	//filter := prepareTagsFilter(gcpScanScope.InstanceTagSelector, gcpScanScope.InstanceTagExclusion)

	// TODO (erezf) unfortunately the tags filter is broken in the google api, so I will need to fetch all instances in required zones, and filter by tags in memory.
	for _, zone := range zones {
		targets, err := c.listInstances(ctx, nil, zone, gcpScanScope)
		if err != nil {
			return nil, fmt.Errorf("failed to list instances: %w", err)
		}

		ret = append(ret, targets...)
	}

	return ret, nil
}

// [https://www.googleapis.com/compute/v1/projects/gcp-etigcp-nprd-12855/zones/us-central1-c, https://www.googleapis.com/compute/v1/projects/gcp-etigcp-nprd-12855/zones/us-central1-a] -> [us-central1-c, us-central1-a]
func getZonesLastPart(zones []string) []string {
	var ret []string

	for _, zone := range zones {
		ret = append(ret, getLastURLPart(&zone))
	}
	return ret
}

func getInstanceBootDisk(vm *computepb.Instance) (*computepb.AttachedDisk, error) {
	for _, disk := range vm.Disks {
		if disk.Boot != nil && *disk.Boot {
			return disk, nil
		}
	}
	return nil, fmt.Errorf("failed to find instance boot disk")
}

//func prepareTagsFilter(includeTags *[]models.Tag, excludeTags *[]models.Tag) string {
//	filter := ""
//
//	if includeTags != nil {
//		for _, tag := range *includeTags {
//			filter += fmt.Sprintf("tags.items=%v AND ", tag.Key)
//		}
//	}
//
//	if excludeTags != nil {
//		for _, tag := range *excludeTags {
//			filter += fmt.Sprintf("-tags.items != %v AND ", tag.Key)
//		}
//	}
//
//	return strings.TrimSuffix(filter, " AND ")
//}

func (c *Client) getZonesToScanFromScanScope(ctx context.Context, scope models.GcpScanScope) ([]string, error) {
	var zones []string
	if scope.AllRegions != nil && *scope.AllRegions {
		regions, err := c.listAllRegions(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list regions: %w", err)
		}
		for _, region := range regions {
			zones = append(zones, getZonesLastPart(region.Zones)...)
		}
		return zones, nil
	}

	if scope.Regions == nil {
		return nil, fmt.Errorf("no regions specifies in scan scope")
	}

	// list zones for specified regions
	for _, region := range *scope.Regions {
		if region.Zones == nil {
			// user did not specify zones, meaning he wants to scan all zones in that region
			// first, get the region information from gcp
			retRegion, err := c.regionsClient.Get(ctx, &computepb.GetRegionRequest{
				Project: c.gcpConfig.ProjectID,
				Region:  region.Name,
			})
			if err != nil {
				_, err := handleGcpRequestError(err, "get region")
				return nil, err
			}

			zones = append(zones, getZonesLastPart(retRegion.Zones)...)
		} else {
			// user specified specific zones to scan in that region
			zones = append(zones, *region.Zones...)
		}
	}
	return zones, nil
}

func (c *Client) listInstances(ctx context.Context, filter *string, zone string, scanScope models.GcpScanScope) ([]models.TargetType, error) {
	var ret []models.TargetType

	it := c.instancesClient.List(ctx, &computepb.ListInstancesRequest{
		Filter:     filter,
		MaxResults: utils.PointerTo[uint32](maxResults),
		Project:    c.gcpConfig.ProjectID,
		Zone:       zone,
	})
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			_, err := handleGcpRequestError(err, "list vms")
			return nil, err
		}
		if !isInScopeByTags(resp, scanScope.InstanceTagSelector, scanScope.InstanceTagExclusion) {
			logrus.Debugf("Ignoring vm %v due to tags filters", *resp.Name)
			continue
		}
		info, err := c.getVMInfoFromVirtualMachine(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to get vminfo from virtual machine: %w", err)
		}
		ret = append(ret, info)
	}

	return ret, nil
}

func isInScopeByTags(vm *computepb.Instance, includeTags *[]models.Tag, excludeTags *[]models.Tag) bool {
	if !hasIncludeTags(vm, includeTags) {
		return false
	}
	if hasExcludeTags(vm, excludeTags) {
		return false
	}
	return true
}

// AND logic - if tags = {tag1:val1, tag2:val2},
// then a vm will be excluded/included only if it has ALL of these tags ({tag1:val1, tag2:val2}).
func hasIncludeTags(vm *computepb.Instance, tags *[]models.Tag) bool {
	if tags == nil {
		return true
	}
	if len(*tags) == 0 {
		return true
	}
	if len(vm.Tags.Items) == 0 {
		return false
	}

	return hasTags(vm.Tags, tags)
}

// AND logic - if tags = {tag1:val1, tag2:val2},
// then a vm will be excluded/included only if it has ALL of these tags ({tag1:val1, tag2:val2}).
func hasExcludeTags(vm *computepb.Instance, tags *[]models.Tag) bool {
	if tags == nil {
		return false
	}
	if len(*tags) == 0 {
		return false
	}
	if len(vm.Tags.Items) == 0 {
		return false
	}

	return hasTags(vm.Tags, tags)
}

func hasTags(vmTags *computepb.Tags, modelsTags *[]models.Tag) bool {
	instanceTags := convertTagsToMap(vmTags)

	for _, tag := range *modelsTags {
		val, ok := instanceTags[tag.Key]
		if !ok {
			return false
		}
		if !(strings.Compare(val, tag.Value) == 0) {
			return false
		}
	}
	return true
}

func (c *Client) listAllRegions(ctx context.Context) ([]*computepb.Region, error) {
	var ret []*computepb.Region

	it := c.regionsClient.List(ctx, &computepb.ListRegionsRequest{
		MaxResults: utils.PointerTo[uint32](maxResults),
		Project:    c.gcpConfig.ProjectID,
	})
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			_, err := handleGcpRequestError(err, "list regions")
			return nil, err
		}

		ret = append(ret, resp)
	}
	return ret, nil
}

func (c *Client) getVMInfoFromVirtualMachine(vm *computepb.Instance) (models.TargetType, error) {
	targetType := models.TargetType{}
	launchTime, err := time.Parse(time.RFC3339, *vm.CreationTimestamp)
	if err != nil {
		return models.TargetType{}, fmt.Errorf("failed to parse time: %v", *vm.CreationTimestamp)
	}
	// get boot disk name
	diskName := getLastURLPart(vm.Disks[0].Source)

	var platform string
	var image string

	// get disk from gcp
	disk, err := c.disksClient.Get(context.TODO(), &computepb.GetDiskRequest{
		Disk:    diskName,
		Project: c.gcpConfig.ProjectID,
		Zone:    getLastURLPart(vm.Zone),
	})
	if err != nil {
		logrus.Warnf("failed to get disk %v: %v", diskName, err)
	} else {
		platform = *disk.Architecture
		image = getLastURLPart(disk.SourceImage)
	}

	err = targetType.FromVMInfo(models.VMInfo{
		InstanceProvider: utils.PointerTo(models.GCP),
		InstanceID:       *vm.Name,
		Image:            image,
		InstanceType:     getLastURLPart(vm.MachineType),
		LaunchTime:       launchTime,
		Location:         getLastURLPart(vm.Zone),
		Platform:         platform,
		SecurityGroups:   &[]models.SecurityGroup{},
		Tags:             convertTags(vm.Tags),
	})
	if err != nil {
		return models.TargetType{}, provider.FatalErrorf("failed to create TargetType from VMInfo: %w", err)
	}

	return targetType, nil
}

// convertTags converts gcp instance tags in the form []string{key1=val1} into models.Tag{Key: key1, Value: val1}
// in case the tag does not contain equal sign, the Key will be the tag and the Value will be empty
func convertTags(tags *computepb.Tags) *[]models.Tag {
	ret := make([]models.Tag, 0, len(tags.Items))
	for _, item := range tags.Items {
		key, val := getKeyValue(item)
		ret = append(ret, models.Tag{
			Key:   key,
			Value: val,
		})
	}
	return &ret
}

func convertTagsToMap(tags *computepb.Tags) map[string]string {
	ret := make(map[string]string, len(tags.Items))
	for _, item := range tags.Items {
		key, val := getKeyValue(item)
		ret[key] = val
	}
	return ret
}

func getKeyValue(str string) (key, value string) {
	spl := strings.Split(str, "=")
	key = spl[0]
	if len(spl) > 1 {
		value = spl[1]
	}
	return
}
