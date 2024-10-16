package main

import (
	"net/http"
	"testing"

	"github.com/MaiMee1/go-apispec/fluent/operation"
	"github.com/MaiMee1/go-apispec/fluent/parameter"
	"github.com/MaiMee1/go-apispec/fluent/schema"
	"github.com/MaiMee1/go-apispec/fluent/security"
	"github.com/MaiMee1/go-apispec/fluent/specs"
)

type Pet struct{}

func TestFluent_PetStore(t *testing.T) {
	api, err := specs.New(
		specs.WithTitle("Swagger Petstore - OpenAPI 3.0"),
		specs.WithDescription("This is a sample Pet Store Server based on the OpenAPI 3.0 specification.  You can find out more about\nSwagger at [https://swagger.io](https://swagger.io). In the third iteration of the pet store, we've switched to the design first approach!\nYou can now help us improve the API whether it's by making changes to the definition itself or to the code.\nThat way, with time, we can improve the API in general, and expose some of the new features in OAS3.\n\n_If you're looking for the Swagger 2.0/OAS 2.0 version of Petstore, then click [here](https://editor.swagger.io/?url=https://petstore.swagger.io/v2/swagger.yaml). Alternatively, you can load via the `Edit > Load Petstore OAS 2.0` menu option!_\n\nSome useful links:\n- [The Pet Store repository](https://github.com/swagger-api/swagger-petstore)\n- [The source API definition for the Pet Store](https://github.com/swagger-api/swagger-petstore/blob/master/src/main/resources/openapi.yaml)"),
		specs.WithTOS("http://swagger.io/terms/"),
		specs.WithContact("", "", "apiteam@swagger.io"),
		specs.WithLicense("Apache 2.0", "http://www.apache.org/licenses/LICENSE-2.0.html"),
		specs.WithVersion("1.0.11"),
		specs.WithExternalDocs("Find out more about Swagger", "http://swagger.io"),
		specs.WithServer("https", "petstore3.swagger.io", 0, "/api/v3"),
		specs.WithTag("pet", "Everything about your Pets"),
		specs.WithTag("store", "Access to Petstore orders"),
		specs.WithTag("user", "Operations about user"),
		specs.WithOperation("updatePet", http.MethodPut, "/pet",
			operation.WithSummary("Update an existing pet"),
			operation.WithDescription("Update an existing pet by Id"),
			operation.WithTags("pet"),
			operation.WithBody("Update an existent pet in the store", true,
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
				"application/x-www-form-urlencoded", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid ID supplied"),
			operation.WithResponse(http.StatusNotFound, "Pet not found"),
			operation.WithResponse(http.StatusUnprocessableEntity, "Validation exception"),
			operation.WithSecurity(security.All{
				"petstore_auth": []string{"write:pets", "read:pets"},
			}),
		),
		specs.WithOperation("addPet", http.MethodPost, "/pet",
			operation.WithSummary("Add a new pet to the store"),
			operation.WithDescription("Add a new pet to the store"),
			operation.WithTags("pet"),
			operation.WithBody("Create a new pet in the store", true,
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
				"application/x-www-form-urlencoded", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid ID supplied"),
			operation.WithResponse(http.StatusUnprocessableEntity, "Validation exception"),
			operation.WithSecurity(security.All{
				"petstore_auth": []string{"write:pets", "read:pets"},
			}),
		),
		specs.WithOperation("findPetsByStatus", http.MethodPost, "/pet/findByStatus",
			operation.WithSummary("Finds Pets by status"),
			operation.WithDescription("Multiple status values can be provided with comma separated strings"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Query("status", "Status values that need to be considered for filter", false, parameter.WithSchemaFor[string](
					schema.WithDefault("available"),
					schema.WithEnum("available", "pending", "sold"),
				), parameter.WithFormStyle(true))),
			operation.WithBody("Create a new pet in the store", true,
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
				"application/x-www-form-urlencoded", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid ID supplied"),
			operation.WithResponse(http.StatusUnprocessableEntity, "Validation exception"),
			operation.WithSecurity(security.All{
				"petstore_auth": []string{"write:pets", "read:pets"},
			}),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	_ = api
	t.Log(api.Json())
}
