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

package gorm

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database/types"
	"gorm.io/gorm"
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

type ScopesTableHandler struct {
	scopesTable *gorm.DB
}

func (db *Handler) ScopesTable() types.ScopesTable {
	return &ScopesTableHandler{
		scopesTable: db.DB.Table(scopeTableName),
	}
}

func (s ScopesTableHandler) GetScopes() (models.ScopeType, error) {
	var scopes models.ScopeType

	return scopes, fmt.Errorf("GetScopes not implemented")
}

func (s ScopesTableHandler) SetScopes(scopes models.ScopeType) (models.ScopeType, error) {
	return scopes, fmt.Errorf("GetScopes not implemented")
}
