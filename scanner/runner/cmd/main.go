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
	"context"
	"fmt"
	"time"

	rr "github.com/openclarity/vmclarity/scanner/runner"
)

// Test start scanner function.
func main() {
	config := rr.PluginConfig{
		Name:          "",
		ImageName:     "", // TODO Add image name
		InputDir:      "", // TODO Add input directory
		OutputDir:     "", // TODO Add output directory
		ScannerConfig: "",
	}

	runner, err := rr.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepare proxy
	err = runner.CreateProxyContainer(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start scanner
	fmt.Printf("Starting scanner %s\n", runner.Name)
	err = runner.StartScanner()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer runner.StopScanner() //nolint:errcheck

	//// block forever
	//for {
	//
	//}

	fmt.Printf("Waiting for scanner %s to be ready\n", runner.Name)
	err = runner.WaitScannerReady(time.Second, time.Minute*2) //nolint:gomnd
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Running scanner %s\n", runner.Name)
	err = runner.RunScanner()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Waiting for scanner %s to finish\n", runner.Name)
	err = runner.WaitScannerDone(time.Second, time.Minute*2) //nolint:gomnd
	if err != nil {
		fmt.Println(err)
		return
	}
}
