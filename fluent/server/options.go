package server

import "github.com/MaiMee1/go-apispec/oas/v3"

// An Option configures a Logger.
type Option interface {
	apply(*oas.Server)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*oas.Server)

func (f optionFunc) apply(o *oas.Server) {
	f(o)
}

func WithTitle(title string) Option {
	return optionFunc(func(server *oas.Server) {
		server.Description = oas.RichText(title)
	})
}

func WithDescription(description oas.RichText) Option {
	return optionFunc(func(server *oas.Server) {
		server.Description = description
	})
}

func WithVariable(name string, defaultValue string, description oas.RichText, examples ...string) Option {
	if name == "" {
		panic("variable name is empty")
	}
	return optionFunc(func(server *oas.Server) {
		server.Variables[name] = oas.ServerVariable{
			Enum:        nil,
			Default:     defaultValue,
			Description: description,
			Extensions: oas.SpecificationExtension{
				"examples": examples,
			},
		}
	})
}

func WithEnumVariable(name string, defaultValue string, description oas.RichText, enum []string, examples ...string) Option {
	if name == "" {
		panic("variable name is empty")
	}
	return optionFunc(func(server *oas.Server) {
		server.Variables[name] = oas.ServerVariable{
			Enum:        enum,
			Default:     defaultValue,
			Description: description,
			Extensions: oas.SpecificationExtension{
				"examples": examples,
			},
		}
	})
}
