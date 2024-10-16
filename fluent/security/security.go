package security

import "github.com/MaiMee1/go-apispec/oas/v3"

func None() oas.SecurityRequirement {
	return oas.SecurityRequirement{}
}

func Scheme(name string, scopes ...string) oas.SecurityRequirement {
	return oas.SecurityRequirement{
		name: scopes,
	}
}

func NewApiKeyScheme(name string, opts ...Option) oas.SecurityScheme {
	security := &oas.SecurityScheme{
		Type: oas.ApiKeyScheme,
		Name: name,
		In:   oas.HeaderLocation,
	}
	for _, opt := range opts {
		opt.apply(security)
	}
	return *security
}

func NewHttpScheme(opts ...Option) oas.SecurityScheme {
	security := &oas.SecurityScheme{
		Type: oas.HttpScheme,
	}
	for _, opt := range opts {
		opt.apply(security)
	}
	return *security
}

func NewOAuth2Scheme(opts ...Option) oas.SecurityScheme {
	security := &oas.SecurityScheme{
		Type:  oas.OAuth2Scheme,
		Flows: &oas.OAuthFlows{},
	}
	for _, opt := range opts {
		opt.apply(security)
	}
	return *security
}

func NewOpenIdConnectScheme(openIdConnectUrl string, opts ...Option) oas.SecurityScheme {
	security := &oas.SecurityScheme{
		Type:             oas.OpenIdConnectScheme,
		OpenIdConnectUrl: openIdConnectUrl,
	}
	for _, opt := range opts {
		opt.apply(security)
	}
	return *security
}
