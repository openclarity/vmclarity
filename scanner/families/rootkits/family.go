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

package rootkits

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/core/log"
	familiesinterface "github.com/openclarity/vmclarity/scanner/families/interfaces"
	familiesresults "github.com/openclarity/vmclarity/scanner/families/results"
	"github.com/openclarity/vmclarity/scanner/families/rootkits/common"
	"github.com/openclarity/vmclarity/scanner/families/rootkits/job"
	"github.com/openclarity/vmclarity/scanner/families/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/job_manager"
)

type Rootkits struct {
	conf Config
}

func (r Rootkits) Run(ctx context.Context, _ *familiesresults.Results) (familiesinterface.IsResults, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "rootkits")
	logger.Info("Rootkits Run...")

	manager := job_manager.New(r.conf.ScannersList, r.conf.ScannersConfig, logger, job.Factory)
	processResults, err := manager.Process(ctx, r.conf.Inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to process inputs for rootkits: %w", err)
	}

	mergedResults := NewMergedResults()
	var rootkitsResults Results

	// Merge results.
	for _, result := range processResults {
		logger.Infof("Merging result from %q", result.ScannerName)
		data, ok := result.Result.(*common.Results)
		if !ok {
			return nil, fmt.Errorf("received results of a wrong type: %T", result)
		}
		if familiesutils.ShouldStripInputPath(result.Input.StripPathFromResult, r.conf.StripInputPaths) {
			data = StripPathFromResult(data, result.InputPath)
		}
		mergedResults = mergedResults.Merge(data)
		rootkitsResults.Metadata.InputScans = append(rootkitsResults.Metadata.InputScans, result.InputScanMetadata)
	}

	rootkitsResults.MergedResults = mergedResults

	logger.Info("Rootkits Done...")

	return &rootkitsResults, nil
}

// StripPathFromResult strip input path from results wherever it is found.
func StripPathFromResult(result *common.Results, path string) *common.Results {
	for i := range result.Rootkits {
		result.Rootkits[i].Message = familiesutils.RemoveMountPathSubStringIfNeeded(result.Rootkits[i].Message, path)
	}
	return result
}

func (r Rootkits) GetType() types.FamilyType {
	return types.Rootkits
}

// ensure types implement the requisite interfaces.
var _ familiesinterface.Family = &Rootkits{}

func New(conf Config) *Rootkits {
	return &Rootkits{
		conf: conf,
	}
}
