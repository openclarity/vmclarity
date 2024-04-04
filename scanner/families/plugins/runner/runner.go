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

package runner

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/plugins/runner"
	"github.com/openclarity/vmclarity/scanner/families/plugins/common"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner/config"
	"github.com/openclarity/vmclarity/scanner/job_manager"
	"github.com/openclarity/vmclarity/scanner/utils"
)

type Scanner struct {
	name       string
	logger     *logrus.Entry
	config     config.Config
	resultChan chan job_manager.Result
}

func New(name string, c job_manager.IsConfig, logger *logrus.Entry, resultChan chan job_manager.Result) job_manager.Job {
	conf := *c.(*common.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       name,
		logger:     logger.Dup().WithField("scanner", name),
		config:     conf[name],
		resultChan: resultChan,
	}
}

func (s *Scanner) Run(sourceType utils.SourceType, userInput string) error {
	go func() {
		ctx := context.Background()

		retResults := common.Results{
			ScannedInput: userInput,
			ScannerName:  s.name,
		}

		if !s.isValidInputType(sourceType) {
			retResults.Error = fmt.Errorf("received invalid input type for plugin scanner: %v", sourceType)
			s.sendResults(retResults, nil)
			return
		}

		config := runner.PluginConfig{
			Name:          s.name,
			ImageName:     s.config.ImageName,
			InputDir:      userInput,
			OutputFile:    s.config.OutputDir,
			ScannerConfig: s.config.ScannerConfig,
		}
		rr, err := runner.New(config)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to create plugin runner: %w", err))
			return
		}

		cleanup, err := rr.Start(ctx)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to start plugin runner: %w", err))
			return
		}
		defer cleanup(ctx)

		err = rr.WaitReady(ctx)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to wait for plugin scanner to be ready: %w", err))
			return
		}

		err = rr.Run(ctx)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to run plugin scanner: %w", err))
			return
		}

		err = rr.WaitDone(ctx)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to wait for plugin scanner to finish: %w", err))
			return
		}

		retResults, err = s.parseResults(retResults, s.config.OutputDir)
		if err != nil {
			s.sendResults(retResults, fmt.Errorf("failed to parse plugin scanner results: %w", err))
			return
		}

		s.sendResults(retResults, nil)
	}()

	return nil
}

func (s *Scanner) isValidInputType(sourceType utils.SourceType) bool {
	switch sourceType {
	case utils.ROOTFS:
		return true
	case utils.DIR, utils.IMAGE, utils.DOCKERARCHIVE, utils.OCIARCHIVE, utils.OCIDIR, utils.FILE, utils.SBOM:
		fallthrough
	default:
		s.logger.Infof("source type %v is not supported for plugin, skipping.", sourceType)
	}
	return false
}

func (s *Scanner) parseResults(results common.Results, resultsDir string) (common.Results, error) {
	// TODO parse scanner output files from directory
	return common.Results{}, nil
}

func (s *Scanner) sendResults(results common.Results, err error) {
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
