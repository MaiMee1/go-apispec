package draft2020

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

type ReferenceMixin[S jsonschema.Keyword] struct {
	Ref        string `json:"$ref,omitempty" validate:"uri-reference"`
	DynamicRef string `json:"$dynamicRef,omitempty" validate:"uri-reference"`
}

func (m *ReferenceMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Applicator
}

func (m *ReferenceMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *ReferenceMixin[S]) Validate(v interface{}) error {
	if m.Ref == "" && m.DynamicRef == "" {
		return nil
	}
	if m.Ref != "" {
		//TODO implement me
	}
	panic("implement me")
}

func (m *ReferenceMixin[S]) Resolve(ctx context.Context) S {
	root := ctx.Value("Root")
	if root == nil {
		panic("root not available")
	}

	var schema S
	b, err := json.Marshal(resolve(m.Ref, root))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &schema); err != nil {
		panic(err)
	}
	return schema
}

func resolve(uri string, v interface{}) interface{} {
	parts := strings.Split(uri, "/")
	if parts[0] != "#" {
		panic(fmt.Errorf("invalid reference format: %s", uri))
	}
	for _, part := range parts[1:] {
		t, ok := v.(map[string]interface{})[part]
		if !ok {
			panic(fmt.Sprintf("key not found: %s", part))
		}
		v = t
	}
	return v
}
