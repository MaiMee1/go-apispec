package jsonschema

import (
	"iter"
	"maps"

	"github.com/MaiMee1/go-apispec/oas/internal/flag"
)

type Kind uint8

const (
	Identifier       Kind = 1 << iota // Control schema identification through setting a URI for the schema and/or changing how the base URI is determined
	ReservedLocation                  // Do not directly affect results, but reserve a place for a specific purpose to ensure interoperability
	Applicator                        // Apply one or more subschemas to a particular location in the instance, and combine or modify their results
	Annotation                        // Attach information to an instance for application use
	Assertion                         // Produce a boolean result when applied to an instance
)

var kindToString = map[Kind]string{
	Identifier:       "Identifier",
	ReservedLocation: "Reserved Location",
	Applicator:       "Applicator",
	Annotation:       "Annotation",
	Assertion:        "Assertion",
}

type Keyword interface {
	Kind() Kind
	AppliesTo(t Type) bool
	Validate(v interface{}) error
}

func (k Kind) Has(ands ...Kind) bool {
	return flag.Has(k, maps.Keys(kindToString), ands...)
}

func (k Kind) Range() iter.Seq[Kind] {
	return flag.Range(k, maps.Keys(kindToString))
}
