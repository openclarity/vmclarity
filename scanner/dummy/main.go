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
	"fmt"
	"github.com/openclarity/vmclarity/scanner/plugin"
	"github.com/openclarity/vmclarity/scanner/plugin/cmd/run"
	"github.com/openclarity/vmclarity/scanner/types"
	"os/exec"
	"time"
)

type DummyScanner struct {
	healthz bool
	status  *types.Status
}

func (d *DummyScanner) Healthz() bool {
	return d.healthz
}

func (d *DummyScanner) Start(config *types.Config) error {
	go func() {
		fmt.Printf("Starting scanner with config: %+v\n", config)
		d.SetStatus(types.NewScannerStatus(types.Running, plugin.PointerTo("Scanner is running...")))
		args := []string{"-a", "-l", "-h"}

		time.Sleep(1 * time.Minute)

		cmd := exec.Command("ls", args...)
		stdout, err := cmd.Output()
		if err != nil {
			d.SetStatus(types.NewScannerStatus(types.Failed, plugin.PointerTo(fmt.Sprintf("Failed to run command: %v", err))))
			fmt.Println(err)
			return
		}

		d.SetStatus(types.NewScannerStatus(types.Done, plugin.PointerTo("Scanner finished running.")))
		fmt.Println(string(stdout))
	}()

	return nil
}

func (d *DummyScanner) GetStatus() *types.Status {
	return d.status
}

func (d *DummyScanner) SetStatus(s *types.Status) {
	d.status = types.NewScannerStatus(s.State, s.Message)
}

func main() {
	// Healthz and status initialized to true and ready,
	// since the scanner does not have any dependencies.
	// Otherwise, the scanner would be initialized with
	// healthz = false and status = NotReady.
	d := &DummyScanner{
		healthz: true,
		status:  types.NewScannerStatus(types.Ready, plugin.PointerTo("Starting scanner...")),
	}

	run.Run(d)
}
