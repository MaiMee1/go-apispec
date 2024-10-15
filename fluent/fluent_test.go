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
				parameter.Query("page", "", false),
				parameter.Query("limit", "", false),
			),
			operation.WithBody("", true, "application/json", 1),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[Test]()),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	_ = api
	t.Log(api.Json())
}
