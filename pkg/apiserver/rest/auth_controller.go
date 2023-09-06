package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *ServerImpl) AuthRedirect(ctx echo.Context) error {
	return sendResponse(ctx, http.StatusOK, ctx.Request().Header)
}
