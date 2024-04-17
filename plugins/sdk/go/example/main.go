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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/plugins/sdk/cmd/run"
	"github.com/openclarity/vmclarity/plugins/sdk/types"
)

//nolint:containedctx
type Scanner struct {
	status *types.Status
}

func (s *Scanner) Metadata() *types.Metadata {
	return &types.Metadata{
		ApiVersion: types.Ptr(types.ApiVersion),
		Name:       types.Ptr("Example scanner"),
		Version:    types.Ptr("v0.1.2"),
	}
}

func (s *Scanner) Start(config *types.Config) {
	log.Infof("Starting scanner with config: %+v\n", config)

	go func() {
		// Mark scan started
		log.Infof("Scanner is running...")
		s.SetStatus(types.NewScannerStatus(types.Running, types.Ptr("Scanner is running...")))

		// Do actual scanning here
		time.Sleep(5 * time.Second) //nolint:gomnd

		// Save scan results
		log.Infof("Scanner finished running.")
		s.SetStatus(types.NewScannerStatus(types.Done, types.Ptr("Scanner finished running.")))
	}()
}

func (s *Scanner) GetStatus() *types.Status {
	return s.status
}

func (s *Scanner) SetStatus(newStatus *types.Status) {
	s.status = types.NewScannerStatus(newStatus.State, newStatus.Message)
}

func (s *Scanner) Stop(timeoutSeconds int) {
	// Shutdown logic
}

func (s *Scanner) Healthz() bool { return true }

func main() {
	run.Run(&Scanner{
		status: types.NewScannerStatus(types.Ready, types.Ptr("Scanner ready")),
	})
}
