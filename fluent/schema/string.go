package schema

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

func String(format oas.Format, opts ...Option) (schema oas.Schema) {
	schema.Type = jsonschema.StringType
	schema.Format = format
	for _, opt := range opts {
		opt.apply(&schema)
	}
	return schema
}
