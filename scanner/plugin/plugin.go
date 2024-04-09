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
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	internal "github.com/openclarity/vmclarity/scanner/plugin/internal/plugin"
	"github.com/openclarity/vmclarity/scanner/types"
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

type Scanner interface {
	Healthz() bool
	Start(config *types.Config) error
	GetStatus() *types.Status
	SetStatus(status *types.Status)
}
