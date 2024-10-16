package security

import (
	"github.com/MaiMee1/go-apispec/oas/iana"
	"github.com/MaiMee1/go-apispec/oas/v3"
)

type Option interface {
	apply(*oas.SecurityScheme)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(scheme *oas.SecurityScheme)

func (f optionFunc) apply(o *oas.SecurityScheme) {
	f(o)
}

func WithDescription(description oas.RichText) Option {
	return optionFunc(func(security *oas.SecurityScheme) {
		security.Description = description
	})
}

func WithInQuery() Option {
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.ApiKeyScheme {
			panic("type must be apiKey")
		}
		security.In = oas.QueryLocation
	})
}

func WithInCookie() Option {
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.ApiKeyScheme {
			panic("type must be apiKey")
		}
		security.In = oas.CookieLocation
	})
}

func WithBasic() Option {
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.HttpScheme {
			panic("type must be http")
		}
		security.Scheme = iana.BasicScheme
	})
}

func WithBearer(format string) Option {
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.HttpScheme {
			panic("type must be http")
		}
		security.Scheme = iana.BearerScheme
		security.BearerFormat = format
	})
}

func WithImplicitFlow(authorizationUrl, refreshUrl string, scopeAndDescriptions ...string) Option {
	if len(scopeAndDescriptions)%2 != 0 {
		panic("scopeAndDescriptions must have an even number")
	}
	var flow oas.ImplicitOAuthFlow

	flow.AuthorizationUrl = authorizationUrl
	flow.RefreshUrl = refreshUrl
	flow.Scopes = make(map[string]string)
	for i := 0; i < len(scopeAndDescriptions)/2; i++ {
		key := scopeAndDescriptions[i*2]
		value := scopeAndDescriptions[i*2+1]
		flow.Scopes[key] = value
	}
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.OAuth2Scheme {
			panic("type must be oauth2")
		}
		security.Flows.Implicit = &flow
	})
}

func WithResourceOwnerPasswordFlow(tokenUrl, refreshUrl string, scopeAndDescriptions ...string) Option {
	if len(scopeAndDescriptions)%2 != 0 {
		panic("scopeAndDescriptions must have an even number")
	}
	var flow oas.PasswordOAuthFlow

	flow.TokenUrl = tokenUrl
	flow.RefreshUrl = refreshUrl
	flow.Scopes = make(map[string]string)
	for i := 0; i < len(scopeAndDescriptions)/2; i++ {
		key := scopeAndDescriptions[i*2]
		value := scopeAndDescriptions[i*2+1]
		flow.Scopes[key] = value
	}
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.OAuth2Scheme {
			panic("type must be oauth2")
		}
		security.Flows.Password = &flow
	})
}

func WithClientCredentialsFlow(tokenUrl, refreshUrl string, scopeAndDescriptions ...string) Option {
	if len(scopeAndDescriptions)%2 != 0 {
		panic("scopeAndDescriptions must have an even number")
	}
	var flow oas.ClientCredentialsOAuthFlow

	flow.TokenUrl = tokenUrl
	flow.RefreshUrl = refreshUrl
	flow.Scopes = make(map[string]string)
	for i := 0; i < len(scopeAndDescriptions)/2; i++ {
		key := scopeAndDescriptions[i*2]
		value := scopeAndDescriptions[i*2+1]
		flow.Scopes[key] = value
	}
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.OAuth2Scheme {
			panic("type must be oauth2")
		}
		security.Flows.ClientCredentials = &flow
	})
}

func WithAuthorizationCodeFlow(authorizationUrl, tokenUrl, refreshUrl string, scopeAndDescriptions ...string) Option {
	if len(scopeAndDescriptions)%2 != 0 {
		panic("scopeAndDescriptions must have an even number")
	}
	var flow oas.AuthorizationCodeOAuthFlow

	flow.AuthorizationUrl = authorizationUrl
	flow.TokenUrl = tokenUrl
	flow.RefreshUrl = refreshUrl
	flow.Scopes = make(map[string]string)
	for i := 0; i < len(scopeAndDescriptions)/2; i++ {
		key := scopeAndDescriptions[i*2]
		value := scopeAndDescriptions[i*2+1]
		flow.Scopes[key] = value
	}
	return optionFunc(func(security *oas.SecurityScheme) {
		if security.Type != oas.OAuth2Scheme {
			panic("type must be oauth2")
		}
		security.Flows.AuthorizationCode = &flow
	})
}
