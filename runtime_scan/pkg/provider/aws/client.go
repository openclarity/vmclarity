package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

type Client struct {
	ec2Client *ec2.Client
}

var (
	snapshotDescription = "VMClarity snapshot"
	tagKey              = "Owner"
	tagVal              = "VMClarity"
	vmclarityTags       = []ec2types.Tag{
		{
			Key:   &tagKey,
			Value: &tagVal,
		},
	}
)

func Create() (*Client, error) {
	var awsClient Client

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %v", err)
	}

	awsClient.ec2Client = ec2.NewFromConfig(cfg)

	return &awsClient, nil
}

func (c *Client) Discover(scope *types.ScanScope) ([]types.Instance, error) {
	var ret []types.Instance
	var filters []ec2types.Filter

	regions, err := c.getRegionsToScan(scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions to scan: %v", err)
	}
	if len(regions) == 0 {
		return nil, fmt.Errorf("no regions to scan")
	}
	filters = append(filters, createInclusionTagsFilters(scope.IncludeTags)...)
	filters = append(filters, createInstanceStateFilters(scope.ScanStopped)...)

	for _, region := range regions {
		// if no vpcs, that mean that we don't need any vpc filters
		if len(region.VPCs) == 0 {
			instances, err := c.GetInstances(filters, scope.ExcludeTags, region.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %v", err)
			}
			ret = append(ret, instances...)
			continue
		}

		// need to do a per vpc call for DescribeInstances
		for _, vpc := range region.VPCs {
			vpcFilters := append(filters, createVPCFilters(vpc)...)

			instances, err := c.GetInstances(vpcFilters, scope.ExcludeTags, region.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %v", err)
			}
			ret = append(ret, instances...)
		}
	}
	return ret, nil
}

func (c *Client) GetInstances(filters []ec2types.Filter, excludeTags []*types.Tag, regionID string) ([]types.Instance, error) {
	var ret []types.Instance

	out, err := c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: utils.Int32Ptr(50), // TODO what will be a good number?
	}, func(options *ec2.Options) {
		options.Region = regionID
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}
	ret = append(ret, getInstancesFromDescribeInstancesOutput(out, excludeTags, regionID)...)

	// use pagination
	for out.NextToken != nil {
		out, err = c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters:    filters,
			MaxResults: utils.Int32Ptr(50),
			NextToken:  out.NextToken,
		}, func(options *ec2.Options) {
			options.Region = regionID
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		ret = append(ret, getInstancesFromDescribeInstancesOutput(out, excludeTags, regionID)...)
	}

	return ret, nil
}

func (c *Client) CreateSnapshot(volume types.Volume) (types.Snapshot, error) {
	params := ec2.CreateSnapshotInput{
		VolumeId:    &volume.ID,
		Description: &snapshotDescription,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSnapshot,
				Tags:         vmclarityTags,
			},
		},
	}
	out, err := c.ec2Client.CreateSnapshot(context.TODO(), &params, func(options *ec2.Options) {
		options.Region = volume.Region
	})
	if err != nil {
		return types.Snapshot{}, fmt.Errorf("failed to create snapshot: %v", err)
	}
	return types.Snapshot{
		ID:     *out.SnapshotId,
		Region: volume.Region,
	}, nil
}

func (c *Client) WaitForSnapshotReady(snapshot types.Snapshot) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	timeout := time.After(3 * time.Minute)

	for {
		select {
		case <-ticker.C:
			out, err := c.ec2Client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
				SnapshotIds: []string{snapshot.ID},
			}, func(options *ec2.Options) {
				options.Region = snapshot.Region
			})
			if err != nil {
				return fmt.Errorf("failed to describe snapshot. snapshotID=%v: %v", snapshot.ID, err)
			}
			if len(out.Snapshots) != 1 {
				return fmt.Errorf("got unexcpected number of snapshots (%v) with snapshot id %v. excpecting 1", len(out.Snapshots), snapshot.ID)
			}
			if out.Snapshots[0].State == ec2types.SnapshotStateCompleted {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout")
		}
	}
}

