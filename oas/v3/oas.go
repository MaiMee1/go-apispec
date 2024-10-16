package oas

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaiMee1/go-apispec/oas/iana"
)

type (
	SemanticVersion string
	Type            int8
	Format          string
	Location        int8
	Style           int8
	Scheme          int8
)

const (
	IntegerType Type = iota + 1 // JSON number without a fraction or exponent part
	NumberType
	StringType
	BooleanType
	ObjectType
	ArrayType
)

func (t Type) String() string {
	switch t {
	case IntegerType:
		return "integer"
	case NumberType:
		return "number"
	case StringType:
		return "string"
	case BooleanType:
		return "boolean"
	case ObjectType:
		return "object"
	case ArrayType:
		return "array"
	case 0:
		return "<0>"
	default:
		panic(t)
	}
}

const (
	NoFormat       Format = ""
	Int32Format    Format = "int32" // signed 32 bits
	Int64Format    Format = "int64" // signed 64 bits (a.k.a long)
	FloatFormat    Format = "float"
	DoubleFormat   Format = "double"
	Base64Format   Format = "base64" // base64 encoded characters
	BinaryFormat   Format = "binary" // octet-stream
	DateFormat     Format = "binary" // As defined by full-date - RFC3339
	DateTimeFormat Format = "binary" // As defined by date-time - RFC3339
	PasswordFormat Format = "binary" // A hint to UIs to obscure input.
)

const (
	QueryLocation Location = iota + 1
	HeaderLocation
	PathLocation
	CookieLocation
)

func (l Location) String() string {
	switch l {
	case QueryLocation:
		return "query"
	case HeaderLocation:
		return "header"
	case PathLocation:
		return "path"
	case CookieLocation:
		return "cookie"
	case 0:
		return "<0>"
	default:
		panic(l)
	}
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

func (s Style) String() string {
	switch s {
	case MatrixStyle:
		return "matrix"
	case LabelStyle:
		return "label"
	case FormStyle:
		return "form"
	case SimpleStyle:
		return "simple"
	case SpaceDelimitedStyle:
		return "spaceDelimited"
	case PipeDelimitedStyle:
		return "pipeDelimited"
	case DeepObjectStyle:
		return "deepObject"
	case 0:
		return "<0>"
	default:
		panic(s)
	}
}

const (
	ApiKeyScheme Scheme = iota + 1
	HttpScheme
	MutualTLSScheme
	OAuth2Scheme
	OpenIdConnectScheme
)

func (s Scheme) String() string {
	switch s {
	case ApiKeyScheme:
		return "apiKey"
	case HttpScheme:
		return "http"
	case MutualTLSScheme:
		return "mutualTLS"
	case OAuth2Scheme:
		return "oauth2"
	case OpenIdConnectScheme:
		return "openIdConnect"
	case 0:
		return "<0>"
	default:
		panic(s)
	}
}

func (v SemanticVersion) Validate() error {
	// TODO:
	return nil
}

func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

//goland:noinspection GoMixedReceiverTypes
func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		//panic(fmt.Errorf("%s %v", b, err))
		return err
	}
	switch s {
	case IntegerType.String():
		*t = IntegerType
	case NumberType.String():
		*t = NumberType
	case StringType.String():
		*t = StringType
	case BooleanType.String():
		*t = BooleanType
	case ObjectType.String():
		*t = ObjectType
	case ArrayType.String():
		*t = ArrayType
	default:
		return fmt.Errorf("invalid type %q", s)
	}
	return nil
}

func (f Format) Validate() error {
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
	switch s {
	case QueryLocation.String():
		*l = QueryLocation
	case HeaderLocation.String():
		*l = HeaderLocation
	case PathLocation.String():
		*l = PathLocation
	case CookieLocation.String():
		*l = CookieLocation
	default:
		return fmt.Errorf("invalid location %q", s)
	}
	return nil
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
	switch st {
	case MatrixStyle.String():
		*s = MatrixStyle
	case LabelStyle.String():
		*s = LabelStyle
	case FormStyle.String():
		*s = FormStyle
	case SimpleStyle.String():
		*s = SimpleStyle
	case SpaceDelimitedStyle.String():
		*s = SpaceDelimitedStyle
	case PipeDelimitedStyle.String():
		*s = PipeDelimitedStyle
	case DeepObjectStyle.String():
		*s = DeepObjectStyle
	default:
		return fmt.Errorf("invalid style %q", st)
	}
	return nil
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
	switch st {
	case ApiKeyScheme.String():
		*s = ApiKeyScheme
	case HttpScheme.String():
		*s = HttpScheme
	case MutualTLSScheme.String():
		*s = MutualTLSScheme
	case OAuth2Scheme.String():
		*s = OAuth2Scheme
	case OpenIdConnectScheme.String():
		*s = OpenIdConnectScheme
	default:
		return fmt.Errorf("invalid scheme %q", st)
	}
	return nil
}

