// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
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

package gitleaks

import (
	"context"
	"encoding/json"
	"fmt"
	job_manager2 "github.com/openclarity/vmclarity/scanner/internal/job_manager"
	types2 "github.com/openclarity/vmclarity/scanner/types"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/secrets/gitleaks/config"
	"github.com/openclarity/vmclarity/scanner/families/secrets/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/utils"
)

const (
	ScannerName    = "gitleaks"
	GitleaksBinary = "gitleaks"
)

type Scanner struct {
	name       string
	logger     *log.Entry
	config     config.Config
	resultChan chan job_manager2.Result
}

func New(_ string, c job_manager2.IsConfig, logger *log.Entry, resultChan chan job_manager2.Result) job_manager2.Job {
	conf := c.(*types.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("scanner", ScannerName),
		config:     conf.Gitleaks,
		resultChan: resultChan,
	}
}

func (a *Scanner) Run(ctx context.Context, sourceType types2.InputType, userInput string) error {
	go func(ctx context.Context) {
		retResults := types.ScannerResult{
			Source:      userInput,
			ScannerName: ScannerName,
		}
		if !a.isValidInputType(sourceType) {
			a.sendResults(retResults, nil)
			return
		}

		// Locate gitleaks binary
		if a.config.BinaryPath == "" {
			a.config.BinaryPath = GitleaksBinary
		}

		gitleaksBinaryPath, err := exec.LookPath(a.config.BinaryPath)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to lookup executable %s: %w", a.config.BinaryPath, err))
			return
		}
		a.logger.Debugf("found gitleaks binary at: %s", gitleaksBinaryPath)

		file, err := os.CreateTemp("", "gitleaks")
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to create temp file. %w", err))
			return
		}
		defer func() {
			_ = os.Remove(file.Name())
		}()
		reportPath := file.Name()

		fsPath, cleanup, err := familiesutils.ConvertInputToFilesystem(ctx, sourceType, userInput)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to convert input to filesystem: %w", err))
			return
		}
		defer cleanup()

		// gitleaks detect --source <source> --no-git -r <report-path> -f json --exit-code 0 --max-target-megabytes 50
		// nolint:gosec
		args := []string{
			"detect",
			"--source",
			fsPath,
			"--no-git",
			"-r",
			reportPath,
			"-f",
			"json",
			"--exit-code",
			"0",
			"--max-target-megabytes",
			"50",
		}
		cmd := exec.Command(gitleaksBinaryPath, args...)
		a.logger.Infof("Running gitleaks command: %v", cmd.String())
		_, err = utils.RunCommand(cmd)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to run gitleaks command: %w", err))
			return
		}

		out, err := os.ReadFile(reportPath)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to read report file from path %v: %w", reportPath, err))
			return
		}

		if err := json.Unmarshal(out, &retResults.Findings); err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to unmarshal results. out: %s. err: %w", out, err))
			return
		}
		a.sendResults(retResults, nil)
	}(ctx)

	return nil
}

func (a *Scanner) isValidInputType(sourceType types2.InputType) bool {
	switch sourceType {
	case types2.DIR, types2.ROOTFS, types2.IMAGE, types2.DOCKERARCHIVE, types2.OCIARCHIVE, types2.OCIDIR:
		return true
	case types2.FILE, types2.SBOM:
		fallthrough
	default:
		a.logger.Infof("source type %v is not supported for gitleaks, skipping.", sourceType)
	}
	return false
}

func (a *Scanner) sendResults(results types.ScannerResult, err error) {
	if err != nil {
		a.logger.Error(err)
		results.Error = err
	}
	select {
	case a.resultChan <- &results:
	default:
		a.logger.Error("Failed to send results on channel")
	}
}
