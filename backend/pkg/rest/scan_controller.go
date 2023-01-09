package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/db_to_rest"
	"github.com/openclarity/vmclarity/backend/pkg/rest/convert/rest_to_db"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"gorm.io/gorm"
)

func (s *ServerImpl) GetScans(ctx echo.Context, params models.GetScansParams) error {
	dbScans, total, err := s.dbHandler.ScansTable().GetScansAndTotal(params)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scans from db: %v", err))
	}

	converted, err := db_to_rest.ConvertScans(dbScans, total)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scans: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PostScans(ctx echo.Context) error {
	var scan models.Scan
	err := ctx.Bind(&scan)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check if scan already exists.
	sc, exist, err := s.dbHandler.ScansTable().CheckExist(*scan.ScanConfigId, *scan.StartTime)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to check scan in db. id=%v: %v", *scan.Id, err))
	}
	if exist {
		return sendResponse(ctx, http.StatusConflict, &sc)
	}

	convertedDB, err := rest_to_db.ConvertScan(&scan)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan: %v", err))
	}
	createdScan, err := s.dbHandler.ScansTable().CreateScan(convertedDB)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to create scan in db: %v", err))
	}

	converted, err := db_to_rest.ConvertScan(createdScan)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan: %v", err))
	}
	return sendResponse(ctx, http.StatusCreated, converted)
}

func (s *ServerImpl) DeleteScansScanID(ctx echo.Context, scanID models.ScanID) error {
	success := models.Success{
		Message: utils.StringPtr(fmt.Sprintf("scan %v deleted", scanID)),
	}

	if err := s.dbHandler.ScansTable().DeleteScan(scanID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusNoContent, &success)
}

func (s *ServerImpl) GetScansScanID(ctx echo.Context, scanID models.ScanID) error {
	scan, err := s.dbHandler.ScansTable().GetScan(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan from db. id=%v: %v", scanID, err))
	}

	converted, err := db_to_rest.ConvertScan(scan)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}

func (s *ServerImpl) PatchScansScanID(ctx echo.Context, scanID models.ScanID) error {
	//TODO implement me
	panic("implement me")
}

func (s *ServerImpl) PutScansScanID(ctx echo.Context, scanID models.ScanID) error {
	var scan models.Scan
	err := ctx.Bind(&scan)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// check that a scan with that id exists.
	_, err = s.dbHandler.ScansTable().GetScan(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(ctx, http.StatusNotFound, fmt.Sprintf("scan was not found in db. scanID=%v: %v", scanID, err))
		}
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scan from db. scanID=%v: %v", scanID, err))
	}

	convertedDB, err := rest_to_db.ConvertScan(&scan)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan: %v", err))
	}
	updatedScan, err := s.dbHandler.ScansTable().UpdateScan(convertedDB, scanID)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to update scan in db. scanID=%v: %v", scanID, err))
	}

	converted, err := db_to_rest.ConvertScan(updatedScan)
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to convert scan: %v", err))
	}
	return sendResponse(ctx, http.StatusOK, converted)
}
