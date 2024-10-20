package specs

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

type API struct {
	document oas.OpenAPI
	opts     []Option
}

func New(options ...Option) (*API, error) {
	api := new(API)
	api.document = oas.Default()
	api.opts = append(api.opts, options...)
	for _, opt := range options {
		opt.apply(api)
	}
	return api, api.document.Validate()
}

func (api *API) clone() *API {
	clone, _ := New(api.opts...)
	return clone
}

func (api *API) WithOptions(opts ...Option) *API {
	c := api.clone()
	c.opts = append(c.opts, opts...)
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (api *API) Json() string {
	b, err := json.Marshal(api.document)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func WithSchemaDefinitions(definitions map[string]*oas.Schema) Option {
	return optionFunc(func(api *API) {
		if api.document.Components.Schemas == nil {
			api.document.Components.Schemas = make(map[string]oas.Schema)
		}
		var needed []string
		for schema := range api.document.IterSchema() {
			if schema.Ref != "" {
				needed = append(needed, schema.Ref)
			}
			for schema.Type.Has(jsonschema.ArrayType) {
				if schema.Items != nil && schema.Items.Y != nil {
					schema = schema.Items.Y
				} else {
					break
				}
			}
			if schema.Type.Has(jsonschema.ObjectType) {
				name, ok := schema.Extensions["Name"].(string)
				if ok {
					needed = append(needed, fmt.Sprintf("#/components/schemas/%v", name))
				}
			}
		}
		for name, schema := range definitions {
			if slices.Contains(needed, fmt.Sprintf("#/components/schemas/%v", name)) {
				api.document.Components.Schemas[name] = *schema
			}
		}
		for schema := range api.document.IterSchema() {
			for schema.Type.Has(jsonschema.ArrayType) {
				if schema.Items != nil && schema.Items.Y != nil {
					schema = schema.Items.Y
				} else {
					break
				}
			}
			if schema.Type.Has(jsonschema.ObjectType) {
				name, ok := schema.Extensions["Name"].(string)
				if ok {
					var s oas.Schema
					s.Ref = fmt.Sprintf("#/components/schemas/%v", name)
					*schema = s
				}
			}
		}
	})
}
