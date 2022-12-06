// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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
	"os"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/server"

	"github.com/openclarity/vmclarity/backend/pkg/common"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

type ServerImpl struct{}

type Server struct {
	port       int
	echoServer *echo.Echo
	dbHandler  database.Database
}

func CreateRESTServer(port int) (*Server, error) {
	swagger, err := server.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec: %v", err)
		os.Exit(1)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())
	// Use oapi-codegen validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// Register server above as the handler for the interface
	server.RegisterHandlers(e, &ServerImpl{})

	return &Server{
		port:       port,
		echoServer: e,
	}, nil
}

func (s *Server) Start(errChan chan struct{}) {
	log.Infof("Starting REST server")
	go func() {
		if err := s.echoServer.Start(fmt.Sprintf("0.0.0.0:%d", s.port)); err != nil {
			log.Errorf("Failed to serve REST server: %v", err)
			errChan <- common.Empty
		}
	}()
	log.Infof("REST server is running")

}

func (s *Server) Stop() {
	log.Infof("Stopping REST server")
	if s.echoServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.echoServer.Shutdown(ctx); err != nil {
			log.Errorf("Failed to shutdown REST server: %v", err)
		}
	}
}
