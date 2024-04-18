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
	"log/slog"
	"net/http"
	"time"

	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	internal "github.com/openclarity/vmclarity/plugins/sdk/internal/plugin"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
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

func (s *Server) Start(address string) error {
	err := s.echo.Start(address)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:gomnd
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}

// Echo returns the underlying Echo server that can be used to e.g. add
// middlewares or register new routes. This should happen before Start.
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

//nolint:wrapcheck
func (s *Server) GetHealthz(ctx echo.Context) error {
	slog.Info("Received GetHealthz request")

	if s.scanner.Healthz() {
		return ctx.JSON(http.StatusOK, nil)
	}

	return ctx.JSON(http.StatusServiceUnavailable, nil)
}

//nolint:wrapcheck
func (s *Server) GetMetadata(ctx echo.Context) error {
	slog.Info("Received GetMetadata request")

	return ctx.JSON(http.StatusOK, s.scanner.Metadata())
}

//nolint:wrapcheck
func (s *Server) PostConfig(ctx echo.Context) error {
	slog.Info("Received PostConfig request")

	var config types.Config
	if err := ctx.Bind(&config); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: types.Ptr("failed to bind request"),
		})
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: types.Ptr("failed to validate request"),
		})
	}

	if s.scanner.GetStatus().State != types.Ready {
		return ctx.JSON(http.StatusConflict, &types.ErrorResponse{
			Message: types.Ptr("scanner is not in ready state"),
		})
	}

	s.scanner.Start(&config)

	return ctx.JSON(http.StatusCreated, nil)
}

//nolint:wrapcheck
func (s *Server) GetStatus(ctx echo.Context) error {
	slog.Info("Received GetStatus request")

	return ctx.JSON(http.StatusOK, s.scanner.GetStatus())
}

//nolint:wrapcheck
func (s *Server) PostStop(ctx echo.Context) error {
	slog.Info("Received StopScanner request")

	var requestBody types.Stop
	if err := ctx.Bind(&requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: types.Ptr("failed to bind request"),
		})
	}

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, &types.ErrorResponse{
			Message: types.Ptr("failed to validate request"),
		})
	}

	s.scanner.Stop(requestBody.TimeoutSeconds)

	return ctx.JSON(http.StatusCreated, nil)
}
