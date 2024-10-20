package draft2020

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
)

type ApplicatorMixin[S jsonschema.Keyword] struct {
	AllOf []S `json:"allOf,omitempty"`
	AnyOf []S `json:"anyOf,omitempty"`
	OneOf []S `json:"oneOf,omitempty"`
	If    S   `json:"if,omitempty"`
	Then  S   `json:"then,omitempty"`
	Else  S   `json:"else,omitempty"`
	Not   S   `json:"not,omitempty"`
}

func (m *ApplicatorMixin[S]) Kind() jsonschema.Kind {
	return jsonschema.Applicator
}

func (m *ApplicatorMixin[S]) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *ApplicatorMixin[S]) Validate(v interface{}) error {
	if len(m.AllOf) > 0 {
		for i, schema := range m.AllOf {
			if err := schema.Validate(v); err != nil {
				return fmt.Errorf("allOff[%d]: %w", i, err)
			}
		}
	}
	if len(m.AnyOf) > 0 {
		slices.ContainsFunc(m.AnyOf, func(schema S) bool {
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
	if !reflect.DeepEqual(m.If, nil) && (!reflect.DeepEqual(m.Then, nil) || !reflect.DeepEqual(m.Else, nil)) {
		eval := m.If.Validate(v) == nil
		if eval {
			if !reflect.DeepEqual(m.Then, nil) {
				if err := m.Then.Validate(v); err != nil {
					return fmt.Errorf("then: %w", err)
				}
			}
		} else {
			if !reflect.DeepEqual(m.Else, nil) {
				if err := m.Else.Validate(v); err != nil {
					return fmt.Errorf("else: %w", err)
				}
			}
		}
	}
	if !reflect.DeepEqual(m.Not, nil) {
		if m.Not.Validate(v) == nil {
			return errors.New("not: schema matched")
		}
	}
	return nil
}
