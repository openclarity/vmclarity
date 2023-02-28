package odatasql

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/CiscoM31/godata"
)

var fixSelectToken sync.Once

// nolint:cyclop
func BuildSQLQuery(schemaMetas map[string]SchemaMeta, schema string, filterString, selectString, expandString *string, top, skip *int) (string, error) {
	// Fix GlobalExpandTokenizer so that it allows for `-` characters in the Literal tokens
	fixSelectToken.Do(func() {
		godata.GlobalExpandTokenizer.Add("^[a-zA-Z0-9_\\'\\.:\\$ \\*-]+", godata.ExpandTokenLiteral)
	})

	// Parse top level $filter and create the top level "WHERE"
	var where string
	if filterString != nil && *filterString != "" {
		filterQuery, err := godata.ParseFilterString(context.TODO(), *filterString)
		if err != nil {
			return "", fmt.Errorf("failed to parse $filter: %w", err)
		}

		// Build the WHERE conditions based on the $filter tree
		conditions, err := buildWhereFromFilter("Data", filterQuery.Tree)
		if err != nil {
			return "", fmt.Errorf("failed to build DB query from $filter: %w", err)
		}

		where = fmt.Sprintf("WHERE %s", conditions)
	}

	var selectQuery *godata.GoDataSelectQuery
	if selectString != nil && *selectString != "" {
		// NOTE(sambetts):
		// For now we'll won't parse the data here and instead pass
		// just the raw value into the selectTree. The select tree will
		// parse the select query using the ExpandParser. If we can
		// update the GoData select parser to handle paths properly and
		// nest query params then we can switch to parsing select once
		// here before passing it to the selectTree.
		selectQuery = &godata.GoDataSelectQuery{RawValue: *selectString}
	}

	var expandQuery *godata.GoDataExpandQuery
	if expandString != nil && *expandString != "" {
		var err error
		expandQuery, err = godata.ParseExpandString(context.TODO(), "Eggs($filter=Name eq 'Egg1')")
		if err != nil {
			return "", fmt.Errorf("failed to parse $expand ")
		}
	}

	// Turn the select and expand query params into a tree that can be used
	// to build nested select statements for the whole schema.
	//
	// TODO(sambetts) This should probably also validate that all the
	// selected/expanded fields are part of the schema.
	selectTree := newSelectTree()
	err := selectTree.insert(nil, nil, selectQuery, expandQuery, false)
	if err != nil {
		return "", fmt.Errorf("failed to parse select and expand: %w", err)
	}

	// Build query selecting fields based on the selectTree

	// For now all queries must start with a root "object" so we create a
	// complex field meta to represent that object
	rootObject := FieldMeta{FieldType: ComplexFieldType, ComplexFieldSchemas: []string{schema}}
	selectFields := buildSelectFields(schemaMetas, rootObject, schema, "Data", "$", selectTree)

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

	table := schemaMetas[schema].Table
	if table == "" {
		return "", fmt.Errorf("trying to query complex type schema %s with no source table", schema)
	}

	return fmt.Sprintf("SELECT ID, %s AS Data FROM %s %s %s", selectFields, table, where, limitStm), nil
}

