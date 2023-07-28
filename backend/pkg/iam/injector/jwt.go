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

package injector

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/backend/pkg/iam"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"golang.org/x/oauth2"
	"net/http"
)

type jwtInjector struct {
	tokenSource oauth2.TokenSource
}

// newJWTInjector creates a client Injector which creates OAuth2 token source
// from key file to generate tokens that will be injected in requests.
//
// TODO: Enable support for creating token sources from file data string and Personal Access Tokens.
// TODO: This can be achieved using functional options, e.g. WithKey(string), WithKeyFile(filepath), WithAccessToken(string).
// TODO: Test against different OIDCs to check if this works. Tested against: Zitadel.
func newJWTInjector(issuer, keyPath string, extraScopes []string) (iam.Injector, error) {
	// Get token source
	tokenSource, err := middleware.JWTProfileFromPath(keyPath)(issuer, append(extraScopes, oidc.ScopeOpenID))
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC token source: %w", err)
	}

	// Return Injector with reusable token source to prevent request spikes
	return &jwtInjector{
		tokenSource: oauth2.ReuseTokenSource(nil, tokenSource),
	}, nil
}

func (injector *jwtInjector) Inject(_ context.Context, request *http.Request) error {
	token, err := injector.tokenSource.Token()
	if err != nil {
		return err
	}

	token.SetAuthHeader(request)
	return nil
}
