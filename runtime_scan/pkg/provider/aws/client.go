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

package aws

import (
	"context"
	"fmt"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/config/aws"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

type Client struct {
	ec2Client *ec2.Client
	awsConfig *aws.Config
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

func Create(ctx context.Context, config *aws.Config) (*Client, error) {
	var awsClient = Client{
		awsConfig: config,
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %v", err)
	}

	awsClient.ec2Client = ec2.NewFromConfig(cfg)

	return &awsClient, nil
}

func (c *Client) Discover(ctx context.Context, scope *types.ScanScope) ([]provider.Instance, error) {
	var ret []provider.Instance
	var filters []ec2types.Filter

	regions, err := c.getRegionsToScan(ctx, scope)
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
			instances, err := c.GetInstances(ctx, filters, scope.ExcludeTags, region.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %v", err)
			}
			ret = append(ret, instances...)
			continue
		}

		// need to do a per vpc call for DescribeInstances
		for _, vpc := range region.VPCs {
			vpcFilters := append(filters, createVPCFilters(vpc)...)

			instances, err := c.GetInstances(ctx, vpcFilters, scope.ExcludeTags, region.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get instances: %v", err)
			}
			ret = append(ret, instances...)
		}
	}
	return ret, nil
}

func (c *Client) LaunchInstance(ctx context.Context, snapshot provider.Snapshot) (provider.Instance, error) {
	out, err := c.ec2Client.RunInstances(ctx, &ec2.RunInstancesInput{
		MaxCount: utils.Int32Ptr(1),
		MinCount: utils.Int32Ptr(1),
		ImageId:  &c.awsConfig.AmiID,
		BlockDeviceMappings: []ec2types.BlockDeviceMapping{
			{
				// attach the snapshot to the instance at launch (a new volume will be created)
				DeviceName: &c.awsConfig.DeviceName,
				Ebs: &ec2types.EbsBlockDevice{
					DeleteOnTermination: utils.BoolPtr(true),
					Encrypted:           nil, // ?
					SnapshotId:          utils.StringPtr(snapshot.GetID()),
					VolumeSize:          nil,                    // default is snapshot size
					VolumeType:          ec2types.VolumeTypeGp2, // TODO need to decide volume type
				},
			},
		},
		InstanceType:   ec2types.InstanceTypeT2Large, // TODO need to decide instance type
		SecurityGroups: nil,                          // use default for now
		SubnetId:       &c.awsConfig.SubnetID,
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags:         vmclarityTags,
			},
		},
		UserData: nil, // TODO put launch script here
	}, func(options *ec2.Options) {
		options.Region = snapshot.GetRegion()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to run instances: %v", err)
	}

	return &InstanceImpl{
		ec2Client: c.ec2Client,
		id:        *out.Instances[0].InstanceId,
		region:    snapshot.GetRegion(),
	}, nil
}

func (c *Client) GetInstances(ctx context.Context, filters []ec2types.Filter, excludeTags []*provider.Tag, regionID string) ([]provider.Instance, error) {
	var ret []provider.Instance

	out, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: utils.Int32Ptr(50), // TODO what will be a good number?
	}, func(options *ec2.Options) {
		options.Region = regionID
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}
	ret = append(ret, c.getInstancesFromDescribeInstancesOutput(out, excludeTags, regionID)...)

	// use pagination
	for out.NextToken != nil {
		out, err = c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
			Filters:    filters,
			MaxResults: utils.Int32Ptr(50),
			NextToken:  out.NextToken,
		}, func(options *ec2.Options) {
			options.Region = regionID
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe instances: %v", err)
		}
		ret = append(ret, c.getInstancesFromDescribeInstancesOutput(out, excludeTags, regionID)...)
	}

	return ret, nil
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

func (c *Client) getInstancesFromDescribeInstancesOutput(result *ec2.DescribeInstancesOutput, excludeTags []*provider.Tag, regionID string) []provider.Instance {
	var ret []provider.Instance

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if hasExcludeTags(excludeTags, instance.Tags) {
				continue
			}
			ret = append(ret, &InstanceImpl{
				ec2Client: c.ec2Client,
				id:        *instance.InstanceId,
				region:    regionID,
			})
		}
	}
	return ret
}

func getVPCSecurityGroupsIDs(vpc provider.VPC) []string {
	var sgs []string
	for _, sg := range vpc.SecurityGroups {
		sgs = append(sgs, sg.ID)
	}
	return sgs
}

const (
	vpcIDFilterName         = "vpc-id"
	sgIDFilterName          = "instance.group-id"
	instanceStateFilterName = "instance-state-name"
)

func createVPCFilters(vpc provider.VPC) []ec2types.Filter {
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

func createInclusionTagsFilters(tags []*provider.Tag) []ec2types.Filter {
	var filters []ec2types.Filter

	for _, tag := range tags {
		filters = append(filters, ec2types.Filter{
			Name:   utils.StringPtr("tag:" + tag.Key),
			Values: []string{tag.Val},
		})
	}

	return filters
}

func (c *Client) getRegionsToScan(ctx context.Context, scope *types.ScanScope) ([]provider.Region, error) {
	if scope.All {
		return c.ListAllRegions(ctx)
	}

	return scope.Regions, nil
}

func (c *Client) ListAllRegions(ctx context.Context) ([]provider.Region, error) {
	var ret []provider.Region
	out, err := c.ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: nil, // display also disabled regions?
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %v", err)
	}
	for _, region := range out.Regions {
		ret = append(ret, provider.Region{
			ID: *region.RegionName,
		})
	}
	return ret, nil
}

func hasExcludeTags(excludeTags []*provider.Tag, instanceTags []ec2types.Tag) bool {
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
