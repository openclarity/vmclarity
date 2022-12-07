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
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func TestGetTargetsTargetIDScanresults(t *testing.T) {
	restServer := createTestRestServer()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().List(gomock.Any(), gomock.Any()).Return([]models.ScanResults{}, nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Get("/targets/1/scanresults?page=1&pageSize=1").Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}

func TestPostTargetsTargetIDScanresults(t *testing.T) {
	restServer := createTestRestServer()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.ScanResultsSummary{}, nil)
	restServer.RegisterHandlers(mockHandler)

	scanResID := "testID"
	newScanResults := models.ScanResults{
		Id: &scanResID,
		Sboms: &models.SbomScan{
			Packages: &[]models.Package{},
		},
		Vulnerabilities: &models.VulnerabilityScan{
			Vulnerabilities: &[]models.Vulnerability{},
		},
		Malwares: &models.MalwareScan{
			Malwares: &[]models.MalwareInfo{},
		},
		Misconfigurations: &models.MisconfigurationScan{
			Misconfigurations: &[]models.MisconfigurationInfo{},
		},
		Secrets: &models.SecretScan{
			Secrets: &[]models.SecretInfo{},
		},
		Rootkits: &models.RootkitScan{
			Rootkits: &[]models.RootkitInfo{},
		},
		Exploits: &models.ExploitScan{
			Exploits: &[]models.ExploitInfo{},
		},
	}
	result := testutil.NewRequest().Post("/targets/1/scanresults").WithJsonBody(newScanResults).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusCreated, result.Code())
}

func TestGetTargetsTargetIDScanresultsScanID(t *testing.T) {
	restServer := createTestRestServer()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.ScanResultsSummary{}, nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Get("/targets/1/scanresults/1").Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}

func TestPutTargetsTargetIDScanresultsScanID(t *testing.T) {
	restServer := createTestRestServer()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.ScanResultsSummary{}, nil)
	restServer.RegisterHandlers(mockHandler)

	scanResID := "testID"
	newScanResults := models.ScanResults{
		Id: &scanResID,
		Sboms: &models.SbomScan{
			Packages: &[]models.Package{},
		},
		Vulnerabilities: &models.VulnerabilityScan{
			Vulnerabilities: &[]models.Vulnerability{},
		},
		Malwares: &models.MalwareScan{
			Malwares: &[]models.MalwareInfo{},
		},
		Misconfigurations: &models.MisconfigurationScan{
			Misconfigurations: &[]models.MisconfigurationInfo{},
		},
		Secrets: &models.SecretScan{
			Secrets: &[]models.SecretInfo{},
		},
		Rootkits: &models.RootkitScan{
			Rootkits: &[]models.RootkitInfo{},
		},
		Exploits: &models.ExploitScan{
			Exploits: &[]models.ExploitInfo{},
		},
	}
	result := testutil.NewRequest().Put("/targets/1/scanresults/1").WithJsonBody(newScanResults).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}
