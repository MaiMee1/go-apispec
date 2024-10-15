package specs

import (
	"encoding/json"

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
