package schema

import (
	"fmt"
	"reflect"

	"github.com/MaiMee1/go-apispec/fluent/schema/encoder"
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/draft2020"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

var enc = encoder.New()

func WithEncoder(opts ...encoder.Option) {
	enc = encoder.New(opts...)
}

func New(typ reflect.Type, opts ...Option) oas.Schema {
	schema := enc.Encode(typ)
	for _, opt := range opts {
		opt.apply(&schema)
	}
	return schema
}

func For[T any](opts ...Option) oas.Schema {
	typ := reflect.TypeFor[T]()
	return New(typ, opts...)
}

func RefFor[T any](opts ...Option) oas.Schema {
	typ := reflect.TypeFor[T]()
	schema := enc.Encode(typ)
	for _, opt := range opts {
		opt.apply(&schema)
	}
	if schema.Type.Has(jsonschema.ObjectType) {
		return oas.Schema{
			ReferenceMixin: draft2020.ReferenceMixin[oas.Schema]{
				Ref: fmt.Sprintf("#/components/schemas/%s", schema.Extensions["Name"].(string)),
			},
		}
	}
	return schema
}

func Cached() map[string]*oas.Schema {
	return enc.Cache()
}
