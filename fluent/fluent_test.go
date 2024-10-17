package main

import (
	"net/http"
	"testing"

	"github.com/MaiMee1/go-apispec/fluent/operation"
	"github.com/MaiMee1/go-apispec/fluent/parameter"
	"github.com/MaiMee1/go-apispec/fluent/schema"
	"github.com/MaiMee1/go-apispec/fluent/specs"
)

type Bug struct {
	X int         `json:"x"`
	Y []Bug       `json:"bugs"`
	Z interface{} `json:"z"`
}

type Test struct {
	Ant string `json:"ant"`
	Bug Bug    `json:"bug"`
}

func TestFluent(t *testing.T) {
	api, err := specs.New(
		specs.WithTitle("Test API"),
		specs.WithDescription("Test Description"),
		specs.WithOperation("getAsdsd", http.MethodGet, "/asfaf",
			operation.WithSummary(""),
			operation.WithParams(
				parameter.Query("page", "", false, parameter.WithSchemaFor[int]()),
				parameter.Query("limit", "", false, parameter.WithSchemaFor[int]()),
			),
			operation.WithBody("", true, "application/json", schema.For[int]()),
			operation.WithResponse(http.StatusOK, "successful operation", "application/json", schema.RefFor[Test]()),
			//operation.WithResponse(http.StatusBadRequest, "bad operation", "application/json", schema.For[Test]()),
		),
		specs.WithSchemaDefinitions(schema.Cached()),
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(api.Json())
}
