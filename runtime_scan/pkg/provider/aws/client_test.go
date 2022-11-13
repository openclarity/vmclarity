package aws

import (
	"reflect"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

//
//func TestClient_ListAllRegions(t *testing.T) {
//	//cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
//	//if err != nil {
//	//	t.Fatalf("%v", err)
//	//}
//	//ec2Client := ec2.NewFromConfig(cfg)
//
//	c, err := Create()
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//	//instance := types.Instance {
//	//	ID:     "i-0f70b335ea12b2853",
//	//	Region: "us-east-1",
//	//}
//	scope := types.ScanScope{
//		All:         false,
//		Regions:     []types.Region{
//			{
//				ID:   "us-east-2",
//				VPCs: []types.VPC{
//					{
//						ID:             "vpc-32ea7c59",
//						SecurityGroups: []types.SecurityGroup{
//							{
//								ID: "sg-4d6b853a",
//							},
//						},
//					},
//				},
//			},
//			{
//				ID:   "us-east-1",
//				VPCs: []types.VPC{
//					{
//						ID:             "vpc-0c41450ba658eed00",
//						SecurityGroups: nil,
//					},
//					{
//						ID:             "vpc-9ca32ce1",
//						SecurityGroups: []types.SecurityGroup{
//							{
//								ID: "sg-030918e8e73254d42",
//							},
//						},
//					},
//				},
//			},
//		},
//		ScanStopped: true,
//		IncludeTags: []*types.Tag{
//			{
//				Key: "Name",
//				Val: "diff-vpc",
//			},
//		},
//		ExcludeTags: nil,
//	}
//	instances, err := c.Discover(&scope)
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//
//	instance := instances[0]
//
//	rootVolume, err := c.GetInstanceRootVolume(instance)
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//	// create a snapshot of that vm
//	srcSnapshot, err := c.CreateSnapshot(rootVolume)
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//	if err := c.WaitForSnapshotReady(srcSnapshot); err != nil {
//		t.Fatalf("%v", err)
//	}
//	//copy the snapshot to the scanner region
//	cpySnapshot, err := c.CopySnapshot(srcSnapshot, "us-east-2")
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//	if err := c.WaitForSnapshotReady(cpySnapshot); err != nil {
//		t.Fatalf("%v", err)
//	}
//	// create the scanner job (vm) with a boot script
//	launchedInstance, err := c.LaunchInstance("ami-0568773882d492fc8", "xvdh", cpySnapshot)
//	if err != nil {
//		t.Fatalf("%v", err)
//	}
//
//	t.Logf("res: %v", launchedInstance.ID)
//}

func Test_createVPCFilters(t *testing.T) {
	var (
		vpcID = "vpc-1"
		sgID1 = "sg-1"
		sgID2 = "sg-2"

		vpcIDFilterName = vpcIDFilterName
		sgIDFilterName  = sgIDFilterName
	)

	type args struct {
		vpc types.VPC
	}
	tests := []struct {
		name string
		args args
		want []ec2types.Filter
	}{
		{
			name: "vpc with no security group",
			args: args{
				vpc: types.VPC{
					ID:             vpcID,
					SecurityGroups: nil,
				},
			},
			want: []ec2types.Filter{
				{
					Name:   &vpcIDFilterName,
					Values: []string{vpcID},
				},
			},
		},
		{
			name: "vpc with one security group",
			args: args{
				vpc: types.VPC{
					ID: vpcID,
					SecurityGroups: []types.SecurityGroup{
						{
							ID: sgID1,
						},
					},
				},
			},
			want: []ec2types.Filter{
				{
					Name:   &vpcIDFilterName,
					Values: []string{vpcID},
				},
				{
					Name:   &sgIDFilterName,
					Values: []string{sgID1},
				},
			},
		},
		{
			name: "vpc with two security groups",
			args: args{
				vpc: types.VPC{
					ID: vpcID,
					SecurityGroups: []types.SecurityGroup{
						{
							ID: sgID1,
						},
						{
							ID: sgID2,
						},
					},
				},
			},
			want: []ec2types.Filter{
				{
					Name:   &vpcIDFilterName,
					Values: []string{vpcID},
				},
				{
					Name:   &sgIDFilterName,
					Values: []string{sgID1, sgID2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createVPCFilters(tt.args.vpc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createVPCFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createInclusionTagsFilters(t *testing.T) {
	var (
		tagName       = "foo"
		filterTagName = "tag:" + tagName
		tagVal        = "bar"
	)

	type args struct {
		tags []*types.Tag
	}
	tests := []struct {
		name string
		args args
		want []ec2types.Filter
	}{
		{
			name: "no tags",
			args: args{
				tags: nil,
			},
			want: nil,
		},
		{
			name: "1 tag",
			args: args{
				tags: []*types.Tag{
					{
						Key: tagName,
						Val: tagVal,
					},
				},
			},
			want: []ec2types.Filter{
				{
					Name:   &filterTagName,
					Values: []string{tagVal},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createInclusionTagsFilters(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createInclusionTagsFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasExcludedTags(t *testing.T) {
	var (
		tagName1 = "foo1"
		tagName2 = "foo2"
		tagVal1  = "bar1"
		tagVal2  = "bar2"
	)

	type args struct {
		excludeTags  []*types.Tag
		instanceTags []ec2types.Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "instance has no tags",
			args: args{
				excludeTags: []*types.Tag{
					{
						Key: tagName1,
						Val: tagVal1,
					},
					{
						Key: "stam1",
						Val: "stam2",
					},
				},
				instanceTags: nil,
			},
			want: false,
		},
		{
			name: "empty excluded tags",
			args: args{
				excludeTags: nil,
				instanceTags: []ec2types.Tag{
					{
						Key:   &tagName1,
						Value: &tagVal1,
					},
					{
						Key:   &tagName2,
						Value: &tagVal2,
					},
				},
			},
			want: false,
		},
		{
			name: "instance has excluded tags",
			args: args{
				excludeTags: []*types.Tag{
					{
						Key: tagName1,
						Val: tagVal1,
					},
					{
						Key: "stam1",
						Val: "stam2",
					},
				},
				instanceTags: []ec2types.Tag{
					{
						Key:   &tagName1,
						Value: &tagVal1,
					},
					{
						Key:   &tagName2,
						Value: &tagVal2,
					},
				},
			},
			want: true,
		},
		{
			name: "instance does not have excluded tags",
			args: args{
				excludeTags: []*types.Tag{
					{
						Key: "stam1",
						Val: "stam2",
					},
					{
						Key: "stam3",
						Val: "stam4",
					},
				},
				instanceTags: []ec2types.Tag{
					{
						Key:   &tagName1,
						Value: &tagVal1,
					},
					{
						Key:   &tagName2,
						Value: &tagVal2,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasExcludeTags(tt.args.excludeTags, tt.args.instanceTags); got != tt.want {
				t.Errorf("hasExcludeTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInstancesFromDescribeInstancesOutput(t *testing.T) {
	type args struct {
		result      *ec2.DescribeInstancesOutput
		excludeTags []*types.Tag
		regionID    string
	}
	tests := []struct {
		name string
		args args
		want []types.Instance
	}{
		{
			name: "no reservations found",
			args: args{
				result: &ec2.DescribeInstancesOutput{
					Reservations: []ec2types.Reservation{},
				},
				excludeTags: nil,
				regionID:    "region-1",
			},
			want: nil,
		},
		{
			name: "no excluded tags",
			args: args{
				result: &ec2.DescribeInstancesOutput{
					Reservations: []ec2types.Reservation{
						{
							Instances: []ec2types.Instance{
								{
									InstanceId: utils.StringPtr("instance-1"),
									Tags: []ec2types.Tag{
										{
											Key:   utils.StringPtr("key-1"),
											Value: utils.StringPtr("val-1"),
										},
									},
								},
								{
									InstanceId: utils.StringPtr("instance-2"),
									Tags: []ec2types.Tag{
										{
											Key:   utils.StringPtr("key-2"),
											Value: utils.StringPtr("val-2"),
										},
									},
								},
							},
						},
						{
							Instances: []ec2types.Instance{
								{
									InstanceId: utils.StringPtr("instance-3"),
								},
							},
						},
					},
				},
				excludeTags: nil,
				regionID:    "region-1",
			},
			want: []types.Instance{
				{
					ID:     "instance-1",
					Region: "region-1",
				},
				{
					ID:     "instance-2",
					Region: "region-1",
				},
				{
					ID:     "instance-3",
					Region: "region-1",
				},
			},
		},
		{
			name: "one excluded instance",
			args: args{
				result: &ec2.DescribeInstancesOutput{
					Reservations: []ec2types.Reservation{
						{
							Instances: []ec2types.Instance{
								{
									InstanceId: utils.StringPtr("instance-1"),
									Tags: []ec2types.Tag{
										{
											Key:   utils.StringPtr("key-1"),
											Value: utils.StringPtr("val-1"),
										},
									},
								},
								{
									InstanceId: utils.StringPtr("instance-2"),
									Tags: []ec2types.Tag{
										{
											Key:   utils.StringPtr("key-2"),
											Value: utils.StringPtr("val-2"),
										},
									},
								},
							},
						},
						{
							Instances: []ec2types.Instance{
								{
									InstanceId: utils.StringPtr("instance-3"),
								},
							},
						},
					},
				},
				excludeTags: []*types.Tag{
					{
						Key: "key-1",
						Val: "val-1",
					},
				},
				regionID: "region-1",
			},
			want: []types.Instance{
				{
					ID:     "instance-2",
					Region: "region-1",
				},
				{
					ID:     "instance-3",
					Region: "region-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInstancesFromDescribeInstancesOutput(tt.args.result, tt.args.excludeTags, tt.args.regionID)

			sort.Slice(got, func(i, j int) bool {
				return got[i].ID > got[j].ID
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].ID > tt.want[j].ID
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInstancesFromDescribeInstancesOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
