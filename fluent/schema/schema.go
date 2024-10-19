package schema

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

func New(typ reflect.Type, opts ...Option) oas.Schema {
	schema := objectSchema(typ)
	for _, opt := range opts {
		opt.apply(&schema)
	}
	return schema
}

func For[T any](opts ...Option) oas.Schema {
	typ := reflect.TypeFor[T]()
	return New(typ, opts...)
}

func RefFor[T any](opts ...Option) oas.Reference {
	typ := reflect.TypeFor[T]()
	schema := objectSchema(typ)
	for _, opt := range opts {
		opt.apply(&schema)
	}
	if schema.Type.Has(jsonschema.ObjectType) {
		ref := oas.Reference{
			Ref: fmt.Sprintf("#/components/schemas/%s", schema.Extensions["Name"].(string)),
		}
		return ref
	}
	panic(fmt.Errorf("%v should be convertable to an object", typ))
}

func Cached() *sync.Map {
	return &cache
}
