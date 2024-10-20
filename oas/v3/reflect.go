package oas

import (
	"iter"
	"reflect"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/jsonschema/oas31"
)

var valueOrReferenceOfPrefix = strings.TrimSuffix(reflect.TypeOf(ValueOrReferenceOf[bool]{}).Name(), "[bool]")
var valueOrReferenceOfSchema = reflect.TypeOf(oas31.Schema{}).Name()

// setRoot recursively find ValueOrReferenceOf fields or elements and sets its Root to root.
func setRoot(v reflect.Value, root interface{}) {
	s := reflect.ValueOf(root)
	for v := range iterValueOrReference(v, true) {
		f := v.FieldByName("Root")
		f.Set(s)
	}
}

// iterValueOrReference recursively find ValueOrReferenceOf fields or elements and sets its Root to root.
//
// If makeCanSet is true, creates a pointer to the existing map value to make its value settable.
func iterValueOrReference(v reflect.Value, makeCanSet bool) iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		//fmt.Println(">> ", v.Type(), fmt.Sprintf("%q", fmt.Sprint(v.Interface())))
		switch v.Kind() {
		case reflect.Invalid:
			panic(v)
		case reflect.Ptr:
			if !v.IsNil() {
				for f := range iterValueOrReference(v.Elem(), makeCanSet) {
					if !yield(f) {
						return
					}
				}
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				for f := range iterValueOrReference(v.Field(i), makeCanSet) {
					if !yield(f) {
						return
					}
				}
			}
			if strings.HasPrefix(v.Type().Name(), valueOrReferenceOfPrefix) {
				//fmt.Println(">>>", v.Type())
				if !yield(v) {
					return
				}
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				for f := range iterValueOrReference(v.Index(i).Addr(), makeCanSet) {
					if !yield(f) {
						return
					}
				}
			}
		case reflect.Map:
			it := v.MapRange()
			for it.Next() {
				value := it.Value()
				if makeCanSet && !value.CanSet() {
					p := reflect.New(value.Type())
					p.Elem().Set(value)
					for f := range iterValueOrReference(p, makeCanSet) {
						if !yield(f) {
							return
						}
					}
					v.SetMapIndex(it.Key(), p.Elem())
				} else {
					for f := range iterValueOrReference(value, makeCanSet) {
						if !yield(f) {
							return
						}
					}
				}
			}
		default:
		}
	}
}

func (doc *OpenAPI) IterRef() iter.Seq[*Reference] {
	return func(yield func(*Reference) bool) {
		v := reflect.ValueOf(doc)
		for v := range iterValueOrReference(v, false) {
			ref := v.FieldByName("Reference").Interface().(*Reference)
			if ref != nil {
				if !yield(ref) {
					return
				}
			}
		}
	}
}

func (doc *OpenAPI) IterSchemaOrRef() iter.Seq[*oas31.Schema] {
	return func(yield func(*oas31.Schema) bool) {
		for v := range iterValueOrReference(reflect.ValueOf(doc.Paths), true) {
			if v.Type().Name() == valueOrReferenceOfSchema {
				p := reflect.New(v.Type())
				p.Elem().Set(v)
				if or, ok := p.Interface().(*oas31.Schema); ok {
					if !yield(or) {
						return
					}
				}
				v.Set(p.Elem())
			}
		}
		for v := range iterValueOrReference(reflect.ValueOf(doc.Webhooks), true) {
			if v.Type().Name() == valueOrReferenceOfSchema {
				p := reflect.New(v.Type())
				p.Elem().Set(v)
				if or, ok := p.Interface().(*oas31.Schema); ok {
					if !yield(or) {
						return
					}
				}
				v.Set(p.Elem())
			}
		}
		// do not yield Components
	}
}
