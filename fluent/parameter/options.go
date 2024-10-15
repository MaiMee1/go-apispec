package parameter

import (
	"strconv"

	"github.com/MaiMee1/go-apispec/oas/v3"
)

// An Option configures a Logger.
type Option interface {
	apply(*oas.Parameter)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*oas.Parameter)

func (f optionFunc) apply(o *oas.Parameter) {
	f(o)
}

func WithExample(value interface{}) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Example = value
	})
}

func WithExamples(summary string, description oas.RichText, value interface{}) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.Examples == nil {
			parameter.Examples = make(map[string]oas.ValueOrReferenceOf[oas.Example])
		}
		k := next(parameter.Examples)
		parameter.Examples[k] = oas.ValueOrReferenceOf[oas.Example]{
			Value: oas.Example{
				Summary:     summary,
				Description: description,
				Value:       value,
			},
		}
	})
}

func next(examples map[string]oas.ValueOrReferenceOf[oas.Example]) string {
	const limit = 100
	var k string
	for i := 0; i < limit; i++ {
		k = strconv.Itoa(i)
		_, ok := examples[k]
		if !ok {
			break
		}
	}
	if k == strconv.Itoa(limit) {
		panic("too many examples")
	}
	return k
}

func WithDeprecated() Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Deprecated = true
	})
}

func WithSerialization(style oas.Style, explode bool) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Style = style
		parameter.Explode = &explode
	})
}

func WithSchema(schema oas.Schema) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Schema = &oas.ValueOrReferenceOf[oas.Schema]{
			Value: schema,
		}
	})
}

func WithSchemaReference(ref oas.Reference) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Schema = &oas.ValueOrReferenceOf[oas.Schema]{
			Ref: ref,
		}
	})
}
