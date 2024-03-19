package server

import "github.com/labstack/echo/v4"

func (s *Server) IsAlive(ctx echo.Context) error {
	return sendResponse(ctx, 200, "ok")
}

func (s *Server) IsReady(ctx echo.Context) error {
	return sendResponse(ctx, 200, "ok")
}
