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
	FieldType          FieldType
	CollectionItemMeta    *FieldMeta
	ComplexFieldSchemas []string
	RelationshipSchema string
	DescriminatorProperty string
}

type SchemaMeta struct {
	Table  string
	Fields map[string]FieldMeta
}

type Schema map[string]FieldMeta
