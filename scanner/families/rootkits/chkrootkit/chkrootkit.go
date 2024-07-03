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
	"github.com/openclarity/vmclarity/scanner/common"
	"github.com/openclarity/vmclarity/scanner/families"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/scanner/families/rootkits/chkrootkit/config"
	chkrootkitutils "github.com/openclarity/vmclarity/scanner/families/rootkits/chkrootkit/utils"
	"github.com/openclarity/vmclarity/scanner/families/rootkits/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/utils"
)

const ScannerName = "chkrootkit"

type Scanner struct {
	logger *log.Entry
	config config.Config
}

func New(_ string, config types.ScannersConfig, logger *log.Entry) (families.Scanner[*types.ScannerResult], error) {
	return &Scanner{
		logger: logger.Dup().WithField("scanner", ScannerName),
		config: config.Chkrootkit,
	}, nil
}

func (s *Scanner) Scan(ctx context.Context, sourceType common.InputType, userInput string) (*types.ScannerResult, error) {
	if !s.isValidInputType(sourceType) {
		return nil, fmt.Errorf("received invalid input type for chkrootkit scanner: %v", sourceType)
	}

	chkrootkitBinaryPath, err := exec.LookPath(s.config.GetBinaryPath())
	if err != nil {
		return nil, fmt.Errorf("failed to lookup executable %s: %w", s.config.BinaryPath, err)
	}
	s.logger.Debugf("found chkrootkit binary at: %s", chkrootkitBinaryPath)

	fsPath, cleanup, err := familiesutils.ConvertInputToFilesystem(ctx, sourceType, userInput)
	if err != nil {
		return nil, fmt.Errorf("failed to convert input to filesystem: %w", err)
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
		return nil, fmt.Errorf("failed to run chkrootkit command: %w", err)
	}

	rootkits, err := chkrootkitutils.ParseChkrootkitOutput(out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chkrootkit output: %w", err)
	}
	rootkits = filterResults(rootkits)
	resultRootkits := toResultsRootkits(rootkits)

	return &types.ScannerResult{
		Rootkits:     resultRootkits,
		ScannedInput: userInput,
		ScannerName:  ScannerName,
	}, nil
}

func (s *Scanner) isValidInputType(sourceType common.InputType) bool {
	switch sourceType {
	case common.DIR, common.ROOTFS, common.IMAGE, common.DOCKERARCHIVE, common.OCIARCHIVE, common.OCIDIR:
		return true
	case common.FILE, common.SBOM:
		fallthrough
	default:
		s.logger.Infof("source type %v is not supported for chkrootkit, skipping.", sourceType)
	}
	return false
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

func init() {
	types.FactoryRegister(ScannerName, New)
}
