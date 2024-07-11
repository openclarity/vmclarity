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
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/to"
	"github.com/openclarity/vmclarity/plugins/runner"
	"github.com/openclarity/vmclarity/plugins/runner/types"
	plugintypes "github.com/openclarity/vmclarity/plugins/sdk-go/types"
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

//nolint:cyclop
func (s *Scanner) Run(ctx context.Context, sourceType utils.SourceType, userInput string) error {
	if !s.isValidInputType(sourceType) {
		return fmt.Errorf("received invalid input type for plugin scanner: %v", sourceType)
	}

	rr, err := runner.New(ctx, types.PluginConfig{
		Name:          s.name,
		ImageName:     s.config.ImageName,
		InputDir:      userInput,
		ScannerConfig: s.config.ScannerConfig,
		BinaryMode:    s.config.BinaryMode,
	})
	if err != nil {
		return fmt.Errorf("failed to create plugin runner: %w", err)
	}

	shutdownRunner := func(ctx context.Context) {
		shutdownContext := context.WithoutCancel(ctx)
		if err := rr.Stop(shutdownContext); err != nil {
			s.logger.WithError(err).Errorf("failed to stop runner")
		}

		if err := rr.Remove(shutdownContext); err != nil {
			s.logger.WithError(err).Errorf("failed to remove runner")
		}
	} //nolint:errcheck

	type result struct {
		Result common.Results
		Err    error
	}

	resChan := make(chan result)

	go func(ctx context.Context) {
		defer func() {
			if e := recover(); e != nil {
				shutdownRunner(ctx)
				panic(e)
			}
		}()

		res := result{
			Result: common.Results{
				ScannedInput: userInput,
				ScannerName:  s.name,
			},
			Err: nil,
		}

		if err := rr.Start(ctx); err != nil {
			res.Err = fmt.Errorf("failed to start plugin runner: %w", err)
			resChan <- res
			return
		}

		if err := rr.WaitReady(ctx); err != nil {
			res.Err = fmt.Errorf("failed to wait for plugin scanner to be ready: %w", err)
			resChan <- res
			return
		}

		// Get plugin metadata
		metadata, err := rr.Metadata(ctx)
		if err != nil {
			res.Err = fmt.Errorf("failed to get plugin scanner metadata: %w", err)
			resChan <- res
			return
		}

		// Stream logs
		go func() {
			logger := s.logger.WithField("metadata", map[string]interface{}{
				"name":       to.ValueOrZero(metadata.Name),
				"version":    to.ValueOrZero(metadata.Version),
				"apiVersion": to.ValueOrZero(metadata.ApiVersion),
			}).WithField("plugin", s.config.Name)

			logs, err := rr.Logs(ctx)
			if err != nil {
				logger.WithError(err).Warnf("could not listen for logs on plugin runner")
				return
			}
			defer logs.Close()

			for r := bufio.NewScanner(logs); r.Scan(); {
				logger.Info(r.Text())
			}
		}()

		if err := rr.Run(ctx); err != nil {
			res.Err = fmt.Errorf("failed to run plugin scanner: %w", err)
			resChan <- res
			return
		}

		if err := rr.WaitDone(ctx); err != nil {
			res.Err = fmt.Errorf("failed to wait for plugin scanner to finish: %w", err)
			resChan <- res
			return
		}

		findings, pluginResult, err := s.parseResults(ctx, rr)
		if err != nil {
			res.Err = fmt.Errorf("failed to parse plugin scanner results: %w", err)
			resChan <- res
			return
		}

		res.Result.Findings = findings
		res.Result.Output = pluginResult
		resChan <- res
	}(ctx)

	select {
	case <-ctx.Done():
		shutdownRunner(ctx)
		return errors.New("plugin context cancelled")
	case r := <-resChan:
		shutdownRunner(ctx)
		s.sendResults(r.Result, r.Err)
	}

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

func (s *Scanner) parseResults(ctx context.Context, runner types.PluginRunner) ([]apitypes.FindingInfo, *plugintypes.Result, error) {
	result, err := runner.Result(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get plugin scanner result: %w", err)
	}
	defer result.Close()

	b, err := io.ReadAll(result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read plugin scanner output: %w", err)
	}

	var pluginResult plugintypes.Result
	err = json.Unmarshal(b, &pluginResult)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal plugin scanner output: %w", err)
	}

	findings, err := apitypes.DefaultPluginAdapter.Result(pluginResult)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert plugin scanner result to vmclarity findings: %w", err)
	}

	return findings, &pluginResult, nil
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
