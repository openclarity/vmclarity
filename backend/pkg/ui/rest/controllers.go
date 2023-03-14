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