type RuntimeExpression string

type Or[A, B comparable] struct {
	X A
	Y B
}

//goland:noinspection GoMixedReceiverTypes
func (o *Or[A, B]) UnmarshalJSON(b []byte) error {
	var (
		x A
		y B
	)
	err1 := json.Unmarshal(b, &x)
	if err1 == nil {
		o.X = x
	}

	err2 := json.Unmarshal(b, &y)
	if err2 == nil {
		o.Y = y
	}

	if err2 == nil || err1 == nil {
		return nil
	}
	return fmt.Errorf("or[%T, %T]: %w: %w", x, y, err2, err1)
}

//goland:noinspection GoMixedReceiverTypes
func (o Or[A, B]) MarshalJSON() ([]byte, error) {
	var x A
	if o.X != x {
		return json.Marshal(o.X)
	}
	return json.Marshal(o.Y)
}

type ValueOrReferenceOf[T any] struct {
	Value T
	Ref   Reference
	Root  interface{} // a reference to the root document, set via reflection
}

func (r ValueOrReferenceOf[T]) Resolve(v ...interface{}) T {
	root := r.Root
	if root == nil {
		if len(v) > 0 {
			root = v[0]
		} else {
			panic("root not available")
		}
	}

	var t T
	if r.Ref.Ref != "" {
		b, _ := json.Marshal(resolve(r.Ref, root))
		_ = json.Unmarshal(b, &t)
		setRoot(reflect.ValueOf(t), root)
		return t
	}
	return r.Value
}

//goland:noinspection GoMixedReceiverTypes
func (r *ValueOrReferenceOf[T]) UnmarshalJSON(b []byte) error {
	var s Reference
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	r.Value = v
	r.Ref = s
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (r ValueOrReferenceOf[T]) MarshalJSON() ([]byte, error) {
	if r.Ref.Ref != "" {
		return json.Marshal(r.Ref)
	}
	return json.Marshal(r.Value)
}

type DataType struct {
	Type   Type   `json:"type" validate:"required"`
	Format Format `json:"format,omitempty"`
}

// UrlTemplate supports ServerVariable and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
type UrlTemplate string

// RichText supports CommonMark markdown formatting.
type RichText string

// SpecificationExtension properties are implemented as patterned fields that are always prefixed by "x-".
type SpecificationExtension map[string]interface{}

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
	Servers      []Server                                `json:"servers,omitempty"`
	Paths        Paths                                   `json:"paths,omitempty"`
	Webhooks     map[string]ValueOrReferenceOf[PathItem] `json:"webhooks,omitempty"`
	Components   Components                              `json:"components"`
	Security     []SecurityRequirement                   `json:"security,omitempty"`
	Tags         []Tag                                   `json:"tags,omitempty"`
	ExternalDocs *ExternalDocumentation                  `json:"externalDocs,omitempty"`
	Extensions   SpecificationExtension                  `json:"-"`
}

func (o *OpenAPI) Validate() error {
	return validate.Struct(o)
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
	Variables   map[string]ServerVariable `json:"variables,omitempty"`
	Extensions  SpecificationExtension    `json:"-"`
}

type ServerVariable struct {
	Enum        []string               `json:"enum,omitempty" validate:"min=1"`
	Default     string                 `json:"default" validate:"required"`
	Description RichText               `json:"description"`
	Extensions  SpecificationExtension `json:"-"`
}

