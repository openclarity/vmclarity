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

package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

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

var (
	// ResourceStatusReadySet represents all states of the resource when it is ready.
	ResourceStatusReadySet = map[types.ResourceStatus]bool{
		types.ResourceStatusCreateComplete:         true,
		types.ResourceStatusDeleteComplete:         true,
		types.ResourceStatusDeleteSkipped:          true,
		types.ResourceStatusUpdateComplete:         true,
		types.ResourceStatusImportComplete:         true,
		types.ResourceStatusImportRollbackComplete: true,
		types.ResourceStatusUpdateRollbackComplete: true,
		types.ResourceStatusRollbackComplete:       true,
	}
	// ResourceStateNotReadySet represents all states of the resource when it is not ready.
	ResourceStatusNotReadySet = map[types.ResourceStatus]bool{
		types.ResourceStatusCreateInProgress:         true,
		types.ResourceStatusDeleteInProgress:         true,
		types.ResourceStatusUpdateInProgress:         true,
		types.ResourceStatusImportInProgress:         true,
		types.ResourceStatusImportRollbackInProgress: true,
		types.ResourceStatusUpdateRollbackInProgress: true,
		types.ResourceStatusRollbackInProgress:       true,
	}
	// ResourceStateDegradedSet represents all states of the resource when it is degraded.
	ResourceStatusDegradedSet = map[types.ResourceStatus]bool{
		types.ResourceStatusCreateFailed:         true,
		types.ResourceStatusDeleteFailed:         true,
		types.ResourceStatusUpdateFailed:         true,
		types.ResourceStatusImportFailed:         true,
		types.ResourceStatusImportRollbackFailed: true,
		types.ResourceStatusUpdateRollbackFailed: true,
		types.ResourceStatusRollbackFailed:       true,
	}
)

// Convert AWS CloudFormation stack resource status to vmclarity service state.
func convertStateFromAWS(state types.ResourceStatus) envtypes.ServiceState {
	if ResourceStatusReadySet[state] {
		return envtypes.ServiceStateReady
	} else if ResourceStatusNotReadySet[state] {
		return envtypes.ServiceStateNotReady
	} else if ResourceStatusDegradedSet[state] {
		return envtypes.ServiceStateDegraded
	}
	return envtypes.ServiceStateUnknown
}
