package draft2020

import (
	"errors"
	"fmt"
	"slices"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type ApplicatorMixin struct {
	AllOf []Schema `json:"allOf,omitempty"`
	AnyOf []Schema `json:"anyOf,omitempty"`
	OneOf []Schema `json:"oneOf,omitempty"`
	If    *Schema  `json:"if,omitempty"`
	Then  *Schema  `json:"then,omitempty"`
	Else  *Schema  `json:"else,omitempty"`
	Not   *Schema  `json:"not,omitempty"`
}

func (m *ApplicatorMixin) Kind() jsonschema.Kind {
	return jsonschema.Applicator
}

func (m *ApplicatorMixin) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *ApplicatorMixin) Validate(v interface{}) error {
	if len(m.AllOf) > 0 {
		for i, schema := range m.AllOf {
			if err := schema.Validate(v); err != nil {
				return fmt.Errorf("allOff[%d]: %w", i, err)
			}
		}
	}
	if len(m.AnyOf) > 0 {
		slices.ContainsFunc(m.AnyOf, func(schema Schema) bool {
			return schema.Validate(v) == nil
		})
	}
	if len(m.OneOf) > 0 {
		var indices []int
		for i, schema := range m.AllOf {
			if schema.Validate(v) == nil {
				indices = append(indices, i)
			}
		}
		if len(indices) != 1 {
			return fmt.Errorf("oneOf: found %d matching at %v, want %d", len(indices), indices, 1)
		}
	}
	if m.If != nil && (m.Then != nil || m.Else != nil) {
		eval := m.If.Validate(v) == nil
		if eval {
			if m.Then != nil {
				if err := m.Then.Validate(v); err != nil {
					return fmt.Errorf("then: %w", err)
				}
			}
		} else {
			if m.Else != nil {
				if err := m.Else.Validate(v); err != nil {
					return fmt.Errorf("else: %w", err)
				}
			}
		}
	}
	if m.Not != nil {
		if m.Not.Validate(v) == nil {
			return errors.New("not: schema matched")
		}
	}
	return nil
}
