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

func (s *ServerImpl) GetTargets(ctx echo.Context, params models.GetTargetsParams) error {
	targets, err := s.dbHandler.TargetsTable().ListTargets(params)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, targets)
}

func (s *ServerImpl) PostTargets(ctx echo.Context) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	newTarget := database.CreateDBTargetFromModel(&target)
	createdTarget, err := s.dbHandler.TargetsTable().CreateTarget(newTarget)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
	return sendResponse(ctx, http.StatusCreated, createdTarget)
}

func (s *ServerImpl) GetTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	targets, err := s.dbHandler.TargetsTable().GetTarget(targetID)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusNotFound, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, targets)
}

func (s *ServerImpl) PutTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	newTarget := database.CreateDBTargetFromModel(&target)
	updatedTarget, err := s.dbHandler.TargetsTable().UpdateTarget(newTarget, targetID)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, updatedTarget)
}

func (s *ServerImpl) DeleteTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	err := s.dbHandler.TargetsTable().DeleteTarget(targetID)
	if err != nil {
		// TODO check errors for status code
		return sendError(ctx, http.StatusNotFound, err.Error())
	}
	return sendResponse(ctx, http.StatusNoContent, "deleted")
}
