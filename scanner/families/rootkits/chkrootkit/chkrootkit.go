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

package chkrootkit

import (
	"context"
	"fmt"
	job_manager2 "github.com/openclarity/vmclarity/scanner/internal/job_manager"
	types2 "github.com/openclarity/vmclarity/scanner/types"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/rootkits/chkrootkit/config"
	chkrootkitutils "github.com/openclarity/vmclarity/scanner/families/rootkits/chkrootkit/utils"
	"github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/utils"
)

const (
	ScannerName      = "chkrootkit"
	ChkrootkitBinary = "chkrootkit"
)

type Scanner struct {
	name       string
	logger     *log.Entry
	config     config.Config
	resultChan chan job_manager2.Result
}

func (s *Scanner) Run(ctx context.Context, sourceType types2.InputType, userInput string) error {
	go func(ctx context.Context) {
		retResults := types.ScannerResult{
			ScannedInput: userInput,
			ScannerName:  ScannerName,
		}

		if !s.isValidInputType(sourceType) {
			retResults.Error = fmt.Errorf("received invalid input type for chkrootkit scanner: %v", sourceType)
			s.sendResults(retResults, nil)
			return
		}

		// Locate chkrootkit binary
		if s.config.BinaryPath == "" {
			s.config.BinaryPath = ChkrootkitBinary
		}

		chkrootkitBinaryPath, err := exec.LookPath(s.config.BinaryPath)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to lookup executable %s: %w", s.config.BinaryPath, err))
			return
		}
		s.logger.Debugf("found chkrootkit binary at: %s", chkrootkitBinaryPath)

		fsPath, cleanup, err := familiesutils.ConvertInputToFilesystem(ctx, sourceType, userInput)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to convert input to filesystem: %w", err))
			return
		}
		defer cleanup()

		args := []string{
			"-r", // Set userInput as the path to the root volume
			fsPath,
		}

		// nolint:gosec
		cmd := exec.Command(chkrootkitBinaryPath, args...)
		s.logger.Infof("running chkrootkit command: %v", cmd.String())
		out, err := utils.RunCommand(cmd)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to run chkrootkit command: %w", err))
			return
		}

		rootkits, err := chkrootkitutils.ParseChkrootkitOutput(out)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to parse chkrootkit output: %w", err))
			return
		}
		rootkits = filterResults(rootkits)

		retResults.Rootkits = toResultsRootkits(rootkits)

		s.sendResults(retResults, nil)
	}(ctx)

	return nil
}

func filterResults(rootkits []chkrootkitutils.Rootkit) []chkrootkitutils.Rootkit {
	// nolint:prealloc
	var ret []chkrootkitutils.Rootkit
	for _, rootkit := range rootkits {
		if rootkit.RkName == "suspicious files and dirs" {
			// This causes many false positives on every VM, as it's just checks for:
			// files=`${find} ${DIR} -name ".[A-Za-z]*" -o -name "...*" -o -name ".. *"`
			// dirs=`${find} ${DIR} -type d -name ".*"`
			continue
		}
		ret = append(ret, rootkit)
	}
	return ret
}

func toResultsRootkits(rootkits []chkrootkitutils.Rootkit) []types.Rootkit {
	// nolint:prealloc
	var ret []types.Rootkit
	for _, rootkit := range rootkits {
		if !rootkit.Infected {
			continue
		}

		ret = append(ret, types.Rootkit{
			Message:     rootkit.Message,
			RootkitName: rootkit.RkName,
			RootkitType: rootkit.RkType,
		})
	}

	return ret
}

func New(_ string, c job_manager2.IsConfig, logger *log.Entry, resultChan chan job_manager2.Result) job_manager2.Job {
	conf := c.(*types.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("scanner", ScannerName),
		config:     config.Config{BinaryPath: conf.Chkrootkit.BinaryPath},
		resultChan: resultChan,
	}
}

func (s *Scanner) isValidInputType(sourceType types2.InputType) bool {
	switch sourceType {
	case types2.DIR, types2.ROOTFS, types2.IMAGE, types2.DOCKERARCHIVE, types2.OCIARCHIVE, types2.OCIDIR:
		return true
	case types2.FILE, types2.SBOM:
		fallthrough
	default:
		s.logger.Infof("source type %v is not supported for chkrootkit, skipping.", sourceType)
	}
	return false
}

func (s *Scanner) sendResults(results types.ScannerResult, err error) {
	if err != nil {
		s.logger.Error(err)
		results.Error = err
	}
	select {
	case s.resultChan <- &results:
	default:
		s.logger.Error("Failed to send results on channel")
	}
}
