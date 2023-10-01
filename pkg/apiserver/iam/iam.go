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
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/types"
	"strings"
)

const userCtxKey = "user"

// GetUserFromContext returns models.User from context.
func GetUserFromContext(ctx echo.Context) *models.User {
	ctxData := ctx.Get(userCtxKey)
	if ctxData == nil {
		return nil
	}
	user, _ := ctxData.(*models.User)
	return user
}

func setUserToContext(eCtx echo.Context, user *models.User) {
	if eCtx == nil {
		return
	}
	eCtx.Set(userCtxKey, user)
}

// NewMiddleware creates an OpenAPI middleware filter which handles request
// authentication and authorization using a specific Provider. Logic flow:
//
// request -> |Authenticator.Authenticate| -> |RoleSyncer.Sync| -> |Authorizer.CanPerform| -> success
//
// Provider will first authenticate the client from request data, synchronize
// user roles on success, and finally try to authorize the request. If
// successful, User will be available in context.
func NewMiddleware(authn types.Authenticator, authz types.Authorizer, store types.AuthStore) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// Remove user from request context
		eCtx := middleware.GetEchoContext(ctx)
		setUserToContext(eCtx, nil)

		// Authenticate
		userInfo, err := authn.Introspect(ctx, input.RequestValidationInput.Request)
		if err != nil {
			return fmt.Errorf("failed to authenticate request: %w", err)
		}

		// Add auth user to request context
		user, err := store.GetUserFromInfo(userInfo)
		if err != nil {
			return err
		}
		setUserToContext(eCtx, &user)

		// TODO: Check RBAC
		for _, ruleDelim := range input.Scopes {
			// For example: "api:update:asset"
			ruleSlice := strings.SplitN(ruleDelim, ":", 3)
			ok, err := authz.CanPerform(user, ruleSlice[0], ruleSlice[1], ruleSlice[2])
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("does not have permissions")
			}
		}

		return nil
	}
}
