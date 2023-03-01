package odatasql

type FieldType string

const (
	PrimitiveFieldType              FieldType = "primitive"
	CollectionFieldType             FieldType = "collection"
	ComplexFieldType                FieldType = "complex"
	RelationshipFieldType           FieldType = "relationship"
	RelationshipCollectionFieldType FieldType = "relationshipCollection"
)

type FieldMeta struct {
	FieldType FieldType

	// Field for collection field types
	CollectionItemMeta *FieldMeta

	// Fields for complex field types
	ComplexFieldSchemas   []string
	DescriminatorProperty string

	// Fields for relationship and relationship collection types
	RelationshipSchema   string
	RelationshipProperty string
}

type SchemaMeta struct {
	Table  string
	Fields map[string]FieldMeta
}

type Schema map[string]FieldMeta
