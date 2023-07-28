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

package role_syncer

import (
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/openclarity/vmclarity/backend/pkg/iam"
)

// NewRoleSyncer creates a new iam.RoleSyncer from config.
// TODO: Use Factory pattern when this supports multiple iam.RoleSyncer.
func NewRoleSyncer(config config.AuthRoleSynchronization) (iam.RoleSyncer, error) {
	return newJwtRoleSyncer(config.JWTRoleClaim), nil
}
