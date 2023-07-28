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

package authorizer

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/openclarity/vmclarity/backend/pkg/iam"
)

//go:embed rbac_model.conf
var rbacModel string

// localRBACAuthorizer enforces RBAC rules loaded from a local file, e.g.
// https://www.aserto.com/blog/building-rbac-in-go
type localRBACAuthorizer struct {
	enforcer *casbin.Enforcer
}

func newLocalRBACAuthorizer(csvRuleFilePath string) (iam.Authorizer, error) {
	enforcer, err := casbin.NewEnforcerSafe(casbin.NewModel(rbacModel), fileadapter.NewAdapter(csvRuleFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create local authorizer: %w", err)
	}

	return &localRBACAuthorizer{
		enforcer: enforcer,
	}, nil
}

func (authorizer *localRBACAuthorizer) CanPerform(_ context.Context, user *iam.User, asset, action string) (bool, error) {
	for _, role := range user.Roles {
		allowed, err := authorizer.enforcer.EnforceSafe(role, asset, action)
		if err != nil {
			return false, err
		}
		if allowed {
			return true, nil
		}
	}
	return false, nil
}
