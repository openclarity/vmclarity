package rest

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/dbtorest"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/resttodb"
)

func (s *ServerImpl) GetDiscoveryScopes(ctx echo.Context, params models.GetDiscoveryScopesParams) error {
	var dbScopes *database.Scopes
	if _, ok := os.LookupEnv("FAKE_SCOPES"); ok {
		dbScopes = getFakeScopes(params.Select, params.Filter)
	} else {
		var err error
		dbScopes, err = s.dbHandler.ScopesTable().GetScopes()
		if err != nil {
			return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scopes from db: %v", err))
		}
	}

	converted, err := dbtorest.ConvertScopes(dbScopes)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scopes: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func getFakeScopes(odataSelect, odataFilter *string) *database.Scopes {
	dbScopes := &database.Scopes{
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
		var regionData []database.AwsScopesRegion
		if len(matches) == 2 {
			filterRegion := matches[1]
			for _, scopesRegion := range dbScopes.AwsScopesRegions {
				if scopesRegion.RegionID == filterRegion {
					regionData = []database.AwsScopesRegion{scopesRegion}
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
		var regionData []database.AwsScopesRegion
		if len(matches) == 3 {
			filterRegion := matches[1]
			filterVpc := matches[2]
			for _, scopesRegion := range dbScopes.AwsScopesRegions {
				if scopesRegion.RegionID == filterRegion {
					for _, vpc := range scopesRegion.AwsRegionVpcs {
						if vpc.VpcID == filterVpc {
							regionData = []database.AwsScopesRegion{scopesRegion}
							regionData[0].AwsRegionVpcs = []database.AwsRegionVpc{vpc}
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

	convertedDB, err := resttodb.ConvertScopes(&scopes)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scopes: %v", err))
	}

	_, err = s.dbHandler.ScopesTable().SetScopes(convertedDB)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to set scopes in db: %v", err).Error())
	}

	return sendResponse(ctx, http.StatusOK, &scopes)
}
