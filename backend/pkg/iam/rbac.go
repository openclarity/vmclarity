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

package iam

import (
	_ "embed"
	"github.com/casbin/casbin"
	scas "github.com/qiangmzsx/string-adapter"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed rbac_model.conf
	rbacModel string
	//go:embed rbac_policy.csv
	rbacPolicy string
	// rbacEnforcer enforces RBAC rules for CanPerform, e.g. https://www.aserto.com/blog/building-rbac-in-go
	rbacEnforcer = casbin.NewEnforcer(casbin.NewModel(rbacModel), scas.NewAdapter(rbacPolicy))
)

// CanPerform checks if User is allowed to perform an action on an asset.
func CanPerform(user *User, asset, action string) bool {
	for _, role := range user.Roles {
		allowed, err := rbacEnforcer.EnforceSafe(role, asset, action)
		if err != nil {
			log.Warnf("Authorization failed for (%s, %s, %s) with error: %v", role, asset, action, err.Error())
		}
		if allowed {
			return true
		}
	}
	return false
}
