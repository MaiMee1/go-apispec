package specs

import (
	"fmt"
	"net/http"

	"github.com/MaiMee1/go-apispec/fluent/operation"
	"github.com/MaiMee1/go-apispec/fluent/server"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

type Option interface {
	apply(*API)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*API)

func (f optionFunc) apply(o *API) {
	f(o)
}

func WithTitle(title string) Option {
	return optionFunc(func(api *API) {
		api.document.Info.Title = title
	})
}

func WithDescription(description oas.RichText) Option {
	return optionFunc(func(api *API) {
		api.document.Info.Description = description
	})
}

func WithTOS(url string) Option {
	return optionFunc(func(api *API) {
		api.document.Info.TermsOfService = url
	})
}

func WithContact(name string, url string, email string) Option {
	return optionFunc(func(api *API) {
		if api.document.Info.Contact == nil {
			api.document.Info.Contact = &oas.Contact{}
		}
		api.document.Info.Contact.Name = name
		api.document.Info.Contact.Url = url
		api.document.Info.Contact.Email = email
	})
}

func WithLicense(name string, url string) Option {
	if name == "" {
		panic("name must not be empty")
	}
	return optionFunc(func(api *API) {
		if api.document.Info.License == nil {
			api.document.Info.License = &oas.License{}
		}
		api.document.Info.License.Name = name
		api.document.Info.License.Url = url
	})
}

func WithVersion(version string) Option {
	return optionFunc(func(api *API) {
		api.document.Info.Version = version
	})
}

func WithServer(protocol, hostname string, port uint16, pathname string, opts ...server.Option) Option {
	srv, err := server.New(protocol, hostname, port, pathname, opts...)
	if err != nil {
		panic(err)
	}
	return optionFunc(func(api *API) {
		api.document.Servers = append(api.document.Servers, *srv)
	})
}

func WithOperation(id string, method string, path string, opts ...operation.Option) Option {
	if id == "" {
		panic("id must not be empty")
	}
	if path == "" {
		panic("path must not be empty")
	}

	op := operation.New(id, opts...)
	return optionFunc(func(api *API) {
		item, ok := api.document.Paths[path]
		if !ok {
			item = oas.PathItem{}
		}
		switch method {
		case http.MethodGet:
			item.Get = op
		case http.MethodHead:
			item.Head = op
		case http.MethodPost:
			item.Post = op
		case http.MethodPut:
			item.Put = op
		case http.MethodPatch:
			item.Patch = op
		case http.MethodDelete:
			item.Delete = op
		case http.MethodOptions:
			item.Options = op
		case http.MethodTrace:
			item.Trace = op
		default:
			panic(fmt.Errorf("invalid http method: %s", method))
		}
		api.document.Paths[path] = item
	})
}

func WithWebhook(name string, method string, opts ...operation.Option) Option {
	if name == "" {
		panic("name must not be empty")
	}

	op := operation.New("", opts...)
	return optionFunc(func(api *API) {
		itemOrRef, ok := api.document.Webhooks[name]
		if !ok {
			itemOrRef = oas.ValueOrReferenceOf[oas.PathItem]{
				Value: oas.PathItem{},
			}
		}
		switch method {
		case http.MethodGet:
			itemOrRef.Value.Get = op
		case http.MethodHead:
			itemOrRef.Value.Head = op
		case http.MethodPost:
			itemOrRef.Value.Post = op
		case http.MethodPut:
			itemOrRef.Value.Put = op
		case http.MethodPatch:
			itemOrRef.Value.Patch = op
		case http.MethodDelete:
			itemOrRef.Value.Delete = op
		case http.MethodOptions:
			itemOrRef.Value.Options = op
		case http.MethodTrace:
			itemOrRef.Value.Trace = op
		default:
			panic(fmt.Errorf("invalid http method: %s", method))
		}
		api.document.Webhooks[name] = itemOrRef
	})
}

func WithComponents(keyAndValues ...interface{}) Option {
	if len(keyAndValues)%2 != 0 {
		panic("keyAndValues must have an even number")
	}
	return optionFunc(func(api *API) {
		api.document.Components = oas.Components{
			Schemas:         make(map[string]oas.Schema),
			Responses:       make(map[string]oas.Response),
			Parameters:      make(map[string]oas.Parameter),
			Examples:        make(map[string]oas.Example),
			RequestBodies:   make(map[string]oas.RequestBody),
			Headers:         make(map[string]oas.Header),
			SecuritySchemes: make(map[string]oas.SecurityScheme),
			Links:           make(map[string]oas.Link),
			Callbacks:       make(map[string]oas.Callback),
			PathItems:       make(map[string]oas.PathItem),
			Extensions:      make(oas.SpecificationExtension),
		}
		for i := 0; i < len(keyAndValues)/2; i++ {
			key := keyAndValues[i*2].(string)
			value := keyAndValues[i*2+1]
			switch v := value.(type) {
			case oas.Schema:
				api.document.Components.Schemas[key] = v
			case oas.Response:
				api.document.Components.Responses[key] = v
			case oas.Parameter:
				api.document.Components.Parameters[key] = v
			case oas.Example:
				api.document.Components.Examples[key] = v
			case oas.RequestBody:
				api.document.Components.RequestBodies[key] = v
			case oas.Header:
				api.document.Components.Headers[key] = v
			case oas.SecurityScheme:
				api.document.Components.SecuritySchemes[key] = v
			case oas.Link:
				api.document.Components.Links[key] = v
			case oas.Callback:
				api.document.Components.Callbacks[key] = v
			case oas.PathItem:
				api.document.Components.PathItems[key] = v
			}
		}
	})
}

func WithSecurity(requirements ...oas.SecurityRequirement) Option {
	return optionFunc(func(api *API) {
		api.document.Security = requirements
	})
}

func WithTag(name string, description oas.RichText) Option {
	return optionFunc(func(api *API) {
		api.document.Tags = append(api.document.Tags, oas.Tag{
			Name:         name,
			Description:  description,
			ExternalDocs: nil,
		})
	})
}

func WithExternalDocs(description oas.RichText, url string) Option {
	return optionFunc(func(api *API) {
		api.document.ExternalDocs = &oas.ExternalDocumentation{
			Description: description,
			Url:         url,
		}
	})
}