func (c *Client) WaitForInstanceReady(instance types.Instance) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	timeout := time.After(3 * time.Minute)

	for {
		select {
		case <-ticker.C:
			out, err := c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
				InstanceIds: []string{instance.ID},
			}, func(options *ec2.Options) {
				options.Region = instance.Region
			})
			if err != nil {
				return fmt.Errorf("failed to describe instance. instanceID=%v: %v", instance.ID, err)
			}
			state := getInstanceState(out, instance.ID)
			if state == ec2types.InstanceStateNameRunning {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout")
		}
	}
}

func (c *Client) CopySnapshot(snapshot types.Snapshot, dstRegion string) (types.Snapshot, error) {
	snap, err := c.ec2Client.CopySnapshot(context.TODO(), &ec2.CopySnapshotInput{
		SourceRegion:     &snapshot.Region,
		SourceSnapshotId: &snapshot.ID,
		Description:      &snapshotDescription,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeSnapshot,
				Tags:         vmclarityTags,
			},
		},
	}, func(options *ec2.Options) {
		options.Region = dstRegion
	})
	if err != nil {
		return types.Snapshot{}, fmt.Errorf("failed to copy snapshot: %v", err)
	}

	return types.Snapshot{
		ID:     *snap.SnapshotId,
		Region: dstRegion,
	}, nil
}

func (c *Client) GetInstanceRootVolume(instance types.Instance) (types.Volume, error) {
	out, err := c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instance.ID},
	}, func(options *ec2.Options) {
		options.Region = instance.Region
	})
	if err != nil {
		return types.Volume{}, fmt.Errorf("failed to describe instances: %v", err)
	}

	if len(out.Reservations) == 0 {
		return types.Volume{}, fmt.Errorf("no reservations were found")
	}
	if len(out.Reservations) > 1 {
		return types.Volume{}, fmt.Errorf("more than one reservations were found")
	}
	if len(out.Reservations[0].Instances) == 0 {
		return types.Volume{}, fmt.Errorf("no instances were found")
	}
	if len(out.Reservations[0].Instances) > 1 {
		return types.Volume{}, fmt.Errorf("more than one instances were found")
	}

	outInstance := out.Reservations[0].Instances[0]
	rootDeviceName := *outInstance.RootDeviceName

	// find root volume of the instance
	for _, blkDevice := range outInstance.BlockDeviceMappings {
		if strings.Compare(*blkDevice.DeviceName, rootDeviceName) == 0 {
			return types.Volume{
				ID:     *blkDevice.Ebs.VolumeId,
				Name:   rootDeviceName,
				Region: instance.Region,
			}, nil
		}
	}
	return types.Volume{}, fmt.Errorf("failed to find root device volume")
}

func (c *Client) LaunchInstance(ami, deviceName string, snapshot types.Snapshot) (types.Instance, error) {
	out, err := c.ec2Client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		MaxCount: utils.Int32Ptr(1),
		MinCount: utils.Int32Ptr(1),
		ImageId:  &ami,
		BlockDeviceMappings: []ec2types.BlockDeviceMapping{
			{
				// attach the snapshot to the instance at launch (a new volume will be created)
				DeviceName: &deviceName,
				Ebs: &ec2types.EbsBlockDevice{
					DeleteOnTermination: utils.BoolPtr(true),
					Encrypted:           nil, // ?
					SnapshotId:          &snapshot.ID,
					VolumeSize:          nil,                    // default is snapshot size
					VolumeType:          ec2types.VolumeTypeGp2, // TODO need to decide volume type
				},
			},
		},
		InstanceType:   ec2types.InstanceTypeT2Large, // TODO need to decide instance type
		SecurityGroups: nil,                          // use default for now
		SubnetId:       nil,                          // use default for now
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags:         vmclarityTags,
			},
		},
		UserData: nil, // TODO put launch script here
	}, func(options *ec2.Options) {
		options.Region = snapshot.Region
	})
	if err != nil {
		return types.Instance{}, fmt.Errorf("failed to run instances: %v", err)
	}

	return types.Instance{
		ID:     *out.Instances[0].InstanceId,
		Region: snapshot.Region,
	}, nil
}

func (c *Client) DeleteInstance(instance types.Instance) error {
	_, err := c.ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instance.ID},
	}, func(options *ec2.Options) {
		options.Region = instance.Region
	})
	if err != nil {
		return fmt.Errorf("failed to terminate instances: %v", err)
	}

	return nil
}

