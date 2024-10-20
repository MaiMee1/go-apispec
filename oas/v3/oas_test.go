package oas

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/MaiMee1/go-apispec/oas/internal/validate"
)

func TestOpenAPI_UnmarshalJSON(t *testing.T) {
	ctx := context.TODO()

	file, err := os.ReadFile("testdata/petstore.json")
	if err != nil {
		t.Fatal(err)
	}

	var root interface{}
	if err := json.Unmarshal(file, &root); err != nil {
		t.Fatal(err)
	}
	ctx = context.WithValue(ctx, "Root", root)

	var document OpenAPI
	if err = json.Unmarshal(file, &document); err != nil {
		t.Fatal(err)
	}
	v := reflect.ValueOf(&document)
	setRoot(v, root)

	if err := validate.Struct(document); err != nil {
		t.Error(err)
	}

	if document.Version != "3.0.3" {
		t.Error(document.Version)
	}
	for k, m := range document.Paths["/pet"].Put.RequestBody.Value.Content {
		b, err := json.Marshal(m.Schema.Resolve(ctx))
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v, %s", k, b)
	}
	t.Log(document.Paths["/pet/findByStatus"].Get.Summary)
}

func FuzzOpenAPI(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		var document OpenAPI
		err := json.Unmarshal([]byte(fmt.Sprintf(`{"openapi":"%s"}`, s)), &document)
		if err != nil {
			t.Fatal(err)
		}
		if document.Version.Validate() != nil {
			t.Error("expected valid openapi version")
		}
	})
}
