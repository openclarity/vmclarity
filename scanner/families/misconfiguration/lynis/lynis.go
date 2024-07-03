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

package lynis

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/lynis/config"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
	scannertypes "github.com/openclarity/vmclarity/scanner/types"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/misconfiguration/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/utils"
)

const ScannerName = "lynis"

func init() {
	types.FactoryRegister(ScannerName, New)
}

type Scanner struct {
	logger *log.Entry
	config config.Config
}

func New(_ string, config types.ScannersConfig, logger *log.Entry) familiestypes.Scanner[*types.ScannerResult] {
	return &Scanner{
		logger: logger.Dup().WithField("scanner", ScannerName),
		config: config.Lynis,
	}
}

// nolint: cyclop
func (a *Scanner) Scan(ctx context.Context, sourceType scannertypes.InputType, userInput string) (*types.ScannerResult, error) {
	// Validate this is an input type supported by the scanner,
	// otherwise return skipped.
	if !a.isValidInputType(sourceType) {
		return nil, fmt.Errorf("unsupported source type for %s: %s", ScannerName, sourceType)
	}

	lynisBinaryPath, err := exec.LookPath(a.config.GetBinaryPath())
	if err != nil {
		return nil, fmt.Errorf("failed to lookup executable %s: %w", a.config.BinaryPath, err)
	}
	a.logger.Debugf("found lynis binary at: %s", lynisBinaryPath)

	reportDir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		err := os.RemoveAll(reportDir)
		if err != nil {
			a.logger.Warningf("failed to remove temp directory: %v", err)
		}
	}()

	reportPath := path.Join(reportDir, "lynis.dat")

	fsPath, cleanup, err := familiesutils.ConvertInputToFilesystem(ctx, sourceType, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to filesystem: %w", err)
	}
	defer cleanup()

	// Build command:
	// lynis audit system \
	//     --report-file <reportDir>/report.dat \
	//     --no-log \
	//     --forensics \
	//     --rootdir <source>
	args := []string{
		"audit",
		"system",
		"--report-file",
		reportPath,
		"--no-log",
		"--forensics",
		"--tests",
		strings.Join(testsToRun, ","),
		"--rootdir",
		fsPath,
	}
	cmd := exec.Command(lynisBinaryPath, args...) // nolint:gosec

	a.logger.Infof("Running command: %v", cmd.String())
	_, err = utils.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	// Get Lynis DB directory
	cmd = exec.Command(lynisBinaryPath, []string{"show", "dbdir"}...) // nolint:gosec
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}
	lynisDBDir := filepath.Clean(strings.TrimSpace(string(out)))

	testDB, err := NewTestDB(a.logger, lynisDBDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load lynis test DB: %w", err)
	}

	reportParser := NewReportParser(testDB)
	misconfigurations, err := reportParser.ParseLynisReport(userInput, reportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report file %v: %w", reportPath, err)
	}

	return &types.ScannerResult{
		ScannerName:       ScannerName,
		Misconfigurations: misconfigurations,
	}, nil
}

func (a *Scanner) isValidInputType(sourceType scannertypes.InputType) bool {
	switch sourceType {
	case scannertypes.ROOTFS, scannertypes.IMAGE, scannertypes.DOCKERARCHIVE, scannertypes.OCIARCHIVE, scannertypes.OCIDIR:
		return true
	case scannertypes.DIR, scannertypes.FILE, scannertypes.SBOM:
		a.logger.Infof("source type %v is not supported for lynis, skipping.", sourceType)
	default:
		a.logger.Infof("unknown source type %v, skipping.", sourceType)
	}
	return false
}
