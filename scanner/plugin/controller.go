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
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/types"
)

//nolint:wrapcheck
func (s *Server) GetHealthz(ctx echo.Context) error {
	log.Info("Received GetHealthz request")

	if s.scanner.Healthz() {
		return ctx.JSON(http.StatusOK, nil)
	}

	return ctx.JSON(http.StatusServiceUnavailable, nil)
}

//nolint:wrapcheck
func (s *Server) GetMetadata(ctx echo.Context) error {
	log.Info("Received GetMetadata request")

	return ctx.JSON(http.StatusOK, &types.Metadata{ApiVersion: types.Ptr("1.0")})
}

//nolint:wrapcheck
func (s *Server) PostConfig(ctx echo.Context) error {
	log.Info("Received PostConfig request")

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
	log.Info("Received GetStatus request")

	return ctx.JSON(http.StatusOK, s.scanner.GetStatus())
}

//nolint:wrapcheck
func (s *Server) PostStop(ctx echo.Context) error {
	log.Info("Received StopScanner request")

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
