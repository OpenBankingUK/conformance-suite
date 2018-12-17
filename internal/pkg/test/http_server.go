package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func MockHTTPServer(code int, body string, headers map[string]string, err error) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addHeaders(w, headers)
		w.WriteHeader(code)
		fmt.Fprint(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			urlStr, _ := url.Parse(server.URL)
			return urlStr, err
		},
	}

	httpClient := &http.Client{Transport: transport}
	return server, httpClient
}

func addHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}
}
