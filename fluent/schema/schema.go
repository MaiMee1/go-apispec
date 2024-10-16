package schema

import (
	"reflect"

	"github.com/MaiMee1/go-apispec/oas/v3"
)

func New(typ reflect.Type, opts ...Option) oas.Schema {
	_, schema := defineObject(typ, false)
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
	name, schema := defineObject(typ, false)
	for _, opt := range opts {
		opt.apply(&schema)
	}
	ref := oas.Reference{
		Ref: makeRef(name),
	}
	return ref
}

func Cached() map[string]*oas.Schema {
	return cache
}
