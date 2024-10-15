package schema

import "github.com/MaiMee1/go-apispec/oas/v3"

func Number(format oas.Format, opts ...Option) oas.Schema {
	schema := &oas.Schema{
		Type:   oas.NumberType,
		Format: format,
	}
	for _, opt := range opts {
		opt.apply(schema)
	}
	return *schema
}
