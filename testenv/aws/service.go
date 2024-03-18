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

package aws

import (
	envtypes "github.com/openclarity/vmclarity/testenv/types"
)

type Service struct {
	ID          string
	Namespace   string
	Application string
	Component   string
	State       envtypes.ServiceState
}

func (s *Service) GetID() string {
	return s.ID
}

func (s *Service) GetNamespace() string {
	return s.Namespace
}

func (s *Service) GetApplicationName() string {
	return s.Application
}

func (s *Service) GetComponentName() string {
	return s.Component
}

func (s *Service) GetState() envtypes.ServiceState {
	return s.State
}

func (s Service) String() string {
	return s.ID
}

// Convert Docker container state to vmclarity service state.
func convertStateFromDocker(state string) envtypes.ServiceState {
	switch state {
	case "running":
		return envtypes.ServiceStateReady
	case "exited":
		return envtypes.ServiceStateNotReady
	}
	return envtypes.ServiceStateUnknown
}
