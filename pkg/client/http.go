package client

import (
	"net/http"
	"time"
)

const (
	DefaultTimeout = time.Duration(time.Second * 25)
)

// NewHTTPClient returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClient(timeout time.Duration) *http.Client {
	return NewHTTPClientWithTransport(timeout, http.DefaultTransport)
}

// NewHTTPClientWithTransport returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClientWithTransport(timeout time.Duration, transport http.RoundTripper) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}
