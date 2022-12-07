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

package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
)

// TODO after db design
type Target struct {
	ID string
}

//go:generate $GOPATH/bin/mockgen --build_flags=--mod=mod -destination=./mock_target.go -package=database github.com/openclarity/vmclarity/backend/pkg/database TargetTable
type TargetTable interface {
	List(params models.GetTargetsParams) ([]models.Target, error)
	Get(targetID models.TargetID) (models.Target, error)
	Create(target *Target) (models.Target, error)
	Update(target *Target, targetID models.TargetID) (models.Target, error)
	Delete(targetID models.TargetID) error
}

type TargetTableHandler struct {
	db *gorm.DB
}

func (db *Handler) TargetTable() TargetTable {
	return &TargetTableHandler{
		db: db.DB,
	}
}

func (t *TargetTableHandler) List(params models.GetTargetsParams) ([]models.Target, error) {
	return []models.Target{}, fmt.Errorf("not implemented")
}

func (t *TargetTableHandler) Get(targetID models.TargetID) (models.Target, error) {
	return models.Target{}, fmt.Errorf("not implemented")
}

func (t *TargetTableHandler) Create(target *Target) (models.Target, error) {
	return models.Target{}, fmt.Errorf("not implemented")
}

func (t *TargetTableHandler) Update(target *Target, targetID models.TargetID) (models.Target, error) {
	return models.Target{}, fmt.Errorf("not implemented")
}

func (t *TargetTableHandler) Delete(targetID models.TargetID) error {
	return fmt.Errorf("not implemented")
}

// TODO after db design
func CreateTarget(target *models.Target) *Target {
	return &Target{
		ID: *target.Id,
	}
}
