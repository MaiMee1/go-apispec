package operation

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MaiMee1/go-apispec/oas/jsonschema/draft2020"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

type Option interface {
	apply(*oas.Operation)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*oas.Operation)

func (f optionFunc) apply(o *oas.Operation) {
	f(o)
}

func WithSummary(summary string) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Summary = summary
	})
}

func WithDescription(description oas.RichText) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Description = description
	})
}

func WithTags(tags ...string) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Tags = tags
	})
}

func WithParams(parameters ...oas.Parameter) Option {
	return optionFunc(func(operation *oas.Operation) {
		for _, param := range parameters {
			operation.Parameters = append(operation.Parameters, param)
		}
	})
}

func WithParamReference(ref string) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Parameters = append(operation.Parameters, oas.Parameter{
			ReferenceMixin: draft2020.ReferenceMixin[oas.Parameter]{
				Ref: ref,
			},
		})
	})
}

func WithBody(description oas.RichText, required bool, keyAndValues ...interface{}) Option {
	if len(keyAndValues)%2 != 0 {
		panic("keyAndValues must have an even number")
	}
	body := oas.RequestBody{
		Description: description,
		Required:    required,
		Content:     make(map[string]oas.MediaType),
	}
	for i := 0; i < len(keyAndValues)/2; i++ {
		key := keyAndValues[i*2].(string)
		value := keyAndValues[i*2+1]
		switch v := value.(type) {
		case oas.Schema:
			body.Content[key] = oas.MediaType{
				Schema:   v,
				Example:  nil,
				Examples: nil,
			}
		}
	}
	return optionFunc(func(operation *oas.Operation) {
		operation.RequestBody = &body
	})
}

func WithBodyReference(ref string) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.RequestBody = &oas.RequestBody{
			ReferenceMixin: draft2020.ReferenceMixin[oas.RequestBody]{
				Ref: ref,
			},
		}
	})
}

func WithResponse(code int, description oas.RichText, keyAndValues ...interface{}) Option {
	if len(keyAndValues)%2 != 0 {
		panic("keyAndValues must have an even number")
	}

	status := strconv.Itoa(code)
	if code == 0 {
		status = "default"
	}
	response := oas.Response{
		Description: description,
		Headers:     nil,
		Content:     make(map[string]oas.MediaType),
	}

	for i := 0; i < len(keyAndValues)/2; i++ {
		key := keyAndValues[i*2].(string)
		value := keyAndValues[i*2+1]
		switch v := value.(type) {
		case oas.Schema:
			response.Content[key] = oas.MediaType{
				Schema:   v,
				Example:  nil,
				Examples: nil,
			}
		default:
			panic(fmt.Errorf("value is invalid type %v", v))
		}
	}
	return optionFunc(func(operation *oas.Operation) {
		if operation.Responses == nil {
			operation.Responses = make(oas.Responses)
		}
		operation.Responses[status] = response
	})
}

func WithResponseReference(code int, ref string) Option {
	status := strconv.Itoa(code)
	if code == 0 {
		status = "default"
	}
	return optionFunc(func(operation *oas.Operation) {
		if operation.Responses == nil {
			operation.Responses = make(oas.Responses)
		}
		operation.Responses[status] = oas.Response{
			ReferenceMixin: draft2020.ReferenceMixin[oas.Response]{
				Ref: ref,
			},
		}
	})
}

func WithCallback(name string, method string, url oas.RuntimeExpression, opts ...Option) Option {
	op := New("", opts...)
	return optionFunc(func(operation *oas.Operation) {
		if operation.Callbacks == nil {
			operation.Callbacks = make(map[string]oas.Callback)
		}
		callback, ok := operation.Callbacks[name]
		if !ok {
			callback = oas.Callback{}
		}
		itemOrRef, ok := callback.Value[url]
		if !ok {
			itemOrRef = oas.PathItem{}
		}
		switch method {
		case http.MethodGet:
			itemOrRef.Get = op
		case http.MethodHead:
			itemOrRef.Head = op
		case http.MethodPost:
			itemOrRef.Post = op
		case http.MethodPut:
			itemOrRef.Put = op
		case http.MethodPatch:
			itemOrRef.Patch = op
		case http.MethodDelete:
			itemOrRef.Delete = op
		case http.MethodOptions:
			itemOrRef.Options = op
		case http.MethodTrace:
			itemOrRef.Trace = op
		default:
			panic(fmt.Errorf("invalid http method: %s", method))
		}
		callback.Value[url] = itemOrRef
		operation.Callbacks[name] = callback
	})
}

//func WithCallbackPathItemReference(name string, url oas.RuntimeExpression, ref string) Option {
//	return optionFunc(func(operation *oas.Operation) {
//		if operation.Callbacks == nil {
//			operation.Callbacks = make(map[string]oas.Callback)
//		}
//		callback, ok := operation.Callbacks[name]
//		if !ok {
//			callback = oas.Callback{}
//		}
//		itemOrRef, ok := callback.Value[url]
//		if !ok {
//			itemOrRef = oas.PathItem{
//				ReferenceMixin: draft2020.ReferenceMixin[oas.PathItem]{
//					Ref: ref,
//				},
//			}
//		}
//		callback.Value[url] = itemOrRef
//		operation.Callbacks[name] = callback
//	})
//}

func WithCallbackReference(name string, ref string) Option {
	return optionFunc(func(operation *oas.Operation) {
		if operation.Callbacks == nil {
			operation.Callbacks = make(map[string]oas.Callback)
		}
		callback, ok := operation.Callbacks[name]
		if !ok {
			callback = oas.Callback{
				ReferenceMixin: draft2020.ReferenceMixin[oas.Callback]{
					Ref: ref,
				},
			}
		}
		operation.Callbacks[name] = callback
	})
}

func WithDeprecated() Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Deprecated = true
	})
}

func WithSecurity(requirements ...oas.SecurityRequirement) Option {
	return optionFunc(func(operation *oas.Operation) {
		operation.Security = requirements
	})
}
