package oas

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/iana"
	"github.com/MaiMee1/go-apispec/oas/internal/validate"
	"github.com/MaiMee1/go-apispec/oas/jsonschema"
	"github.com/MaiMee1/go-apispec/oas/jsonschema/oas31"
)

type (
	Schema                = oas31.Schema
	ExternalDocumentation = oas31.ExternalDocumentation
	Discriminator         = oas31.Discriminator
	XML                   = oas31.XML
	// SpecificationExtension properties are implemented as patterned fields that are always prefixed by "x-".
	SpecificationExtension = oas31.SpecificationExtension
)

type (
	SemanticVersion string
	Type            = jsonschema.Type
	Format          = jsonschema.Format
	Location        int8
	Style           int8
	Scheme          int8
)

const (
	Int32Format  Format = "int32"
	Int64Format  Format = "int64"
	FloatFormat  Format = "float"
	DoubleFormat Format = "double"

	Base64Format   Format = "base64"   // base64 encoded characters
	BinaryFormat   Format = "binary"   // octet-stream
	PasswordFormat Format = "password" // A hint to UIs to obscure input.
)

const (
	QueryLocation Location = iota + 1
	HeaderLocation
	PathLocation
	CookieLocation
)

var locationToString = []string{
	0:              "<0>",
	QueryLocation:  "query",
	HeaderLocation: "header",
	PathLocation:   "path",
	CookieLocation: "cookie",
}

func (l Location) String() string {
	return locationToString[l]
}

const (
	MatrixStyle Style = iota + 1
	LabelStyle
	FormStyle
	SimpleStyle
	SpaceDelimitedStyle
	PipeDelimitedStyle
	DeepObjectStyle
)

var styleToString = []string{
	0:                   "<0>",
	MatrixStyle:         "matrix",
	LabelStyle:          "label",
	FormStyle:           "form",
	SimpleStyle:         "simple",
	SpaceDelimitedStyle: "spaceDelimited",
	PipeDelimitedStyle:  "pipeDelimited",
	DeepObjectStyle:     "deepObject",
}

func (s Style) String() string {
	return styleToString[s]
}

const (
	ApiKeyScheme Scheme = iota + 1
	HttpScheme
	MutualTLSScheme
	OAuth2Scheme
	OpenIdConnectScheme
)

var schemeToString = []string{
	0:                   "<0>",
	ApiKeyScheme:        "apiKey",
	HttpScheme:          "http",
	MutualTLSScheme:     "mutualTLS",
	OAuth2Scheme:        "oauth2",
	OpenIdConnectScheme: "openIdConnect",
}

func (s Scheme) String() string {
	return schemeToString[s]
}

func (v SemanticVersion) Validate() error {
	// TODO:
	return nil
}

func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

//goland:noinspection GoMixedReceiverTypes
func (l *Location) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if i := slices.Index(locationToString, s); i != -1 {
		*l = Location(i)
		return nil
	}
	return fmt.Errorf("invalid location %q", s)
}

func (s Style) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

//goland:noinspection GoMixedReceiverTypes
func (s *Style) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}

	if i := slices.Index(styleToString, st); i != -1 {
		*s = Style(i)
		return nil
	}
	return fmt.Errorf("invalid style %q", st)
}

func (s Scheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

//goland:noinspection GoMixedReceiverTypes
func (s *Scheme) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}

	if i := slices.Index(schemeToString, st); i != -1 {
		*s = Scheme(i)
		return nil
	}
	return fmt.Errorf("invalid scheme %q", st)
}

type RuntimeExpression string

type ValueOrReferenceOf[T any] struct {
	Value     T
	Reference *Reference
	Root      interface{} // a reference to the root document, set via reflection
}

func (o ValueOrReferenceOf[T]) Resolve(v ...interface{}) T {
	root := o.Root
	if root == nil {
		if len(v) > 0 {
			root = v[0]
		} else {
			panic("root not available")
		}
	}

	var t T
	if o.Reference != nil {
		b, _ := json.Marshal(resolve(o.Reference, root))
		_ = json.Unmarshal(b, &t)
		setRoot(reflect.ValueOf(t), root)
		return t
	}
	return o.Value
}

