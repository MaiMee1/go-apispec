package oas31

import (
	"iter"
	"reflect"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/abc"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/draft2020"
)

var zero Schema

type Schema struct {
	draft2020.MetaSchemaMixin[*Schema]
	draft2020.ReferenceMixin[Schema]
	draft2020.MetaDataMixin
	draft2020.ValidationMixin
	draft2020.StringMixin
	draft2020.NumericMixin
	draft2020.ObjectMixin[*Schema]
	draft2020.ArrayMixin[*Schema]
	draft2020.UnevaluatedMixin[*Schema]
	draft2020.ApplicatorMixin[*Schema]
	OASMixin
}

func (m *Schema) Keywords() iter.Seq[jsonschema.Keyword] {
	return func(yield func(jsonschema.Keyword) bool) {
		if !reflect.DeepEqual(m.MetaSchemaMixin, zero.MetaDataMixin) {
			if !yield(&m.MetaSchemaMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.ReferenceMixin, zero.ReferenceMixin) {
			if !yield(&m.ReferenceMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.MetaDataMixin, zero.MetaDataMixin) {
			if !yield(&m.MetaDataMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.ValidationMixin, zero.ValidationMixin) {
			if !yield(&m.ValidationMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.StringMixin, zero.StringMixin) {
			if !yield(&m.StringMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.NumericMixin, zero.NumericMixin) {
			if !yield(&m.NumericMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.ObjectMixin, zero.ObjectMixin) {
			if !yield(&m.ObjectMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.ArrayMixin, zero.ArrayMixin) {
			if !yield(&m.ArrayMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.UnevaluatedMixin, zero.UnevaluatedMixin) {
			if !yield(&m.UnevaluatedMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.ApplicatorMixin, zero.ApplicatorMixin) {
			if !yield(&m.ApplicatorMixin) {
				return
			}
		}
		if !reflect.DeepEqual(m.OASMixin, zero.OASMixin) {
			if !yield(&m.OASMixin) {
				return
			}
		}
	}
}

func (m *Schema) Kind() jsonschema.Kind {
	return abc.Kind(m)
}

func (m *Schema) AppliesTo(t jsonschema.Type) bool {
	return abc.AppliesTo(m, t)
}

func (m *Schema) Validate(v interface{}) error {
	return abc.Validate(m, v)
}
