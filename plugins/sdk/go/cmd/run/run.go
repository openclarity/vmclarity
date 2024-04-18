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
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/openclarity/vmclarity/plugins/sdk/plugin"
)

// TODO(ramizpolic): allow usage of custom slog

type options struct {
	config *Config
}

func WithConfig(config Config) func(*options) {
	return func(o *options) {
		o.config = &config
	}
}

// Run starts an HTTP server based on given config data. Logger always writes
// data to standard output in order to be able to collect logs.
func Run(scanner plugin.Scanner, opts ...func(*options)) {
	// Load option data
	options := &options{
		config: NewConfig(),
	}
	for _, opt := range opts {
		opt(options)
	}

	// Init logger
	initLogger(options.config.LogLevel)

	// Start server
	server, err := plugin.NewServer(scanner)
	if err != nil {
		slog.Error("failed to create HTTP server", slog.Any("error", err))
	}

	go func() {
		slog.Info("Plugin HTTP server starting...")
		if err = server.Start(options.config.ListenAddress); err != nil {
			slog.Error("failed to start HTTP server", slog.Any("error", err))
		}
	}()

	defer func() {
		slog.Info("Plugin HTTP server stopping...")
		if err = server.Stop(); err != nil {
			slog.Error("failed to stop HTTP server", slog.Any("error", err))
			return
		}
		slog.Info("Plugin HTTP server stopped")
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sig
	slog.Warn(fmt.Sprintf("Received a termination signal: %v", s))
}
