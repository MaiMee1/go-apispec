package oas

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
)

func New(filename string) (*OpenAPI, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var root interface{}
	if err := json.Unmarshal(file, &root); err != nil {
		return nil, err
	}
	var document OpenAPI
	if err = json.Unmarshal(file, &document); err != nil {
		return nil, err
	}

	v := reflect.ValueOf(&document)
	setRoot(v, root)

	if err := validate.Struct(document); err != nil {
		return &document, err
	}
	return &document, nil
}

var valueOrReferenceOfPrefix = strings.TrimSuffix(reflect.TypeOf(ValueOrReferenceOf[bool]{}).Name(), "[bool]")

// setRoot recursively find ValueOrReferenceOf fields or elements and sets its Root to root.
func setRoot(v reflect.Value, root interface{}) {
	//fmt.Println(">> ", v.Type(), fmt.Sprintf("%q", fmt.Sprint(v.Interface())))
	switch v.Kind() {
	case reflect.Invalid:
		panic(v)
	case reflect.Ptr:
		if !v.IsNil() {
			setRoot(v.Elem(), root)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			setRoot(f, root)
		}
		if strings.HasPrefix(v.Type().Name(), valueOrReferenceOfPrefix) {
			//fmt.Println(">>>", v.Type())
			f := v.FieldByName("Root")
			f.Set(reflect.ValueOf(root))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			setRoot(v.Index(i).Addr(), root)
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			p := reflect.New(iter.Value().Type())
			p.Elem().Set(iter.Value())
			setRoot(p, root)
			v.SetMapIndex(iter.Key(), p.Elem())
		}
	default:
	}
}
