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

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func (s *ServerImpl) GetTargetsTargetIDScanResults(
	ctx echo.Context,
	targetID models.TargetID,
	params models.GetTargetsTargetIDScanResultsParams,
) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	targets, err := s.dbHandler.ScanResultsTable().ListScanResults(targetID, params)
	if err != nil {
		// TODO check errors and for status code
		return sendError(ctx, http.StatusNotFound, oopsMsg)
	}
	return sendResponse(ctx, http.StatusOK, targets)
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

	s.lock.Lock()
	defer s.lock.Unlock()

	newScanResults := database.CreateDBScanResultsFromModel(&scanResults)
	scanResultsSummary, err := s.dbHandler.ScanResultsTable().CreateScanResults(targetID, newScanResults)
	if err != nil {
		// TODO check errors and for status code
		return sendError(ctx, http.StatusConflict, oopsMsg)
	}
	return sendResponse(ctx, http.StatusCreated, scanResultsSummary)
}

//nolint:cyclop
func (s *ServerImpl) GetTargetsTargetIDScanResultsScanID(
	ctx echo.Context,
	targetID models.TargetID,
	scanID models.ScanID,
	params models.GetTargetsTargetIDScanResultsScanIDParams,
) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var result interface{}
	var err error

	if params.ScanType == nil {
		result, err = s.dbHandler.ScanResultsTable().GetScanResultsSummary(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
		return sendResponse(ctx, http.StatusOK, result)
	}
	switch *params.ScanType {
	case models.SBOM:
		result, err = s.dbHandler.ScanResultsTable().GetSBOM(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.VULNERABILITY:
		result, err = s.dbHandler.ScanResultsTable().GetVulnerabilities(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.MALWARE:
		result, err = s.dbHandler.ScanResultsTable().GetMalwares(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.ROOTKIT:
		result, err = s.dbHandler.ScanResultsTable().GetRootkits(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.SECRET:
		result, err = s.dbHandler.ScanResultsTable().GetSecrets(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.MISCONFIGURATION:
		result, err = s.dbHandler.ScanResultsTable().GetMisconfigurations(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	case models.EXPLOIT:
		result, err = s.dbHandler.ScanResultsTable().GetExploits(targetID, scanID)
		if err != nil {
			// TODO check errors and for status code
			return sendError(ctx, http.StatusNotFound, oopsMsg)
		}
	default:
		return sendError(ctx, http.StatusBadRequest, oopsMsg)
	}
	return sendResponse(ctx, http.StatusOK, result)
}

func (s *ServerImpl) PutTargetsTargetIDScanResultsScanID(
	ctx echo.Context,
	targetID models.TargetID,
	scanID models.ScanID,
) error {
	var scanResults models.ScanResults
	err := ctx.Bind(&scanResults)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	newScanResults := database.CreateDBScanResultsFromModel(&scanResults)
	scanResultsSummary, err := s.dbHandler.ScanResultsTable().UpdateScanResults(targetID, scanID, newScanResults)
	if err != nil {
		// TODO check errors and for status code
		return sendError(ctx, http.StatusConflict, oopsMsg)
	}
	return sendResponse(ctx, http.StatusOK, scanResultsSummary)
}
