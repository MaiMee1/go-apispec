package schema

import "github.com/MaiMee1/go-apispec/oas/v3"

type Option interface {
	apply(*oas.Schema)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*oas.Schema)

func (f optionFunc) apply(o *oas.Schema) {
	f(o)
}

func WithExample(value interface{}) Option {
	return optionFunc(func(schema *oas.Schema) {
		schema.Example = value
	})
}

func WithTitle(title string) Option {
	return optionFunc(func(schema *oas.Schema) {
		schema.Title = title
	})
}

func WithDescription(description oas.RichText) Option {
	return optionFunc(func(schema *oas.Schema) {
		schema.Description = string(description)
	})
}

func WithDefault(value interface{}) Option {
	return optionFunc(func(schema *oas.Schema) {
		schema.Default = value
	})
}

func WithEnum(value ...interface{}) Option {
	return optionFunc(func(schema *oas.Schema) {
		schema.Enum = value
	})
}
