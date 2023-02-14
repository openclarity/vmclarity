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
package gorm

import (
	"context"
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/CiscoM31/godata"

	"gorm.io/datatypes"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var fixSelectToken sync.Once

type ODataObject struct {
	gorm.Model
	Data datatypes.JSON
}

var schemaMeta = map[string]schema{
	"ScanConfig": {
		"id":                 fieldMeta{collection: false, primitive: true},
		"name":               fieldMeta{collection: false, primitive: true},
		"scanFamiliesConfig": fieldMeta{collection: false, primitive: false, ty: []string{"ScanFamiliesConfig"}},
		"scheduled":          fieldMeta{collection: false, primitive: false, ty: []string{"SingleScheduleScanConfig"}},
		"scope":              fieldMeta{collection: false, primitive: false, ty: []string{"AwsScanScope"}},
	},
	"ScanFamiliesConfig": {
		"exploits": fieldMeta{collection: false, primitive: false, ty: []string{"Exploits"}},
		"sbom":     fieldMeta{collection: false, primitive: false, ty: []string{"Sbom"}},
	},
	"Exploits": {
		"enabled": fieldMeta{collection: false, primitive: true},
	},
	"Sbom": {
		"enabled": fieldMeta{collection: false, primitive: true},
	},
	"SingleScheduleScanConfig": {
		"operationTime": fieldMeta{collection: false, primitive: true},
	},
	"AwsScanScope": {
		"all":                        fieldMeta{primitive: true},
		"instanceTagExclusion":       fieldMeta{collection: true, ty: []string{"Tag"}},
		"instanceTagSelector":        fieldMeta{collection: true, ty: []string{"Tag"}},
		"regions":                    fieldMeta{collection: true, ty: []string{"AwsRegion"}},
		"shouldScanStoppedInstances": fieldMeta{primitive: true},
	},
	"Tag": {
		"key":   fieldMeta{primitive: true},
		"value": fieldMeta{primitive: true},
	},
	"AwsRegion": {
		"id":   fieldMeta{primitive: true},
		"name": fieldMeta{primitive: true},
		"vpcs": fieldMeta{collection: true, ty: []string{"AwsVPC"}},
	},
	"AwsVPC": {
		"id":             fieldMeta{primitive: true},
		"name":           fieldMeta{primitive: true},
		"securityGroups": fieldMeta{collection: true, ty: []string{"AwsSecurityGroup"}},
	},
	"AwsSecurityGroup": {
		"id":   fieldMeta{primitive: true},
		"name": fieldMeta{primitive: true},
	},
}

// nolint:cyclop
func ODataQuery(db *gorm.DB, table string, schema string, filter *string, selectString *string, collection bool, result interface{}) error {
	// Fix GlobalExpandTokenizer so that it allows for `-` characters in the Literal tokens
	fixSelectToken.Do(func() {
		godata.GlobalExpandTokenizer.Add("^[a-zA-Z0-9_\\'\\.:\\$ \\*-]+", godata.ExpandTokenLiteral)
	})

	// Parse top level $filter and create the top level "WHERE"
	var where string
	if filter != nil && *filter != "" {
		filterQuery, err := godata.ParseFilterString(context.TODO(), *filter)
		if err != nil {
			return fmt.Errorf("failed to parse $filter: %w", err)
		}

		// Build the WHERE conditions based on the $filter tree
		conditions, err := buildWhereFromFilter("Data", filterQuery.Tree)
		if err != nil {
			return fmt.Errorf("failed to build DB query from $filter: %w", err)
		}

		where = fmt.Sprintf("WHERE %s", conditions)
	}

	selectItems := []*godata.ExpandItem{}
	if selectString != nil && *selectString != "" {
		// Parse $select using ParseExpandString because godata.ParseSelectString
		// is a nieve implementation and doesn't handle query options properly
		expandQuery, err := godata.ParseExpandString(context.TODO(), *selectString)
		if err != nil {
			return fmt.Errorf("failed to parse $select: %w", err)
		}
		selectItems = expandQuery.ExpandItems
	}

	// Turn the select query into a tree that can be used to build nested
	// select queries for all the embedded types.
	selectTree := buildSelectTreeFromSelect(selectItems)

	// Build query selecting fields based on the selectTree
	selectFields := buildSelectFields("$", "Data", schema, schemaMeta["ScanConfig"], selectTree)

	query := fmt.Sprintf("SELECT ID, %s AS Data FROM %s %s", selectFields, table, where)
	if collection {
		if err := db.Raw(query).Find(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	} else {
		query := fmt.Sprintf("%s LIMIT 1", query)
		if err := db.Raw(query).Find(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	}
	return nil
}

type fieldMeta struct {
	collection bool
	primitive  bool
	ty         []string
}

type schema map[string]fieldMeta

type selectNode struct {
	children map[string]*selectNode
	filter   *godata.GoDataFilterQuery
}

func newSelectTree() *selectNode {
	return &selectNode{
		children: map[string]*selectNode{},
	}
}

func (st *selectNode) insert(ts []*godata.Token, filter *godata.GoDataFilterQuery) {
	if len(ts) == 0 {
		return
	}

	childName := ts[0].Value

	child, ok := st.children[childName]
	if !ok {
		st.children[childName] = newSelectTree()
		child = st.children[childName]
	}

	if len(ts) == 1 {
		child.filter = filter
	}

	child.insert(ts[1:], filter)
}

func buildSelectTreeFromSelect(si []*godata.ExpandItem) *selectNode {
	tree := newSelectTree()
	for _, s := range si {
		tree.insert(s.Path, s.Filter)
	}
	return tree
}

// nolint:cyclop,gocognit
func buildSelectFields(path string, source string, schemaName string, sch schema, st *selectNode) string {
	parts := []string{fmt.Sprintf("'objectType', '%s'", schemaName)}
	for key, meta := range sch {
		var where string
		var sel *selectNode
		if st != nil {
			var ok bool
			sel, ok = st.children[key]
			// If any children specified, but this key isn't one of them
			// then skip this key.
			if len(st.children) > 0 && !ok {
				continue
			}
		}

		var extract string
		if meta.primitive {
			// Primitive
			extract = fmt.Sprintf("%s -> '%s.%s'", source, path, key)
		} else if meta.collection {
			// List
			newSource := fmt.Sprintf("%s.value", key)
			if sel != nil && sel.filter != nil {
				conditions, _ := buildWhereFromFilter(newSource, sel.filter.Tree)
				where = fmt.Sprintf("WHERE %s", conditions)
			}

			var subQuery string
			if len(meta.ty) == 1 {
				subQuery = buildSelectFields("$", newSource, meta.ty[0], schemaMeta[meta.ty[0]], sel)
			} else {
				objects := []string{}
				for _, t := range meta.ty {
					objects = append(objects, buildSelectFields("$", newSource, t, schemaMeta[t], sel))
				}
				subQuery = fmt.Sprintf("(SELECT %sOptions.value FROM JSON_EACH(JSON_ARRAY(%s)) AS %sOptions WHERE %sOptions.value -> '$.objectType' = %s -> '%s.objectType')", key, strings.Join(objects, ","), key, key, newSource, path)
			}

			extract = fmt.Sprintf("(SELECT JSON_GROUP_ARRAY(%s) FROM JSON_EACH(%s, '%s.%s') AS %s %s)", subQuery, source, path, key, key, where)
		} else {
			// Struct
			if len(meta.ty) == 1 {
				extract = buildSelectFields(fmt.Sprintf("%s.%s", path, key), source, meta.ty[0], schemaMeta[meta.ty[0]], sel)
			} else {
				objects := []string{}
				for _, t := range meta.ty {
					objects = append(objects, buildSelectFields(fmt.Sprintf("%s.%s", path, key), source, t, schemaMeta[t], sel))
				}
				extract = fmt.Sprintf("(SELECT %s.value FROM JSON_EACH(JSON_ARRAY(%s)) AS %s WHERE %s.value -> '$.objectType' = %s -> '%s.%s.objectType')", key, strings.Join(objects, ","), key, key, source, path, key)
			}
		}
		part := fmt.Sprintf("'%s', %s", key, extract)
		parts = append(parts, part)
	}

	return fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(parts, ","))
}

var sqlOperators = map[string]string{
	"eq":         "=",
	"ne":         "!=",
	"gt":         ">",
	"ge":         ">=",
	"lt":         "<",
	"le":         "<=",
	"or":         "or",
	"contains":   "%%%s%%",
	"endswith":   "%%%s",
	"startswith": "%s%%",
}

// nolint:cyclop
func buildWhereFromFilter(source string, node *godata.ParseNode) (string, error) {
	operator := node.Token.Value

	var query string
	switch operator {
	case "eq", "ne", "gt", "ge", "lt", "le":
		queryField := node.Children[0].Token.Value
		// TODO Possibly convert "slash paths" to "dot paths"
		queryPath := fmt.Sprintf("$.%s", queryField)

		right := node.Children[1].Token.Value
		var value string
		switch node.Children[1].Token.Type {
		case godata.ExpressionTokenString:
			value = strings.ReplaceAll(right, "'", "\"")
		case godata.ExpressionTokenBoolean:
			value = right
		}

		query = fmt.Sprintf("%s -> '%s' %s '%s'", source, queryPath, sqlOperators[operator], value)
	case "and":
		left, err := buildWhereFromFilter(source, node.Children[0])
		if err != nil {
			return query, err
		}
		right, err := buildWhereFromFilter(source, node.Children[1])
		if err != nil {
			return query, err
		}
		query = fmt.Sprintf("(%s AND %s)", left, right)
	case "or":
		left, err := buildWhereFromFilter(source, node.Children[0])
		if err != nil {
			return query, err
		}
		right, err := buildWhereFromFilter(source, node.Children[1])
		if err != nil {
			return query, err
		}
		query = fmt.Sprintf("(%s OR %s)", left, right)
	case "contains", "endswith", "startswith":
		queryField := node.Children[0].Token.Value
		queryPath := fmt.Sprintf("$.%s", queryField)

		right := node.Children[1].Token.Value
		var value interface{}
		switch node.Children[1].Token.Type {
		case godata.ExpressionTokenString:
			r := strings.ReplaceAll(right, "'", "")
			value = fmt.Sprintf(sqlOperators[operator], r)
		default:
			return query, fmt.Errorf("unsupported token type")
		}
		query = fmt.Sprintf("%s -> '%s' LIKE '%s'", source, queryPath, value)
	}

	return query, nil
}
