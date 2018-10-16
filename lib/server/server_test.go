package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// Test the NewServer constructor method returns a server.
func TestServer_NewServer(t *testing.T) {
	t.Parallel()

	server := NewServer()
	assert.NotNil(t, server)
}

// Test the GET `/api/health` endpoint.
func TestServer_HealthHandler(t *testing.T) {
	t.Parallel()

	server := NewServer()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	// Make the request
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

// Test that POST `/api/validation-runs` endpoint
func TestServer_ValidationRunsHandler(t *testing.T) {
	t.Parallel()

	server := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/api/validation-runs", nil)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Make the request
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(
		t,
		echo.MIMEApplicationJSONCharsetUTF8,
		rec.HeaderMap.Get(echo.HeaderContentType),
	)

	assert.NotNil(t, rec.Body)
	var body ValidationRunsResponse
	json.Unmarshal(rec.Body.Bytes(), &body)

	id, err := uuid.Parse(body.ID)
	assert.NoError(t, err)

	assert.Equal(t, id.String(), body.ID)
}

// Test the GET `/api/validation-runs/:id` endpoint.
func TestServer_ValidationRunsIDHandler(t *testing.T) {
	t.Parallel()

	server := NewServer()

	id := "c243d5b6-32f0-45ce-a516-1fc6bb6c3c9a"
	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/validation-runs/%s", id),
		nil,
	)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Make the request
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(
		t,
		echo.MIMEApplicationJSONCharsetUTF8,
		rec.HeaderMap.Get(echo.HeaderContentType),
	)

	assert.NotNil(t, rec.Body)
	var body ValidationRunsIDResponse
	json.Unmarshal(rec.Body.Bytes(), &body)

	assert.Equal(t, id, body.Status)
}
