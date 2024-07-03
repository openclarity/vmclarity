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

package plugins

import (
	"context"
	"fmt"

	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner"
	"github.com/openclarity/vmclarity/scanner/families/plugins/types"
	"github.com/openclarity/vmclarity/scanner/internal/scan_manager"
)

type Plugins struct {
	conf types.Config
}

func New(conf types.Config) families.Family[*types.Result] {
	return &Plugins{
		conf: conf,
	}
}

func (p *Plugins) GetType() families.FamilyType {
	return families.Plugins
}

func (p *Plugins) Run(ctx context.Context, _ *families.Results) (*types.Result, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "plugins")
	logger.Info("Plugins Run...")

	// Register plugins dynamically instead of registering the runner itself
	for _, n := range p.conf.ScannersList {
		types.FactoryRegister(n, runner.New)
	}

	// Top level BinaryMode overrides the individual scanner BinaryMode if set
	if p.conf.BinaryMode != nil {
		for name := range *p.conf.ScannersConfig {
			// for _, config := range *p.conf.ScannersConfig {
			config := (*p.conf.ScannersConfig)[name]
			config.BinaryMode = *p.conf.BinaryMode
			(*p.conf.ScannersConfig)[name] = config
		}
	}

	manager := job_manager.New(p.conf.ScannersList, p.conf.ScannersConfig, logger, factory)
	processResults, err := manager.Process(ctx, p.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for plugins: %w", err)
	}

	pluginsResults := types.NewResult()

	// Merge results
	for _, result := range results {
		logger.Infof("Merging result from %q", result.Metadata.ScannerName)
		pluginsResults.Merge(result.Metadata, result.ScanResult)
	}

	logger.Info("Plugins Done...")

	return pluginsResults, nil
}
