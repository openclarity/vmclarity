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
	Output   apitypes.PluginOutput  `json:"Output"`
}

func (*Results) IsResults() {}

func (r *Results) GetTotal() int {
	total := 0
	if r.Output.Exploits != nil {
		total += len(*r.Output.Exploits)
	}
	if r.Output.InfoFinder != nil {
		total += len(*r.Output.InfoFinder)
	}
	if r.Output.Malware != nil {
		total += len(*r.Output.Malware)
	}
	if r.Output.Misconfigurations != nil {
		total += len(*r.Output.Misconfigurations)
	}
	if r.Output.Packages != nil {
		total += len(*r.Output.Packages)
	}
	if r.Output.Rootkits != nil {
		total += len(*r.Output.Rootkits)
	}
	if r.Output.Secrets != nil {
		total += len(*r.Output.Secrets)
	}
	if r.Output.Vulnerabilities != nil {
		total += len(*r.Output.Vulnerabilities)
	}

	return total
}
