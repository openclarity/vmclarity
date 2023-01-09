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

	echo "github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/db_to_rest"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/rest_to_db"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"gorm.io/gorm"
)

func (s *ServerImpl) GetTargets(ctx echo.Context, params models.GetTargetsParams) error {
	dbTargets, total, err := s.dbHandler.TargetsTable().GetTargetsAndTotal(params)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get targets from db: %v", err))
	}

	converted, err := db_to_rest.ConvertTargets(dbTargets, total)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert targets: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PostTargets(ctx echo.Context) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	disc, err := target.TargetInfo.Discriminator()
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}
	switch disc {
	case "VMInfo":
		vminfo, err := target.TargetInfo.AsVMInfo()
		targetFromDB, exists, err := s.dbHandler.TargetsTable().CheckVMInfoExists(*vminfo.InstanceID, *vminfo.Location)
		if err != nil {
			return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to check vminfo existence: %v", err))
		}
		if exists {
			converted, err := db_to_rest.ConvertTarget(targetFromDB)
			if err != nil {
				return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
			}
			return sendResponse(ctx, http.StatusConflict, converted)
		}
	default:
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("unknown target type %v", disc))
	}

	convertedDB, err := rest_to_db.ConvertTarget(&target)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
	}
	createdTarget, err := s.dbHandler.TargetsTable().CreateTarget(convertedDB)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create target in db: %v", err))
	}

	converted, err := db_to_rest.ConvertTarget(createdTarget)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
	}
	return sendResponse(ctx, http.StatusCreated, converted)
}

func (s *ServerImpl) GetTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	target, err := s.dbHandler.TargetsTable().GetTarget(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to get target from db. targetID=%v: %v", targetID, err).Error())
	}

	converted, err := db_to_rest.ConvertTarget(target)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PutTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	var target models.Target
	err := ctx.Bind(&target)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Errorf("failed to bind request: %v", err).Error())
	}

	_, err = s.dbHandler.TargetsTable().GetTarget(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("target with was not found in db. targetID=%v", targetID))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get target from db: %v", err))
	}

	convertedDB, err := rest_to_db.ConvertTarget(&target)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
	}
	updatedTarget, err := s.dbHandler.TargetsTable().UpdateTarget(convertedDB, targetID)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Errorf("failed to update target in db. targetID=%v: %v", targetID, err).Error())
	}

	converted, err := db_to_rest.ConvertTarget(updatedTarget)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert target: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) DeleteTargetsTargetID(ctx echo.Context, targetID models.TargetID) error {
	success := models.Success{
		Message: utils.StringPtr(fmt.Sprintf("target %v deleted", targetID)),
	}

	if err := s.dbHandler.TargetsTable().DeleteTarget(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusNoContent, &success)
}
