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
	"context"
	_ "embed"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	scas "github.com/qiangmzsx/string-adapter"
	"net/http"
	"strings"
)

const userCtxKey = "user"

var (
	//go:embed rbac_model.conf
	rbacModel string
	//go:embed rbac_policy.csv
	rbacPolicy string
	// rbacEnforcer enforces RBAC rules for CanPerform, e.g. https://www.aserto.com/blog/building-rbac-in-go
	rbacEnforcer = casbin.NewEnforcer(casbin.NewModel(rbacModel), scas.NewAdapter(rbacPolicy))
)

// User defines an authenticated user.
type User struct {
	ID    string          `json:"id"`
	Roles map[string]bool `json:"roles"`
}

// Provider implements server-side IAM synchronization policy.
type Provider interface {
	// Authenticate validates and verifies user auth details from request against
	// some auth provider. It should also be able to fetch permissions associated
	// with that key from some location.
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
}

// OapiAuthenticatorForProvider creates an OpenAPI authenticator using a specific Provider.
func OapiAuthenticatorForProvider(m Provider) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// Authenticate
		user, err := m.Authenticate(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return err
		}

		// Update request context with user data
		if eCtx := middleware.GetEchoContext(ctx); eCtx != nil {
			eCtx.Set(userCtxKey, user)
		}

		// Authorize
		if len(input.Scopes) == 0 {
			return nil
		}
		for _, scope := range input.Scopes {
			reqScope := strings.Split(scope, ":")
			if len(reqScope) != 2 {
				return fmt.Errorf("unknown asset:action defined for route %s", input.RequestValidationInput.Route.Path)
			}
			asset, action := reqScope[0], reqScope[1]
			if CanPerform(user, asset, action) {
				return nil
			}
		}
		return fmt.Errorf("not allowed, missing required permissions")
	}
}

// GetUserFromContext returns User from context.
func GetUserFromContext(ctx echo.Context) *User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}

	user, _ := ctxData.(*User)
	return user
}

// CanPerform checks if User is allowed to perform an action on an asset.
func CanPerform(user *User, asset, action string) bool {
	for role := range user.Roles {
		if ok, _ := rbacEnforcer.EnforceSafe(role, asset, action); ok {
			return true
		}
	}
	return false
}
