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

package types

import (
	"github.com/openclarity/vmclarity/scanner/families/types"
	"time"
)

var _ types.FamilyResult = &FamilyResult{}

type FamilyResult struct {
	Metadata      types.Metadata
	MergedResults *MergedResults `yaml:"merged_results"`
}

func NewFamilyResult() *FamilyResult {
	return &FamilyResult{
		Metadata: types.Metadata{
			Timestamp: time.Now(),
			Scanners:  []string{},
		},
		MergedResults: NewMergedResults(),
	}
}

func (*FamilyResult) IsResult() {}

type MergedResults struct {
	Results []ScannerResult
}

func NewMergedResults() *MergedResults {
	return &MergedResults{}
}

func (m *MergedResults) Merge(other ScannerResult) *MergedResults {
	m.Results = append(m.Results, other)
	return m
}
