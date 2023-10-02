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

package authz

import (
	_ "embed"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"

	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
)

//go:embed rbac_model.conf
var rbacModel string

// New creates an authorizer which will use a local CSV file to configure role rules.
func New() (types.Authorizer, error) {
	// Load config
	config := LoadConfig()

	// Create enforcer
	enforcer, err := casbin.NewEnforcerSafe(casbin.NewModel(rbacModel), fileadapter.NewAdapter(config.RuleFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create rbac model for Authorizer=localrbac: %w", err)
	}

	// Return local RBAC Authorizer
	return &localRBAC{
		enforcer: enforcer,
	}, nil
}

// localRBAC enforces RBAC rules loaded from a local file.
//
// Check examples at: https://www.aserto.com/blog/building-rbac-in-go
type localRBAC struct {
	enforcer *casbin.Enforcer
}

func (authorizer *localRBAC) CanPerform(user models.User, category, action, asset string) (bool, error) {
	// TODO: Use better matcher
	for _, role := range *user.Roles {
		allowed, err := authorizer.enforcer.EnforceSafe(role, category, action, asset)
		if err != nil {
			return false, fmt.Errorf("failed checking auth role for Authorizer=localrbac: %w", err)
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}
