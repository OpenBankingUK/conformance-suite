package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"
	"github.com/stretchr/testify/require"
)

// TestServerSwaggerHandlers - paths (e.g., /swagger/account-transaction-v3.0/v3.0/docs) are mapped to handlers.
func TestServerSwaggerHandlers(t *testing.T) {
	require := test.NewRequire(t)

	expectedSwaggerUIPaths := expectedSwaggerUIPaths()

	handlersMap := swaggerHandlers(nullLogger())
	require.Len(handlersMap, len(expectedSwaggerUIPaths))

	for path := range expectedSwaggerUIPaths {
		handler, handlerFound := handlersMap[path]
		require.True(handlerFound, fmt.Sprintf("path=%s has no handler", path))
		require.NotNil(handler, fmt.Sprintf("path=%s has no handler", path))
	}
}

// TestServerSwaggerHandlers - paths (e.g., /swagger/account-transaction-v3.0/v3.0/docs) serve the swagger ui.
func TestServerSwaggerHandlersServesUI(t *testing.T) {
	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(t, server.Shutdown(context.TODO()))
	}()
	require.NotNil(t, server)

	expectedSwaggerUIPaths := expectedSwaggerUIPaths()
	for path, schemaVersion := range expectedSwaggerUIPaths {
		path, schemaVersion := path, schemaVersion
		t.Run(path, func(t *testing.T) {
			require := test.NewRequire(t)

			code, body, headers := request(http.MethodGet, path, nil, server)

			// do assertions
			require.NotNil(body)
			bodyExpected := expectedSwaggerUIHTMLResponse(schemaVersion)
			bodyActual := body.String()
			require.Equal(bodyExpected, bodyActual)

			require.Equal(http.StatusOK, code)
			require.Len(headers, 2)
			require.Equal("text/html; charset=utf-8", headers["Content-Type"][0])
		})
	}
}

// expectedSwaggerUIPaths - returns map like below. The key is the full path at which
// the value (the swagger definition) will be served at in the Echo server. E.g.,
// /swagger/account-transaction-v3.0/v3.0/docs = https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json
func expectedSwaggerUIPaths() map[string]string {
	expectedSwaggerUIPaths := map[string]string{}

	specs := model.Specifications()
	for _, spec := range specs {
		fullPath := fmt.Sprintf("/swagger/%s/%s/docs", spec.Identifier, spec.Version) // /swagger/account-transaction-v3.0/v3.0/docs
		specURL := spec.SchemaVersion.String()

		expectedSwaggerUIPaths[fullPath] = specURL
	}

	return expectedSwaggerUIPaths
}

func expectedSwaggerUIHTMLResponse(specURL string) string {
	spaces := "    "
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>API documentation</title>
%s
    <meta name="viewport" content="width=device-width, initial-scale=1">

%s
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='%s'></redoc>
    <script src="/static/redoc/bundles/redoc.standalone.js"> </script>
  </body>
</html>
`, spaces, spaces, specURL)
}
