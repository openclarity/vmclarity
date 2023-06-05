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

package models

func (s *TargetScanState) GetState() (TargetScanStateState, bool) {
	var state TargetScanStateState
	var ok bool

	if s.State != nil {
		state, ok = *s.State, true
	}
	return state, ok
}

func (s *TargetScanState) GetErrors() []string {
	var errs []string

	if s.Errors != nil {
		errs = *s.Errors
	}
	return errs
}
