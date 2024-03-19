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
)

func (s *Server) GetScanFindingsForScan(ctx echo.Context, scanID types.ScanID) error {
	// Check if scan exists
	scan, err := s.store.Scans().Get(scanID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// We don't want to return findings for a scan that is not in succeeded state to
	// prevent the queries to findings DB, but also to allow consumers to poll this
	// endpoint directly without having to poll GET /scan/{scanID}.
	switch state := scan.Status.State; state {
	case types.ScanStatusStateInProgress, types.ScanStatusStatePending:
		return sendResponse(ctx, http.StatusAccepted,
			fmt.Sprintf("scan in %s state, check later", state))

	case types.ScanStatusStateAborted, types.ScanStatusStateFailed:
		return sendError(ctx, http.StatusMethodNotAllowed,
			fmt.Sprintf("scan in %s state, no results", state))

	case types.ScanStatusStateDone:
		// continue
	}

	// Get scan findings
	findings, err := s.store.ScanFindings().GetAll(store.GetScanFindingsRequest{
		ScanID: &scanID,
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, &types.ScanFindings{
		Count: len(findings),
		Items: findings,
	})
}
