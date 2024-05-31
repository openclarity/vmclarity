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

package plugins

import (
	apitypes "github.com/openclarity/vmclarity/api/types"
	familiestypes "github.com/openclarity/vmclarity/scanner/families/types"
)

type Results struct {
	Metadata familiestypes.Metadata `json:"Metadata"`
	Output   []apitypes.FindingInfo `json:"Output"`
	RawData  map[string]interface{} `json:"RawData"`
}

func (*Results) IsResults() {}

func (r *Results) GetTotal() int {
	return len(r.Output)
}
