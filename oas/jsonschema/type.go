package jsonschema

import (
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"reflect"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/internal/flag"
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

var kindToType = []Type{
	reflect.Invalid:       0,
	reflect.Bool:          BooleanType,
	reflect.Int:           IntegerType,
	reflect.Int8:          IntegerType,
	reflect.Int16:         IntegerType,
	reflect.Int32:         IntegerType,
	reflect.Int64:         IntegerType,
	reflect.Uint:          IntegerType,
	reflect.Uint8:         IntegerType,
	reflect.Uint16:        IntegerType,
	reflect.Uint32:        IntegerType,
	reflect.Uint64:        IntegerType,
	reflect.Uintptr:       IntegerType,
	reflect.Float32:       NumberType,
	reflect.Float64:       NumberType,
	reflect.Complex64:     0,
	reflect.Complex128:    0,
	reflect.Array:         ArrayType,
	reflect.Chan:          0,
	reflect.Func:          0,
	reflect.Interface:     AnyType,
	reflect.Map:           ObjectType,
	reflect.Pointer:       0,
	reflect.Slice:         ArrayType,
	reflect.String:        StringType,
	reflect.Struct:        ObjectType,
	reflect.UnsafePointer: 0,
}

func TypeOf(v any) Type {
	return kindToType[reflect.TypeOf(v).Kind()]
}

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
	return flag.Has(t, maps.Keys(typeToString), ands...)
}

//goland:noinspection GoMixedReceiverTypes
func (t Type) Range() iter.Seq[Type] {
	return flag.Range(t, maps.Keys(typeToString))
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
