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

package _default

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) IsAlive(ctx echo.Context) error {
	return sendResponse(ctx, 200, nil)
}

func (s *Server) IsReady(ctx echo.Context) error {
	return sendResponse(ctx, 200, nil)
}

func (s *Server) GetScannerInfo(ctx echo.Context) error {
	info, err := s.manager.Scanner().GetInfo(context.Background())
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
	return sendResponse(ctx, http.StatusOK, info)
}
