package oas

import (
	"context"
	"encoding/json"
	"os"
	"reflect"

	"github.com/MaiMee1/go-apispec/oas/internal/validate"
)

func New(filename string) (*OpenAPI, error) {
	ctx := context.TODO()

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var document OpenAPI
	if err = json.Unmarshal(file, &document); err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, "Root", document)
	setContext(reflect.ValueOf(&document), ctx)

	if err := validate.Struct(document); err != nil {
		return &document, err
	}
	return &document, nil
}
