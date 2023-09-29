package zitadel

import (
	"context"
	"net/http"
)

type AuthUser struct {
	ID string

	// Data from auth providers
	FromOIDC *AuthFromOIDC
}

type AuthFromOIDC struct {
	Claims map[string]interface{}
}

type Authenticator interface {
	Authenticate(ctx context.Context, req *http.Request) (*AuthUser, error)
}

type User struct {
	ID      string
	Name    string
	Email   string
	Roles   []string
	Banned  bool
	Machine bool
}

type UserStore interface {
	GetUserFromAuth(ctx context.Context, authUser *AuthUser) (*User, error)

	CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error

	RevokeUserAccess(ctx context.Context, user *User) error
	RevokeUserAccessToken(ctx context.Context, user *User, token string) error
}