func (c *Client) DeleteSnapshot(snapshot types.Snapshot) error {
	_, err := c.ec2Client.DeleteSnapshot(context.TODO(), &ec2.DeleteSnapshotInput{
		SnapshotId: &snapshot.ID,
	}, func(options *ec2.Options) {
		options.Region = snapshot.Region
	})
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %v", err)
	}

	return nil
}

func getInstanceState(result *ec2.DescribeInstancesOutput, instanceID string) ec2types.InstanceStateName {
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if strings.Compare(*instance.InstanceId, instanceID) == 0 {
				if instance.State != nil {
					return instance.State.Name
				}
			}
		}
	}
	return ec2types.InstanceStateNamePending
}

func getInstancesFromDescribeInstancesOutput(result *ec2.DescribeInstancesOutput, excludeTags []*types.Tag, regionID string) []types.Instance {
	var ret []types.Instance

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if hasExcludeTags(excludeTags, instance.Tags) {
				continue
			}
			ret = append(ret, types.Instance{
				ID:     *instance.InstanceId,
				Region: regionID,
			})
		}
	}
	return ret
}

func getVPCSecurityGroupsIDs(vpc types.VPC) []string {
	var sgs []string
	for _, sg := range vpc.SecurityGroups {
		sgs = append(sgs, sg.ID)
	}
	return sgs
}

const (
	vpcIDFilterName         = "vpc-id"
	sgIDFilterName          = "instance.group-id" // TODO is this the right one?
	instanceStateFilterName = "instance-state-name"
)

func createVPCFilters(vpc types.VPC) []ec2types.Filter {
	var ret = make([]ec2types.Filter, 0)
	var sgs []string

	// create per vpc filters
	ret = append(ret, ec2types.Filter{
		Name:   utils.StringPtr(vpcIDFilterName),
		Values: []string{vpc.ID},
	})
	sgs = getVPCSecurityGroupsIDs(vpc)
	if len(sgs) > 0 {
		ret = append(ret, ec2types.Filter{
			Name:   utils.StringPtr(sgIDFilterName),
			Values: sgs,
		})
	}

	log.Infof("VPC filter created: %+v", ret)

	return ret
}

func createInstanceStateFilters(scanStopped bool) []ec2types.Filter {
	var filters []ec2types.Filter
	var states = []string{"running"}
	if scanStopped {
		states = append(states, "stopped")
	}

	// TODO these are the states: pending | running | shutting-down | terminated | stopping | stopped
	// Do we want to scan any other state (other than running and stopped)
	filters = append(filters, ec2types.Filter{
		Name:   utils.StringPtr(instanceStateFilterName),
		Values: states,
	})
	return filters
}

func createInclusionTagsFilters(tags []*types.Tag) []ec2types.Filter {
	var filters []ec2types.Filter

	for _, tag := range tags {
		filters = append(filters, ec2types.Filter{
			Name:   utils.StringPtr("tag:" + tag.Key),
			Values: []string{tag.Val},
		})
	}

	return filters
}

func (c *Client) getRegionsToScan(scope *types.ScanScope) ([]types.Region, error) {
	if scope.All {
		return c.ListAllRegions()
	}

	var ret []types.Region
	for _, region := range scope.Regions {
		ret = append(ret, region)
	}

	return ret, nil
}

func (c *Client) ListAllRegions() ([]types.Region, error) {
	var ret []types.Region
	out, err := c.ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{
		AllRegions: nil, // display also disabled regions?
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %v", err)
	}
	for _, region := range out.Regions {
		ret = append(ret, types.Region{
			ID: *region.RegionName,
		})
	}
	return ret, nil
}

func hasExcludeTags(excludeTags []*types.Tag, instanceTags []ec2types.Tag) bool {
	var excludedTagsMap = make(map[string]string)

	for _, tag := range excludeTags {
		excludedTagsMap[tag.Key] = tag.Val
	}
	for _, instanceTag := range instanceTags {
		if val, ok := excludedTagsMap[*instanceTag.Key]; ok {
			if strings.Compare(val, *instanceTag.Value) == 0 {
				return true
			}
		}
	}
	return false
}
