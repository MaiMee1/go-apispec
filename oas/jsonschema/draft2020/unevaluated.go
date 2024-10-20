package draft2020

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/ser"
)

type UnevaluatedMixin[S jsonschema.Keyword] struct {
	UnevaluatedItems      *ser.Or[bool, S] `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *ser.Or[bool, S] `json:"unevaluatedProperties,omitempty"`
}

func (m *UnevaluatedMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Applicator | jsonschema.Annotation
}

func (m *UnevaluatedMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.ObjectType | jsonschema.ArrayType)
}

func (m *UnevaluatedMixin[S]) Validate(v interface{}) error {
	//TODO implement me
	panic("implement me")
}
