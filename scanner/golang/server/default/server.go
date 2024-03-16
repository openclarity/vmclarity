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
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	server2 "github.com/openclarity/vmclarity/scanner/server"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
	"time"
)

type Server struct {
	echo    *echo.Echo
	manager types.ScanManager
}

func NewServer(scanner types.Scanner) (*Server, error) {
	// Get swagger specs
	swagger, err := server2.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}

	// Create server instance
	server := &Server{
		echo:    echo.New(),
		manager: types.NewScannerManager(scanner),
	}

	// Log all requests
	server.echo.Use(echomiddleware.Logger())

	// Recover any panics into HTTP 500
	server.echo.Use(echomiddleware.Recover())

	// Use oapi-codegen validation middleware to validate
	// the API group against the OpenAPI schema.
	server.echo.Use(middleware.OapiRequestValidator(swagger))

	// Register paths with the server implementation
	server2.RegisterHandlers(server.echo, server)

	return server, nil
}

// Start starts the server and blocks until the server exits or returns an error
func (s *Server) Start(address string) error {
	err := s.echo.Start(address)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stops the server with 10 second timeout or returns an error
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
