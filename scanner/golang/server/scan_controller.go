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
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
	"time"
)

func (s *Server) GetScans(ctx echo.Context, params types.GetScansParams) error {
	scans, err := s.store.Scans().GetAll(types.GetScansRequest{
		State: params.State,
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	count := len(scans)
	return sendResponse(ctx, http.StatusOK, types.Scans{
		Count: &count,
		Items: &scans,
	})
}

func (s *Server) GetScan(ctx echo.Context, scanID types.ScanID) error {
	result, err := s.store.Scans().Get(scanID)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, result)
}

func (s *Server) CreateScan(ctx echo.Context) error {
	// Load request
	var scan types.Scan
	if err := ctx.Bind(&scan); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// Create scan
	inputsCount := len(scan.Inputs)
	scan, err := s.store.Scans().Create(types.Scan{
		Inputs:        scan.Inputs,
		JobsLeftToRun: &inputsCount,
		Status: &types.ScanStatus{
			LastTransitionTime: time.Now(),
			State:              types.ScanStatusStatePending,
		},
		TimeoutSeconds: scan.TimeoutSeconds,
	})
	if err != nil {
		var checkErr *types.PreconditionFailedError
		if errors.As(err, &checkErr) {
			return sendError(ctx, http.StatusBadRequest, checkErr.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}

func (s *Server) DeleteScan(ctx echo.Context, scanID types.ScanID) error {
	// Delete scan
	err := s.store.Scans().Delete(scanID)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Clear findings for the given scan
	_ = s.store.ScanFindings().Delete(types.DeleteScanFindingsRequest{
		ScanID: &scanID,
	})

	// TODO: send stop signal via orchestrator in case scan is running

	return sendResponse(ctx, http.StatusOK, "scan successfully deleted")
}

func (s *Server) GetScanResult(ctx echo.Context, scanID types.ScanID) error {
	// Get scan
	scan, err := s.store.Scans().Get(scanID)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Check scan state
	switch state := scan.Status.State; state {
	case types.ScanStatusStateInProgress:
		msg := "scan still in progress"
		return sendResponse(ctx, http.StatusAccepted, types.ScanProgressResponse{
			JobsDone:      scan.JobsCompleted,
			JobsRemaining: scan.JobsLeftToRun,
			Message:       &msg,
		})

	case types.ScanStatusStatePending, types.ScanStatusStateAborted, types.ScanStatusStateFailed:
		return sendError(ctx, http.StatusPreconditionFailed,
			fmt.Sprintf("cannot get scan result due to scan state: %s", state))
	}

	// Get scan results
	scanResult, err := s.store.ScanFindings().GetAll(types.GetScanFindingsRequest{
		ScanID: &scanID,
	})
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, "scan findings not found or created yet")
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	count := len(scanResult)
	return sendResponse(ctx, http.StatusOK, types.ScanFindings{
		Count: &count,
		Items: &scanResult,
	})
}

func (s *Server) StopScan(ctx echo.Context, scanID types.ScanID) error {
	err := s.orchestrator.StartScan(scanID)
	if err != nil {
		if errors.Is(err, types.ErrNotRunning) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, "scan stoppe successfuly")
}
