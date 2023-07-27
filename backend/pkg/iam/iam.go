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
	"net/http"
	"strings"
)

const userCtxKey = "user"

// User defines an authenticated user.
type User struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
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

// Injector implements client-side request authentication data injection.
type Injector interface {
	// Inject fetches the authentication data from a source and injects them into
	// request data (e.g. into Authorization header).
	Inject(ctx context.Context, request *http.Request) error
}

// Provider implements server-side IAM synchronization policy.
// TODO: Consider separating auth and authz, e.g. permission fetcher can be passed to Provider.
type Provider interface {
	// Authenticate validates and verifies user auth details from request against an
	// auth provider. It should also be able to fetch permissions associated with the
	// user from some location (e.g. directly from the token, db,...).
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
}

// OapiFilterForProvider creates an OpenAPI filter function which handles request
// authentication and authorization using a specific Provider.
func OapiFilterForProvider(m Provider) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// TODO: Explore caching options to reduce checks against identity server

		// Authenticate
		user, err := m.Authenticate(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return err
		}

		// Update request context with user data
		if eCtx := middleware.GetEchoContext(ctx); eCtx != nil {
			eCtx.Set(userCtxKey, user)
		}

		// Authorize. Route permissions are defined as "asset:action", so we
		// extract the asset and action from the requested permission scope
		if len(input.Scopes) == 0 {
			return nil
		}
		for _, scope := range input.Scopes {
			reqScope := strings.Split(scope, ":")
			if len(reqScope) != 2 {
				return fmt.Errorf("unknown asset:action found: %s", scope)
			}
			asset, action := reqScope[0], reqScope[1]
			if CanPerform(user, asset, action) {
				return nil
			}
		}
		return fmt.Errorf("not allowed, missing required permissions")
	}
}
