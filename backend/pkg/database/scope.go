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
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	scopeTableName = "scopes"
)

type Scopes struct {
	Base
	Type string `json:"type,omitempty" gorm:"column:type" faker:"-"`

	// AWS Scope
	AwsScopesRegions []AwsScopesRegion `gorm:"foreignKey:RegionID"`
}

type AwsScopesRegion struct {
	RegionID      string         `gorm:"primarykey" faker:"-"`
	AwsRegionVpcs []AwsRegionVpc `gorm:"foreignKey:VpcID"`
}

type AwsRegionVpc struct {
	VpcID                string                `gorm:"primarykey" faker:"-"`
	AwsVpcSecurityGroups []AwsVpcSecurityGroup `gorm:"foreignKey:GroupID"`
}

type AwsVpcSecurityGroup struct {
	GroupID string `gorm:"primarykey" faker:"-"`
}

type ScopesTable interface {
	GetScopes() (*Scopes, error)
	SetScopes(scopes *Scopes) (*Scopes, error)
}

type ScopesTableHandler struct {
	scopesTable *gorm.DB
}

func (db *Handler) ScopesTable() ScopesTable {
	return &ScopesTableHandler{
		scopesTable: db.DB.Table(scopeTableName),
	}
}

func (s ScopesTableHandler) GetScopes() (*Scopes, error) {
	var scopes []Scopes

	if err := s.scopesTable.Preload("AwsScopesRegions.AwsRegionVpcs.AwsVpcSecurityGroups").Preload(clause.Associations).Find(&scopes).Error; err != nil {
		return nil, fmt.Errorf("failed to get scopes: %w", err)
	}

	scopesB, err := json.Marshal(scopes)
	if err != nil {
		log.Errorf("Failed to marshal scopes: %v", err)
	} else {
		fmt.Printf("scopesB=%s\n\n", scopesB)
	}

	return &scopes[0], nil
}

func (s ScopesTableHandler) SetScopes(scopes *Scopes) (*Scopes, error) {
	if err := s.scopesTable.Session(&gorm.Session{FullSaveAssociations: true}).Save(scopes).Error; err != nil {
		return nil, fmt.Errorf("failed to save scopes in db: %w", err)
	}

	return scopes, nil
}
