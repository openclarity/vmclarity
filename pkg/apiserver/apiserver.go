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

package apiserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Portshift/go-utils/healthz"

	"github.com/openclarity/vmclarity/pkg/apiserver/database"
	databaseTypes "github.com/openclarity/vmclarity/pkg/apiserver/database/types"
	"github.com/openclarity/vmclarity/pkg/apiserver/rest"
	"github.com/openclarity/vmclarity/utils/log"
)

func createDatabaseConfig(config *Config) databaseTypes.DBConfig {
	return databaseTypes.DBConfig{
		DriverType:     config.DatabaseDriver,
		EnableInfoLogs: config.EnableDBInfoLogs,
		DBPassword:     config.DBPassword,
		DBUser:         config.DBUser,
		DBHost:         config.DBHost,
		DBPort:         config.DBPort,
		DBName:         config.DBName,
		LocalDBPath:    config.LocalDBPath,
	}
}

const defaultChanSize = 100

func Run(ctx context.Context, config *Config) {
	logger := log.GetLoggerFromContextOrDiscard(ctx)

	errChan := make(chan struct{}, defaultChanSize)

	healthServer := healthz.NewHealthServer(config.HealthCheckAddress)
	healthServer.Start()
	healthServer.SetIsReady(false)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger.Info("VMClarity backend is running")

	dbConfig := createDatabaseConfig(config)
	dbHandler, err := database.InitializeDatabase(dbConfig)
	if err != nil {
		logger.Fatalf("Failed to initialise database: %v", err)
	}

	if config.EnableDBFakeData {
		go database.CreateDemoData(ctx, dbHandler)
	}

	// nolint:contextcheck
	restServer, err := rest.CreateRESTServer(config.ListenAddress, dbHandler)
	if err != nil {
		logger.Fatalf("Failed to create REST server: %v", err)
	}
	restServer.Start(ctx, errChan)
	defer restServer.Stop(ctx)

	healthServer.SetIsReady(true)
	logger.Info("VMClarity backend is ready")

	// Wait for deactivation
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-errChan:
		cancel()
		logger.Errorf("Received an error - shutting down")
	case s := <-sig:
		cancel()
		logger.Warningf("Received a termination signal: %v", s)
	}
}