//goland:noinspection GoMixedReceiverTypes
func (o *ValueOrReferenceOf[T]) UnmarshalJSON(b []byte) error {
	var r Reference
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	o.Value = v
	if r.Ref != "" {
		o.Reference = &r
	}
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (o ValueOrReferenceOf[T]) MarshalJSON() ([]byte, error) {
	if o.Reference != nil {
		return json.Marshal(o.Reference)
	}
	return json.Marshal(o.Value)
}

type DataType struct {
	Type   jsonschema.Type `json:"type" validate:"required"`
	Format Format          `json:"format,omitempty"`
}

// UrlTemplate supports ServerVariable and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
type UrlTemplate string

// RichText supports CommonMark markdown formatting.
type RichText string

// Default returns a minimal starting OpenAPI specs.
func Default() OpenAPI {
	document := OpenAPI{}
	document.Version = "3.1.0"
	document.Info.Title = "Unnamed API"
	document.Info.Version = "1.0.0"
	document.Paths = make(Paths)
	return document
}

type OpenAPI struct {
	Version      SemanticVersion                         `json:"openapi" validate:"required"`
	Info         Info                                    `json:"info" validate:"required"`
	Servers      []Server                                `json:"servers,omitempty" validate:"dive"`
	Paths        Paths                                   `json:"paths,omitempty" validate:"dive"`
	Webhooks     map[string]ValueOrReferenceOf[PathItem] `json:"webhooks,omitempty" validate:"dive"`
	Components   Components                              `json:"components"`
	Security     []SecurityRequirement                   `json:"security,omitempty" validate:"dive"`
	Tags         []Tag                                   `json:"tags,omitempty" validate:"dive"`
	ExternalDocs *oas31.ExternalDocumentation            `json:"externalDocs,omitempty"`
	Extensions   SpecificationExtension                  `json:"-"`
}

func (doc *OpenAPI) Validate() error {
	return validate.Struct(doc)
}

// Info provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
type Info struct {
	Title          string                 `json:"title" validate:"required"`
	Description    RichText               `json:"description"`
	TermsOfService string                 `json:"termsOfService,omitempty" validate:"omitempty,url"`
	Contact        *Contact               `json:"contact,omitempty"`
	License        *License               `json:"license,omitempty"`
	Version        string                 `json:"version" validate:"required"`
	Extensions     SpecificationExtension `json:"-"`
}

type Contact struct {
	Name       string                 `json:"name,omitempty"`
	Url        string                 `json:"url,omitempty" validate:"omitempty,url"`
	Email      string                 `json:"email,omitempty" validate:"omitempty,email"`
	Extensions SpecificationExtension `json:"-"`
}

type License struct {
	Name       string                 `json:"name" validate:"required"`
	Url        string                 `json:"url,omitempty" validate:"url,omitempty"`
	Extensions SpecificationExtension `json:"-"`
}

type Server struct {
	Url         UrlTemplate               `json:"url" validate:"required"`
	Description RichText                  `json:"description"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" validate:"dive"`
	Extensions  SpecificationExtension    `json:"-"`
}

type ServerVariable struct {
	Enum        []string               `json:"enum,omitempty" validate:"min=1"`
	Default     string                 `json:"default" validate:"required"`
	Description RichText               `json:"description"`
	Extensions  SpecificationExtension `json:"-"`
}

type Components struct {
	Schemas         map[string]oas31.Schema   `json:"schemas,omitempty" validate:"dive"`
	Responses       map[string]Response       `json:"responses,omitempty" validate:"dive"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty" validate:"dive"`
	Examples        map[string]Example        `json:"examples,omitempty" validate:"dive"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty" validate:"dive"`
	Headers         map[string]Header         `json:"headers,omitempty" validate:"dive"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty" validate:"dive"`
	Links           map[string]Link           `json:"links,omitempty" validate:"dive"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty" validate:"dive"`
	PathItems       map[string]PathItem       `json:"pathItems,omitempty" validate:"dive"`
	Extensions      SpecificationExtension    `json:"-"`
}

var fieldNameRe = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)

func isValidComponentsKey(key string) bool {
	return fieldNameRe.Match([]byte(key))
}

type Paths map[string]PathItem

type PathItem struct {
	Ref         string                          `json:"$ref,omitempty" validate:"omitempty,uri"`
	Summary     string                          `json:"summary,omitempty"`
	Description RichText                        `json:"description,omitempty"`
	Get         *Operation                      `json:"get,omitempty"`
	Put         *Operation                      `json:"put,omitempty"`
	Post        *Operation                      `json:"post,omitempty"`
	Delete      *Operation                      `json:"delete,omitempty"`
	Options     *Operation                      `json:"options,omitempty"`
	Head        *Operation                      `json:"head,omitempty"`
	Patch       *Operation                      `json:"patch,omitempty"`
	Trace       *Operation                      `json:"trace,omitempty"`
	Servers     []Server                        `json:"servers,omitempty"`
	Parameters  []ValueOrReferenceOf[Parameter] `json:"parameters,omitempty"`
	Extensions  SpecificationExtension          `json:"-"`
}

func (i *PathItem) Range() map[string]Operation {
	m := make(map[string]Operation)
	if i.Get != nil {
		m["GET"] = *i.Get
	}
	if i.Put != nil {
		m["PUT"] = *i.Put
	}
	if i.Post != nil {
		m["POST"] = *i.Post
	}
	if i.Delete != nil {
		m["DELETE"] = *i.Delete
	}
	if i.Options != nil {
		m["OPTIONS"] = *i.Options
	}
	if i.Head != nil {
		m["HEAD"] = *i.Head
	}
	if i.Patch != nil {
		m["PATCH"] = *i.Patch
	}
	if i.Trace != nil {
		m["TRACE"] = *i.Trace
	}
	return m
}

