package schema

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

func Array(item oas.Schema, opts ...Option) (schema oas.Schema) {
	schema.Type = jsonschema.ArrayType
	schema.Items.Y = &item
	for _, opt := range opts {
		opt.apply(&schema)
	}
	return schema
}
