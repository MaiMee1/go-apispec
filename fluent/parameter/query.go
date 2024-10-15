package parameter

import "github.com/MaiMee1/go-apispec/oas/v3"

func Query(name string, description oas.RichText, required bool, opts ...Option) oas.Parameter {
	param := &oas.Parameter{
		In:          oas.QueryLocation,
		Style:       oas.FormStyle,
		Name:        name,
		Description: description,
		Required:    required,
	}
	for _, opt := range opts {
		opt.apply(param)
	}
	return *param
}

func WithAllowedReserved() Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.In != oas.QueryLocation {
			panic("in must be query")
		}
		parameter.AllowReserved = true
	})
}
