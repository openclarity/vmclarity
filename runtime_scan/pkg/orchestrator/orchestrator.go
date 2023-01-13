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

package orchestrator

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/orchestrator/scanconfig_watcher"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/orchestrator/scheduler"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
)

type Orchestrator struct {
	config            *_config.OrchestratorConfig
	scanConfigWatcher *scanconfig_watcher.ScanConfigWatcher
	scheduler         *scheduler.Scheduler
}

func Create(config *_config.OrchestratorConfig, providerClient provider.Client) (*Orchestrator, error) {
	backendClient, err := client.NewClientWithResponses(
		fmt.Sprintf("%s:%d/%s", config.BackendAddress, config.BackendRestPort, config.BackendBaseURL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create a backend client: %v", err)
	}
	scanConfigChan := make(chan *map[string]models.ScanConfig)
	orc := &Orchestrator{
		config:            config,
		scanConfigWatcher: scanconfig_watcher.CreateScanConfigWatcher(scanConfigChan, backendClient),
		scheduler:         scheduler.CreateScheduler(scanConfigChan, &config.ScannerConfig, providerClient, backendClient),
	}

	return orc, nil
}

func (o *Orchestrator) Start(errChan chan struct{}) {
	log.Infof("Starting Orchestrator server")

	go o.scanConfigWatcher.Start(errChan)
	go o.scheduler.Start(errChan)
}

func (o *Orchestrator) Stop(errChan chan struct{}) {
	// Stop orchestrator server
	o.scanConfigWatcher.Stop()
	o.scheduler.Stop()

	log.Infof("Stop Orchestrator server")
}
