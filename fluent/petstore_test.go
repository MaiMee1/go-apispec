package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/MaiMee1/go-apispec/fluent/operation"
	"github.com/MaiMee1/go-apispec/fluent/parameter"
	"github.com/MaiMee1/go-apispec/fluent/schema"
	"github.com/MaiMee1/go-apispec/fluent/schema/encoder"
	"github.com/MaiMee1/go-apispec/fluent/security"
	"github.com/MaiMee1/go-apispec/fluent/specs"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

type Tag struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Category struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Pet struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `json:"name" validate:"required"`
	Category  Category `json:"category"`
	PhotoUrls []string `json:"photoUrls" validate:"required"`
	Tags      []Tag    `json:"tags,omitempty"`
	Status    string   `json:"status,omitempty"`
}
type ApiResponse struct {
	Code    int32  `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

func TestFluent_PetStore(t *testing.T) {
	schema.WithEncoder(encoder.WithNameFilter(func(s string) string {
		return strings.TrimPrefix(s, "github.com.MaiMee1.go_apispec.fluent.")
	}))

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
		specs.WithComponents(
			"petstore_auth", security.NewOAuth2Scheme(
				security.WithImplicitFlow("https://petstore3.swagger.io/oauth/authorize", "",
					"write:pets", "modify pets in your account",
					"read:pets", "read your pets",
				),
			),
			"api_key", security.NewApiKeyScheme("api_key"),
		),
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
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
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
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("findPetsByStatus", http.MethodGet, "/pet/findByStatus",
			operation.WithSummary("Finds Pets by status"),
			operation.WithDescription("Multiple status values can be provided with comma separated strings"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Query("status", "Status values that need to be considered for filter", false, parameter.WithSchemaFor[string](
					schema.WithDefault("available"),
					schema.WithEnum("available", "pending", "sold"),
				), parameter.WithFormStyle(true)),
			),
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
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("findPetsByTags", http.MethodGet, "/pet/findByTags",
			operation.WithSummary("Finds Pets by tags"),
			operation.WithDescription("Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing."),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Query("tags", "Tags to filter by", false, parameter.WithSchemaFor[[]string](), parameter.WithFormStyle(true)),
			),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[[]Pet](),
				"application/xml", schema.For[[]Pet](),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid tag value"),
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("getPetById", http.MethodGet, "/pet/{pedId}",
			operation.WithSummary("Find pet by ID"),
			operation.WithDescription("Returns a single pet"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Path("petId", "ID of pet to return", true, parameter.WithSchemaFor[int64]()),
			),
			operation.WithResponse(http.StatusOK, "successful operation",
				"application/json", schema.For[Pet](),
				"application/xml", schema.For[Pet](),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid ID supplied"),
			operation.WithResponse(http.StatusNotFound, "Pet not found"),
			operation.WithSecurity(security.Scheme("api_key")),
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("updatePetWithForm", http.MethodPost, "/pet/{pedId}",
			operation.WithSummary("Updates a pet in the store with form data"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Path("petId", "ID of pet that needs to be updated", true, parameter.WithSchemaFor[int64]()),
				parameter.Query("name", "Name of pet that needs to be updated", false, parameter.WithSchemaFor[string]()),
				parameter.Query("status", "Status of pet that needs to be updated", false, parameter.WithSchemaFor[string]()),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid input"),
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("deletePet", http.MethodDelete, "/pet/{pedId}",
			operation.WithSummary("Deletes a pet"),
			operation.WithDescription("deletes a pet"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Header("api_key", "", false, parameter.WithSchemaFor[string]()),
				parameter.Path("petId", "Pet id to delete", true, parameter.WithSchemaFor[int64]()),
			),
			operation.WithResponse(http.StatusBadRequest, "Invalid pet value"),
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("uploadFile", http.MethodPost, "/pet/{pedId}/uploadImage",
			operation.WithSummary("uploads an image"),
			operation.WithTags("pet"),
			operation.WithParams(
				parameter.Path("petId", "ID of pet to update", true, parameter.WithSchemaFor[int64]()),
				parameter.Query("additionalMetadata", "Additional Metadata", false, parameter.WithSchemaFor[string]()),
			),
			operation.WithBody("", false, "application/octet-stream", schema.String(oas.BinaryFormat)),
			operation.WithResponse(http.StatusOK, "successful operation", "application/json", schema.For[ApiResponse]()),
			operation.WithSecurity(security.Scheme("petstore_auth", "write:pets", "read:pets")),
		),
		specs.WithOperation("getInventory", http.MethodPost, "/store/inventory",
			operation.WithSummary("Returns pet inventories by status"),
			operation.WithDescription("Returns a map of status codes to quantities"),
			operation.WithTags("store"),
			operation.WithResponse(http.StatusOK, "successful operation", "application/json", schema.For[map[string]int32]()),
			operation.WithSecurity(security.Scheme("api_key")),
		),
		specs.WithSchemaDefinitions(schema.Cached()),
	)
	if err != nil {
		t.Fatal(err)
	}
	_ = api
	t.Log(api.Json())
}
