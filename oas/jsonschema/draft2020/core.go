package draft2020

import (
	"context"

	"github.com/MaiMee1/go-apispec/oas/jsonpointer"
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type MetaSchemaMixin[S jsonschema.Keyword] struct {
	Schema        string          `json:"$schema,omitempty" validate:"omitempty,uri"`
	Id            string          `json:"$id,omitempty" validate:"omitempty,uri-reference"`
	Comments      string          `json:"$comments,omitempty"`
	Defs          map[string]S    `json:"$defs,omitempty" validate:"dive"`
	Anchor        string          `json:"$anchor,omitempty"`
	DynamicAnchor string          `json:"$dynamicAnchor,omitempty"`
	Vocabulary    map[string]bool `json:"$vocabulary,omitempty" validate:"dive,keys,uri-reference,endkeys"`
}

func (m *MetaSchemaMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Identifier | jsonschema.ReservedLocation
}

func (m *MetaSchemaMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *MetaSchemaMixin[S]) Validate(v interface{}) error {
	_ = v.(*Schema)
	return nil
}

type ReferenceMixin[S any] struct {
	Ref        string `json:"$ref,omitempty" validate:"uri-reference"`
	DynamicRef string `json:"$dynamicRef,omitempty" validate:"uri-reference"`

	ctx context.Context
}

func (m *ReferenceMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Applicator
}

func (m *ReferenceMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *ReferenceMixin[S]) Validate(v interface{}) error {
	return nil
}

func (m *ReferenceMixin[S]) WithContext(ctx context.Context) *ReferenceMixin[S] {
	m.ctx = ctx
	return m
}

func (m *ReferenceMixin[S]) Resolve() S {
	root := m.ctx.Value("Root")
	if root == nil {
		panic("root not available")
	}

	v, err := jsonpointer.UriFragment(m.Ref).Access(root)
	if err != nil {
		panic(err)
	}
	return v.Interface().(S)
}
