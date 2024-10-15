package schema

import "github.com/MaiMee1/go-apispec/oas/v3"

func Integer(format oas.Format, opts ...Option) oas.Schema {
	schema := &oas.Schema{
		Type:   oas.IntegerType,
		Format: format,
	}
	for _, opt := range opts {
		opt.apply(schema)
	}
	return *schema
}
