package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func expectedJSONHeaders() http.Header {
	return http.Header{
		"Content-Security-Policy": []string{"default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' blob:; worker-src blob:; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self' ws: wss:;"},
		"X-Xss-Protection":        []string{"1; mode=block"},
		"X-Content-Type-Options":  []string{"nosniff"},
		"X-Frame-Options":         []string{"SAMEORIGIN"},
		"Content-Type":            []string{"application/json; charset=UTF-8"},
		"Vary":                    []string{"Accept-Encoding"},
	}
}

const (
	marshalIndentPrefix = ``
	marshalIndentindent = `  `
)

// Generic util function for making test requests.
func request(method, path string, body io.Reader, server *Server) (int, *bytes.Buffer, http.Header) {
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	server.ServeHTTP(rec, req)

	return rec.Code, rec.Body, rec.Header()
}

// nullLogger - create a logger that discards output.
func nullLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger.WithField("app", "test")
}
