package encoder

import (
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/oas31"
	"github.com/MaiMee1/go-apispec/oas/ser"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

// fallbackKindToFormat uses Go type names as fallback for JSON serializable values
var fallbackKindToFormat = []oas.Format{
	reflect.Invalid:    "",
	reflect.Bool:       "bool",
	reflect.Int:        "int",
	reflect.Int8:       "int8",
	reflect.Int16:      "int16",
	reflect.Int32:      "int32",
	reflect.Int64:      "int64",
	reflect.Uint:       "uint",
	reflect.Uint8:      "uint8",
	reflect.Uint16:     "uint16",
	reflect.Uint32:     "uint32",
	reflect.Uint64:     "uint64",
	reflect.Uintptr:    "uintptr",
	reflect.Float32:    "float32",
	reflect.Float64:    "float64",
	reflect.Complex64:  "complex64",
	reflect.Complex128: "complex128",
}

// kindToFormat defines mapping from Go [reflect.Kind] to canonical [oas.Format]
var kindToFormat = map[reflect.Kind]oas.Format{
	reflect.Int32:   oas.Int32Format,
	reflect.Int64:   oas.Int64Format,
	reflect.Float32: oas.FloatFormat,
	reflect.Float64: oas.DoubleFormat,
}

func format2(t reflect.Type, fallback bool) oas.Format {
	if format, ok := kindToFormat[t.Kind()]; ok {
		return format
	}
	if fallback {
		return fallbackKindToFormat[t.Kind()]
	}
	return ""
}

func dataType2(t reflect.Type) oas.Type {
	switch t.Kind() {
	case reflect.Invalid, reflect.Pointer, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func:
		return 0
	case reflect.Bool:
		return jsonschema.BooleanType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return jsonschema.IntegerType
	case reflect.Float32, reflect.Float64:
		return jsonschema.NumberType
	case reflect.Array, reflect.Slice:
		return jsonschema.ArrayType
	case reflect.Interface, reflect.Map, reflect.Struct:
		return jsonschema.ObjectType
	case reflect.String:
		return jsonschema.StringType
	}
	panic("unreachable")
}

func name2(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		panic(t)
	}
	return t.PkgPath() + "." + name
}

var defaultNameFilter = strings.NewReplacer("/", ".", "-", "_").Replace

type Encoder struct {
	cache         sync.Map // map[string]*oas.Schema
	nameFilter    StringFilter
	nullableMap   bool
	nullableSlice bool
}

func New(opts ...Option) *Encoder {
	enc := new(Encoder)
	enc.cache = sync.Map{}
	enc.nameFilter = defaultNameFilter
	for _, opt := range opts {
		opt.apply(enc)
	}
	return enc
}

func (enc *Encoder) Encode(t reflect.Type) oas.Schema {
	return enc.objectSchema(t)
}

func (enc *Encoder) Cache() map[string]*oas.Schema {
	m := make(map[string]*oas.Schema)
	enc.cache.Range(func(k, v interface{}) bool {
		m[k.(string)] = v.(*oas.Schema)
		return true
	})
	return m
}

func (enc *Encoder) format(t reflect.Type) oas.Format {
	return format2(t, false)
}

func (enc *Encoder) dataType(t reflect.Type) oas.Type {
	return dataType2(t)
}

func (enc *Encoder) makeName(t reflect.Type) string {
	return enc.nameFilter(name2(t))
}

func isRequired(sf reflect.StructField) bool {
	tag, ok := sf.Tag.Lookup("validate")
	if !ok {
		return false
	}
	if tag == "-" {
		return false
	}
	tags := strings.Split(tag, ",")
	return slices.Contains(tags, "required")
}

func (enc *Encoder) diveStruct(t reflect.Type) (properties map[string]*oas.Schema, required map[string]struct{}) {
	required = make(map[string]struct{})
	properties = make(map[string]*oas.Schema)
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// handle embedded fields
		if sf.Anonymous {
			t := sf.Type
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}
			if !sf.IsExported() && t.Kind() != reflect.Struct {
				// skip embedded fields of unexported non-struct types
				continue
			}
			// dive into embedded fields of un/exported struct types
			props, req := enc.diveStruct(t)
			for name, prop := range props {
				if _, ok := properties[name]; ok {
					panic("promoted fields in conflict")
				}
				properties[name] = prop
			}
			for name := range req {
				required[name] = struct{}{}
			}
			continue
		} else if !sf.IsExported() {
			// skip unexported non-embedded fields
			continue
		}

		tag := sf.Tag.Get("json")
		if tag == "-" {
			// skip non-serialized fields
			continue
		}

		name, jsonOpts := parseTag(tag)
		if !isValidTag(name) {
			// fallback to field name
			name = sf.Name
		}

		// check for duplicate names
		if _, ok := properties[name]; ok {
			panic("fields in conflict")
		}

		schema := oas.Schema{
			OASMixin: oas31.OASMixin{
				Extensions: oas.SpecificationExtension{
					"GoType": sf.Type,
				},
			},
		}
		if jsonOpts.Contains("string") {
			// encoding/json only add quotes to strings, floats, integers, and booleans
			switch sf.Type.Kind() {
			case reflect.String:
				schema.Type = jsonschema.StringType
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
				reflect.Float32, reflect.Float64:
				schema.Type = jsonschema.StringType
				schema.Format = format2(sf.Type, true)
			default:
				schema = enc.objectSchema(sf.Type)
			}
		} else {
			schema = enc.objectSchema(sf.Type)
		}

		if isRequired(sf) {
			required[name] = struct{}{}
			// JSON Schema's "required" does not mean the value cannot be null (just that the key must be present)
			// but validate tag's "required" expects not nil value
			schema.Type = schema.Type & ^jsonschema.NullType
		}

		properties[name] = &schema
	}
	return
}

func (enc *Encoder) objectSchema(t reflect.Type) (schema oas.Schema) {
	var nullable bool
	// some types allow nil values
	switch t.Kind() {
	case reflect.Pointer:
		nullable = true
	case reflect.Map:
		nullable = enc.nullableMap
	case reflect.Array, reflect.Slice:
		nullable = enc.nullableSlice
	default:
	}

	// unwrap pointer once
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema.Type = enc.dataType(t)
	schema.Format = enc.format(t)
	schema.Extensions = oas.SpecificationExtension{
		"GoType": t,
	}
	switch t.Kind() {
	case reflect.Struct:
		name := enc.makeName(t)
		if c, ok := enc.cache.Load(name); ok {
			return *c.(*oas.Schema)
		}

		schema.Extensions["Name"] = name
		enc.cache.Store(name, &schema)

		properties, required := enc.diveStruct(t)
		schema.Required = slices.Collect(maps.Keys(required))
		schema.Properties = properties
	case reflect.Interface:
		schema.Type = 0
	case reflect.Map:
		item := enc.objectSchema(t.Elem())
		schema.AdditionalProperties = &ser.Or[bool, *oas.Schema]{
			Y: &item,
		}
	case reflect.Slice, reflect.Array:
		item := enc.objectSchema(t.Elem())
		schema.Items = &ser.Or[bool, *oas.Schema]{
			Y: &item,
		}
	default:
	}
	if nullable {
		schema.Type = schema.Type | jsonschema.NullType
	}
	return
}
