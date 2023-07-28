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

package rolesyncer

import (
	"context"
	"fmt"

	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
)

var RoleSyncerTypeJwt iam.RoleSyncerType = "jwt"

type jwtRoleSyncer struct {
	roleClaim string
}

func newJwtRoleSyncer(roleClaim string) iam.RoleSyncer {
	return &jwtRoleSyncer{
		roleClaim: roleClaim,
	}
}

func (roleSyncer *jwtRoleSyncer) Type() iam.RoleSyncerType {
	return RoleSyncerTypeJwt
}

func (roleSyncer *jwtRoleSyncer) Sync(ctx context.Context, user *iam.User) error {
	// Check user
	if user == nil {
		return fmt.Errorf("no user provider")
	}
	if user.JwtClaims == nil {
		return fmt.Errorf("no user jwt claims found")
	}

	// Get user roles from token
	tokenRoles, ok := user.JwtClaims[roleSyncer.roleClaim]
	if !ok {
		return fmt.Errorf("cannot get user roles from token claim")
	}

	// Get user roles from token roles
	var userRoles []string
	switch tokenRoles := tokenRoles.(type) {
	case map[string]interface{}:
		index := 0
		userRoles = make([]string, len(tokenRoles))
		for roleClaim := range tokenRoles {
			userRoles[index] = roleClaim
			index++
		}
	case []string:
		userRoles = tokenRoles
	default:
		return fmt.Errorf("cannot extract roles from token roles type %T", tokenRoles)
	}

	// Sync user roles
	user.Roles = userRoles
	return nil
}
