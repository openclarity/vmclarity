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
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database/gorm"
)

func (s *ServerImpl) GetDiscoveryScopes(ctx echo.Context, params models.GetDiscoveryScopesParams) error {
	var dbScopes *gorm.Scopes
	if _, ok := os.LookupEnv("FAKE_SCOPES"); ok {
		dbScopes = getFakeScopes(params.Select, params.Filter)
	} else {
		//var err error
		//dbScopes, err = s.dbHandler.ScopesTable().GetScopes()
		//if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scopes from db: not implemented"))
		//}
	}

	converted, err := gorm.ConvertToRestScopes(dbScopes)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scopes: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func getFakeScopes(odataSelect, odataFilter *string) *gorm.Scopes {
	dbScopes := &gorm.Scopes{
		Type: "AwsScope",
		AwsScopesRegions: []gorm.AwsScopesRegion{
			{
				RegionID: "eu-central-1",
				AwsRegionVpcs: []gorm.AwsRegionVpc{
					{
						VpcID: "vpc-1-from-eu-central-1",
						AwsVpcSecurityGroups: []gorm.AwsVpcSecurityGroup{
							{
								GroupID: "sg-1-from-vpc-1-from-eu-central-1",
							},
						},
					},
					{
						VpcID: "vpc-2-from-eu-central-1",
						AwsVpcSecurityGroups: []gorm.AwsVpcSecurityGroup{
							{
								GroupID: "sg-2-from-vpc-1-from-eu-central-1",
							},
						},
					},
				},
			},
			{
				RegionID: "us-east-1",
				AwsRegionVpcs: []gorm.AwsRegionVpc{
					{
						VpcID: "vpc-1-from-us-east-1",
						AwsVpcSecurityGroups: []gorm.AwsVpcSecurityGroup{
							{
								GroupID: "sg-1-from-vpc-1-from-us-east-1",
							},
						},
					},
					{
						VpcID: "vpc-2-from-us-east-1",
						AwsVpcSecurityGroups: []gorm.AwsVpcSecurityGroup{
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

	if odataSelect != nil {
		fmt.Printf("\nodataSelect=%v\n", *odataSelect)
	} else {
		fmt.Printf("odataSelect is empty\n\n")
	}
	if odataFilter != nil {
		fmt.Printf("\nodataFilter=%+v\n\n", *odataFilter)
	} else {
		fmt.Printf("\nodataFilter is empty\n\n")
	}

	if odataSelect != nil && *odataSelect == "AwsScope.Regions" && (odataFilter == nil || *odataFilter == "") {
		// use case 1, only regions: discovery/scopes?$select=AwsScope.Regions
		fmt.Printf("\nuse case 1, only regions\n")
		dbScopes.AwsScopesRegions[0].AwsRegionVpcs = nil
		dbScopes.AwsScopesRegions[1].AwsRegionVpcs = nil
	} else if odataSelect != nil && *odataSelect == "AwsScope.Regions.Vpcs" && odataFilter != nil {
		// use case 2 only vpcs for a specific region: /discovery/scopes?$select=AwsScope.Regions.Vpcs&$filter=AwsScope.Regions.ID eq 'specific-region'
		fmt.Printf("\nuse case 2 only vpcs for a specific region\n")
		r := regexp.MustCompile(`AwsScope\.Regions\.ID eq '(?P<First>[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*)'`)
		matches := r.FindStringSubmatch(*odataFilter)
		var regionData []gorm.AwsScopesRegion
		if len(matches) == 2 {
			filterRegion := matches[1]
			for _, scopesRegion := range dbScopes.AwsScopesRegions {
				if scopesRegion.RegionID == filterRegion {
					regionData = []gorm.AwsScopesRegion{scopesRegion}
					for i, _ := range regionData[0].AwsRegionVpcs {
						regionData[0].AwsRegionVpcs[i].AwsVpcSecurityGroups = nil
					}
					break
				}
			}
		}
		if len(regionData) == 0 {
			dbScopes.AwsScopesRegions = nil
		} else {
			dbScopes.AwsScopesRegions = regionData
		}
	} else if odataSelect != nil && *odataSelect == "AwsScope.Regions.Vpcs.securityGroups" && odataFilter != nil {
		// use case 3 only security groups for a specific region and vpc: /discovery/scopes?$select=AwsScope.Regions.Vpcs.securityGroups&$filter=AwsScope.Regions.ID eq 'specific-region' and AwsScope.Regions.Vpcs.ID eq 'specific vpc'
		fmt.Printf("\nuse case 3 only security groups for a specific region and vpc\n")
		r := regexp.MustCompile(`AwsScope\.Regions\.ID eq '(?P<First>[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*)' and AwsScope\.Regions\.Vpcs.ID eq '(?P<Second>[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*[a-z0-9]*[-]*)'`)
		matches := r.FindStringSubmatch(*odataFilter)
		var regionData []gorm.AwsScopesRegion
		if len(matches) == 3 {
			filterRegion := matches[1]
			filterVpc := matches[2]
			for _, scopesRegion := range dbScopes.AwsScopesRegions {
				if scopesRegion.RegionID == filterRegion {
					for _, vpc := range scopesRegion.AwsRegionVpcs {
						if vpc.VpcID == filterVpc {
							regionData = []gorm.AwsScopesRegion{scopesRegion}
							regionData[0].AwsRegionVpcs = []gorm.AwsRegionVpc{vpc}
							break
						}
					}

				}
				if len(regionData) > 0 {
					break
				}
			}
		}
		if len(regionData) == 0 {
			dbScopes.AwsScopesRegions = nil
		} else {
			dbScopes.AwsScopesRegions = regionData
		}
	}

	return dbScopes
}

func (s *ServerImpl) PutDiscoveryScopes(ctx echo.Context) error {
	var scopes models.ScopeType
	err := ctx.Bind(&scopes)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Errorf("failed to bind request: %v", err).Error())
	}

	_, err = s.dbHandler.ScopesTable().SetScopes(scopes)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to set scopes in db: %v", err).Error())
	}

	return sendResponse(ctx, http.StatusOK, &scopes)
}
