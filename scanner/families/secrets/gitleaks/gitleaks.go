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
	"github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/secrets/gitleaks/config"
	"github.com/openclarity/vmclarity/scanner/families/secrets/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/utils"
)

const ScannerName = "gitleaks"

type Scanner struct {
	logger *log.Entry
	config config.Config
}

func New(_ string, config types.ScannersConfig, logger *log.Entry) (families.Scanner[*types.ScannerResult], error) {
	return &Scanner{
		logger: logger.Dup().WithField("scanner", ScannerName),
		config: config.Gitleaks,
	}, nil
}

func (a *Scanner) Scan(ctx context.Context, sourceType common.InputType, userInput string) (*types.ScannerResult, error) {
	if !a.isValidInputType(sourceType) {
		return nil, fmt.Errorf("received invalid input type for gitleaks scanner: %v", sourceType)
	}

	gitleaksBinaryPath, err := exec.LookPath(a.config.GetBinaryPath())
	if err != nil {
		return nil, fmt.Errorf("failed to lookup executable %s: %w", a.config.BinaryPath, err)
	}
	a.logger.Debugf("found gitleaks binary at: %s", gitleaksBinaryPath)

	file, err := os.CreateTemp("", "gitleaks")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file. %w", err)
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()
	reportPath := file.Name()

	fsPath, cleanup, err := familiesutils.ConvertInputToFilesystem(ctx, sourceType, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to filesystem: %w", err)
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
		return nil, fmt.Errorf("failed to run gitleaks command: %w", err)
	}

	out, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file from path %v: %w", reportPath, err)
	}

	var findings []types.Finding
	if err := json.Unmarshal(out, &findings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results. out: %s. err: %w", out, err)
	}

	return &types.ScannerResult{
		Findings:    findings,
		Source:      userInput,
		ScannerName: ScannerName,
	}, nil
}

func (a *Scanner) isValidInputType(sourceType common.InputType) bool {
	switch sourceType {
	case common.DIR, common.ROOTFS, common.IMAGE, common.DOCKERARCHIVE, common.OCIARCHIVE, common.OCIDIR:
		return true
	case common.FILE, common.SBOM:
		fallthrough
	default:
		a.logger.Infof("source type %v is not supported for gitleaks, skipping.", sourceType)
	}
	return false
}

func init() {
	types.FactoryRegister(ScannerName, New)
}
