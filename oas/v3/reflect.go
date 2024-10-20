package oas

import (
	"context"
	"iter"
	"reflect"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/jsonschema/draft2020"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/oas31"
)

var referenceMixinPrefix = strings.TrimSuffix(reflect.TypeOf(draft2020.ReferenceMixin[bool]{}).Name(), "[bool]")
var trueWhenReferenceMixinPrefix = func(v reflect.Value) bool {
	return strings.Contains(v.Type().Name(), referenceMixinPrefix)
}

// setContext recursively find ReferenceMixin fields or elements and sets its ctx to ctx.
func setContext(v reflect.Value, ctx context.Context) {
	param1 := reflect.ValueOf(ctx)
	for v := range iterValueOrReference(v, trueWhenReferenceMixinPrefix, true) {
		f := v.MethodByName("WithContext")
		f.Call([]reflect.Value{param1})
	}
}

// iterValueOrReference recursively find ValueOrReferenceOf fields or elements and sets its Root to root.
//
// If makeCanSet is true, creates a pointer to the existing map value to make its value settable.
func iterValueOrReference(v reflect.Value, f func(v reflect.Value) bool, makeCanSet bool) iter.Seq[reflect.Value] {
	return func(yield func(reflect.Value) bool) {
		//fmt.Println(">> ", v.Type(), fmt.Sprintf("%q", fmt.Sprint(v.Interface())))
		switch v.Kind() {
		case reflect.Invalid:
			panic(v)
		case reflect.Ptr:
			if !v.IsNil() {
				for f := range iterValueOrReference(v.Elem(), f, makeCanSet) {
					if !yield(f) {
						return
					}
				}
			}
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				if f(v.Field(i)) {
					//fmt.Println(">>>", v.Field(i))
					if !yield(v.Field(i).Addr()) {
						return
					}
				}
				for f := range iterValueOrReference(v.Field(i), f, makeCanSet) {
					if !yield(f) {
						return
					}
				}
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				for f := range iterValueOrReference(v.Index(i).Addr(), f, makeCanSet) {
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
					for f := range iterValueOrReference(p, f, makeCanSet) {
						if !yield(f) {
							return
						}
					}
					v.SetMapIndex(it.Key(), p.Elem())
				} else {
					for f := range iterValueOrReference(value, f, makeCanSet) {
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

var schema = reflect.TypeOf(oas31.Schema{}).Name()
var trueWhenSchema = func(v reflect.Value) bool {
	return v.Type().Name() == schema
}

func (doc *OpenAPI) IterSchema() iter.Seq[*oas31.Schema] {
	return func(yield func(*oas31.Schema) bool) {
		for v := range iterValueOrReference(reflect.ValueOf(doc.Paths), trueWhenSchema, true) {
			p := reflect.New(v.Type())
			p.Elem().Set(v)
			if or, ok := p.Interface().(*oas31.Schema); ok {
				if !yield(or) {
					return
				}
			}
			v.Set(p.Elem())
		}
		for v := range iterValueOrReference(reflect.ValueOf(doc.Webhooks), trueWhenSchema, true) {
			p := reflect.New(v.Type())
			p.Elem().Set(v)
			if or, ok := p.Interface().(*oas31.Schema); ok {
				if !yield(or) {
					return
				}
			}
			v.Set(p.Elem())
		}
		// do not yield Components
	}
}
