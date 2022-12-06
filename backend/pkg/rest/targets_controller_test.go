// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package rest

import (
	"net/http"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTargetsController(t *testing.T) {
	restServer, err := CreateRESTServer(8080)
	if err != nil {
		log.Fatalf("Failed to create REST server: %v", err)
	}

	result := testutil.NewRequest().Get("/targets?page=1&pageSize=1").Go(t, restServer.echoServer)
	assert.Equal(t, http.StatusOK, result.Code())
}
