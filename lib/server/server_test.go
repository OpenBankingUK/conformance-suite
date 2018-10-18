package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/google/uuid"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/lib/server"
)

func TestServer(t *testing.T) {
	RunSpecs(t, "Server Suite")
}

var _ bool = Describe("Server", func() {
	var (
		the_server *echo.Echo
	)

	BeforeEach(func() {
		the_server = server.NewServer()
	})

	It("is not nil", func() {
		assert.NotNil(GinkgoT(), the_server)
	})

	Describe("GET /api/health", func() {
		Context("when successful", func() {
			It("returns OK", func() {
				req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
				rec := httptest.NewRecorder()

				the_server.ServeHTTP(rec, req) // Make the request

				assert.Equal(GinkgoT(), http.StatusOK, rec.Code)
				assert.Equal(GinkgoT(), "OK", rec.Body.String())
			})
		})
	})

	Describe("POST /api/validation-runs", func() {
		Context("when successful", func() {
			It("returns validation run ID in JSON and nil error", func() {
				req := httptest.NewRequest(http.MethodPost, "/api/validation-runs", nil)
				rec := httptest.NewRecorder()

				the_server.ServeHTTP(rec, req)
				assert.Equal(GinkgoT(), http.StatusAccepted, rec.Code)
				assert.Equal(GinkgoT(), echo.MIMEApplicationJSONCharsetUTF8, rec.HeaderMap.Get(echo.HeaderContentType))

				assert.NotNil(GinkgoT(), rec.Body)
				var body server.ValidationRunsResponse
				json.Unmarshal(rec.Body.Bytes(), &body)

				id, err := uuid.Parse(body.ID)
				assert.NoError(GinkgoT(), err)

				assert.Equal(GinkgoT(), id.String(), body.ID)
			})
		})
	})
	Describe("GET /api/validation-runs/${id}", func() {
		Context("when succesful", func() {
			It("returns OK ", func() {
				id := "c243d5b6-32f0-45ce-a516-1fc6bb6c3c9a"
				req := httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("/api/validation-runs/%s", id),
					nil,
				)
				// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				the_server.ServeHTTP(rec, req)

				assert.Equal(GinkgoT(), http.StatusOK, rec.Code)
				assert.Equal(GinkgoT(), echo.MIMEApplicationJSONCharsetUTF8, rec.HeaderMap.Get(echo.HeaderContentType))

				assert.NotNil(GinkgoT(), rec.Body)
				var body server.ValidationRunsIDResponse
				json.Unmarshal(rec.Body.Bytes(), &body)

				assert.Equal(GinkgoT(), id, body.Status)
			})
		})
	})
})
