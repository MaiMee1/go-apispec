package draft2020

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/ser"
	"github.com/go-playground/validator/v10"
)

const validateTagSep = ','

var validate = validator.New(validator.WithRequiredStructEnabled())

type ArrayMixin struct {
	MaxItems    int                     `json:"maxItems,omitempty" validate:"omitempty,gte=0"`
	MinItems    int                     `json:"minItems,omitempty" validate:"omitempty,gte=0"`
	UniqueItems bool                    `json:"uniqueItems,omitempty"`
	PrefixItems []Schema                `json:"prefixItems,omitempty"`
	Items       *ser.Or[*bool, *Schema] `json:"items,omitempty"`
	Contains    *Schema                 `json:"contains,omitempty"`
	MaxContains int                     `json:"maxContains,omitempty" validate:"omitempty,gte=0"`
	MinContains *int                    `json:"minContains,omitempty" validate:"omitempty,gte=0"`
}

func (m *ArrayMixin) Kind() jsonschema.Kind {
	return jsonschema.Assertion
}

func (m *ArrayMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.ArrayType)
}

func (m *ArrayMixin) Validate(v interface{}) error {
	arr := v.([]interface{})
	if err := validate.Var(arr, m.validateTag()); err != nil {
		return err
	}
	if err := m.checkItems(arr); err != nil {
		return err
	}
	if err := m.checkContains(arr); err != nil {
		return err
	}
	return nil
}

func (m *ArrayMixin) validateTag() string {
	b := strings.Builder{}
	if m.MaxItems != 0 {
		b.WriteString("max=")
		b.WriteString(strconv.Itoa(m.MaxItems))
		b.WriteRune(validateTagSep)
	}
	if m.MinItems != 0 {
		b.WriteString("min=")
		b.WriteString(strconv.Itoa(m.MinItems))
		b.WriteRune(validateTagSep)
	}
	if m.UniqueItems {
		b.WriteString("unique")
		b.WriteRune(validateTagSep)
	}
	return b.String()
}

func (m *ArrayMixin) checkItems(arr []interface{}) error {
	for i, e := range arr {
		var valid func(v interface{}) error
		if len(m.PrefixItems) > i {
			valid = m.PrefixItems[i].Validate
		} else {
			if m.Items != nil {
				if m.Items.X != nil {
					if *m.Items.X {
						valid = func(v interface{}) error { return nil }
					} else {
						valid = func(v interface{}) error { return errors.New("items: false") }
					}
				} else {
					valid = m.Items.Y.Validate
				}
			} else {
				valid = func(v interface{}) error { return nil }
			}
		}
		err := valid(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ArrayMixin) checkContains(arr []interface{}) error {
	if m.Contains == nil {
		return nil
	}
	found := 0
	for _, e := range arr {
		if m.Contains.Validate(e) == nil {
			found++
		}
	}
	wantMin := 1 // default MinContains is 1: https://json-schema.org/draft/2020-12/json-schema-validation#section-6.4.5
	if m.MinContains != nil {
		wantMin = *m.MinContains
	}
	if wantMin > found {
		return fmt.Errorf("contains: found %d matching items, want at least %d", found, wantMin)
	}
	if m.MaxContains != 0 && found > m.MaxContains {
		return fmt.Errorf("contains: found %d matching items, want at most %d", found, m.MaxContains)
	}
	return nil
}

type NumericMixin struct {
	MultipleOf       int  `json:"multipleOf,omitempty" validate:"omitempty,gt=0"`
	Maximum          int  `json:"maximum,omitempty"`
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty"`
	Minimum          int  `json:"minimum,omitempty"`
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty"`
}

func (m *NumericMixin) Kind() jsonschema.Kind {
	return jsonschema.Assertion
}

func (m *NumericMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.IntegerType | jsonschema.NumberType)
}

func (m *NumericMixin) Validate(v interface{}) error {
	return validate.Var(v, m.validateTag())
}

func (m *NumericMixin) validateTag() string {
	b := strings.Builder{}
	if m.MultipleOf != 0 {
		b.WriteString("multipleOf=")
		b.WriteString(strconv.Itoa(m.MultipleOf))
		b.WriteRune(validateTagSep)
	}
	if m.Maximum != 0 {
		b.WriteString("lt")
		if !m.ExclusiveMaximum {
			b.WriteRune('e')
		}
		b.WriteRune('=')
		b.WriteString(strconv.Itoa(m.Maximum))
		b.WriteRune(validateTagSep)
	}
	if m.Minimum != 0 {
		b.WriteString("gt")
		if !m.ExclusiveMinimum {
			b.WriteRune('e')
		}
		b.WriteRune('=')
		b.WriteString(strconv.Itoa(m.Minimum))
		b.WriteRune(validateTagSep)
	}
	return b.String()
}

type ObjectMixin struct {
	MaxProperties        int                    `json:"maxProperties,omitempty" validate:"omitempty,gte=0"`
	MinProperties        int                    `json:"minProperties,omitempty" validate:"omitempty,gte=0"`
	PropertyNames        *Schema                `json:"propertyNames,omitempty" validate:"omitempty,dive"`
	Required             []string               `json:"required,omitempty" validate:"omitempty,min=1,unique"`
	Properties           map[string]Schema      `json:"properties,omitempty" validate:"dive"`
	PatternProperties    map[string]Schema      `json:"patternProperties,omitempty" validate:"dive,keys,regex,endkeys"`
	AdditionalProperties *ser.Or[bool, *Schema] `json:"additionalProperties,omitempty"`
}

func (m *ObjectMixin) Kind() jsonschema.Kind {
	return jsonschema.Assertion
}

func (m *ObjectMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.ObjectType)
}

