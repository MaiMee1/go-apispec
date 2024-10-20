package abc

import (
	"iter"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type KeywordCollection interface {
	Keywords() iter.Seq[jsonschema.Keyword]
}

func Kind(m KeywordCollection) jsonschema.Kind {
	var kind jsonschema.Kind
	for mixin := range m.Keywords() {
		kind = kind | mixin.Kind()
	}
	return kind
}

func AppliesTo(m KeywordCollection, t jsonschema.Type) bool {
	for mixin := range m.Keywords() {
		if mixin.AppliesTo(t) {
			return true
		}
	}
	return true
}

func Validate(m KeywordCollection, v interface{}) error {
	var err error // TODO: create validation errors
	t := jsonschema.TypeOf(v)
	for mixin := range m.Keywords() {
		if mixin.AppliesTo(t) && mixin.Kind().Has(jsonschema.Assertion) {
			err = mixin.Validate(v)
			if err != nil {
				return err
			}
		}
	}
	return err
}
