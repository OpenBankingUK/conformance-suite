package server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedirectHandlersFragmentOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	bodyExpected := `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`
	// read the file we expect to be served.
	code, body, headers := request(
		http.MethodPost,
		`/api/redirect/fragment/ok`,
		strings.NewReader(bodyExpected),
		server,
	)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}

func TestRedirectHandlersQueryOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	bodyExpected := `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`
	// read the file we expect to be served.
	code, body, headers := request(
		http.MethodPost,
		`/api/redirect/query/ok`,
		strings.NewReader(bodyExpected),
		server,
	)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}

func TestRedirectHandlersError(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	bodyExpected := `
{
    "error_description": "JWT invalid. Expiration time incorrect.",
    "error": "invalid_request",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`
	// read the file we expect to be served.
	code, body, headers := request(
		http.MethodPost,
		`/api/redirect/error`,
		strings.NewReader(bodyExpected),
		server,
	)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 2)
	require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)
}
