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
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/openclarity/vmclarity/api/models"
)

const (
	targetsTableName = "targets"
)

type Target struct {
	gorm.Model

	Type     string `json:"type" json:"type,omitempty" gorm:"column:type"`
	Location string `json:"location,omitempty" gorm:"column:location"`

	// VMInfo
	InstanceID       string `json:"instance_id,omitempty" gorm:"column:instance_id"`
	InstanceProvider string `json:"instance_provider,omitempty" gorm:"column:instance_provider"`

	// PodInfo
	PodName string `json:"pod_name,omitempty" gorm:"column:pod_name"`

	// DirInfo
	DirName string `json:"dir_name,omitempty" gorm:"column:dir_name"`
}

type TargetsTable interface {
	GetTargetsAndTotal(params models.GetTargetsParams) ([]*Target, int64, error)
	GetTarget(targetID models.TargetID) (*Target, error)
	CheckVMInfoExists(instanceID string, location string) (*Target, bool, error)
	CreateTarget(target *Target) (*Target, error)
	UpdateTarget(target *Target, targetID models.TargetID) (*Target, error)
	DeleteTarget(targetID models.TargetID) error
}

type TargetsTableHandler struct {
	targetsTable *gorm.DB
}

func (db *Handler) TargetsTable() TargetsTable {
	return &TargetsTableHandler{
		targetsTable: db.DB.Table(targetsTableName),
	}
}

func (t *TargetsTableHandler) CheckVMInfoExists(instanceID string, location string) (*Target, bool, error) {
	var target *Target

	if err := t.targetsTable.Where("instance_id = ? AND location = ?", instanceID, location).First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return target, true, nil
}

func (t *TargetsTableHandler) GetTargetsAndTotal(params models.GetTargetsParams) ([]*Target, int64, error) {
	var count int64
	var targets []*Target

	tx := t.targetsTable

	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count total: %w", err)
	}

	if err := tx.Find(&targets).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find targets: %w", err)
	}

	return targets, count, nil
}

func (t *TargetsTableHandler) GetTarget(targetID models.TargetID) (*Target, error) {
	var target *Target

	if err := t.targetsTable.Where("id = ?", targetID).First(&target).Error; err != nil {
		return nil, fmt.Errorf("failed to get target by id %q: %w", targetID, err)
	}

	return target, nil
}

func (t *TargetsTableHandler) CreateTarget(target *Target) (*Target, error) {
	if err := t.targetsTable.Create(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (t *TargetsTableHandler) UpdateTarget(target *Target, targetID models.TargetID) (*Target, error) {
	id, err := strconv.Atoi(targetID)
	if err != nil {
		return nil, err
	}
	target.ID = uint(id)
	t.targetsTable.Save(target)

	return target, err
}

func (t *TargetsTableHandler) DeleteTarget(targetID models.TargetID) error {
	if err := t.targetsTable.Delete(&Scan{}, targetID).Error; err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}
	return nil
}
