package server

import (
	"fmt"
	"net/url"

	"github.com/MaiMee1/go-apispec/oas/v3"
)

func New(protocol, hostname, port, pathname string, opts ...Option) (*oas.Server, error) {
	if protocol == "" {
		protocol = "http"
	}
	if hostname == "" {
		hostname = "localhost"
	}
	if port == "" {
		port = "80"
	}
	if pathname == "" {
		pathname = "/"
	}
	uri, err := url.Parse(fmt.Sprintf("%s://%s:%s%s", protocol, hostname, port, pathname))
	if err != nil {
		return nil, err
	}
	server := &oas.Server{
		Url:         oas.UrlTemplate(uri.String()),
		Description: "",
		Variables:   nil,
		Extensions:  nil,
	}
	for _, opt := range opts {
		opt.apply(server)
	}
	return server, nil
}
