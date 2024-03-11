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
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
	"scanner/scanner"
)

var sknr = scanner.Scanner{}

func (s *Server) GetScanResult(ctx echo.Context) error {
	result, err := sknr.GetScanResult()
	if err != nil {
		return sendError(ctx, 400, err.Error())
	}
	return sendResponse(ctx, 200, result)
}

func (s *Server) StartScan(ctx echo.Context) error {
	// TODO: check that the provided scan and asset IDs are valid
	var scanTemplate types.ScanTemplate
	err := ctx.Bind(&scanTemplate)
	if err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	result, err := sknr.StartScan(scanTemplate)
	if err != nil {
		return sendError(ctx, 400, err.Error())
	}
	return sendResponse(ctx, 200, result)
}

func (s *Server) GetScan(ctx echo.Context) error {
	result, err := sknr.GetScan()
	if err != nil {
		return sendError(ctx, 400, err.Error())
	}
	return sendResponse(ctx, 200, result)
}

func (s *Server) StopScan(ctx echo.Context, params types.StopScanParams) error {
	err := sknr.StopScan()
	if err != nil {
		return sendError(ctx, 400, err.Error())
	}
	return sendResponse(ctx, 200, "")
}
