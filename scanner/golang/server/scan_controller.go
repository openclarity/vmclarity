// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
)

func (s *Server) CreateScan(ctx echo.Context) error {
	// Load scan
	var scan types.Scan
	if err := ctx.Bind(&scan); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// Start scan
	scan, err := s.store.Scans().Create(scan)
	if err != nil {
		//if errors.Is(err, ErrScanAlreadyExists) {
		//	return sendError(ctx, http.StatusConflict, err.Error())
		//}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}

func (s *Server) GetScan(ctx echo.Context, scanID types.ScanID) error {
	result, err := s.store.Scans().Get(scanID)
	if err != nil {
		//if errors.Is(err, ErrScanNotFound) {
		//	return sendError(ctx, http.StatusNotFound, err.Error())
		//}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, result)
}

func (s *Server) DeleteScan(ctx echo.Context, scanID types.ScanID) error {
	err := s.store.Scans().Delete(scanID)
	if err != nil {
		//if errors.Is(err, ErrScanNotFound) {
		//	return sendError(ctx, http.StatusNotFound, err.Error())
		//}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, "")
}

func (s *Server) GetScanResult(ctx echo.Context, scanID types.ScanID) error {
	// check scan status
	scan, err := s.store.Scans().Get(scanID)
	if err != nil {
		//if errors.Is(err, ErrScanInProgress) {
		//	return sendError(ctx, http.StatusAccepted, err.Error())
		//}
		//if errors.Is(err, ErrScanNotFound) {
		//	return sendError(ctx, http.StatusNotFound, err.Error())
		//}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// get scan result
	scanResult, err := s.store.ScanFindings().GetAll(types.GetScanFindingsRequest{
		ScanID: scanID,
	})
	if err != nil {
		//if errors.Is(err, ErrScanInProgress) {
		//	return sendError(ctx, http.StatusAccepted, err.Error())
		//}
		//if errors.Is(err, ErrScanNotFound) {
		//	return sendError(ctx, http.StatusNotFound, err.Error())
		//}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, scanResult)
}

func (s *Server) StopScan(ctx echo.Context, scanID types.ScanID) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetScans(ctx echo.Context, params types.GetScansParams) error {
	scans, err := s.store.Scans().GetAll(types.GetScansRequest{
		State: params.State,
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, scans)
}
