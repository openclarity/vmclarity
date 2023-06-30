// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
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
	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/common"
	databaseTypes "github.com/openclarity/vmclarity/backend/pkg/database/types"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func (s *ServerImpl) GetAssetScans(ctx echo.Context, params models.GetAssetScansParams) error {
	dbAssetScans, err := s.dbHandler.AssetScansTable().GetAssetScans(params)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scans results from db: %v", err))
	}

	return sendResponse(ctx, http.StatusOK, dbAssetScans)
}

func (s *ServerImpl) PostAssetScans(ctx echo.Context) error {
	var assetScan models.AssetScan
	err := ctx.Bind(&assetScan)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	createdAssetScan, err := s.dbHandler.AssetScansTable().CreateAssetScan(assetScan)
	if err != nil {
		var conflictErr *common.ConflictError
		if errors.As(err, &conflictErr) {
			existResponse := &models.AssetScanExists{
				Message:   utils.PointerTo(conflictErr.Reason),
				AssetScan: &createdAssetScan,
			}
			return sendResponse(ctx, http.StatusConflict, existResponse)
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create asset scan in db: %v", err))
	}

	return sendResponse(ctx, http.StatusCreated, createdAssetScan)
}

func (s *ServerImpl) GetAssetScansAssetScanID(ctx echo.Context, assetScanID models.AssetScanID, params models.GetAssetScansAssetScanIDParams) error {
	dbAssetScan, err := s.dbHandler.AssetScansTable().GetAssetScan(assetScanID, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get asset scan from db. assetScanID=%v: %v", assetScanID, err))
	}

	return sendResponse(ctx, http.StatusOK, dbAssetScan)
}

// nolint:cyclop
func (s *ServerImpl) PatchAssetScansAssetScanID(ctx echo.Context, assetScanID models.AssetScanID, params models.PatchAssetScansAssetScanIDParams) error {
	// TODO: check that the provided scan and asset IDs are valid
	var assetScan models.AssetScan
	err := ctx.Bind(&assetScan)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that an asset scan with that id exists.
	_, err = s.dbHandler.AssetScansTable().GetAssetScan(assetScanID, models.GetAssetScansAssetScanIDParams{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("asset scan was not found. assetScanID=%v: %v", assetScanID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get asset scan. assetScanID=%v: %v", assetScanID, err))
	}

	// PATCH request might not contain the ID in the body, so set it from
	// the URL field so that the DB layer knows which object is being updated.
	if assetScan.Id != nil && *assetScan.Id != assetScanID {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("id in body %s does not match object %s to be updated", *assetScan.Id, assetScanID))
	}
	assetScan.Id = &assetScanID

	updatedAssetScan, err := s.dbHandler.AssetScansTable().UpdateAssetScan(assetScan, params)
	if err != nil {
		var validationErr *common.BadRequestError
		var conflictErr *common.ConflictError
		var preconditionFailedErr *databaseTypes.PreconditionFailedError
		switch true {
		case errors.As(err, &conflictErr):
			existResponse := &models.AssetScanExists{
				Message:   utils.PointerTo(conflictErr.Reason),
				AssetScan: &updatedAssetScan,
			}
			return sendResponse(ctx, http.StatusConflict, existResponse)
		case errors.As(err, &validationErr):
			return sendError(ctx, http.StatusBadRequest, err.Error())
		case errors.As(err, &preconditionFailedErr):
			return sendError(ctx, http.StatusPreconditionFailed, err.Error())
		default:
			return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update asset scan in db. assetScanID=%v: %v", assetScanID, err))
		}
	}

	return sendResponse(ctx, http.StatusOK, updatedAssetScan)
}

// nolint:cyclop
func (s *ServerImpl) PutAssetScansAssetScanID(ctx echo.Context, assetScanID models.AssetScanID, params models.PutAssetScansAssetScanIDParams) error {
	// TODO: check that the provided scan and asset IDs are valid
	var assetScan models.AssetScan
	err := ctx.Bind(&assetScan)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that an asset scan with that id exists.
	_, err = s.dbHandler.AssetScansTable().GetAssetScan(assetScanID, models.GetAssetScansAssetScanIDParams{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("asset scan was not found. assetScanID=%v: %v", assetScanID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get asset scan. assetScanID=%v: %v", assetScanID, err))
	}

	// PUT request might not contain the ID in the body, so set it from
	// the URL field so that the DB layer knows which object is being updated.
	if assetScan.Id != nil && *assetScan.Id != assetScanID {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("id in body %s does not match object %s to be updated", *assetScan.Id, assetScanID))
	}
	assetScan.Id = &assetScanID

	updatedAssetScan, err := s.dbHandler.AssetScansTable().SaveAssetScan(assetScan, params)
	if err != nil {
		var validationErr *common.BadRequestError
		var conflictErr *common.ConflictError
		var preconditionFailedErr *databaseTypes.PreconditionFailedError
		switch true {
		case errors.As(err, &conflictErr):
			existResponse := &models.AssetScanExists{
				Message:   utils.PointerTo(conflictErr.Reason),
				AssetScan: &updatedAssetScan,
			}
			return sendResponse(ctx, http.StatusConflict, existResponse)
		case errors.As(err, &validationErr):
			return sendError(ctx, http.StatusBadRequest, err.Error())
		case errors.As(err, &preconditionFailedErr):
			return sendError(ctx, http.StatusPreconditionFailed, err.Error())
		default:
			return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update asset scan in db. assetScanID=%v: %v", assetScanID, err))
		}
	}

	return sendResponse(ctx, http.StatusOK, updatedAssetScan)
}
