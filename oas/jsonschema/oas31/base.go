package oas31

import (
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type SpecificationExtension map[string]interface{}

type OASMixin struct {
	Example       interface{}            `json:"examples,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Discriminator *Discriminator         `json:"discriminator,omitempty"`
	Xml           *XML                   `json:"xml,omitempty"`
	Extensions    SpecificationExtension `json:"-"`
}

func (s OASMixin) Kind() jsonschema.Kind {
	return jsonschema.Annotation
}

func (s OASMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.ObjectType)
}

func (s OASMixin) Validate(v interface{}) error {
	return nil
}

type ExternalDocumentation struct {
	Description string                 `json:"description,omitempty"`
	Url         string                 `json:"url,omitempty" validate:"required,url"`
	Extensions  SpecificationExtension `json:"-"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName,omitempty" validate:"required"`
	Mapping      map[string]string `json:"mapping,omitempty" validate:"dive,uri-reference"`
}

type XML struct {
	Name       string                 `json:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty" validate:"omitempty,uri"`
	Prefix     string                 `json:"prefix,omitempty"`
	Attribute  bool                   `json:"attribute,omitempty"`
	Wrapped    bool                   `json:"wrapped,omitempty"`
	Extensions SpecificationExtension `json:"-"`
}
