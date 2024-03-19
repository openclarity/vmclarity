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
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
	"time"
)

func (s *Server) GetScans(ctx echo.Context, params types.GetScansParams) error {
	scans, err := s.store.Scans().GetAll(store.GetScansRequest{
		State:        params.State,
		MetaSelector: params.MetaSelectors,
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, &types.Scans{
		Count: len(scans),
		Items: scans,
	})
}

func (s *Server) CreateScan(ctx echo.Context) error {
	// Load request
	var scan types.Scan
	if err := ctx.Bind(&scan); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// Create scan
	now := time.Now()
	scan, err := s.store.Scans().Create(types.Scan{
		Annotations:              scan.Annotations,
		InProgressTimeoutSeconds: scan.InProgressTimeoutSeconds,
		Inputs:                   scan.Inputs,
		PendingTimeoutSeconds:    scan.PendingTimeoutSeconds,
		Status: &types.ScanStatus{
			LastTransitionTime: time.Now(),
			State:              types.ScanStatusStatePending,
		},
		SubmitTime: &now,
		Summary:    &types.ScanSummary{},
	})
	if err != nil {
		var checkErr *store.PreconditionFailedError
		if errors.As(err, &checkErr) {
			return sendError(ctx, http.StatusBadRequest, checkErr.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, &scan)
}

func (s *Server) GetScan(ctx echo.Context, scanID types.ScanID) error {
	scan, err := s.store.Scans().Get(scanID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, &scan)
}

func (s *Server) MarkScanAborted(ctx echo.Context, scanID types.ScanID) error {
	statusMsg := "scan was manually aborted via API"
	scan, err := s.store.Scans().Update(scanID, types.Scan{
		Status: &types.ScanStatus{
			LastTransitionTime: time.Now(),
			Message:            &statusMsg,
			State:              types.ScanStatusStateAborted,
		},
	})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, scan)
}
