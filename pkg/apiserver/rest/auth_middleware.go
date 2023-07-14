package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/zitadel-go/v2/pkg/api/middleware"
	"net/http"
)

// newAuthMiddleware creates a middleware which intercepts every call and checks
// for a correct Bearer token using OAuth2 introspection by sending the token to
// the introspection endpoint. Relies on third-party Zitadel https://zitadel.com.
func newAuthMiddleware(issuer, appKeyPath string) (echo.MiddlewareFunc, error) {
	resourceServer, err := rs.NewResourceServerFromKeyFile(issuer, appKeyPath)
	if err != nil {
		return nil, err
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			err := middleware.Introspect(req.Context(), req.Header.Get("authorization"), resourceServer)
			if err != nil {
				return sendError(c, http.StatusUnauthorized, err.Error())
			}
			return next(c)
		}
	}, nil
}
