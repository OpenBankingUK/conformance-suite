package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const (
	DefaultTimeout = time.Second * 25
)

// go 1.3 - can't just clone default transport as in go 1.4
var (
	defaultInsecureTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	defaultNoProxyInsecureTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
)

// NewHTTPClient returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClient(timeout time.Duration) *http.Client {

	// Golang 1.4
	// transport := http.DefaultTransport.(*http.Transport).Clone()
	// transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return NewHTTPClientWithTransport(timeout, defaultInsecureTransport)
}

// NewHTTPClientWithoutProxy returns a HTTP client without Proxy. Use it when you would like to skip proxy setting for some requests (e.g. version checks)
func NewHTTPClientWithoutProxy(timeout time.Duration) *http.Client {
	return NewHTTPClientWithTransport(timeout, defaultNoProxyInsecureTransport)
}

// NewHTTPClientWithTransport returns a more appropriate HTTP client as opposed the default provided by `net/http`
func NewHTTPClientWithTransport(timeout time.Duration, transport *http.Transport) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}
