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
	apitypes "github.com/openclarity/vmclarity/api/types"
	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/core/to"
	plugintypes "github.com/openclarity/vmclarity/plugins/sdk-go/types"
	"github.com/openclarity/vmclarity/scanner/families/interfaces"
	"github.com/openclarity/vmclarity/scanner/families/plugins/common"
	"github.com/openclarity/vmclarity/scanner/families/plugins/runner"
	"github.com/openclarity/vmclarity/scanner/families/results"
	"github.com/openclarity/vmclarity/scanner/families/types"
	"github.com/openclarity/vmclarity/scanner/job_manager"
)

type Plugins struct {
	conf Config
}

var _ interfaces.Family = &Plugins{}

func (p *Plugins) Run(ctx context.Context, res *results.Results) (interfaces.IsResults, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "plugins")
	logger.Info("Plugins Run...")

	factory := job_manager.NewJobFactory()
	for _, n := range p.conf.ScannersList {
		factory.Register(n, runner.New)
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

	// Merge results from all plugins into the same output
	var pluginsResults Results

	var mergedResults []apitypes.FindingInfo
	mergedPluginResult := make(map[string]plugintypes.Result)

	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.ScannerName)
		data, ok := result.Result.(*common.Results)
		if !ok {
			return nil, fmt.Errorf("received results of a wrong type: %T", result)
		}
		mergedResults = append(mergedResults, data.Findings...)
		mergedPluginResult[result.ScannerName] = to.ValueOrZero(data.Output)
		pluginsResults.Metadata.InputScans = append(pluginsResults.Metadata.InputScans, result.InputScanMetadata)
	}

	pluginsResults.Findings = mergedResults
	pluginsResults.PluginOutputs = mergedPluginResult

	logger.Info("Plugins Done...")

	return &pluginsResults, nil
}

func (p *Plugins) GetType() types.FamilyType {
	return types.Plugins
}

func New(conf Config) *Plugins {
	return &Plugins{
		conf: conf,
	}
}
