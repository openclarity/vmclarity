// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
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

package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	backendmodels "github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/api/ui_backend/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func (s *ServerImpl) GetDashboardRiskiestRegions(ctx echo.Context, params models.GetDashboardRiskiestRegionsParams) error {
	scans, err := s.BackendClient.GetScans(context.TODO(), backendmodels.GetScansParams{})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, fmt.Sprintf("failed to get scans from backend: %v", err))
	}
	log.Errorf("got scans from backend: %+v", scans)
	return sendResponse(ctx, http.StatusOK, models.RiskiestRegions{
		Message: utils.StringPtr("riskiest regions!"),
	})
}
