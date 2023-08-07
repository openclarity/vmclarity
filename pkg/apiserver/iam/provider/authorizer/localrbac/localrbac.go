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

package localrbac

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
)

//go:embed rbac_model.conf
var rbacModel string

// New creates an authorizer which will use a local CSV file to configure role rules.
func New() (iam.Authorizer, error) {
	// Load config
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("localrbac: failed to load config: %w", err)
	}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("localrbac: failed to validate config: %w", err)
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcerSafe(casbin.NewModel(rbacModel), fileadapter.NewAdapter(config.RuleFilePath))
	if err != nil {
		return nil, fmt.Errorf("localrbac: failed to create local rbac authorizer: %w", err)
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

func (authorizer *localRBAC) CanPerform(_ context.Context, user *iam.User, asset, action string) (bool, error) {
	for _, role := range user.GetRoles() {
		allowed, err := authorizer.enforcer.EnforceSafe(role, asset, action)
		if err != nil {
			return false, fmt.Errorf("localrbac: failed checking auth role: %w", err)
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}
