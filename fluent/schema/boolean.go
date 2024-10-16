package schema

import "github.com/MaiMee1/go-apispec/oas/v3"

func Boolean(opts ...Option) oas.Schema {
	schema := &oas.Schema{
		Type: oas.BooleanType,
	}
	for _, opt := range opts {
		opt.apply(schema)
	}
	return *schema
}
