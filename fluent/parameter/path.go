package parameter

import "github.com/MaiMee1/go-apispec/oas/v3"

func Path(name string, description oas.RichText, required bool, opts ...Option) oas.Parameter {
	param := &oas.Parameter{
		In:          oas.PathLocation,
		Style:       0, // oas.SimpleStyle
		Name:        name,
		Description: description,
		Required:    required,
	}
	for _, opt := range opts {
		opt.apply(param)
	}
	return *param
}
