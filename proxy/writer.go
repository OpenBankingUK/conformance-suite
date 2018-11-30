package proxy

import (
	"bytes"
	"net/http"
)

// Response - generic response interface
type Response interface {
	Status() int
	Header() http.Header
	Body() []byte
}

// WriterRecorder - struct to access call responses
type WriterRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

// WriteHeader - sorts out the response header status
func (w *WriterRecorder) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

// Write - the body of the response
func (w *WriterRecorder) Write(body []byte) (n int, err error) {
	if n, err := w.body.Write(body); err != nil {
		return n, err
	}

	return w.ResponseWriter.Write(body)
}

// Body - return the response body
func (w *WriterRecorder) Body() []byte {
	return w.body.Bytes()
}

// Status - returns the response status
func (w *WriterRecorder) Status() int {
	if w.status == 0 {
		return 200
	}
	return w.status
}
