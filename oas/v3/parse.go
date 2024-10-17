package oas

import (
	"encoding/json"
	"os"
	"reflect"
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
