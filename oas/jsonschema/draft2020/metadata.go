package draft2020

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type MetaDataMixin struct {
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	Examples    []interface{} `json:"examples,omitempty"`
	Deprecated  bool          `json:"deprecated,omitempty"`
	ReadOnly    bool          `json:"readOnly,omitempty"`
	WriteOnly   bool          `json:"writeOnly,omitempty"`
}

func (m *MetaDataMixin) Kind() jsonschema.Kind {
	return jsonschema.Annotation
}

func (m *MetaDataMixin) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *MetaDataMixin) Validate(v interface{}) error {
	return nil
}
