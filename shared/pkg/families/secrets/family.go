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

package secrets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/results"
)

type Secrets struct {
	conf   Config
	logger *log.Entry
}

func (s Secrets) Run(res *results.Results) (_interface.IsResults, error) {
	s.logger.Info("Secrets Run...")

	// validate that gitleaks binary exists
	if _, err := os.Stat(s.conf.GitleaksConfig.BinaryPath); err != nil {
		return nil, fmt.Errorf("failed to find binary in %v: %v", s.conf.GitleaksConfig.BinaryPath, err)
	}

	// ./gitleaks detect -v --source=<source> --no-git -r <report-path> -f json --exit-code 0
	cmd := exec.Command(s.conf.GitleaksConfig.BinaryPath, "detect", fmt.Sprintf("--source=%v", s.conf.GitleaksConfig.Source), "--no-git", "-r", s.conf.GitleaksConfig.ReportPath, "-f", "json", "--exit-code", "0")
	_, err := runCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run gitleaks command: %v", err)
	}
	out, err := os.ReadFile(s.conf.GitleaksConfig.ReportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file from path: %v. %v", s.conf.GitleaksConfig.ReportPath, err)
	}

	log.Infof("gitleaks results: %s", out)

	var retResults Results
	if err := json.Unmarshal(out, &retResults.Findings); err != nil {
		return nil, err
	}

	s.logger.Info("Secrets Done...")
	return &retResults, nil
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	//cmd := exec.Command(name, arg)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		err = errors.New(fmt.Sprintf("%v. %v", err, errb.String()))
		return nil, fmt.Errorf("failed to run command: %v. %v", cmd.String(), err)
	}
	return outb.Bytes(), nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Secrets{}

func New(logger *log.Entry, conf Config) *Secrets {
	return &Secrets{
		conf:   conf,
		logger: logger.Dup().WithField("family", "secrets"),
	}
}
