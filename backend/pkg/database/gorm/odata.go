// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gorm

import (
	"fmt"

	"github.com/openclarity/vmclarity/backend/pkg/database/odatasql"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ODataObject struct {
	gorm.Model
	Data datatypes.JSON
}

var schemaMetas = map[string]odatasql.SchemaMeta{
	"ScanConfig": {
		Table: "scan_configs",
		Fields: odatasql.Schema{
			"id":                 odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"name":               odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"scanFamiliesConfig": odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"ScanFamiliesConfig"},
			},
			"scheduled":          odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"SingleScheduleScanConfig"},
			},
			"scope":              odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"AwsScanScope"},
				DescriminatorProperty: "objectType",
			},
		},
	},
	"ScanFamiliesConfig": {
		Fields: odatasql.Schema{
			"exploits":          odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"ExploitsConfig"},
			},
			"malware":           odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"MalwareConfig"},
			},
			"misconfigurations": odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"MisconfigurationsConfig"},
			},
			"rootkits":          odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"RootkitsConfig"},
			},
			"sbom":              odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"SBOMConfig"},
			},
			"secrets":           odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"SecretsConfig"},
			},
			"vulnerabilties":    odatasql.FieldMeta{
				FieldType: odatasql.ComplexFieldType,
				ComplexFieldSchemas: []string{"VulnerabiltiesConfig"},
			},
		},
	},
	"ExploitsConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"MalwareConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"MisconfigurationsConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"RootkitsConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"SBOMConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"SecretsConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"VulnerabilitiesConfig": {
		Fields: odatasql.Schema{
			"enabled": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"SingleScheduleScanConfig": {
		Fields: odatasql.Schema{
			"operationTime": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"AwsScanScope": {
		Fields: odatasql.Schema{
			"objectType": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"all":                        odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"shouldScanStoppedInstances": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"instanceTagExclusion":       odatasql.FieldMeta{
				FieldType: odatasql.CollectionFieldType,
				CollectionItemMeta: &odatasql.FieldMeta{
					FieldType: odatasql.ComplexFieldType,
					ComplexFieldSchemas: []string{"Tag"},
				},
			},
			"instanceTagSelector":        odatasql.FieldMeta{
				FieldType: odatasql.CollectionFieldType,
				CollectionItemMeta: &odatasql.FieldMeta{
					FieldType: odatasql.ComplexFieldType,
					ComplexFieldSchemas: []string{"Tag"},
				},
			},
			"regions":                    odatasql.FieldMeta{
				FieldType: odatasql.CollectionFieldType,
				CollectionItemMeta: &odatasql.FieldMeta{
					FieldType: odatasql.ComplexFieldType,
					ComplexFieldSchemas: []string{"AwsRegion"},
				},
			},
		},
	},
	"Tag": {
		Fields: odatasql.Schema{
			"key":   odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"value": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
	"AwsRegion": {
		Fields: odatasql.Schema{
			"id":   odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"name": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"vpcs": odatasql.FieldMeta{
				FieldType: odatasql.CollectionFieldType,
				CollectionItemMeta: &odatasql.FieldMeta{
					FieldType: odatasql.ComplexFieldType,
					ComplexFieldSchemas: []string{"AwsVPC"},
				},
			},
		},
	},
	"AwsVPC": {
		Fields: odatasql.Schema{
			"id":             odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"name":           odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"securityGroups": odatasql.FieldMeta{
				FieldType: odatasql.CollectionFieldType,
				CollectionItemMeta: &odatasql.FieldMeta{
					FieldType: odatasql.ComplexFieldType,
					ComplexFieldSchemas: []string{"AwsSecurityGroup"},
				},
			},
		},
	},
	"AwsSecurityGroup": {
		Fields: odatasql.Schema{
			"id":   odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
			"name": odatasql.FieldMeta{FieldType: odatasql.PrimitiveFieldType},
		},
	},
}

func ODataQuery(db *gorm.DB, schema string, filterString, selectString, expandString *string, top, skip *int, collection bool, result interface{}) error {
	// If we're not getting a collection, make sure the result is limited
	// to 1 item.
	if !collection {
		top = utils.IntPtr(1)
		skip = nil
	}

	// Build the raw SQL query using the odatasql library, this will also
	// parse and validate the ODATA query params.
	query, err := odatasql.BuildSQLQuery(schemaMetas, schema, filterString, selectString, expandString, top, skip)
	if err != nil {
		return fmt.Errorf("failed to build query for DB: %w", err)
	}

	// Use the query to populate "result" using the gorm finalisers so that
	// the gorm error handling processes things like no results found.
	if collection {
		if err := db.Raw(query).Find(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	} else {
		if err := db.Raw(query).First(result).Error; err != nil {
			return fmt.Errorf("failed to query DB: %w", err)
		}
	}
	return nil
}

/*
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

		subQuery := buildSelectFields(*field.CollectionItemMeta, newIdentifier, newSource, "$", newSelectNode)
		return fmt.Sprintf("(SELECT JSON_GROUP_ARRAY(%s) FROM JSON_EACH(%s, '%s') AS %s %s)", subQuery, source, path, identifier, where)
	case odatasql.ComplexFieldType:
		objects := []string{}
		for _, schemaName := range field.ComplexFieldSchemas {
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
*/
