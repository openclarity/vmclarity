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

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func (s *ServerImpl) GetTargetsTargetIDScanResults(
	ctx echo.Context,
	targetID models.TargetID,
	params models.GetTargetsTargetIDScanResultsParams,
) error {
	results, err := s.dbHandler.ScanResultsTable().ListScanResults(targetID, params)
	if err != nil {
		// TODO check errors for status code
		log.Errorf("%v", err)
		return sendError(ctx, http.StatusInternalServerError, oops)
	}
	resultsModel := []models.ScanResults{}
	for _, result := range results {
		result := result
		resultModel := database.CreateModelScanResultsFromDB(&result)
		resultsModel = append(resultsModel, *resultModel)
	}
	return sendResponse(ctx, http.StatusOK, &resultsModel)
}

func (s *ServerImpl) PostTargetsTargetIDScanResults(
	ctx echo.Context,
	targetID models.TargetID,
) error {
	var scanResults models.ScanResults
	err := ctx.Bind(&scanResults)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	newScanResults := database.CreateDBScanResultsFromModel(&scanResults)
	scanResultsSummary, err := s.dbHandler.ScanResultsTable().CreateScanResults(targetID, newScanResults)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
	return sendResponse(ctx, http.StatusCreated, scanResultsSummary)
}

//nolint:cyclop
func (s *ServerImpl) GetTargetsTargetIDScanResultsScanID(
	ctx echo.Context,
	targetID models.TargetID,
	scanID models.ScanID,
) error {
	result, err := s.dbHandler.ScanResultsTable().GetScanResults(targetID, scanID)
	if err != nil {
		// TODO check errors for status code
		log.Errorf("%v", err)
		return sendError(ctx, http.StatusNotFound, oops)
	}
	return sendResponse(ctx, http.StatusOK, database.CreateModelScanResultsFromDB(result))
}

func (s *ServerImpl) PutTargetsTargetIDScanResultsScanID(
	ctx echo.Context,
	targetID models.TargetID,
	scanID models.ScanID,
) error {
	var scanResults models.ScanResults
	err := ctx.Bind(&scanResults)
	if err != nil {
		log.Errorf("%v", err)
		return sendError(ctx, http.StatusBadRequest, oops)
	}

	newScanResults := database.CreateDBScanResultsFromModel(&scanResults)
	scanResultsSummary, err := s.dbHandler.ScanResultsTable().UpdateScanResults(targetID, scanID, newScanResults)
	if err != nil {
		// TODO check errors for status code
		log.Errorf("%v", err)
		return sendError(ctx, http.StatusInternalServerError, oops)
	}
	return sendResponse(ctx, http.StatusOK, scanResultsSummary)
}
