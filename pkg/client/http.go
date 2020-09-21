package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	DefaultTimeout = time.Second * 25
)

// NewHTTPClient returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClient(timeout time.Duration) *http.Client {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return NewHTTPClientWithTransport(timeout, transport)
}

// NewHTTPClientWithTransport returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClientWithTransport(timeout time.Duration, transport *http.Transport) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}
