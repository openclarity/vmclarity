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
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
)

func (s *Server) GetFindings(ctx echo.Context, params types.GetFindingsParams) error {
	findings, err := s.store.ScanFindings().GetAll(store.GetScanFindingsRequest{
		MetaSelector: types.MetaSelectorsToMap(params.MetaSelectors),
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusOK, &types.ScanFindings{
		Count: len(findings),
		Items: findings,
	})
}

func (s *Server) GetScanFindingsForScan(ctx echo.Context, scanID types.ScanID) error {
	// Check if scan exists
	// TODO: use DB query here instead of asking from the store
	if _, err := s.store.Scans().Get(scanID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
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
