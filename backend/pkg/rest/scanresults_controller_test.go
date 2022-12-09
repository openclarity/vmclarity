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
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

var scanResultsPath = fmt.Sprintf("%s/targets/1/scanResults", baseURL)

func TestGetTargetsTargetIDScanResults(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().List(gomock.Any(), gomock.Any()).Return(&[]models.ScanResults{}, nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Get(fmt.Sprintf("%s?page=1&pageSize=1", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}

func TestPostTargetsTargetIDScanResults(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&models.ScanResultsSummary{}, nil)
	restServer.RegisterHandlers(mockHandler)

	scanResID := testID
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
	result := testutil.NewRequest().Post(scanResultsPath).WithJsonBody(newScanResults).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusCreated, result.Code())
}

func TestGetTargetsTargetIDScanResultsScanID(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable).AnyTimes()
	restServer.RegisterHandlers(mockHandler)

	gomock.InOrder(
		mockScanResultsTable.EXPECT().GetSummary(gomock.Any(), gomock.Any()).Return(&models.ScanResultsSummary{}, nil),
		mockScanResultsTable.EXPECT().GetSBOM(gomock.Any(), gomock.Any()).Return(&models.SbomScan{}, nil),
		mockScanResultsTable.EXPECT().GetVulnerabilities(gomock.Any(), gomock.Any()).Return(&models.VulnerabilityScan{}, nil),
		mockScanResultsTable.EXPECT().GetMalwares(gomock.Any(), gomock.Any()).Return(&models.MalwareScan{}, nil),
		mockScanResultsTable.EXPECT().GetRootkits(gomock.Any(), gomock.Any()).Return(&models.RootkitScan{}, nil),
		mockScanResultsTable.EXPECT().GetSecrets(gomock.Any(), gomock.Any()).Return(&models.SecretScan{}, nil),
		mockScanResultsTable.EXPECT().GetMisconfigurations(gomock.Any(), gomock.Any()).Return(&models.MisconfigurationScan{}, nil),
		mockScanResultsTable.EXPECT().GetExploits(gomock.Any(), gomock.Any()).Return(&models.ExploitScan{}, nil),
	)

	result := testutil.NewRequest().Get(fmt.Sprintf("%s/1", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=SBOM", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=VULNERABILITY", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=MALWARE", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=ROOTKIT", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=SECRET", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=MISCONFIGURATION", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())

	result = testutil.NewRequest().Get(fmt.Sprintf("%s/1?scanType=EXPLOIT", scanResultsPath)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}

func TestPutTargetsTargetIDScanResultsScanID(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockScanResultsTable := database.NewMockScanResultsTable(mockCtrl)
	mockHandler.EXPECT().ScanResultsTable().Return(mockScanResultsTable)
	mockScanResultsTable.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.ScanResultsSummary{}, nil)
	restServer.RegisterHandlers(mockHandler)

	scanResID := testID
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
	result := testutil.NewRequest().Put(fmt.Sprintf("%s/1", scanResultsPath)).WithJsonBody(newScanResults).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}
