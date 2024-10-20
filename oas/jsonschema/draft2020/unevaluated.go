package draft2020

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/ser"
)

type UnevaluatedMixin struct {
	UnevaluatedItems      *ser.Or[bool, *Schema] `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties *ser.Or[bool, *Schema] `json:"unevaluatedProperties,omitempty"`
}

func (m *UnevaluatedMixin) Kind() jsonschema.Kind {
	return jsonschema.Applicator | jsonschema.Annotation
}

func (m *UnevaluatedMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.ObjectType | jsonschema.ArrayType)
}

func (m *UnevaluatedMixin) Validate(v interface{}) error {
	//TODO implement me
	panic("implement me")
}
