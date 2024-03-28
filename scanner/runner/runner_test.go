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

package runner

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Test start scanner function
func TestStartScanner(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	config := PluginConfig{
		Name:          "test-scanner",
		ImageName:     "alpine:latest",
		InputDir:      wd + "/input",
		OutputDir:     wd + "/output",
		ScannerConfig: "plugin.json",
	}

	runner, err := New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Starting scanner %s\n", runner.Name)
	err = runner.StartScanner()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer runner.StopScanner()

	time.Sleep(10 * time.Second)

	fmt.Printf("Waiting for scanner %s to be ready\n", runner.Name)
	err = runner.WaitScannerReady(time.Second, time.Minute*2)
	if err != nil {
		fmt.Println(err)
		return
	}
}
