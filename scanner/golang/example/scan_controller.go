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

package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
)

func (s *Server) StartScan(ctx echo.Context) error {
	// Load and validate scan template
	var scanTemplate types.ScanTemplate
	if err := ctx.Bind(&scanTemplate); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}
	if err := scanTemplate.Validate(); err != nil {
		return sendError(ctx, http.StatusBadRequest, err.Error())
	}

	// Start scan
	scan, err := s.Scanner.StartScan(scanTemplate)
	if err != nil {
		if errors.Is(err, types.ErrScanAlreadyExists) {
			return sendError(ctx, http.StatusConflict, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}

func (s *Server) GetScan(ctx echo.Context) error {
	result, err := s.Scanner.GetScan()
	if err != nil {
		if errors.Is(err, types.ErrScanNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, result)
}

func (s *Server) StopScan(ctx echo.Context) error {
	err := s.Scanner.StopScan()
	if err != nil {
		if errors.Is(err, types.ErrScanNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, "")
}

func (s *Server) GetScanResult(ctx echo.Context) error {
	scanResult, err := s.Scanner.GetScanResult()
	if err != nil {
		if errors.Is(err, types.ErrScanInProgress) {
			return sendError(ctx, http.StatusProcessing, err.Error())
		}
		if errors.Is(err, types.ErrScanNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, scanResult)
}
