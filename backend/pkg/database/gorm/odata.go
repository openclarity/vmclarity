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
	"fmt"
	"strings"
	"sync"

	"github.com/CiscoM31/godata"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var fixSelectToken sync.Once

type ODataObject struct {
	gorm.Model
	Data datatypes.JSON
}

type fieldType int

const (
	primitiveFieldType fieldType = iota
	complexFieldType
	collectionFieldType
)

type fieldMeta struct {
	fieldType           fieldType
	collectionItemMeta  *fieldMeta
	complexFieldSchemas []string
}

type schema map[string]fieldMeta

var schemaMeta = map[string]schema{
	"ScanConfig": {
		"id":                 fieldMeta{fieldType: primitiveFieldType},
		"name":               fieldMeta{fieldType: primitiveFieldType},
		"scanFamiliesConfig": fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"ScanFamiliesConfig"}},
		"scheduled":          fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"SingleScheduleScanConfig"}},
		"scope":              fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"AwsScanScope"}},
	},
	"ScanFamiliesConfig": {
		"exploits":          fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"ExploitsConfig"}},
		"malware":           fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"MalwareConfig"}},
		"misconfigurations": fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"MisconfigurationsConfig"}},
		"rootkits":          fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"RootkitsConfig"}},
		"sbom":              fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"SBOMConfig"}},
		"secrets":           fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"SecretsConfig"}},
		"vulnerabilties":    fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"VulnerabiltiesConfig"}},
	},
	"ExploitsConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"MalwareConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"MisconfigurationsConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"RootkitsConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"SBOMConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"SecretsConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"VulnerabilitiesConfig": {
		"enabled": fieldMeta{fieldType: primitiveFieldType},
	},
	"SingleScheduleScanConfig": {
		"operationTime": fieldMeta{fieldType: primitiveFieldType},
	},
	"AwsScanScope": {
		"all":                        fieldMeta{fieldType: primitiveFieldType},
		"instanceTagExclusion":       fieldMeta{fieldType: collectionFieldType, collectionItemMeta: &fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"Tag"}}},
		"instanceTagSelector":        fieldMeta{fieldType: collectionFieldType, collectionItemMeta: &fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"Tag"}}},
		"regions":                    fieldMeta{fieldType: collectionFieldType, collectionItemMeta: &fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"AwsRegion"}}},
		"shouldScanStoppedInstances": fieldMeta{fieldType: primitiveFieldType},
	},
	"Tag": {
		"key":   fieldMeta{fieldType: primitiveFieldType},
		"value": fieldMeta{fieldType: primitiveFieldType},
	},
	"AwsRegion": {
		"id":   fieldMeta{fieldType: primitiveFieldType},
		"name": fieldMeta{fieldType: primitiveFieldType},
		"vpcs": fieldMeta{fieldType: collectionFieldType, collectionItemMeta: &fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"AwsVPC"}}},
	},
	"AwsVPC": {
		"id":             fieldMeta{fieldType: primitiveFieldType},
		"name":           fieldMeta{fieldType: primitiveFieldType},
		"securityGroups": fieldMeta{fieldType: collectionFieldType, collectionItemMeta: &fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{"AwsSecurityGroup"}}},
	},
	"AwsSecurityGroup": {
		"id":   fieldMeta{fieldType: primitiveFieldType},
		"name": fieldMeta{fieldType: primitiveFieldType},
	},
}

// nolint:cyclop
func ODataQuery(db *gorm.DB, table string, schema string, filter *string, selectString *string, top, skip *int, collection bool, result interface{}) error {
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
	selectFields := buildSelectFields(fieldMeta{fieldType: complexFieldType, complexFieldSchemas: []string{schema}}, schema, "Data", "$", selectTree)

	// Build paging statement
	var limitStm string
	if top != nil || skip != nil {
		limitVal := -1 // Negative means no limit, if no "$top" is specified this is what we want
		if top != nil {
			limitVal = *top
		}
		limitStm = fmt.Sprintf("LIMIT %d", limitVal)

		if skip != nil {
			limitStm = fmt.Sprintf("%s OFFSET %d", limitStm, *skip)
		}
	}

	query := fmt.Sprintf("SELECT ID, %s AS Data FROM %s %s %s", selectFields, table, where, limitStm)
	if collection {
		if err := db.Raw(query).Find(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	} else {
		query := fmt.Sprintf("%s LIMIT 1", query)
		if err := db.Raw(query).First(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	}
	return nil
}

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

// nolint:cyclop
func buildSelectFields(field fieldMeta, identifier, source, path string, st *selectNode) string {
	switch field.fieldType {
	case collectionFieldType:
		newIdentifier := fmt.Sprintf("%sOptions", identifier)
		newSource := fmt.Sprintf("%s.value", identifier)

		var where string
		var newSelectNode *selectNode
		if st != nil {
			if st.filter != nil {
				conditions, _ := buildWhereFromFilter(newSource, st.filter.Tree)
				where = fmt.Sprintf("WHERE %s", conditions)
			}
			newSelectNode = &selectNode{children: st.children}
		}

		subQuery := buildSelectFields(*field.collectionItemMeta, newIdentifier, newSource, "$", newSelectNode)
		return fmt.Sprintf("(SELECT JSON_GROUP_ARRAY(%s) FROM JSON_EACH(%s, '%s') AS %s %s)", subQuery, source, path, identifier, where)
	case complexFieldType:
		objects := []string{}
		for _, schemaName := range field.complexFieldSchemas {
			schema := schemaMeta[schemaName]
			parts := []string{fmt.Sprintf("'objectType', '%s'", schemaName)}
			for key, fm := range schema {
				var sel *selectNode
				if st != nil && len(st.children) > 0 {
					var ok bool
					sel, ok = st.children[key]
					if !ok {
						continue
					}
				}

				extract := buildSelectFields(fm, fmt.Sprintf("%s%s", identifier, key), source, fmt.Sprintf("%s.%s", path, key), sel)
				part := fmt.Sprintf("'%s', %s", key, extract)
				parts = append(parts, part)
			}
			objects = append(objects, fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(parts, ",")))
		}
		if len(objects) == 1 {
			return objects[0]
		}
		return fmt.Sprintf("(SELECT %s.value FROM JSON_EACH(JSON_ARRAY(%s)) AS %s WHERE %s.value -> '$.objectType' = %s -> '%s.objectType')", identifier, strings.Join(objects, ","), identifier, identifier, source, path)
	case primitiveFieldType:
		fallthrough
	default:
		// If root of source (path is just $) is primitive just return the source
		if path == "$" {
			return source
		}
		return fmt.Sprintf("%s -> '%s'", source, path)
	}
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
