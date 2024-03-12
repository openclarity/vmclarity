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

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	scannerserver "github.com/openclarity/vmclarity/scanner/server"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo *echo.Echo
}

func NewServer() (*Server, error) {
	// Get swagger specs
	swagger, err := scannerserver.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load swagger spec: %w", err)
	}

	// Create server instance
	server := &Server{
		echo: echo.New(),
	}

	// Log all requests
	server.echo.Use(echomiddleware.Logger())

	// Recover any panics into HTTP 500
	server.echo.Use(echomiddleware.Recover())

	// Use oapi-codegen validation middleware to validate
	// the API group against the OpenAPI schema.
	server.echo.Use(middleware.OapiRequestValidator(swagger))

	// Register paths with the server implementation
	scannerserver.RegisterHandlers(server.echo, server)

	return server, nil
}

// Start starts the server in a goroutine or exits with fatal error
func (s *Server) Start(address string) {
	log.Infof("Starting scanner server")
	go func() {
		if err := s.echo.Start(address); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Scanner server errored: %v", err)
		}
	}()
	log.Infof("Scanner server ready to start")
}

// Stop stops the server with 10 second timeout
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Infof("Terminating scanner server")
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Errorf("Failed to terminate server: %v", err)
		return
	}
	log.Infof("Scanner server sucessfully terminated")
}

func main() {
	// Load components
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := initLogger(config.LogLevel, os.Stderr); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Run server
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create scanner server: %v", err)
	}
	server.Start(config.ListenAddress)
	defer server.Stop()

	// Wait for shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	log.Warningf("Received a termination signal: %v", s)
}
