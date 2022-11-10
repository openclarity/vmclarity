package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

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

type Client struct {
	ec2Client *ec2.Client
}

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

	regions, err := c.getRegionsToScan(scope)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions to scan: %v", err)
	}
	if len(regions) == 0 {
		return nil, fmt.Errorf("no regions to scan")
	}

	for _, region := range regions {
		describeFilters := createDescribeFilters(scope, region)

		out, err := c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters:    describeFilters,
			MaxResults: nil, // TODO
			NextToken:  nil, // TODO
		}, func(options *ec2.Options) {
			options.Region = region
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		// get all returned instances:
		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				if instance.InstanceId != nil {
					ret = append(ret, types.Instance{
						ID:     *instance.InstanceId,
						Region: region,
					})
				}
			}
		}
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
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(2 * time.Minute)

	for {
		select {
		case <-ticker.C:
			out, err := c.ec2Client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
				SnapshotIds: []string{snapshot.ID},
			}, func(options *ec2.Options) {
				options.Region = snapshot.Region
			})
			if err != nil {
				return fmt.Errorf("failed to descrive snapshot: %v", err)
			}
			if out.Snapshots[0].State == ec2types.SnapshotStateCompleted {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout")
		}
	}
}

func (c *Client) CopySnapshot(snapshot types.Snapshot, dstRegion string) (types.Snapshot, error) {
	snap, err := c.ec2Client.CopySnapshot(context.TODO(), &ec2.CopySnapshotInput{
		SourceRegion:      &snapshot.Region,
		SourceSnapshotId:  &snapshot.ID,
		Description:       &snapshotDescription,
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
	var max = int32(1)
	var min = int32(1)

	out, err := c.ec2Client.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		MaxCount: &max,
		MinCount: &min,
		ImageId:  &ami,
		BlockDeviceMappings: []ec2types.BlockDeviceMapping{
			{
				DeviceName: &deviceName,
				Ebs: &ec2types.EbsBlockDevice{
					DeleteOnTermination: nil, // ?
					Encrypted:           nil, // ?
					SnapshotId:          &snapshot.ID,
					VolumeSize:          nil,                    // default is snapshot size
					VolumeType:          ec2types.VolumeTypeGp2, // ?
				},
			},
		},
		InstanceType:      ec2types.InstanceTypeT2Large,
		MetadataOptions:   nil, // ?
		SecurityGroups:    nil, // use default for now
		SubnetId:          nil, // use default for now
		TagSpecifications: nil, // need to specify tags
		UserData:          nil, // put launch script here
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

func findRegionByID(regions []types.Region, regionID string) types.Region {
	for _, region := range regions {
		if strings.Compare(region.ID, regionID) == 0 {
			return region
		}
	}
	return types.Region{}
}

func getRegionVPCsAndSecurityGroups(region types.Region) (vpcs []string, sgs []string) {
	for _, vpc := range region.VPCs {
		vpcs = append(vpcs, vpc.ID)
		for _, sc := range vpc.SecurityGroups {
			sgs = append(sgs, sc.ID)
		}
	}
	return
}

func createDescribeFilters(scopes *types.ScanScope, regionID string) []ec2types.Filter {
	var ret []ec2types.Filter
	var vpcs []string
	var sgs []string

	if !scopes.All {
		region := findRegionByID(scopes.Regions, regionID)
		vpcs, sgs = getRegionVPCsAndSecurityGroups(region)
	}

	vpcID := "vpc-id"
	sgID := "instance.group-id" // TODO

	if len(sgs) > 0 {
		ret = append(ret, ec2types.Filter{
			Name:   &sgID,
			Values: sgs,
		})
	} else if len(vpcs) > 0 {
		ret = append(ret, ec2types.Filter{
			Name:   &vpcID,
			Values: vpcs,
		})
	}

	if len(scopes.IncludeTags) > 0 {
		for _, tag := range scopes.IncludeTags {
			name := "tag:" + tag.Key
			ret = append(ret, ec2types.Filter{
				Name:   &name,
				Values: []string{tag.Val},
			})
		}
	}

	return ret
}

func (c *Client) getRegionsToScan(scope *types.ScanScope) ([]string, error) {
	if scope.All {
		return c.ListAllRegions()
	}

	var ret []string
	for _, region := range scope.Regions {
		ret = append(ret, region.ID)
	}

	return ret, nil
}

func (c *Client) ListAllRegions() ([]string, error) {
	var ret []string
	out, err := c.ec2Client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	for _, region := range out.Regions {
		ret = append(ret, *region.RegionName)
	}
	return ret, nil
}
