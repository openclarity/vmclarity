// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package common

import "github.com/openclarity/vmclarity/cli/families/rootkits/types"

type Results struct {
	Rootkits     []Rootkit
	ScannedInput string
	ScannerName  string
	Error        error
}

type Rootkit struct {
	Message     string            `json:"message,omitempty"`
	RootkitName string            `json:"RootkitName,omitempty"`
	RootkitType types.RootkitType `json:"RootkitType,omitempty"`
}

func (r *Results) GetError() error {
	return r.Error
}
