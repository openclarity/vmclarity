package rest

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
	"github.com/openclarity/vmclarity/pkg/apiserver/iam/authn"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
)

func (s *ServerImpl) Login(ctx echo.Context) error {
	codeVerifier, verifierErr := randomBytesInHex(32) // 64 character string here
	if verifierErr != nil {
		return ctx.String(http.StatusInternalServerError, verifierErr.Error())
	}
	sha2 := sha256.New()
	_, _ = io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	state, stateErr := randomBytesInHex(24)
	if stateErr != nil {
		return ctx.String(http.StatusInternalServerError, stateErr.Error())
	}

	redirectUrl := s.authn.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	return ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (s *ServerImpl) Callback(ctx echo.Context) error {
	sess, _ := session.Get("authenticate-sessions", ctx)
	sess.Options = &sessions.Options{
		Path:     "/secret",
		MaxAge:   5 * 60, // 5min
		HttpOnly: true,
	}

	// Exchange an authorization code for a token.
	code := ctx.QueryParam("code")
	token, err := s.authn.Exchange(ctx.Request().Context(),
		code,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
	}

	_, err = s.authn.Verify(ctx.Request().Context(), token)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
	}

	// Redirect to logged in page.
	// Set user as authenticated
	sess.Values["authenticated"] = true
	err = sess.Save(ctx.Request(), ctx.Response())
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Failed to save session.")
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, "/user")
}

func (s *ServerImpl) Logout(ctx echo.Context) error {
	logoutUrl, err := url.Parse(authn.LoadConfig().Issuer + "/v2/logout")
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	scheme := "http"
	if ctx.Request().TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request().Host)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	logoutUrl.RawQuery = parameters.Encode()

	// Revoke users authentication
	sess, _ := session.Get("authenticate-sessions", ctx)
	sess.Values["authenticated"] = false

	return ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

// TODO: done.
func (s *ServerImpl) GetCurrentUser(ctx echo.Context) error {
	user := iam.GetUserFromContext(ctx)
	if user == nil {
		return sendError(ctx, http.StatusUnauthorized, "not authenticated")
	}
	return sendResponse(ctx, http.StatusOK, user)
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("Could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}
