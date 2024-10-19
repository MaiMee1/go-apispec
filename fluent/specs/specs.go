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
		var needed []*oas.ValueOrReferenceOf[oas.Schema]
		for schema := range api.document.IterSchemaOrRef() {
			if schema.Reference != nil {
				needed = append(needed, schema)
			}
			for schema.Value.Type.Has(jsonschema.ArrayType) {
				if schema.Value.Items == nil {
					// should not but may happen
					break
				} else {
					schema = schema.Value.Items
				}
			}
			if schema.Value.Type.Has(jsonschema.ObjectType) {
				name, ok := schema.Value.Extensions["Name"].(string)
				if ok {
					schema.Reference = &oas.Reference{
						Ref: fmt.Sprintf("#/components/schemas/%v", name),
					}
					needed = append(needed, schema)
				}
			}
		}
		for name, schema := range definitions {
			hasThisName := func(schema *oas.ValueOrReferenceOf[oas.Schema]) bool {
				if schema.Reference != nil && schema.Reference.Ref == fmt.Sprintf("#/components/schemas/%v", name) {
					return true
				}
				return false
				//return name == schema.Value.Extensions["Name"]
			}
			if slices.ContainsFunc(needed, hasThisName) {
				api.document.Components.Schemas[name] = *schema
			}
		}
	})
}
