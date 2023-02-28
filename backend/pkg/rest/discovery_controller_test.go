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

package rest

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/backend/pkg/database"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func Test_getFakeScopes(t *testing.T) {
	allScopes := &database.Scopes{
		Type: "AwsScope",
		AwsScopesRegions: []database.AwsScopesRegion{
			{
				RegionID: "eu-central-1",
				AwsRegionVpcs: []database.AwsRegionVpc{
					{
						VpcID: "vpc-1-from-eu-central-1",
						AwsVpcSecurityGroups: []database.AwsVpcSecurityGroup{
							{
								GroupID: "sg-1-from-vpc-1-from-eu-central-1",
							},
						},
					},
					{
						VpcID: "vpc-2-from-eu-central-1",
						AwsVpcSecurityGroups: []database.AwsVpcSecurityGroup{
							{
								GroupID: "sg-2-from-vpc-1-from-eu-central-1",
							},
						},
					},
				},
			},
			{
				RegionID: "us-east-1",
				AwsRegionVpcs: []database.AwsRegionVpc{
					{
						VpcID: "vpc-1-from-us-east-1",
						AwsVpcSecurityGroups: []database.AwsVpcSecurityGroup{
							{
								GroupID: "sg-1-from-vpc-1-from-us-east-1",
							},
						},
					},
					{
						VpcID: "vpc-2-from-us-east-1",
						AwsVpcSecurityGroups: []database.AwsVpcSecurityGroup{
							{
								GroupID: "sg-1-from-vpc-2-from-us-east-1",
							},
							{
								GroupID: "sg-2-from-vpc-2-from-us-east-1",
							},
						},
					},
				},
			},
		},
	}

	type args struct {
		odataSelect *string
		odataFilter *string
	}
	tests := []struct {
		name string
		args args
		want *database.Scopes
	}{
		{
			name: "no filter and select",
			args: args{
				odataSelect: nil,
				odataFilter: nil,
			},
			want: allScopes,
		},
		{
			name: "only regions",
			args: args{
				odataSelect: utils.StringPtr("AwsScope.Regions"),
				odataFilter: nil,
			},
			want: &database.Scopes{
				Type: "AwsScope",
				AwsScopesRegions: []database.AwsScopesRegion{
					{
						RegionID: "eu-central-1",
					},
					{
						RegionID: "us-east-1",
					},
				},
			},
		},
		{
			name: "only vpcs by region",
			args: args{
				odataSelect: utils.StringPtr("AwsScope.Regions.Vpcs"),
				odataFilter: utils.StringPtr("AwsScope.Regions.ID eq 'eu-central-1'"),
			},
			want: &database.Scopes{
				Type: "AwsScope",
				AwsScopesRegions: []database.AwsScopesRegion{
					{
						RegionID: "eu-central-1",
						AwsRegionVpcs: []database.AwsRegionVpc{
							{
								VpcID: "vpc-1-from-eu-central-1",
							},
							{
								VpcID: "vpc-2-from-eu-central-1",
							},
						},
					},
				},
			},
		},
		{
			name: "security groups for a specific vpc and region",
			args: args{
				odataSelect: utils.StringPtr("AwsScope.Regions.Vpcs.securityGroups"),
				odataFilter: utils.StringPtr("AwsScope.Regions.ID eq 'us-east-1' and AwsScope.Regions.Vpcs.ID eq 'vpc-2-from-us-east-1'"),
			},
			want: &database.Scopes{
				Type: "AwsScope",
				AwsScopesRegions: []database.AwsScopesRegion{
					{
						RegionID: "us-east-1",
						AwsRegionVpcs: []database.AwsRegionVpc{
							{
								VpcID: "vpc-2-from-us-east-1",
								AwsVpcSecurityGroups: []database.AwsVpcSecurityGroup{
									{
										GroupID: "sg-1-from-vpc-2-from-us-east-1",
									},
									{
										GroupID: "sg-2-from-vpc-2-from-us-east-1",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFakeScopes(tt.args.odataSelect, tt.args.odataFilter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFakeScopes() = %v, want %v", got, tt.want)
			}
		})
	}
}