func (m *ObjectMixin) Validate(v interface{}) error {
	_ = v.(*Schema)
	return nil
}

type StringMixin struct {
	MaxLength int    `json:"maxLength,omitempty" validate:"omitempty,gte=0"`
	MinLength int    `json:"minLength,omitempty" validate:"omitempty,gte=0"`
	Pattern   string `json:"pattern,omitempty"  validate:"regex"`
}

func (m *StringMixin) Kind() jsonschema.Kind {
	return jsonschema.Assertion
}

func (m *StringMixin) AppliesTo(t jsonschema.Type) bool {
	return t.Has(jsonschema.StringType)
}

func (m *StringMixin) Validate(v interface{}) error {
	str := v.(string)
	return validate.Var(str, m.validateTag())
}

func (m *StringMixin) validateTag() string {
	b := strings.Builder{}
	if m.MaxLength != 0 {
		b.WriteString("max=")
		b.WriteString(strconv.Itoa(m.MaxLength))
		b.WriteRune(validateTagSep)
	}
	if m.MinLength != 0 {
		b.WriteString("min=")
		b.WriteString(strconv.Itoa(m.MinLength))
		b.WriteRune(validateTagSep)
	}
	if m.Pattern != "" {
		b.WriteString("regex_ecma=")
		b.WriteString(m.Pattern)
		b.WriteRune(validateTagSep)
	}
	return b.String()
}

type ValidationMixin struct {
	Type   jsonschema.Type   `json:"type,omitempty"`
	Format jsonschema.Format `json:"format,omitempty"`
	Enum   []interface{}     `json:"enum,omitempty"`
	Const  interface{}       `json:"const,omitempty"`
}

func (m *ValidationMixin) Kind() jsonschema.Kind {
	return jsonschema.Annotation | jsonschema.Assertion
}

func (m *ValidationMixin) AppliesTo(t jsonschema.Type) bool {
	return true
}

func (m *ValidationMixin) Validate(v interface{}) error {
	if m.Const != nil {
		if v != m.Const {
			return errors.New("const: invalid value")
		}
	}
	if len(m.Enum) != 0 {
		for _, enum := range m.Enum {
			if enum == v {
				return nil
			}
		}
		return errors.New("enum: invalid value")
	}
	if !m.Type.Has(jsonschema.TypeOf(v)) {
		return fmt.Errorf("type: want %s, got %s", m.Type, jsonschema.TypeOf(v))
	}
	// TODO: Format
	return nil
}
