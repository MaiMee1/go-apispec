package schema

import (
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/MaiMee1/go-apispec/oas/jsonschema"
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
	return oas.NoFormat
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
		switch t.Kind() {
		case reflect.Interface:
			name = "any"
		case reflect.Pointer:
			name = "ptr"
		default:
			panic(t)
		}
	}
	return t.PkgPath() + name
}

type StringFilter = func(string) string

var DefaultEncoder = Encoder{
	nameFilter: func(qualifiedName string) string {
		qualifiedName = strings.ReplaceAll(qualifiedName, "/", ".")
		qualifiedName = strings.ReplaceAll(qualifiedName, "-", "_")
		return qualifiedName
	},
}

type Encoder struct {
	nameFilter StringFilter
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

func diveStruct(t reflect.Type) (properties map[string]oas.ValueOrReferenceOf[oas.Schema], required map[string]struct{}) {
	required = make(map[string]struct{})
	properties = make(map[string]oas.ValueOrReferenceOf[oas.Schema])
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
			props, req := diveStruct(t)
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
			Extensions: oas.SpecificationExtension{
				"GoType": sf.Type,
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
				schema = objectSchema(sf.Type)
			}
		} else {
			schema = objectSchema(sf.Type)
		}

		if isRequired(sf) {
			required[name] = struct{}{}
			// JSON Schema's "required" does not mean the value cannot be null (just that the key must be present)
			// but validate tag's "required" expects not nil value
			schema.Type = schema.Type & ^jsonschema.NullType
		}

		properties[name] = oas.ValueOrReferenceOf[oas.Schema]{
			Value: schema,
		}
	}

	return
}

var cache = sync.Map{} // map[string]*oas.Schema

func objectSchema(t reflect.Type) (schema oas.Schema) {
	var nullable bool
	// some types allow nil values
	switch t.Kind() {
	case reflect.Pointer,
		reflect.Interface, reflect.Map,
		reflect.Array, reflect.Slice:
		nullable = true
	default:
	}

	// unwrap pointer once
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		name := DefaultEncoder.makeName(t)
		if c, ok := cache.Load(name); ok {
			return *c.(*oas.Schema)
		}

		schema = oas.Schema{
			Type:   DefaultEncoder.dataType(t),
			Format: DefaultEncoder.format(t),
			Extensions: oas.SpecificationExtension{
				"GoType": t,
				"Name":   name,
			},
		}
		if nullable {
			schema.Type = schema.Type | jsonschema.NullType
		}
		cache.Store(name, &schema)

		properties, required := diveStruct(t)
		schema.Required = slices.Collect(maps.Keys(required))
		schema.Properties = properties
	case reflect.Interface:
		name := DefaultEncoder.makeName(t)
		if c, ok := cache.LoadOrStore(name, &oas.Schema{
			Type: jsonschema.AnyType,
			Extensions: oas.SpecificationExtension{
				"GoType": t,
				"Name":   name,
			},
		}); ok {
			return *c.(*oas.Schema)
		}
	case reflect.Map:
		schema = oas.Schema{
			Type:   DefaultEncoder.dataType(t) | jsonschema.NullType,
			Format: DefaultEncoder.format(t),
			Extensions: oas.SpecificationExtension{
				"GoType": t,
			},
		}

		item := objectSchema(t.Elem())
		schema.AdditionalProperties = &oas.Or[bool, *oas.ValueOrReferenceOf[oas.Schema]]{
			Y: &oas.ValueOrReferenceOf[oas.Schema]{
				Value: item,
			},
		}
	case reflect.Slice, reflect.Array:
		schema = oas.Schema{
			Type:   DefaultEncoder.dataType(t) | jsonschema.NullType,
			Format: DefaultEncoder.format(t),
			Extensions: oas.SpecificationExtension{
				"GoType": t,
			},
		}

		item := objectSchema(t.Elem())
		schema.Items = &oas.ValueOrReferenceOf[oas.Schema]{
			Value: item,
		}
	default:
		schema = oas.Schema{
			Type:   DefaultEncoder.dataType(t),
			Format: DefaultEncoder.format(t),
			Extensions: oas.SpecificationExtension{
				"GoType": t,
			},
		}
		if nullable {
			schema.Type = schema.Type | jsonschema.NullType
		}
	}
	return
}

func define(t reflect.Type) map[string]oas.Schema {
	objMap := make(map[string]oas.Schema)
	obj := objectSchema(t)
	if !obj.Type.Has(jsonschema.ObjectType) {
		return objMap
	}
	name := obj.Extensions["Name"].(string)
	objMap[name] = obj

	dirty := true
	for dirty {
		dirty = false
		for _, d := range objMap {
			if d.Items != nil { // update
				prop := d.Items.Value
				for prop.Items != nil {
					prop = prop.Items.Value
				}
				if prop.Properties != nil {
					if prop.Type.Has(jsonschema.ObjectType) {
						name := prop.Extensions["Name"].(string)
						if _, exists := objMap[name]; !exists {
							child := objectSchema(prop.Extensions["GoType"].(reflect.Type))
							objMap[name] = child
							dirty = true
						}
					}
					if prop.AdditionalProperties != nil && prop.AdditionalProperties.Y != nil {
						prop = prop.AdditionalProperties.Y.Value
					}
				}
			}
			for _, p := range d.Properties {
				prop := &p.Value
				for prop.Items != nil { // update
					prop = &prop.Items.Value
				}
				for prop != nil {
					if prop.Type.Has(jsonschema.ObjectType) {
						name := prop.Extensions["Name"].(string)
						if _, exists := objMap[name]; !exists {
							child := objectSchema(prop.Extensions["GoType"].(reflect.Type))
							objMap[name] = child
							dirty = true
						}
					}
					if prop.AdditionalProperties != nil && prop.AdditionalProperties.Y != nil {
						prop = &prop.AdditionalProperties.Y.Value
					} else {
						prop = nil
					}
				}
			}
		}
	}
	return objMap
}
