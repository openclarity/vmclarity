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

package plugin

import (
	"github.com/openclarity/vmclarity/cli/job_manager"
	"github.com/openclarity/vmclarity/cli/utils"
	"github.com/sirupsen/logrus"
)

type ScannerPlugin struct {
	//TODO Define plugin struct
	name       string
	logger     *logrus.Entry
	config     string // JSON config from asset scan
	resultChan chan job_manager.Result
}

func (p *ScannerPlugin) Run(sourceType utils.SourceType, userInput string) error {
	//TODO Implement plugin logic
	// * validate source type (must be directory or file)
	// * validate source (must exist)
	// * validate userInput
	// * StartScanner()
	// * WaitScannerReady() which includes method to poll healthz
	// * RunScanner() which includes method to post config
	// * WaitScannerDone() which includes method to poll status
	// * StopScanner()
	// * GetFindings() and send results to result channel

	return nil
}
