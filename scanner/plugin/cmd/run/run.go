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
	"github.com/openclarity/vmclarity/scanner/plugin"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"syscall"

	"os"
)

func Run(scanner plugin.Scanner) {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err = InitLogger(config.LogLevel, os.Stderr); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	server, err := plugin.NewServer(scanner)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	go func() {
		log.Infof("Plugin HTTP server starting...")
		if err = server.Start(config.ListenAddress); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	defer func() {
		log.Infof("Plugin HTTP server stopping...")
		if err = server.Stop(); err != nil {
			log.Errorf("Failed to stop HTTP server: %v", err)
			return
		}
		log.Infof("Plugin HTTP server stopped")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	log.Warningf("Received a tremination signal: %v", s)
}
