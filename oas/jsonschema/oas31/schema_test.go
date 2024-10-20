package oas31

import (
	"encoding/json"
	"testing"
)

func TestSchema_UnmarshalJSON(t *testing.T) {
	var s Schema
	data := `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/schema-optional",
  "$dynamicAnchor": "meta",
  "$vocabulary": {
    "https://json-schema.org/draft/2020-12/vocab/core": true,
    "https://example.com/vocab/example-vocab": false
  },
  "allOf": [
    { "$ref": "https://json-schema.org/draft/2020-12/meta/core" },
    { "$ref": "https://example.com/meta/example-vocab" }
  ]
}`
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestSchema_MarshalJSON(t *testing.T) {
	var s Schema
	data := `{"$schema":"https://example.com/schema-required","$id":"https://my-schema.com"}`
	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		t.Fatal(err)
	}
	// Act
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != data {
		t.Error("not equal")
	}
	t.Log(string(b))
}
