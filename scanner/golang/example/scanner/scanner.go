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

package scanner

import (
	"context"
	dockle_run "github.com/Portshift/dockle/pkg"
	"github.com/openclarity/vmclarity/scanner/types"
)

type Scanner struct{}

func (s *Scanner) GetInfo(ctx context.Context) (*types.ScannerInfo, error) {
	return &types.ScannerInfo{
		Name:    "cisdocker",
		Version: "1.23",
	}, nil
}

func (s *Scanner) Scan(ctx context.Context, family types.ScanFamily, input types.ScanObjectInput) (*types.ScanResult, error) {
	// Validate this is an input type supported by the scanner,
	// otherwise return skipped.
	if !s.isValidInputType(input.Type) {
		return nil, nil // skip
	}

	dockleCfg := createDockleConfig(input)
	assessmentMap, err := dockle_run.RunFromConfig(dockleCfg)
	if err != nil {
		return nil, err
	}

	findings := parseDockleReport(input, assessmentMap)

	return &findings, nil
}

func (s *Scanner) isValidInputType(sourceType types.ScanObjectInputType) bool {
	switch sourceType {
	case types.InputTypeImage, types.InputTypeDockerArchive:
		return true
	default:
		return false
	}
}
