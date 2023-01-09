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
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/db_to_rest"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/rest_to_db"
	"gorm.io/gorm"
)

func (s *ServerImpl) GetScanResults(ctx echo.Context, params models.GetScanResultsParams) error {
	dbScanResults, total, err := s.dbHandler.ScanResultsTable().GetScanResultsAndTotal(params)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scans results from db: %v", err))
	}

	converted, err := db_to_rest.ConvertScanResults(dbScanResults, total)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan results: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PostScanResults(ctx echo.Context) error {
	var scanResult models.TargetScanResult
	err := ctx.Bind(&scanResult)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check if scan results already exists for that target and scan.
	dbScanResult, exist, err := s.dbHandler.ScanResultsTable().CheckExist(scanResult.ScanId, scanResult.TargetId)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan result from db. scanID=%v, targetID=%v: %v", scanResult.ScanId, scanResult.TargetId, err))
	}
	if exist {
		return sendResponse(ctx, http.StatusConflict, &dbScanResult)
	}

	convertedDB, err := rest_to_db.ConvertScanResult(&scanResult)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan result: %v", err))
	}
	createdScanResult, err := s.dbHandler.ScanResultsTable().CreateScanResult(convertedDB)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create scan result in db: %v", err))
	}

	converted, err := db_to_rest.ConvertScanResult(createdScanResult)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan result: %v", err))
	}
	return sendResponse(ctx, http.StatusCreated, converted)
}

func (s *ServerImpl) GetScanResultsScanResultID(ctx echo.Context, scanResultID models.ScanResultID, params models.GetScanResultsScanResultIDParams) error {
	dbScanResult, err := s.dbHandler.ScanResultsTable().GetScanResult(scanResultID, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan result from db. scanResultID=%v: %v", scanResultID, err))
	}

	converted, err := db_to_rest.ConvertScanResult(dbScanResult)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan result: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PatchScanResultsScanResultID(ctx echo.Context, scanResultID models.ScanResultID) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerImpl) PutScanResultsScanResultID(ctx echo.Context, scanResultID models.ScanResultID) error {
	var scanResult models.TargetScanResult
	err := ctx.Bind(&scanResult)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that a scan result with that id exists.
	_, err = s.dbHandler.ScanResultsTable().GetScanResult(scanResultID, models.GetScanResultsScanResultIDParams{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("scan result was not found. scanResultID=%v: %v", scanResultID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan result. scanResultID=%v: %v", scanResultID, err))
	}

	convertedDB, err := rest_to_db.ConvertScanResult(&scanResult)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan result: %v", err))
	}
	updatedScanResult, err := s.dbHandler.ScanResultsTable().UpdateScanResult(convertedDB, scanResultID)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update scan result in db. scanResultID=%v: %v", scanResultID, err))
	}

	converted, err := db_to_rest.ConvertScanResult(updatedScanResult)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan result: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}