type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`
	Examples        map[string]Example        `json:"examples,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]Link           `json:"links,omitempty"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty"`
	PathItems       map[string]PathItem       `json:"pathItems,omitempty"`
	Extensions      SpecificationExtension    `json:"-"`
}

var fieldNameRe = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)

func isValidComponentsKey(key string) bool {
	return fieldNameRe.Match([]byte(key))
}

type Paths map[string]PathItem

type PathItem struct {
	Ref         string                          `json:"$ref,omitempty" validate:"uri"`
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
	ExternalDocs *ExternalDocumentation                  `json:"externalDocs,omitempty"`
	OperationId  string                                  `json:"operationId,omitempty"`
	Parameters   []ValueOrReferenceOf[Parameter]         `json:"parameters,omitempty"`
	RequestBody  *ValueOrReferenceOf[RequestBody]        `json:"requestBody,omitempty"`
	Responses    Responses                               `json:"responses,omitempty"`
	Callbacks    map[string]ValueOrReferenceOf[Callback] `json:"callbacks,omitempty"`
	Deprecated   bool                                    `json:"deprecated,omitempty"`
	Security     []SecurityRequirement                   `json:"security,omitempty"`
	Servers      []Server                                `json:"servers,omitempty"`
	Extensions   SpecificationExtension                  `json:"-"`
}

type ExternalDocumentation struct {
	Description RichText               `json:"description"`
	Url         string                 `json:"url" validate:"required,url"`
	Extensions  SpecificationExtension `json:"-"`
}

type Parameter struct {
	Name            string                                 `json:"name" validate:"required"`
	In              Location                               `json:"in" validate:"required"`
	Description     RichText                               `json:"description"`
	Required        bool                                   `json:"required" validate:"required_if=In path"`
	Deprecated      bool                                   `json:"deprecated,omitempty"`
	AllowEmptyValue bool                                   `json:"allowEmptyValue,omitempty"` // Deprecated
	Style           Style                                  `json:"style,omitempty"`
	Explode         *bool                                  `json:"explode,omitempty"`
	AllowReserved   bool                                   `json:"allowReserved,omitempty"`
	Schema          *ValueOrReferenceOf[Schema]            `json:"schema,omitempty" validate:"required_without=Content"`
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
	Schema     ValueOrReferenceOf[Schema]             `json:"schema,omitempty"`
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
	Schema          *ValueOrReferenceOf[Schema]            `json:"schema,omitempty" validate:"required_without=Content"`
	Content         map[string]MediaType                   `json:"content,omitempty" validate:"required_without=Schema"`
	Example         interface{}                            `json:"example,omitempty"`
	Examples        map[string]ValueOrReferenceOf[Example] `json:"examples,omitempty"`
	Extensions      SpecificationExtension                 `json:"-"`
}

type Tag struct {
	Name         string                 `json:"name" validate:"required"`
	Description  RichText               `json:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
	Extensions   SpecificationExtension `json:"-"`
}

// Reference is defined by https://datatracker.ietf.org/doc/html/draft-pbryan-zyp-json-ref-03d follows the same structure, behavior and rules.
type Reference struct {
	Ref string `json:"$ref" validate:"required,uri"`
}