type Operation struct {
	Tags         []string                                `json:"tags,omitempty"`
	Summary      string                                  `json:"summary"`
	Description  RichText                                `json:"description"`
	ExternalDocs *oas31.ExternalDocumentation            `json:"externalDocs,omitempty"`
	OperationId  string                                  `json:"operationId,omitempty"`
	Parameters   []ValueOrReferenceOf[Parameter]         `json:"parameters,omitempty" validate:"dive"`
	RequestBody  *ValueOrReferenceOf[RequestBody]        `json:"requestBody,omitempty"`
	Responses    Responses                               `json:"responses,omitempty" validate:"dive"`
	Callbacks    map[string]ValueOrReferenceOf[Callback] `json:"callbacks,omitempty" validate:"dive"`
	Deprecated   bool                                    `json:"deprecated,omitempty"`
	Security     []SecurityRequirement                   `json:"security,omitempty" validate:"dive"`
	Servers      []Server                                `json:"servers,omitempty" validate:"dive"`
	Extensions   SpecificationExtension                  `json:"-"`
}

type Parameter struct {
	Name            string                                 `json:"name" validate:"required"`
	In              Location                               `json:"in" validate:"required"`
	Description     RichText                               `json:"description"`
	Required        bool                                   `json:"required" validate:"required_if=In 3"`
	Deprecated      bool                                   `json:"deprecated,omitempty"`
	AllowEmptyValue bool                                   `json:"allowEmptyValue,omitempty"` // Deprecated
	Style           Style                                  `json:"style,omitempty"`
	Explode         *bool                                  `json:"explode,omitempty"`
	AllowReserved   bool                                   `json:"allowReserved,omitempty"`
	Schema          *oas31.Schema                          `json:"schema,omitempty" validate:"required_without=Content"`
	Content         map[string]MediaType                   `json:"content,omitempty" validate:"required_without=Schema"`
	Example         interface{}                            `json:"example,omitempty"`
	Examples        map[string]ValueOrReferenceOf[Example] `json:"examples,omitempty"`
	Extensions      SpecificationExtension                 `json:"-"`
}

type RequestBody struct {
	Description RichText               `json:"description"`
	Content     map[string]MediaType   `json:"content" validate:"required"`
	Required    bool                   `json:"required"`
	Extensions  SpecificationExtension `json:"-"`
}

type MediaType struct {
	Schema     *oas31.Schema                          `json:"schema,omitempty"`
	Example    interface{}                            `json:"example,omitempty"`
	Examples   map[string]ValueOrReferenceOf[Example] `json:"examples,omitempty"`
	Encoding   map[string]Encoding                    `json:"encoding,omitempty"`
	Extensions SpecificationExtension                 `json:"-"`
}

type Encoding struct {
	ContentType   string                                `json:"contentType"`
	Headers       map[string]ValueOrReferenceOf[Header] `json:"headers,omitempty"`
	Style         Style                                 `json:"style,omitempty"`
	Explode       *bool                                 `json:"explode,omitempty"`
	AllowReserved bool                                  `json:"allowReserved,omitempty"`
	Extensions    SpecificationExtension                `json:"-"`
}

type Responses map[string]ValueOrReferenceOf[Response]

func (r Responses) Validate() error {
	for key := range maps.Keys(r) {
		if !isValidResponsesKey(key) {
			return fmt.Errorf("invalid response key: %s", key)
		}
	}
	return nil
}

var httpStatusCodeRe = regexp.MustCompile("^[1-5][0-9X]{2}$")

func isValidResponsesKey(key string) bool {
	if key != "default" && !httpStatusCodeRe.Match([]byte(key)) {
		return false
	}
	return true
}

func isValidContentKey(key string) bool {
	// TODO: media type or media type range
	return true
}

type Response struct {
	Description RichText                              `json:"description" validate:"required"`
	Headers     map[string]ValueOrReferenceOf[Header] `json:"headers,omitempty"`
	Content     map[string]MediaType                  `json:"content,omitempty"`
	Links       map[string]ValueOrReferenceOf[Link]   `json:"links,omitempty"`
	Extensions  SpecificationExtension                `json:"-"`
}

type Callback map[RuntimeExpression]ValueOrReferenceOf[PathItem]

func isValidCallbackKey(key string) bool {
	// TODO: valid runtime expression
	return true
}

