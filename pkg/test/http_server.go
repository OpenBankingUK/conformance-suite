package test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// HTTPServer creates a server with a provided response status code and body
func HTTPServer(code int, body string, headers map[string]string) (*httptest.Server, string) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addHeaders(w, headers)
		w.WriteHeader(code)
		fmt.Fprint(w, body)
	}))

	return srv, srv.URL
}

// HTTPSServer creates a TLS server with a provided response status code and body
func HTTPSServer(tlsConfig *tls.Config, code int, body string, headers map[string]string) (*httptest.Server, string) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addHeaders(w, headers)
		w.WriteHeader(code)
		fmt.Fprint(w, body)
	}))
	srv.TLS = tlsConfig
	srv.StartTLS()

	return srv, srv.URL
}

func addHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}
}
