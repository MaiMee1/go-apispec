package draft2020

import (
	"github.com/MaiMee1/go-apispec/oas/iana"
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type Encoding string

type ContentMixin struct {
	ContentEncoding  Encoding       `json:"contentEncoding,omitempty"`
	ContentMediaType iana.MediaType `json:"contentMediaType,omitempty"`
	ContentSchema    Schema         `json:"contentSchema,omitempty"`
}

func (m *ContentMixin) Kind() jsonschema.Kind {
	return jsonschema.Annotation
}

func (m *ContentMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.StringType)
}

func (m *ContentMixin) Validate(v interface{}) error {
	return nil
}
