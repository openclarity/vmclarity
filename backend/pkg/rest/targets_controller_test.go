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

func TestGetTargets(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockTargetTable := database.NewMockTargetsTable(mockCtrl)
	mockHandler.EXPECT().TargetsTable().Return(mockTargetTable)
	mockTargetTable.EXPECT().List(gomock.Any()).Return(&[]models.Target{}, nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Get(fmt.Sprintf("%s/targets?page=1&pageSize=1", baseURL)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
	var resultTarget []models.Target
	result.UnmarshalJsonToObject(&resultTarget)
	assert.Equal(t, []models.Target{}, resultTarget)

}

func TestPostTargets(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockTargetTable := database.NewMockTargetsTable(mockCtrl)
	mockHandler.EXPECT().TargetsTable().Return(mockTargetTable)
	mockTargetTable.EXPECT().Create(gomock.Any()).Return(&models.Target{}, nil)
	restServer.RegisterHandlers(mockHandler)

	targetID := testID
	targetType := models.TargetType("VM")
	scanResults := uint32(1)
	instanceName := "instance"
	instanceProvider := models.CloudProvider("AWS")
	location := "eu-central2"
	vmInfo := models.VMInfo{
		InstanceName:     &instanceName,
		InstanceProvider: &instanceProvider,
		Location:         &location,
	}
	targetInfo := &models.Target_TargetInfo{}
	if err := targetInfo.FromVMInfo(vmInfo); err != nil {
		t.Errorf("failed to create target info")
	}

	newTarget := models.Target{
		Id:          &targetID,
		ScanResults: &scanResults,
		TargetType:  &targetType,
		TargetInfo:  targetInfo,
	}
	result := testutil.NewRequest().Post(fmt.Sprintf("%s/targets", baseURL)).WithJsonBody(newTarget).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusCreated, result.Code())
	var resultTarget models.Target
	result.UnmarshalJsonToObject(&resultTarget)
	assert.Equal(t, models.Target{}, resultTarget)
}

func TestGetTargetTargetID(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockTargetTable := database.NewMockTargetsTable(mockCtrl)
	mockHandler.EXPECT().TargetsTable().Return(mockTargetTable)
	mockTargetTable.EXPECT().Get(gomock.Any()).Return(&models.Target{}, nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Get(fmt.Sprintf("%s/targets/1", baseURL)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
	var resultTarget models.Target
	result.UnmarshalJsonToObject(&resultTarget)
	assert.Equal(t, models.Target{}, resultTarget)
}

func TestPutTargetTargetID(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockTargetTable := database.NewMockTargetsTable(mockCtrl)
	mockHandler.EXPECT().TargetsTable().Return(mockTargetTable)
	mockTargetTable.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&models.Target{}, nil)
	restServer.RegisterHandlers(mockHandler)

	targetID := testID
	targetType := models.TargetType("VM")
	scanResults := uint32(1)
	instanceName := "instance"
	instanceProvider := models.CloudProvider("AWS")
	location := "eu-central2"
	vmInfo := models.VMInfo{
		InstanceName:     &instanceName,
		InstanceProvider: &instanceProvider,
		Location:         &location,
	}
	targetInfo := &models.Target_TargetInfo{}
	if err := targetInfo.FromVMInfo(vmInfo); err != nil {
		t.Errorf("failed to create target info")
	}

	newTarget := models.Target{
		Id:          &targetID,
		ScanResults: &scanResults,
		TargetType:  &targetType,
		TargetInfo:  targetInfo,
	}
	result := testutil.NewRequest().Put(fmt.Sprintf("%s/targets/1", baseURL)).WithJsonBody(newTarget).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
	var resultTarget models.Target
	result.UnmarshalJsonToObject(&resultTarget)
	assert.Equal(t, models.Target{}, resultTarget)
}

func TestDeleteTargetTargetID(t *testing.T) {
	restServer := createTestRestServer(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHandler := database.NewMockDatabase(mockCtrl)
	mockTargetTable := database.NewMockTargetsTable(mockCtrl)
	mockHandler.EXPECT().TargetsTable().Return(mockTargetTable)
	mockTargetTable.EXPECT().Delete(gomock.Any()).Return(nil)
	restServer.RegisterHandlers(mockHandler)

	result := testutil.NewRequest().Delete(fmt.Sprintf("%s/targets/1", baseURL)).Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusNoContent, result.Code())
}
