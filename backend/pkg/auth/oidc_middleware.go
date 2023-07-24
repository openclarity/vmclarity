package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/backend/pkg/common"
	"github.com/openclarity/vmclarity/backend/pkg/config"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/oidc"
	"net/http"
	"strings"
)

// NewOIDCMiddleware creates a middleware which intercepts every call and checks
// for a correct Bearer token using OAuth2 introspection by sending the token to
// the introspection endpoint.
func NewOIDCMiddleware(config *config.Config) (echo.MiddlewareFunc, error) {
	var resourceServer rs.ResourceServer
	var err error
	//TODO: This needs to be fixed; File works normally, but client-secret does not
	if config.OIDCAppFilePath != "" {
		resourceServer, err = rs.NewResourceServerFromKeyFile(config.OIDCIssuer, config.OIDCAppFilePath)
	} else {
		resourceServer, err = rs.NewResourceServerClientCredentials(config.OIDCIssuer, config.OIDCClientID, config.OIDCClientSecret)
	}
	if err != nil {
		return nil, err
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Authenticate
			authHeader := c.Request().Header.Get("authorization")
			if authHeader == "" {
				return common.SendError(c, http.StatusUnauthorized, "auth header missing")
			}
			parts := strings.Split(authHeader, oidc.PrefixBearer)
			if len(parts) != 2 {
				return common.SendError(c, http.StatusUnauthorized, "invalid auth header")
			}
			resp, err := rs.Introspect(ctx, resourceServer, parts[1])
			if err != nil || !resp.IsActive() {
				return common.SendError(c, http.StatusUnauthorized, "invalid token")
			}

			// Inject authorization data
			// roleScope := config.OIDCRolesScope
			// claims := resp.GetClaims()
			return next(c)
		}
	}, nil
}
