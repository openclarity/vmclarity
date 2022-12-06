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

	"github.com/labstack/echo/v4"

	"github.com/openclarity/vmclarity/api/models"
)

func (s *ServerImpl) DeleteTargetTargetID(
	ctx echo.Context,
	targetID models.TargetID,
) error {
	return nil
}

func (s *ServerImpl) GetTargetTargetID(
	ctx echo.Context,
	targetID models.TargetID,
) error {
	return nil
}

func (s *ServerImpl) PutTargetTargetID(
	ctx echo.Context,
	targetID models.TargetID,
) error {
	return nil
}

func (s *ServerImpl) GetTargets(
	ctx echo.Context,
	params models.GetTargetsParams,
) error {
	return ctx.JSON(http.StatusOK, []models.Target{})
}

func (s *ServerImpl) PostTargets(
	ctx echo.Context,
) error {
	return nil
}
