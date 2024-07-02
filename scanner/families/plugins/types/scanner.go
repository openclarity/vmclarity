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

package types

import (
	"github.com/openclarity/vmclarity/api/types"
	plugintypes "github.com/openclarity/vmclarity/plugins/sdk-go/types"
	runnerconfig "github.com/openclarity/vmclarity/scanner/families/plugins/runner/config"
)

type ScannersConfig map[string]runnerconfig.Config

func (ScannersConfig) IsConfig() {}

type ScannerResult struct {
	Findings     []types.FindingInfo
	Output       *plugintypes.Result
	ScannedInput string
	ScannerName  string
	Error        error
}

func (r *ScannerResult) GetError() error {
	return r.Error
}
