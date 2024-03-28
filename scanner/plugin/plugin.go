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

package plugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	internal "github.com/openclarity/vmclarity/scanner/plugin/internal/plugin"
	"github.com/openclarity/vmclarity/scanner/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	echo    *echo.Echo
	scanner Scanner
}

func NewServer(scanner Scanner) (*Server, error) {
	_, err := internal.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}

	server := &Server{
		echo:    echo.New(),
		scanner: scanner,
	}

	server.echo.Use(echomiddleware.Logger())
	server.echo.Use(echomiddleware.Recover())

	internal.RegisterHandlers(server.echo, server)

	return server, nil
}

func (s *Server) GetHealthz(ctx echo.Context) error {
	log.Info("Received GetHealthz request")

	if s.scanner.Healthz() {
		s.scanner.SetStatus(types.NewScannerStatus(types.Ready, nil))
		return ctx.JSON(http.StatusOK, nil)
	}

	return ctx.JSON(http.StatusServiceUnavailable, nil)
}

func (s *Server) GetMetadata(ctx echo.Context) error {
	log.Info("Received GetMetadata request")

	return ctx.JSON(http.StatusOK, &types.Metadata{ApiVersion: PointerTo("1.0")})
}

func (s *Server) PostConfig(ctx echo.Context) error {
	log.Info("Received PostConfig request")

	var config types.Config
	if err := ctx.Bind(&config); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: PointerTo("failed to bind request"),
		})
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: PointerTo("failed to validate request"),
		})
	}

	if err := s.scanner.Start(&config); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &types.ErrorResponse{
			Message: PointerTo(fmt.Sprintf("failed to start scanner: %v", err)),
		})
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (s *Server) GetStatus(ctx echo.Context) error {
	log.Info("Received GetStatus request")
	return ctx.JSON(http.StatusOK, s.scanner.GetStatus())
}

func (s *Server) Start(address string) error {
	err := s.echo.Start(address)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func PointerTo[T any](value T) *T {
	return &value
}

type Scanner interface {
	Healthz() bool
	Start(config *types.Config) error
	GetStatus() *types.Status
	SetStatus(status *types.Status)
}
