package oas

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOpenAPI_UnmarshalJSON(t *testing.T) {
	document, err := New("testdata/petstore.json")
	if err != nil {
		t.Error(err)
	}

	if document.Version != "3.0.3" {
		t.Error(document.Version)
	}
	for k, m := range document.Paths["/pet"].Put.RequestBody.Content {
		b, err := json.Marshal(m.Schema.Resolve())
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
