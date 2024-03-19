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

package run

import (
	scannerserver "github.com/openclarity/vmclarity/scanner/server"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// TODO(ramizpolic): change to use client stubs rather than REST server import

func Run() {
	// Load components
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := InitLogger(config.LogLevel, os.Stderr); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create server
	server, err := scannerserver.NewServer()
	if err != nil {
		log.Fatalf("Failed to create scanner server: %v", err)
	}

	// Handle server start
	go func() {
		log.Infof("Scanner server starting...")
		err := server.Start(config.ListenAddress)
		if err != nil {
			log.Fatalf("Scanner server failed: %v", err)
		}
	}()

	// Handle server shutdown
	defer func() {
		log.Infof("Terminating scanner server...")
		if err := server.Stop(); err != nil {
			log.Errorf("Failed to terminate scanner server: %v", err)
			return
		}
		log.Infof("Scanner server terminated successfully")
	}()

	// Wait for shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	log.Warningf("Received a termination signal: %v", s)
}
