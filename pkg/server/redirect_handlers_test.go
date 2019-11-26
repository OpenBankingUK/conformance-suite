package server

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"
)

type testTableItem struct {
	label                string
	endpoint             string
	httpStatusExpected   int
	requestBody          string
	responseBodyExpected string
}

func TestRedirectHandlersFragmentOK(t *testing.T) {
	require := test.NewRequire(t)

	journey := &MockJourney{}
	journey.On("CollectToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	server := NewServer(journey, nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// Valid code and c_hash combination
	testItemOK := testTableItem{
		label:                "fragment_ok",
		endpoint:             `/api/redirect/fragment/ok`,
		httpStatusExpected:   http.StatusOK,
		responseBodyExpected: "null",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// c_hash will not match calculated c_hash due to invalid `code`
	testItemInvalidCode := testTableItem{
		label:                "fragment_invalid_code",
		endpoint:             "/api/redirect/fragment/ok",
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
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
	testItemInvalidCHash := testTableItem{
		label:                "fragment_invalid_c_hash",
		endpoint:             "/api/redirect/fragment/ok",
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
    "scope": "openid accounts",
    "id_token": "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2p3dC1pZHAuZXhhbXBsZS5jb20iLCJzdWIiOiJtYWlsdG86bWlrZUBleGFtcGxlLmNvbSIsIm5iZiI6MTU0OTU1NjY0MiwiZXhwIjoxNTQ5NTYwMjQyLCJpYXQiOjE1NDk1NTY2NDIsImp0aSI6ImlkMTIzNDU2IiwidHlwIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9yZWdpc3RlciIsICJjX2hhc2giOiAiaW52YWxpZC1jX2hhc2gifQ.CG8bb_wT7EetLsE3SB8W30K3z_be14ZNXsjWQiklXkqImE-aWFvCqruwh-3aAG5xvaQ_u6T5mj7jaK-ZX93v591FsMPmX1MWyUYNfJp5MsPsWUfzZX69Us5UAqOgZ2zxu662prcE8fVqsL-GB-boVR_0e1SUj4NjKhiEHCNVYe-SGclRSZPvjRf0ymQBacmFz84kLqFVZYTXJFkufXd09FUopnNVK-aK2aCc39TaxzxFwwLaAW_iOtJnzHUtnNdF1OUW5MLTeJYd7hPg0Oq5hPUtz2h2XLVl76ERdJYNWNa1yws4gaWE9PaDgNu-mYbfEZVIHnB28XkB7d6BaCW-GQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []testTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		// TODO(mbana): We are not checking the c_hash at the moment, so skip the tests.
		if ttItem.label == "fragment_invalid_code" || ttItem.label == "fragment_invalid_c_hash" {
			continue
		}

		// read the file we expect to be served.
		code, body, headers := request(
			http.MethodPost,
			ttItem.endpoint,
			strings.NewReader(ttItem.requestBody),
			server,
		)

		// do assertions.
		require.NotNil(body)
		bodyActual := body.String()
		require.JSONEq(ttItem.responseBodyExpected, bodyActual, ttItem.label)

		require.Equal(ttItem.httpStatusExpected, code, ttItem.label)
		require.Equal(expectedJsonHeaders(), headers, ttItem.label)
	}
}

func TestRedirectHandlersQueryOK(t *testing.T) {
	journey := &MockJourney{}
	journey.On("CollectToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	server := NewServer(journey, nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()
	require.NotNil(t, server)

	// valid code and c_hash combination
	testItemOK := testTableItem{
		label:                "query_ok",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusOK,
		responseBodyExpected: "null",
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	// invalid value for `code`
	testItemInvalidCode := testTableItem{
		label:                "query_invalid_code",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
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
	testItemInvalidCHash := testTableItem{
		label:                "query_invalid_c_hash",
		endpoint:             `/api/redirect/query/ok`,
		httpStatusExpected:   http.StatusBadRequest,
		responseBodyExpected: `{"error":"c_hash invalid"}`,
		requestBody: `
{
    "code": "a052c795-742d-415a-843f-8a4939d740d1",
    "scope": "openid accounts",
    "id_token": "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2p3dC1pZHAuZXhhbXBsZS5jb20iLCJzdWIiOiJtYWlsdG86bWlrZUBleGFtcGxlLmNvbSIsIm5iZiI6MTU0OTU1NjY0MiwiZXhwIjoxNTQ5NTYwMjQyLCJpYXQiOjE1NDk1NTY2NDIsImp0aSI6ImlkMTIzNDU2IiwidHlwIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9yZWdpc3RlciIsICJjX2hhc2giOiAiaW52YWxpZC1jX2hhc2gifQ.CG8bb_wT7EetLsE3SB8W30K3z_be14ZNXsjWQiklXkqImE-aWFvCqruwh-3aAG5xvaQ_u6T5mj7jaK-ZX93v591FsMPmX1MWyUYNfJp5MsPsWUfzZX69Us5UAqOgZ2zxu662prcE8fVqsL-GB-boVR_0e1SUj4NjKhiEHCNVYe-SGclRSZPvjRf0ymQBacmFz84kLqFVZYTXJFkufXd09FUopnNVK-aK2aCc39TaxzxFwwLaAW_iOtJnzHUtnNdF1OUW5MLTeJYd7hPg0Oq5hPUtz2h2XLVl76ERdJYNWNa1yws4gaWE9PaDgNu-mYbfEZVIHnB28XkB7d6BaCW-GQ",
    "state": "5a6b0d7832a9fb4f80f1170a"
}
	`,
	}

	ttData := []testTableItem{
		testItemOK,
		testItemInvalidCode,
		testItemInvalidCHash,
	}

	for _, ttItem := range ttData {
		ttItem := ttItem
		t.Run(ttItem.label, func(t *testing.T) {
			// TODO(mbana): We are not checking the c_hash at the moment, so skip the tests.
			if ttItem.label == "query_invalid_c_hash" || ttItem.label == "query_invalid_code" {
				t.Skip()
			}

			require := test.NewRequire(t)

			// read the file we expect to be served.
			code, body, headers := request(
				http.MethodPost,
				ttItem.endpoint,
				strings.NewReader(ttItem.requestBody),
				server,
			)

			// do assertions.
			require.NotNil(body)
			bodyActual := body.String()
			require.JSONEq(ttItem.responseBodyExpected, bodyActual)

			require.Equal(ttItem.httpStatusExpected, code)
			require.Equal(expectedJsonHeaders(), headers)
		})
	}
}

func TestRedirectHandlersError(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
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
	require.NotNil(body)
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)

	require.Equal(http.StatusOK, code)
	require.Equal(expectedJsonHeaders(), headers)
}
