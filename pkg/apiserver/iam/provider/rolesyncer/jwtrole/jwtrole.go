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

package jwtrole

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
)

// New creates a role syncer which syncs User roles from JWT token claims.
func New() (iam.RoleSyncer, error) {
	// Load config
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("jwtrole: failed to load config: %w", err)
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("jwtrole: failed to validate config: %w", err)
	}

	// Return JWT RoleSyncer
	return &jwtRoleSyncer{
		roleClaim: config.RoleClaim,
	}, nil
}

type jwtRoleSyncer struct {
	roleClaim string
}

func (roleSyncer *jwtRoleSyncer) Sync(_ context.Context, user *iam.User) error {
	// Check user
	if user == nil {
		return fmt.Errorf("jwtrole: no user provider")
	}
	if user.JwtClaims == nil {
		return fmt.Errorf("jwtrole: no user jwt claims found")
	}

	// Get user roles from token
	tokenRoles, ok := user.JwtClaims[roleSyncer.roleClaim]
	if !ok {
		return fmt.Errorf("jwtrole: cannot get user roles %s from token claim", roleSyncer.roleClaim)
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
		return fmt.Errorf("jwtrole: cannot extract roles from token roles type %T", tokenRoles)
	}

	// Sync user roles
	user.SetRoles(userRoles)
	return nil
}
