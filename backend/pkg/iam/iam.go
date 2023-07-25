package iam

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/openclarity/vmclarity/api/models"
	"net/http"
	"strings"
)

const userCtxKey = "user"

// User defines an authenticated user
type User struct {
	ID    string          `json:"id"`
	Roles map[string]bool `json:"roles"`
}

// Provider implements IAM synchronization policy.
type Provider interface {
	Authenticate(ctx context.Context, request *http.Request) (*User, error)
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

		// Authorize - this can be done somewhere else in the chain by inferring user
		// data from context
		return authorize(user, input.Scopes)
	}
}

// GetAllowedRolesFromContext returns a list of roles from context allowed to
// perform a request.
func GetAllowedRolesFromContext(ctx context.Context) []string {
	requiredRoles, _ := ctx.Value(models.IamPolicyScopes).([]string)
	return requiredRoles
}

// GetUserFromContext returns User from context or throws an error.
func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userCtxKey).(*User)
	if !ok || user == nil {
		return nil, fmt.Errorf("no user found in context")
	}
	return user, nil
}

// authorize authorizes the request by returning nil if the User has at least one
// role from allowedRoles.
func authorize(user *User, allowedRoles []string) error {
	for _, e := range allowedRoles {
		if _, ok := user.Roles[e]; ok {
			return nil
		}
	}
	return fmt.Errorf("not allowed, requires at least one of: %s", strings.Join(allowedRoles, ", "))
}
