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
	"github.com/openclarity/vmclarity/backend/pkg/common"
	databaseTypes "github.com/openclarity/vmclarity/backend/pkg/database/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func (s *ServerImpl) GetTargets(ctx echo.Context, params models.GetTargetsParams) error {
	dbTargets, err := s.dbHandler.TargetsTable().GetTargets(params)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get targets from db: %v", err))
	}
	//scanResults, _, err := s.dbHandler.ScanResultsTable().GetScanResultsAndTotal(database.GetScanResultsParams{})
	//if err != nil {
	//	return err
	//}
	//targetsIDScanCount := make(map[string]int)
	//targetsIDSummary := make(map[string][]byte)
	//for _, result := range scanResults {
	//	targetsIDScanCount[result.TargetID] = targetsIDScanCount[result.TargetID] + 1
	//	targetsIDSummary[result.TargetID] = result.Summary // Temporary override just for demo
	//}

	return sendResponse(ctx, http.StatusOK, dbTargets)
}

// nolint:cyclop
func (s *ServerImpl) PostTargets(ctx echo.Context) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	createdTarget, err := s.dbHandler.TargetsTable().CreateTarget(target)
	if err != nil {
		var conflictErr *common.ConflictError
		if errors.As(err, &conflictErr) {
			existResponse := &models.TargetExists{
				Message: utils.StringPtr(conflictErr.Reason),
				Target:  &createdTarget,
			}
			return sendResponse(ctx, http.StatusConflict, existResponse)
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create target in db: %v", err))
	}

	return sendResponse(ctx, http.StatusCreated, createdTarget)
}

func (s *ServerImpl) GetTargetsTargetID(ctx echo.Context, targetID models.TargetID, params models.GetTargetsTargetIDParams) error {
	target, err := s.dbHandler.TargetsTable().GetTarget(targetID, params)
	if err != nil {
		if errors.Is(err, databaseTypes.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("Target with ID %v not found", targetID))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get target from db. targetID=%v: %v", targetID, err))
	}

	return sendResponse(ctx, http.StatusOK, target)
}

func (s *ServerImpl) PutTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	target.Id = &targetID
	updatedTarget, err := s.dbHandler.TargetsTable().SaveTarget(target)
	if err != nil {
		if errors.Is(err, databaseTypes.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("Target with ID %v not found", targetID))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get target from db. targetID=%v: %v", targetID, err))
	}

	return sendResponse(ctx, http.StatusOK, updatedTarget)
}

func (s *ServerImpl) DeleteTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	success := models.Success{
		Message: utils.StringPtr(fmt.Sprintf("target %v deleted", targetID)),
	}

	if err := s.dbHandler.TargetsTable().DeleteTarget(targetID); err != nil {
		if errors.Is(err, databaseTypes.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("Target with ID %v not found", targetID))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to delete target from db. targetID=%v: %v", targetID, err))
	}

	return sendResponse(ctx, http.StatusNoContent, &success)
}
