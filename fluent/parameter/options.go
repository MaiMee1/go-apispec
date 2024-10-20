package parameter

import (
	"strconv"

	"github.com/MaiMee1/go-apispec/fluent/schema"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

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
			parameter.Examples = make(map[string]oas.Example)
		}
		k := next(parameter.Examples)
		parameter.Examples[k] = oas.Example{
			Summary:     summary,
			Description: description,
			Value:       value,
		}
	})
}

func next(examples map[string]oas.Example) string {
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

// WithMatrixStyle use Path-style parameters defined by RFC6570.
//
// See section 3.2.7. Path-Style Parameter Expansion: {;var}.
//
// Assume a parameter named color has one of the following values:
//
//	"blue"
//	[]string{"blue","black","brown"}
//	map[string]int{ "R": 100, "G": 200, "B": 150 }
//
//	explode=false           	explode=true
//	;color                  	;color
//	;color=blue             	;color=blue
//	;color=blue,black,brown 	;color=blue;color=black;color=brown
//	;color=R,100,G,200,B,150	;R=100;G=200;B=150
func WithMatrixStyle(explode bool) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.In != oas.PathLocation {
			panic("in must be path")
		}
		parameter.Style = oas.MatrixStyle
		parameter.Explode = &explode
	})
}

// WithLabelStyle use Label style parameters defined by RFC6570.
//
// See section 3.2.5. Label Expansion with Dot-Prefix: {.var}.
//
// Assume a parameter named color has one of the following values:
//
//	"blue"
//	[]string{"blue","black","brown"}
//	map[string]int{ "R": 100, "G": 200, "B": 150 }
//
//	explode=false     	explode=true
//	.                 	.
//	.blue             	.blue
//	.blue.black.brown 	.blue.black.brown
//	.R.100.G.200.B.150	.R=100.G=200.B=150
func WithLabelStyle(explode bool) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.In != oas.PathLocation {
			panic("in must be path")
		}
		parameter.Style = oas.LabelStyle
		parameter.Explode = &explode
	})
}

// WithFormStyle use Form style parameters defined by RFC6570.
//
// See section 3.2.8. Form-Style Query Expansion: {?var}
//
// Assume a parameter named color has one of the following values:
//
//	"blue"
//	[]string{"blue","black","brown"}
//	map[string]int{ "R": 100, "G": 200, "B": 150 }
//
//	explode=false         	explode=true
//	color=                 	color=
//	color=blue             	color=blue
//	color=blue,black,brown 	color=blue&color=black&color=brown
//	color=R,100,G,200,B,150	R=100&G=200&B=150
func WithFormStyle(explode bool) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.In != oas.QueryLocation && parameter.In != oas.CookieLocation {
			panic("in must be query or cookie")
		}
		parameter.Style = oas.FormStyle
		parameter.Explode = &explode
	})
}

// WithSimpleStyle use Simple style parameters defined by RFC6570.
//
// See section 3.2.2. Simple String Expansion: {var}
//
// Assume a parameter named color has one of the following values:
//
//	"blue"
//	[]string{"blue","black","brown"}
//	map[string]int{ "R": 100, "G": 200, "B": 150 }
//
//	explode=false    	explode=true
//	blue             	blue
//	blue,black,brown 	blue,black,brown
//	R,100,G,200,B,150	R=100,G=200,B=150
func WithSimpleStyle(explode bool) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		if parameter.In != oas.QueryLocation && parameter.In != oas.CookieLocation {
			panic("in must be query or cookie")
		}
		parameter.Style = oas.FormStyle
		parameter.Explode = &explode
	})
}

func WithSchemaFor[T any](opts ...schema.Option) Option {
	s := schema.For[T](opts...)
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Schema = s
	})
}

func WithSchema(schema oas.Schema) Option {
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Schema = schema
	})
}

func WithSchemaReference(ref string) Option {
	var s oas.Schema
	s.Ref = ref
	return optionFunc(func(parameter *oas.Parameter) {
		parameter.Schema = s
	})
}

func WithComplexSerialization(keyAndValues ...interface{}) Option {
	if len(keyAndValues)%2 != 0 {
		panic("keyAndValues must have an even number")
	}
	return optionFunc(func(parameter *oas.Parameter) {
		for i := 0; i < len(keyAndValues)/2; i++ {
			key := keyAndValues[i*2].(string)
			value := keyAndValues[i*2+1]
			switch v := value.(type) {
			case oas.Schema:
				parameter.Content[key] = oas.MediaType{
					Schema:   v,
					Example:  nil,
					Examples: nil,
				}
			}
		}
	})
}