type Example struct {
	Summary       string                 `json:"summary"`
	Description   RichText               `json:"description"`
	Value         interface{}            `json:"value,omitempty" validate:"excluded_with=ExternalValue"`
	ExternalValue string                 `json:"externalValue,omitempty" validate:"excluded_with=Value"`
	Extensions    SpecificationExtension `json:"-"`
}

type Link struct {
	OperationRef string                 `json:"operationRef,omitempty" validate:"required_without=OperationId"`
	OperationId  string                 `json:"operationId,omitempty" validate:"required_without=OperationRef"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	RequestBody  []interface{}          `json:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Server       string                 `json:"server,omitempty"`
	Extensions   SpecificationExtension `json:"-"`
}

type Header struct {
	Description     RichText                               `json:"description"`
	Required        bool                                   `json:"required" validate:"required_if=In path"`
	Deprecated      bool                                   `json:"deprecated"`
	AllowEmptyValue bool                                   `json:"allowEmptyValue,omitempty"` // Deprecated
	Style           Style                                  `json:"style,omitempty"`
	Explode         *bool                                  `json:"explode,omitempty"`
	AllowReserved   bool                                   `json:"allowReserved,omitempty"`
	Schema          *oas31.Schema                          `json:"schema,omitempty" validate:"required_without=Content"`
	Content         map[string]MediaType                   `json:"content,omitempty" validate:"required_without=Schema"`
	Example         interface{}                            `json:"example,omitempty"`
	Examples        map[string]ValueOrReferenceOf[Example] `json:"examples,omitempty"`
	Extensions      SpecificationExtension                 `json:"-"`
}

type Tag struct {
	Name         string                       `json:"name" validate:"required"`
	Description  RichText                     `json:"description,omitempty"`
	ExternalDocs *oas31.ExternalDocumentation `json:"externalDocs,omitempty"`
	Extensions   SpecificationExtension       `json:"-"`
}

// Reference is defined by https://datatracker.ietf.org/doc/html/draft-pbryan-zyp-json-ref-03d follows the same structure, behavior and rules.
type Reference struct {
	Ref string `json:"$ref" validate:"required,url_fragment"`
}

func resolve(r *Reference, v interface{}) any {
	parts := strings.Split(r.Ref, "/")
	if parts[0] != "#" {
		panic(fmt.Errorf("invalid reference format: %s", r.Ref))
	}
	for _, part := range parts[1:] {
		t, ok := v.(map[string]interface{})[part]
		if !ok {
			panic(fmt.Sprintf("key not found: %s", part))
		}
		v = t
	}
	return v
}

type SecurityScheme struct {
	Type             Scheme                 `json:"type" validate:"required"`
	Description      RichText               `json:"description,omitempty"`
	Name             string                 `json:"name,omitempty"  validate:"required_if=Type 1,excluded_unless=Type 1"`
	In               Location               `json:"in,omitempty"  validate:"required_if=Type 1,excluded_unless=Type 1,omitempty,oneof=1 2 4"`
	Scheme           iana.AuthScheme        `json:"scheme,omitempty" validate:"required_if=Type 2,excluded_unless=Type 2"`
	BearerFormat     string                 `json:"bearerFormat,omitempty" validate:"excluded_unless=Scheme bearer"`
	Flows            *OAuthFlows            `json:"flows,omitempty" validate:"required_if=Type 4,excluded_unless=Type 4"`
	OpenIdConnectUrl string                 `json:"openIdConnectUrl,omitempty" validate:"required_if=Type 5,excluded_unless=Type 5,omitempty,url"`
	Extensions       SpecificationExtension `json:"-"`
}

type OAuthFlows struct {
	Implicit          *ImplicitOAuthFlow          `json:"implicit,omitempty"`
	Password          *PasswordOAuthFlow          `json:"password,omitempty"`
	ClientCredentials *ClientCredentialsOAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *AuthorizationCodeOAuthFlow `json:"authorizationCode,omitempty"`
	Extensions        SpecificationExtension      `json:"-"`
}

type ImplicitOAuthFlow struct {
	authorizationUrlMixin
	oAuthFlow
}

type PasswordOAuthFlow struct {
	tokenUrlMixin
	oAuthFlow
}

type ClientCredentialsOAuthFlow struct {
	tokenUrlMixin
	oAuthFlow
}

type AuthorizationCodeOAuthFlow struct {
	authorizationUrlMixin
	tokenUrlMixin
	oAuthFlow
}

type authorizationUrlMixin struct {
	AuthorizationUrl string `json:"authorizationUrl" validate:"required,url"`
}

type tokenUrlMixin struct {
	TokenUrl string `json:"tokenUrl" validate:"required,url"`
}

type oAuthFlow struct {
	RefreshUrl string                 `json:"refreshUrl,omitempty" validate:"omitempty,url"`
	Scopes     map[string]string      `json:"scopes" validate:"required"`
	Extensions SpecificationExtension `json:"-"`
}

type SecurityRequirement map[string][]string
