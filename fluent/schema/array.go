package schema

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

func Array(item oas.Schema, opts ...Option) oas.Schema {
	schema := &oas.Schema{
		Type: jsonschema.ArrayType,
		Items: &oas.ValueOrReferenceOf[oas.Schema]{
			Value: item,
		},
	}
	for _, opt := range opts {
		opt.apply(schema)
	}
	return *schema
}
