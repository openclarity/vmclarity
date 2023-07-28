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
	"github.com/openclarity/vmclarity/backend/pkg/iam"
	"net/http"
)

type tokenInjector struct {
	accessToken string
}

// newTokenInjector creates an iam.Injector which adds personal access token to request.
func newTokenInjector(accessToken string) (iam.Injector, error) {
	return &tokenInjector{
		accessToken: accessToken,
	}, nil
}

func (injector *tokenInjector) Inject(_ context.Context, request *http.Request) error {
	request.Header.Set("Authorization", "Bearer "+injector.accessToken)
	return nil
}
