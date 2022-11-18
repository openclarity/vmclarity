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

package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type ScannerJobConfig struct {
	DirectoryToScan   string            `json:"directory_to_scan"`
	ServerToReport    string            `json:"server_to_report"`
	VulnerabilityScan VulnerabilityScan `json:"vulnerability_scan"`
	RootkitScan       RootkitScan       `json:"rootkit_scan"`
	MisconfigScan     MisconfigScan     `json:"misconfig_scan"`
	SecretScan        SecretScan        `json:"secret_scan"`
	MalewareScan      MalwareScan       `json:"malawre_scan"`
	ExploitCheck      ExploitCheck      `json:"exploit_check"`
}

type VulnerabilityScan struct {
	Vuls Vuls `json:"vuls"`
}

type RootkitScan struct {
	Chkrootkit Chkrootkit `json:"chkrootkit"`
}

type MisconfigScan struct {
	Lynis Lynis `json:"lynis"`
}

type SecretScan struct {
	Trufflehog Trufflehog `json:"trufflehog"`
}

type MalwareScan struct {
	Clamav Clamav `json:"clamav"`
}

type ExploitCheck struct {
	Vuls Vuls `json:"vuls"`
}

type Vuls struct {
	Config Config `json:"config"`
}

type Chkrootkit struct {
	Config Config `json:"config"`
}

type Lynis struct {
	Config Config `json:"config"`
}

type Trufflehog struct {
	Config Config `json:"config"`
}

type Clamav struct {
	Config Config `json:"config"`
}

type Config struct {
	Someconfig string `json:"someconfig"`
}

func GenerateCloudConfig() error {
	vars := make(map[string]interface{})
	// parse the template
	tmpl, _ := template.ParseFiles("scanner_boot/templates/cloud-config.cfg.tmpl")

	// create a new file
	file, _ := os.Create("cloud-config.cfg")
	defer file.Close()

	confB, err := createScannerConfig()
	if err != nil {
		return err
	}
	vars["Config"] = bytes.NewBuffer(confB).String()
	fmt.Println(vars["Config"])
	return tmpl.Execute(file, vars)
}

func createScannerConfig() ([]byte, error) {
	config := ScannerJobConfig{
		DirectoryToScan: "/test/path",
		ServerToReport:  "127.0.0.1",
		VulnerabilityScan: VulnerabilityScan{
			Vuls: Vuls{
				Config: Config{
					Someconfig: "vuls",
				},
			},
		},
		RootkitScan: RootkitScan{
			Chkrootkit: Chkrootkit{
				Config: Config{
					Someconfig: "chkrootkit",
				},
			},
		},
		MisconfigScan: MisconfigScan{
			Lynis: Lynis{
				Config: Config{
					Someconfig: "lynis",
				},
			},
		},
		SecretScan: SecretScan{
			Trufflehog: Trufflehog{
				Config: Config{
					Someconfig: "trufflehog",
				},
			},
		},
		MalewareScan: MalwareScan{
			Clamav: Clamav{
				Config: Config{
					Someconfig: "clamav",
				},
			},
		},
	}

	scannerConfigB, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("falied to mashal config: %v", err)
	}

	return scannerConfigB, nil
}
