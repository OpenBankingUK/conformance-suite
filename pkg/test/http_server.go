package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// HTTPServer creates a server with a provided response status code and body
func HTTPServer(code int, body string, headers map[string]string) (*httptest.Server, string) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addHeaders(w, headers)
		w.WriteHeader(code)
		fmt.Fprint(w, body)
	}))

	return server, server.URL
}

func addHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}
}