package draft2020

import (
	"github.com/MaiMee1/go-apispec/oas/iana"
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type Encoding string

type ContentMixin[S jsonschema.Keyword] struct {
	ContentEncoding  Encoding       `json:"contentEncoding,omitempty"`
	ContentMediaType iana.MediaType `json:"contentMediaType,omitempty"`
	ContentSchema    S              `json:"contentSchema,omitempty"`
}

func (m *ContentMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Annotation
}

func (m *ContentMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.StringType)
}

func (m *ContentMixin[S]) Validate(v interface{}) error {
	return nil
}
