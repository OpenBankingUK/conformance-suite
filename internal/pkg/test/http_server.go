package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func MockHTTPServer(code int, body string, headers map[string]string, err error) (*httptest.Server, string) {
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