// nolint:cyclop,gocognit
func buildSelectFields(schemaMetas map[string]SchemaMeta, field FieldMeta, identifier, source, path string, st *selectNode) string {
	switch field.FieldType {
	case PrimitiveFieldType:
		// If root of source (path is just $) is primitive just return the source
		if path == "$" {
			return source
		}
		return fmt.Sprintf("%s -> '%s'", source, path)
	case CollectionFieldType:
		newIdentifier := fmt.Sprintf("%sOptions", identifier)
		newSource := fmt.Sprintf("%s.value", identifier)

		var where string
		var newSelectNode *selectNode
		if st != nil {
			if st.filter != nil {
				conditions, _ := buildWhereFromFilter(newSource, st.filter.Tree)
				where = fmt.Sprintf("WHERE %s", conditions)
			}
			newSelectNode = &selectNode{children: st.children, expand: st.expand}
		}

		subQuery := buildSelectFields(schemaMetas, *field.CollectionItemMeta, newIdentifier, newSource, "$", newSelectNode)
		return fmt.Sprintf("(SELECT JSON_GROUP_ARRAY(%s) FROM JSON_EACH(%s, '%s') AS %s %s)", subQuery, source, path, identifier, where)
	case ComplexFieldType:
		objects := []string{}
		for _, schemaName := range field.ComplexFieldSchemas {
			schema := schemaMetas[schemaName]

			parts := []string{}
			if field.DescriminatorProperty != "" {
				parts = append(parts, fmt.Sprintf("'%s', '%s'", field.DescriminatorProperty, schemaName))
			}
			for key, fm := range schema.Fields {
				if field.DescriminatorProperty != "" && key == field.DescriminatorProperty {
					continue
				}

				var sel *selectNode
				if st != nil && len(st.children) > 0 {
					var ok bool
					sel, ok = st.children[key]
					if !ok {
						continue
					}
				}

				extract := buildSelectFields(schemaMetas, fm, fmt.Sprintf("%s%s", identifier, key), source, fmt.Sprintf("%s.%s", path, key), sel)
				part := fmt.Sprintf("'%s', %s", key, extract)
				parts = append(parts, part)
			}
			objects = append(objects, fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(parts, ",")))
		}

		if len(objects) == 1 {
			return objects[0]
		}
		if field.DescriminatorProperty == "" {
			//TODO(sambetts) Error, if multiple schema there must be a descriminator
		}
		return fmt.Sprintf("(SELECT %s.value FROM JSON_EACH(JSON_ARRAY(%s)) AS %s WHERE %s.value -> '$.%s' = %s -> '%s.%s')", identifier, strings.Join(objects, ","), identifier, identifier, field.DescriminatorProperty, source, path, field.DescriminatorProperty)

	case RelationshipFieldType:
		if st == nil || !st.expand {
			return fmt.Sprintf("%s -> '%s'", source, path)
		}

		schemaName := field.RelationshipSchema
		schema := schemaMetas[schemaName]
		newsource := fmt.Sprintf("%s.Data", schema.Table)
		parts := []string{fmt.Sprintf("'ObjectType', '%s'", schemaName)}
		for key, fm := range schema.Fields {
			var sel *selectNode
			if st != nil && len(st.children) > 0 {
				var ok bool
				sel, ok = st.children[key]
				if !ok {
					continue
				}
			}

			extract := buildSelectFields(schemaMetas, fm, fmt.Sprintf("%s%s", identifier, key), newsource, fmt.Sprintf("$.%s", key), sel)
			part := fmt.Sprintf("'%s', %s", key, extract)
			parts = append(parts, part)
		}
		object := fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(parts, ","))

		return fmt.Sprintf("(SELECT %s FROM %s WHERE %s -> '$.Id' == %s -> '%s.Id')", object, schema.Table, newsource, source, path)
	case RelationshipCollectionFieldType:
		if st == nil || !st.expand {
			return fmt.Sprintf("%s -> '%s'", source, path)
		}

		schemaName := field.RelationshipSchema
		schema := schemaMetas[schemaName]
		newSource := fmt.Sprintf("%s.Data", schema.Table)

		where := fmt.Sprintf("WHERE %s -> '$.Id' = %s.value -> '$.Id'", newSource, identifier)
		if st != nil {
			if st.filter != nil {
				conditions, _ := buildWhereFromFilter(newSource, st.filter.Tree)
				where = fmt.Sprintf("%s and %s", where, conditions)
			}
		}

		parts := []string{fmt.Sprintf("'ObjectType', '%s'", schemaName)}
		for key, fm := range schema.Fields {
			var sel *selectNode
			if st != nil && len(st.children) > 0 {
				var ok bool
				sel, ok = st.children[key]
				if !ok {
					continue
				}
			}

			extract := buildSelectFields(schemaMetas, fm, fmt.Sprintf("%s%s", identifier, key), newSource, fmt.Sprintf("$.%s", key), sel)
			part := fmt.Sprintf("'%s', %s", key, extract)
			parts = append(parts, part)
		}
		subQuery := fmt.Sprintf("JSON_OBJECT(%s)", strings.Join(parts, ","))

		return fmt.Sprintf("(SELECT JSON_GROUP_ARRAY(%s) FROM %s,JSON_EACH(%s, '%s') AS %s %s)", subQuery, schema.Table, source, path, identifier, where)
	default:
		return ""
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