func resolve(r Reference, v interface{}) any {
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

type Schema struct {
	Title       string      `json:"title,omitempty"`
	Description RichText    `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`

	MultipleOf           int                                    `json:"multipleOf,omitempty" validate:"gt=0"`
	Maximum              int                                    `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                                   `json:"exclusiveMaximum,omitempty"`
	Minimum              int                                    `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                                   `json:"exclusiveMinimum,omitempty"`
	MaxLength            int                                    `json:"maxLength,omitempty" validate:"gte=0"`
	MinLength            int                                    `json:"minLength,omitempty" validate:"gte=0"`
	Pattern              string                                 `json:"pattern,omitempty"`
	Items                *ValueOrReferenceOf[Schema]            `json:"items,omitempty" validate:"required_if=Type 6"`
	MaxItems             int                                    `json:"maxItems,omitempty" validate:"gte=0"`
	MinItems             int                                    `json:"minItems,omitempty" validate:"gte=0"`
	UniqueItems          bool                                   `json:"uniqueItems,omitempty"`
	MaxProperties        int                                    `json:"maxProperties,omitempty" validate:"gte=0"`
	MinProperties        int                                    `json:"minProperties,omitempty" validate:"gte=0"`
	Required             []string                               `json:"required,omitempty" validate:"min=1,unique"`
	Properties           map[string]ValueOrReferenceOf[Schema]  `json:"properties,omitempty"`
	AdditionalProperties *Or[bool, *ValueOrReferenceOf[Schema]] `json:"additionalProperties,omitempty"`
	Enum                 []interface{}                          `json:"enum,omitempty"`
	Type                 Type                                   `json:"type,omitempty"`
	AllOf                []ValueOrReferenceOf[Schema]           `json:"allOf,omitempty"`
	AnyOf                []ValueOrReferenceOf[Schema]           `json:"anyOf,omitempty"`
	OneOf                []ValueOrReferenceOf[Schema]           `json:"oneOf,omitempty"`
	Not                  *ValueOrReferenceOf[Schema]            `json:"not,omitempty"`
	Format               Format                                 `json:"format,omitempty"`

	Nullable      bool                   `json:"nullable,omitempty"`
	Discriminator *Discriminator         `json:"discriminator,omitempty"`
	ReadOnly      bool                   `json:"readOnly,omitempty" validate:""`
	WriteOnly     bool                   `json:"writeOnly,omitempty" `
	Xml           *XML                   `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
	Example       interface{}            `json:"example,omitempty"`
	Deprecated    bool                   `json:"deprecated,omitempty"`

	Extensions SpecificationExtension `json:"-"`
}

func (s Schema) Validate() error {
	switch s.Type {
	case IntegerType, NumberType, StringType, BooleanType:
		return validate.Var(s, s.ValidationTag())
	case ObjectType:
		return nil
	case ArrayType:
		return nil
	default:
		panic(fmt.Errorf("invalid type %s", s.Type))
	}
}

// ValidationTag returns a tag style validation to use for a value specified by the Schema.
//
//	tag := schema.ValidationTag()
//	err := validate.Var(v, tag)
func (s Schema) ValidationTag() string {
	const sep = ','
	b := strings.Builder{}
	switch s.Type {
	case IntegerType, NumberType:
		if s.MultipleOf != 0 {
			b.WriteString("multipleOf=")
			b.WriteString(strconv.Itoa(s.MultipleOf))
			b.WriteRune(sep)
		}
		if s.Maximum != 0 {
			b.WriteString("lt")
			if !s.ExclusiveMaximum {
				b.WriteRune('e')
			}
			b.WriteRune('=')
			b.WriteString(strconv.Itoa(s.Maximum))
			b.WriteRune(sep)
		}
		if s.Minimum != 0 {
			b.WriteString("gt")
			if !s.ExclusiveMinimum {
				b.WriteRune('e')
			}
			b.WriteRune('=')
			b.WriteString(strconv.Itoa(s.Minimum))
			b.WriteRune(sep)
		}
		if len(s.Enum) > 0 {
			var format string
			if s.Type == IntegerType {
				format = "%d"
			} else {
				format = "%f"
			}
			b.WriteString("one of=")
			for _, e := range s.Enum {
				b.WriteString(fmt.Sprintf(format, e))
				b.WriteRune(' ')
			}
			b.WriteRune(sep)
		}
		return b.String()
	case StringType:
		if s.MaxLength != 0 {
			b.WriteString("max=")
			b.WriteString(strconv.Itoa(s.MaxLength))
			b.WriteRune(sep)
		}
		if s.MinLength != 0 {
			b.WriteString("min=")
			b.WriteString(strconv.Itoa(s.MinLength))
			b.WriteRune(sep)
		}
		if s.Pattern != "" {
			b.WriteString("regex_ecma=")
			b.WriteString(s.Pattern)
			b.WriteRune(sep)
		}
		if len(s.Enum) > 0 {
			b.WriteString("one of=")
			for _, e := range s.Enum {
				b.WriteString(fmt.Sprintf("%s", e))
				b.WriteRune(' ')
			}
			b.WriteRune(sep)
		}
		return b.String()
	case BooleanType:
		if len(s.Enum) > 0 {
			b.WriteString("one of=")
			for _, e := range s.Enum {
				b.WriteString(fmt.Sprintf("%t", e))
				b.WriteRune(' ')
			}
			b.WriteRune(sep)
		}
		return ""
	case ObjectType:
		return ""
	case ArrayType:
		if s.MaxItems != 0 {
			b.WriteString("max=")
			b.WriteString(strconv.Itoa(s.MaxItems))
			b.WriteRune(sep)
		}
		if s.MinItems != 0 {
			b.WriteString("min=")
			b.WriteString(strconv.Itoa(s.MinItems))
			b.WriteRune(sep)
		}
		if s.UniqueItems {
			b.WriteString("unique")
			b.WriteRune(sep)
		}
		return b.String()
	default:
		panic(fmt.Errorf("invalid type %s", s.Type))
	}
}

type Discriminator struct {
	PropertyName string            `json:"propertyName" validate:"required"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

type XML struct {
	// TODO:
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
	Implicit          *OAuthFlow             `json:"implicit,omitempty"`
	Password          *OAuthFlow             `json:"password,omitempty"`
	ClientCredentials *OAuthFlow             `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow             `json:"authorizationCode,omitempty"`
	Extensions        SpecificationExtension `json:"-"`
}

type OAuthFlow struct {
	AuthorizationUrl string                 `json:"authorizationUrl,omitempty" validate:"url"`
	TokenUrl         string                 `json:"tokenUrl,omitempty" validate:"url"`
	RefreshUrl       string                 `json:"refreshUrl,omitempty" validate:"url"`
	Scopes           map[string]string      `json:"scopes,omitempty"`
	Extensions       SpecificationExtension `json:"-"`
}

type SecurityRequirement map[string][]string
