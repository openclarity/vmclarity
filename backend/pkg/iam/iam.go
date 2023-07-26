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
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/api/models"
	"net/http"
)

const userCtxKey = "user"

// User defines an authenticated user
type User struct {
	ID    string          `json:"id"`
	Roles map[string]bool `json:"roles"`
}

// Provider implements server-side IAM synchronization policy.
type Provider interface {
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
}

// Injector implements client-side authentication data injection.
type Injector interface {
	Inject(ctx context.Context, req *http.Request) error
}

// OapiAuthenticatorForProvider creates an OpenAPI authenticator for a given Provider
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

		// Authorize - this can be done somewhere else in the chain by inferring
		// user/role data from context
		return authorize(user, input.Scopes)
	}
}

// GetRequiredRolesFromContext returns a list of roles from context required to
// perform a request.
func GetRequiredRolesFromContext(ctx echo.Context) []string {
	ctxData := ctx.Get(models.IamPolicyScopes)
	if ctxData == nil {
		return nil
	}

	requiredRoles, _ := ctxData.([]string)
	return requiredRoles
}

// GetUserFromContext returns User from context or throws an error.
func GetUserFromContext(ctx echo.Context) *User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}

	user, _ := ctxData.(*User)
	return user
}

// authorize authorizes the request by returning nil if the User has all requiredRoles.
func authorize(user *User, requiredRoles []string) error {
	for _, role := range requiredRoles {
		if _, ok := user.Roles[role]; !ok {
			return fmt.Errorf("not allowed, requires %s", role)
		}
	}
	return nil
}
