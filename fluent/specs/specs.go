package specs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

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

func WithSchemaDefinitions(definitions *sync.Map) Option {
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
				}
				if schema.Value.Items.Reference == nil {
					schema = schema.Value.Items
				}
			}
			if schema.Value.Type.Has(jsonschema.ObjectType) {
				//typ := schema.Value.Extensions["GoType"].(reflect.Type)
				//name := makeName(typ)
				//schema.Value.Extensions["Name"] = name
				//schema.Reference = &oas.Reference{
				//	Ref: fmt.Sprintf("#/components/schemas/%v", name),
				//}
				needed = append(needed, schema)
			}
		}
		definitions.Range(func(key, value interface{}) bool {
			name, schema := key.(string), value.(*oas.Schema)
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
			return true
		})
	})
}

func makeName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		if t.Kind() == reflect.Interface {
			name = "any // exclude"
		} else {
			name = "ptr" + makeName(t.Elem())
		}
	}
	pkgPath := t.PkgPath()
	if pkgPath != "." {
		pkgPath += "."
	}
	fullName := pkgPath + name
	fullName = strings.ReplaceAll(fullName, "/", ".") // updated
	return strings.ReplaceAll(fullName, "-", "_")
}
