package jsonschema

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"strings"
)

type Type uint8

const (
	NullType Type = 1 << iota
	IntegerType
	NumberType
	StringType
	BooleanType
	ObjectType
	ArrayType
)

const AnyType = NullType | IntegerType | NumberType | StringType | BooleanType | ObjectType | ArrayType

var typeToString = map[Type]string{
	NullType:    "null",
	IntegerType: "integer",
	NumberType:  "number",
	StringType:  "string",
	BooleanType: "boolean",
	ObjectType:  "object",
	ArrayType:   "array",
}
var stringToType map[string]Type

//goland:noinspection GoMixedReceiverTypes
func (t Type) Has(ands ...Type) bool {
	for _, and := range ands {
		var ok = false
		for or := range and.Range() {
			if t&or == or {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) Range() iter.Seq[Type] {
	return func(yield func(Type) bool) {
		for typ := range maps.Keys(typeToString) {
			if t&typ == typ {
				if !yield(typ) {
					return
				}
			}
		}
	}
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) String() string {
	if t == 0 {
		return "<0>"
	}
	if s, ok := typeToString[t]; ok {
		return s
	}
	b := strings.Builder{}
	sep := '('
	for typ := range t.Range() {
		b.WriteRune(sep)
		b.WriteString(typeToString[typ])
		sep = '|'
	}
	b.WriteRune(')')
	return b.String()
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) MarshalJSON() ([]byte, error) {
	// if no combination, serialize as string
	if s, ok := typeToString[t]; ok {
		return json.Marshal(s)
	}

	// if is combination, serialize as array of strings
	var arr []string
	for typ := range t.Range() {
		arr = append(arr, typeToString[typ])
	}
	return json.Marshal(arr)
}

//goland:noinspection GoMixedReceiverTypes
func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		typ, ok := stringToType[s]
		if !ok {
			return fmt.Errorf("invalid type %q", s)
		}
		*t = typ
		return nil
	}

	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		var tt Type
		for _, s := range arr {
			typ, ok := stringToType[s]
			if !ok {
				// TODO: return multiple errors
				return fmt.Errorf("invalid type %q", s)
			}
			tt = tt | typ
		}
		*t = tt
		return nil
	}

	return fmt.Errorf("invalid json value %q, expect string or array of strings", string(b))
}

func all2[Map ~map[K]V, K comparable, V any](m Map) iter.Seq2[V, K] {
	return func(yield func(V, K) bool) {
		for k, v := range m {
			if !yield(v, k) {
				return
			}
		}
	}
}

func init() {
	stringToType = maps.Collect(all2(typeToString))
}
