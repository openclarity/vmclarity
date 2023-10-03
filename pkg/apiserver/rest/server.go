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
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authn"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authstore"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authz"
	iamTypes "github.com/openclarity/vmclarity/pkg/apiserver/iam/types"

	"github.com/getkin/kin-openapi/openapi3filter"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/openclarity/vmclarity/api/server"
	"github.com/openclarity/vmclarity/pkg/apiserver/common"
	databaseTypes "github.com/openclarity/vmclarity/pkg/apiserver/database/types"
	"github.com/openclarity/vmclarity/pkg/shared/log"
)

const (
	shutdownTimeoutSec = 10
)

type ServerImpl struct {
	dbHandler databaseTypes.Database
	authStore iamTypes.AuthStore
	authn     iamTypes.Authenticator
}

type Server struct {
	port       int
	echoServer *echo.Echo
}

func CreateRESTServer(port int, dbHandler databaseTypes.Database) (*Server, error) {
	e, err := createEchoServer(dbHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest server: %w", err)
	}
	return &Server{
		port:       port,
		echoServer: e,
	}, nil
}

func createEchoServer(dbHandler databaseTypes.Database) (*echo.Echo, error) {
	swagger, err := server.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}

	// Create server
	e := echo.New()

	// Use store
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// Log all requests
	e.Use(echomiddleware.Logger())

	// Recover any panics into HTTP 500
	e.Use(echomiddleware.Recover())

	// Create IAM service data
	authenticator, err := authn.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}
	authorizer, err := authz.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create authorizer: %w", err)
	}
	authStore, err := authstore.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create authstore: %w", err)
	}

	// Use oapi-codegen validation middleware to validate the API group against the OpenAPI schema.
	// Authenticator function must be defined due to OAPI auth specs.
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: iam.NewMiddleware(authenticator, authorizer, authStore),
		},
	}))

	// Register paths with the backend implementation
	server.RegisterHandlers(e, &ServerImpl{
		dbHandler: dbHandler,
		authStore: authStore,
		authn:     authenticator,
	})

	return e, nil
}

func (s *Server) Start(ctx context.Context, errChan chan struct{}) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	logger.Infof("Starting REST server")
	go func() {
		if err := s.echoServer.Start(fmt.Sprintf("0.0.0.0:%d", s.port)); err != nil {
			logger.Errorf("Failed to start REST server: %v", err)
			errChan <- common.Empty
		}
	}()
	logger.Infof("REST server is running")
}

func (s *Server) Stop(ctx context.Context) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	logger.Infof("Stopping REST server")
	if s.echoServer != nil {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeoutSec*time.Second)
		defer cancel()
		if err := s.echoServer.Shutdown(ctx); err != nil {
			logger.Errorf("Failed to shutdown REST server: %v", err)
		}
	}
}
