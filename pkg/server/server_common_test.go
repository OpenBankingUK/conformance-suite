package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func expectedJsonHeaders() http.Header {
	return http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/json; charset=UTF-8"},
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
