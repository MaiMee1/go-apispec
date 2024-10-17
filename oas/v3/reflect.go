package oas

import (
	"iter"
	"reflect"
	"strings"
)

var valueOrReferenceOfPrefix = strings.TrimSuffix(reflect.TypeOf(ValueOrReferenceOf[bool]{}).Name(), "[bool]")
var valueOrReferenceOfSchema = reflect.TypeOf(ValueOrReferenceOf[Schema]{}).Name()

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
					yield(f)
				}
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				for f := range iterValueOrReference(v.Field(i), makeCanSet) {
					yield(f)
				}
			}
			if strings.HasPrefix(v.Type().Name(), valueOrReferenceOfPrefix) {
				//fmt.Println(">>>", v.Type())
				yield(v)
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				for f := range iterValueOrReference(v.Index(i).Addr(), makeCanSet) {
					yield(f)
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
						yield(f)
					}
					v.SetMapIndex(it.Key(), p.Elem())
				} else {
					for f := range iterValueOrReference(value, makeCanSet) {
						yield(f)
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
				yield(ref)
			}
		}
	}
}

func (doc *OpenAPI) IterSchemaOrRef() iter.Seq[ValueOrReferenceOf[Schema]] {
	return func(yield func(ValueOrReferenceOf[Schema]) bool) {
		v := reflect.ValueOf(doc)
		for v := range iterValueOrReference(v, false) {
			if v.Type().Name() == valueOrReferenceOfSchema {
				if or, ok := v.Interface().(ValueOrReferenceOf[Schema]); ok {
					yield(or)
				}
			}
		}
	}
}
