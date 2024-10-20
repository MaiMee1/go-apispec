package jsonschema

type Format string

const (
	NoFormat Format = ""

	DateTimeFormat            Format = "date-time"             // as defined by the "date-time" ABNF rule in RFC 3339
	DateFormat                Format = "date"                  // as defined by the "full-date" ABNF rule in RFC 3339
	TimeFormat                Format = "time"                  // as defined by the "full-time" ABNF rule in RFC 3339
	DurationFormat            Format = "duration"              // as defined by the "duration" ABNF rule in RFC 3339
	EmailFormat               Format = "email"                 // as defined by the "Mailbox" ABNF rule in RFC 5321
	IdnEmailFormat            Format = "idn-email"             // as defined by the "Mailbox" ABNF rule in RFC 5321 extended by RFC 6531
	Ipv4Format                Format = "ipv4"                  // as defined by the "dotted-quad" ABNF rule in RFC 2673
	Ipv6Format                Format = "ipv6"                  // as defined by RFC 4291
	UriFormat                 Format = "uri"                   // as defined by the "URI" ABNF rule in RFC 3986
	UriReferenceFormat        Format = "uri-reference"         // as defined by the "URI-reference" ABNF rule in RFC 3986
	IriFormat                 Format = "iri"                   // as defined by the "IRI" ABNF rule in RFC 3987
	IriReferenceFormat        Format = "iri-reference"         // as defined by the "IRI-reference" ABNF rule in RFC 3987
	UuidFormat                Format = "uuid"                  // as defined by the "UUID" ABNF rule in RFC 4122
	UriTemplateFormat         Format = "uri-template"          // as defined by the "URI-Template" ABNF rule in RFC 6570
	JsonPointerFormat         Format = "json-pointer"          // as defined by RFC 6901
	RelativeJsonPointerFormat Format = "relative-json-pointer" // as defined by RFC 6901
	RegexFormat               Format = "regex"                 // regular expression in the ECMA-262 dialect
)
