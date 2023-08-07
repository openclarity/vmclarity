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

package envbearer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/openclarity/vmclarity/pkg/apiserver/iam"
)

// New creates an injector which adds Bearer token read from env variable to request.
func New(tokenEnv string) (iam.Injector, error) {
	if tokenEnv == "" {
		return nil, fmt.Errorf("cannot use empty env variable for Injector=envbearer")
	}

	return &bearer{
		tokenEnv: tokenEnv,
	}, nil
}

type bearer struct {
	tokenEnv string
}

func (injector *bearer) Inject(_ context.Context, request *http.Request) error {
	request.Header.Set("Authorization", "Bearer "+os.Getenv(injector.tokenEnv))
	return nil
}
