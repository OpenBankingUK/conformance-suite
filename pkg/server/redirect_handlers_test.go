package server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTableItem holds variables for each test run
type TestTableItem struct {
	label              	 string
	endpoint           	 string
	httpStatusExpected 	 int
	requestBody        	 string
	responseBodyExpected string
}

func TestRedirectHandlersFragmentOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// Valid code and c_hash combination
	testItemOK := TestTableItem{
		label:              "fragment ok",
		endpoint:           `/api/redirect/fragment/ok`,
		httpStatusExpected: http.StatusOK,
		responseBodyExpected: "null",
		requestBody:`
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// c_hash will not match calculated c_hash due to invalid `code`
	testItemInvalidCode := TestTableItem{
		label:              "fragment invalid code",
		endpoint:           "/api/redirect/fragment/ok",
		httpStatusExpected: http.StatusBadRequest,
		responseBodyExpected: "{\"error\":\"calculated c_hash `Scli9Z1BOsMPd3VjeC_2Kg` does not equal expected c_hash `1lkuHAniRCfVMKlDsJqM3A`\"}",
		requestBody: `
{
    "code": "---invalid code---",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// Note c_hash is manipulated (invalid) in JWT, meaning invalid signature as a result
	// if signature validation is implemented this test shall fail.
	testItemInvalidCHash := TestTableItem{
		label:              "fragment invalid c_hash",
		endpoint:           "/api/redirect/fragment/ok",
		httpStatusExpected: http.StatusBadRequest,
		responseBodyExpected: "{\"error\":\"calculated c_hash `1lkuHAniRCfVMKlDsJqM3A` does not equal expected c_hash `bad-c_hash`\"}",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiJiYWQtY19oYXNoIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQ5MzEzODk3LCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MSwianRpIjoiNTY5NDgxMDEtYjY0NS00NjFkLThhNzUtYTgzNmZiMGVmODAzIn0.GGWSIW6SI2XXZkfIZq6QJ5O-j6BHZbqOGvCWxOaNypU",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []TestTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		// read the file we expect to be served.
		code, body, headers := request(
			http.MethodPost,
			ttItem.endpoint,
			strings.NewReader(ttItem.requestBody),
			server,
		)

		// do assertions.
		require.Equal(ttItem.httpStatusExpected, code, ttItem.label)
		require.Len(headers, 2, ttItem.label)
		require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0], ttItem.label)
		require.NotNil(body, ttItem.label)

		bodyActual := body.String()
		require.Equal(ttItem.responseBodyExpected, bodyActual, ttItem.label)
	}
}

func TestRedirectHandlersQueryOK(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// valid code and c_hash combination
	testItemOK := TestTableItem{
		label:              "query ok",
		endpoint:           `/api/redirect/query/ok`,
		httpStatusExpected: http.StatusOK,
		responseBodyExpected: "null",
		requestBody:`
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// invalid value for `code`
	testItemInvalidCode := TestTableItem{
		label:              "query invalid code",
		endpoint:           `/api/redirect/query/ok`,
		httpStatusExpected: http.StatusBadRequest,
		responseBodyExpected: "{\"error\":\"calculated c_hash `Scli9Z1BOsMPd3VjeC_2Kg` does not equal expected c_hash `1lkuHAniRCfVMKlDsJqM3A`\"}",
		requestBody: `
{
    "code": "---invalid code---",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// invalid value for `c_hash` inside the `id_token` field
	testItemInvalidCHash := TestTableItem{
		label:              "query invalid c_hash",
		endpoint:           `/api/redirect/query/ok`,
		httpStatusExpected: http.StatusBadRequest,
		responseBodyExpected: "{\"error\":\"calculated c_hash `1lkuHAniRCfVMKlDsJqM3A` does not equal expected c_hash `bad-c_hash`\"}",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiJiYWQtY19oYXNoIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQ5MzEzODk3LCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MSwianRpIjoiNTY5NDgxMDEtYjY0NS00NjFkLThhNzUtYTgzNmZiMGVmODAzIn0.GGWSIW6SI2XXZkfIZq6QJ5O-j6BHZbqOGvCWxOaNypU",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []TestTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		// read the file we expect to be served.
		code, body, headers := request(
			http.MethodPost,
			ttItem.endpoint,
			strings.NewReader(ttItem.requestBody),
			server,
		)

		// do assertions.
		require.Equal(ttItem.httpStatusExpected, code, ttItem.label)
		require.Len(headers, 2, ttItem.label)
		require.Equal("application/json; charset=UTF-8", headers["Content-Type"][0], ttItem.label)
		require.NotNil(body, ttItem.label)

		bodyActual := body.String()
		require.JSONEq(ttItem.responseBodyExpected, bodyActual, ttItem.label)
	}
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
