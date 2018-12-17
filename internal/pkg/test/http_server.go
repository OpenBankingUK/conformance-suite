package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func MockHTTPServer(code int, body string, headerKey string, headerValue string, err error) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerKey, headerValue)
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


