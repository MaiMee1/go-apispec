package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/v3"
)

func WithMask(keys ...string) Option {
	return optionFunc(func(schema *oas.Schema) {
		for _, key := range keys {
			delete(schema.Properties, key)
		}
	})
}

func makeRef(name string) string {
	return fmt.Sprintf("#/components/schemas/%v", name)
}

func makeName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		if t.Kind() == reflect.Interface {
			name = "any"
		} else {
			name = "ptr" + makeName(t.Elem())
		}
	}
	pkgPath := t.PkgPath()
	if pkgPath != "." {
		pkgPath += "."
	}
	fullName := pkgPath + name
	fullName = strings.ReplaceAll(fullName, "/", ".") // updated
	return strings.ReplaceAll(fullName, "-", "_")
}

func parseFormat(t reflect.Type) oas.Format {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return parseFormat(t.Elem())
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return oas.NoFormat
	case reflect.Int32:
		return oas.Int32Format
	case reflect.Int64:
		return oas.Int64Format
	case reflect.Float32:
		return oas.FloatFormat
	case reflect.Float64:
		return oas.DoubleFormat
	default:
		return oas.NoFormat
	}
}

func parseType(t reflect.Type) oas.Type {
	kind := t.Kind()
	if kind == reflect.Ptr {
		return parseType(t.Elem())
	}

	// can cache here
	//if pt := Get(t.String()); pt != Unknown {
	//	return pt
	//}

	switch kind {
	case reflect.Bool:
		return oas.BooleanType
	case reflect.Float32, reflect.Float64:
		return oas.NumberType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return oas.IntegerType
	case reflect.String:
		return oas.StringType
	case reflect.Array, reflect.Slice:
		return oas.ArrayType
	case reflect.Struct, reflect.Map, reflect.Interface:
		return oas.ObjectType
	default:
		return 0
	}
}

func buildProperty(t reflect.Type) (properties map[string]oas.ValueOrReferenceOf[oas.Schema], required []string) {
	//fmt.Println("buildProperty", t)
	properties = make(map[string]oas.ValueOrReferenceOf[oas.Schema])
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// skip unexported fields
		if !field.IsExported() {
			continue
		}
		if field.Anonymous {
			// 暂不处理匿名结构的required
			ps, _ := buildProperty(field.Type)
			for name, p := range ps {
				properties[name] = p
			}
			continue
		}

		// determine the json name of the field
		name := strings.TrimSpace(field.Tag.Get("json"))
		if name == "" || strings.HasPrefix(name, ",") {
			name = field.Name

		} else {
			// strip out things like , omitempty
			parts := strings.Split(name, ",")
			name = parts[0]
		}

		parts := strings.Split(name, ",") // foo,omitempty => foo
		name = parts[0]
		if name == "-" {
			// honor json ignore tag
			continue
		}

		var p oas.Schema
		jsonTag := field.Tag.Get("json")
		if strings.Contains(jsonTag, ",string") {
			p.Type = oas.StringType
			p.Extensions["GoType"] = field.Type
		} else {
			_, p = defineObject(field.Type, false)
		}

		// determine the extra info of the field
		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			parts := strings.Split(validateTag, ",")
			for _, part := range parts {
				if part == "required" {
					required = append(required, name)
					break
				}
			}
		}
		// TODO: description from comments
		properties[name] = oas.ValueOrReferenceOf[oas.Schema]{
			Value: p,
		}
	}
	return properties, required
}

var cache = make(map[string]*oas.Schema)

func defineObject(t reflect.Type, nullable bool) (string, oas.Schema) {
	//fmt.Println("defineObject", t, nullable)
	kind := t.Kind()
	if kind == reflect.Ptr {
		// unwrap pointer
		return defineObject(t.Elem(), true)
	}

	name := makeName(t)
	// TODO: cache better
	if c, ok := cache[name]; ok {
		//fmt.Println("cache hit", name)
		return name, *c
	}

	p := oas.Schema{
		Type:     parseType(t),
		Format:   parseFormat(t),
		Nullable: nullable,
		Extensions: oas.SpecificationExtension{
			"GoType": t,
		},
	}
	cache[name] = &p

	switch kind {
	case reflect.Struct:
		properties, required := buildProperty(t)
		p.Required = required
		p.Properties = properties
		name = makeName(t)
	case reflect.Interface:
		p.AdditionalProperties = &oas.Or[bool, *oas.ValueOrReferenceOf[oas.Schema]]{
			X: true,
		}
	case reflect.Map:
		_, item := defineObject(t.Elem(), nullable)
		p.AdditionalProperties = &oas.Or[bool, *oas.ValueOrReferenceOf[oas.Schema]]{
			Y: &oas.ValueOrReferenceOf[oas.Schema]{
				Value: item,
			},
		}
		name = makeName(t)
	case reflect.Slice, reflect.Array:
		// unwrap array or slice
		_, item := defineObject(t.Elem(), true)
		p.Items = &oas.ValueOrReferenceOf[oas.Schema]{
			Value: item,
		}
		name = makeName(t)
	default:
	}

	return name, p
}

func define(t reflect.Type) map[string]oas.Schema {
	objMap := make(map[string]oas.Schema)
	name, obj := defineObject(t, false)
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
					if prop.Type == oas.ObjectType {
						name := makeName(prop.Extensions["GoType"].(reflect.Type))
						if _, exists := objMap[name]; !exists {
							name, child := defineObject(prop.Extensions["GoType"].(reflect.Type), false)
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
					if prop.Type == oas.ObjectType {
						name := makeName(prop.Extensions["GoType"].(reflect.Type))
						if _, exists := objMap[name]; !exists {
							name, child := defineObject(prop.Extensions["GoType"].(reflect.Type), false)
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

//func MakeSchema(prototype interface{}) oas.Schema {
//	name, obj := defineObject(prototype, false)
//	//schema := &oas.Schema{
//	//	Prototype: prototype,
//	//	Reference:       makeRef(name),
//	//}
//	_ = makeRef(name)
//	return obj
//}
