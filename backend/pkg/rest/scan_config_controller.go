package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/dbtorest"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/resttodb"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func (s *ServerImpl) GetScanConfigs(ctx echo.Context, params models.GetScanConfigsParams) error {
	dbScanConfigs, total, err := s.dbHandler.ScanConfigsTable().GetScanConfigsAndTotal(resttodb.ConvertGetScanConfigsParams(params))
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan configs from db: %v", err))
	}

	converted, err := dbtorest.ConvertScanConfigs(dbScanConfigs, total)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan configs: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PostScanConfigs(ctx echo.Context) error {
	var scanConfig models.ScanConfig
	err := ctx.Bind(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check if scan config with that name already exists.
	sc, exist, err := s.dbHandler.ScanConfigsTable().CheckExist(*scanConfig.Name)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to check scan config in db. name=%v: %v", *scanConfig.Name, err))
	}
	if exist {
		return sendResponse(ctx, http.StatusConflict, &sc)
	}

	convertedDB, err := resttodb.ConvertScanConfig(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config: %v", err))
	}
	createdScanConfig, err := s.dbHandler.ScanConfigsTable().CreateScanConfig(convertedDB)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create scan config in db: %v", err))
	}

	converted, err := dbtorest.ConvertScanConfig(createdScanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config: %v", err))
	}
	return sendResponse(ctx, http.StatusCreated, converted)
}

func (s *ServerImpl) DeleteScanConfigsScanConfigID(ctx echo.Context, scanConfigID models.ScanConfigID) error {
	success := models.Success{
		Message: utils.StringPtr(fmt.Sprintf("scan config %v deleted", scanConfigID)),
	}

	if err := s.dbHandler.ScanConfigsTable().DeleteScanConfig(scanConfigID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusNoContent, &success)
}

func (s *ServerImpl) GetScanConfigsScanConfigID(ctx echo.Context, scanConfigID models.ScanConfigID) error {
	sc, err := s.dbHandler.ScanConfigsTable().GetScanConfig(scanConfigID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan config from db. scanConfigID=%v: %v", scanConfigID, err))
	}

	converted, err := dbtorest.ConvertScanConfig(sc)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config. scanConfigID=%v: %v", scanConfigID, err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PatchScanConfigsScanConfigID(ctx echo.Context, scanConfigID models.ScanConfigID) error {
	var scanConfig models.ScanConfig
	err := ctx.Bind(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that a scan config with that id exists.
	_, err = s.dbHandler.ScanConfigsTable().GetScanConfig(scanConfigID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("scan config was not found. scanConfigID=%v: %v", scanConfigID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan config from db. scanConfigID=%v: %v", scanConfigID, err))
	}

	convertedDB, err := resttodb.ConvertScanConfig(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config: %v", err))
	}
	updatedScanConfig, err := s.dbHandler.ScanConfigsTable().UpdateScanConfig(convertedDB, scanConfigID)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update scan config in db. scanConfigID=%v: %v", scanConfigID, err))
	}

	converted, err := dbtorest.ConvertScanConfig(updatedScanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config. scanConfigID=%v: %v", scanConfigID, err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PutScanConfigsScanConfigID(ctx echo.Context, scanConfigID models.ScanConfigID) error {
	var scanConfig models.ScanConfig
	err := ctx.Bind(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that a scan config with that id exists.
	_, err = s.dbHandler.ScanConfigsTable().GetScanConfig(scanConfigID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("scan config was not found. scanConfigID=%v: %v", scanConfigID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan config from db. scanConfigID=%v: %v", scanConfigID, err))
	}

	convertedDB, err := resttodb.ConvertScanConfig(&scanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config: %v", err))
	}
	updatedScanConfig, err := s.dbHandler.ScanConfigsTable().SaveScanConfig(convertedDB, scanConfigID)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update scan config in db. scanConfigID=%v: %v", scanConfigID, err))
	}

	converted, err := dbtorest.ConvertScanConfig(updatedScanConfig)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan config. scanConfigID=%v: %v", scanConfigID, err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}
