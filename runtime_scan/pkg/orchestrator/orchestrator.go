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
	"sync"

	log "github.com/sirupsen/logrus"

	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	aws2 "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/aws"
	_scanner "github.com/openclarity/vmclarity/runtime_scan/pkg/scanner"
)

type Orchestrator struct {
	scanner        *_scanner.Scanner
	config         *_config.Config
	providerClient provider.Client
	sync.Mutex
}

//go:generate $GOPATH/bin/mockgen -destination=./mock_orchestrator.go -package=orchestrator github.com/openclarity/kubeclarity/runtime_scan/pkg/orchestrator VulnerabilitiesScanner
type VulnerabilitiesScanner interface {
	Start(errChan chan struct{})
	Scan(scanConfig *_config.ScanConfig, scanDone chan struct{}) error
	//ScanProgress() types.ScanProgress
	///Results() *types.ScanResults
	//Clear()
	Stop()
}

func Create(config *_config.Config) (*Orchestrator, error) {
	// for now will statically create aws client here (until we support more cloud providers)
	awsClient, err := aws2.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create aws client: %v", err)
	}

	orc := &Orchestrator{
		config:         config,
		providerClient: awsClient,
		Mutex:          sync.Mutex{},
	}

	return orc, nil
}

func (o *Orchestrator) Start(errChan chan struct{}) {
	// Start result server
	log.Infof("Starting Orchestrator server")
}

func (o *Orchestrator) Stop() {
	log.Infof("Stopping Orchestrator server")
}

func (o *Orchestrator) Scan(scanConfig *_config.ScanConfig, scanDone chan struct{}) error {
	instances, err := o.providerClient.Discover(&scanConfig.DiscoveryFilters)
	if err != nil {
		return err
	}
	scanConfig.Instances = instances

	if err := o.getScanner().Scan(scanConfig, scanDone); err != nil {
		return fmt.Errorf("failed to scan: %v", err)
	}

	return nil
}

//
//func (o *Orchestrator) ScanProgress() types.ScanProgress {
//	return o.getScanner().ScanProgress()
//}
//
//func (o *Orchestrator) Results() *types.ScanResults {
//	return o.getScanner().Results()
//}

//func (o *Orchestrator) Clear() {
//	o.Lock()
//	defer o.Unlock()
//
//	log.Infof("Clearing Orchestrator")
//	o.scanner.Clear()
//	//o.scanner = _scanner.CreateScanner(o.config, o.clientset)
//}

func (o *Orchestrator) getScanner() *_scanner.Scanner {
	o.Lock()
	defer o.Unlock()

	return o.scanner
}
