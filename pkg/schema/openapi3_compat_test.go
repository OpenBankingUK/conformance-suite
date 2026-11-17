package schema

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compatSpec has what the compatibility shims act on: string formats, a declared
// response header and "allOf" schemas at several depths.
const compatSpec = `{
	"openapi": "3.0.0",
	"info": {"title": "compat", "version": "1"},
	"paths": {"/x": {"get": {"responses": {"200": {
		"description": "ok",
		"headers": {"x-fapi-interaction-id": {"schema": {"type": "string"}}},
		"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Root"}}}
	}}}}},
	"components": {"schemas": {
		"Root": {"type": "object", "properties": {
			"Date": {"type": "string", "format": "date"},
			"DateTime": {"type": "string", "format": "date-time"},
			"Data": {"type": "object", "properties": {
				"Deep": {"type": "object", "properties": {
					"Account": {"allOf": [{"$ref": "#/components/schemas/Account"}, {"description": "must be populated"}]}
				}},
				"Nested": {"allOf": [{"$ref": "#/components/schemas/Holder"}, {"$ref": "#/components/schemas/Zed"}]},
				"Items": {"type": "array", "items": {"allOf": [{"$ref": "#/components/schemas/Account"}, {"description": "item"}]}},
				"Triple": {"allOf": [
					{"type": "object", "required": ["p"]},
					{"type": "object", "required": ["q"]},
					{"type": "object", "required": ["r"]}
				]}
			}}
		}},
		"Account": {"type": "object", "required": ["Identification", "Name"], "properties": {"Identification": {"type": "string"}, "Name": {"type": "string"}}},
		"Holder": {"type": "object", "properties": {"Inner": {"allOf": [{"$ref": "#/components/schemas/Account"}, {"type": "object", "required": ["Extra"]}]}}},
		"Zed": {"type": "object", "required": ["Zed"]}
	}}
}`

func validateCompatResponse(t *testing.T, header http.Header, body string) error {
	t.Helper()

	doc, err := openapi3.NewLoader().LoadFromData([]byte(compatSpec))
	require.NoError(t, err)
	router, err := legacyrouter.NewRouter(doc)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodGet, "http://localhost/x", nil)
	require.NoError(t, err)
	route, pathParams, err := router.FindRoute(req)
	require.NoError(t, err)

	return OpenAPI3Validator{router: router, doc: doc}.validateResponse(validateParams{
		httpReq:    req,
		route:      route,
		pathParams: pathParams,
		statusCode: http.StatusOK,
		header:     header,
		body:       []byte(body),
	})
}

func TestEmptyDeclaredResponseHeaderIsRejected(t *testing.T) {
	header := http.Header{"Content-Type": {"application/json"}, "X-Fapi-Interaction-Id": {""}}

	err := validateCompatResponse(t, header, `{}`)

	require.Error(t, err)
	assert.Contains(t, err.Error(), `response header "x-fapi-interaction-id" doesn't match schema: Value is not nullable`)
	assert.Equal(t, []string{""}, header["X-Fapi-Interaction-Id"], "the caller's header must not be modified")
}

func TestLegacyStringFormatsAreKept(t *testing.T) {
	header := http.Header{"Content-Type": {"application/json"}}

	t.Run("date accepts a zero month and day", func(t *testing.T) {
		assert.NoError(t, validateCompatResponse(t, header, `{"Date":"2021-00-00"}`))
	})
	t.Run("date-time accepts a missing timezone", func(t *testing.T) {
		assert.NoError(t, validateCompatResponse(t, header, `{"DateTime":"2021-01-01T00:00:00"}`))
	})
}

func TestAllOfFailureNamesOnlyTheFirstFailingBranch(t *testing.T) {
	header := http.Header{"Content-Type": {"application/json"}}
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "object nested in objects",
			body: `{"Data":{"Deep":{"Account":{"Name":"n"}}}}`,
			want: `Error at "/Data/Deep/Account": Error at "/Identification": property "Identification" is missing`,
		},
		{
			name: "null where every branch fails",
			body: `{"Data":{"Deep":{"Account":null}}}`,
			want: `Error at "/Data/Deep/Account": Value is not nullable`,
		},
		{
			name: "allOf inside allOf",
			body: `{"Data":{"Nested":{"Inner":{"Name":"n"}}}}`,
			want: `Error at "/Data/Nested": Error at "/Inner": Error at "/Identification": property "Identification" is missing`,
		},
		{
			name: "array item",
			body: `{"Data":{"Items":[{"Identification":"i","Name":"n"},{"Name":"n"}]}}`,
			want: `Error at "/Data/Items/1": Error at "/Identification": property "Identification" is missing`,
		},
		{
			name: "several failing branches",
			body: `{"Data":{"Triple":{}}}`,
			want: `Error at "/Data/Triple": Error at "/p": property "p" is missing`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCompatResponse(t, header, tc.body)

			require.Error(t, err)
			firstLine := strings.SplitN(err.Error(), "\n", 2)[0]
			assert.Equal(t, "response body doesn't match schema #/components/schemas/Root: "+tc.want, firstLine)
		})
	}
}

func TestEmptyHeaderValuesAsNull(t *testing.T) {
	assert.Nil(t, emptyHeaderValuesAsNull(nil))

	in := http.Header{"X-Empty": {""}, "X-Set": {"v"}, "X-None": {}}
	out := emptyHeaderValuesAsNull(in)

	assert.Equal(t, http.Header{"X-Empty": {}, "X-Set": {"v"}, "X-None": {}}, out)
	assert.Equal(t, http.Header{"X-Empty": {""}, "X-Set": {"v"}, "X-None": {}}, in, "the input must not be modified")
}

func TestLegacyAllOfMessageDefersToTheLibrary(t *testing.T) {
	tests := map[string]*openapi3.SchemaError{
		"not an allOf":                         {SchemaField: "type", Origin: errors.New("x")},
		"allOf without a cause":                {SchemaField: "allOf"},
		"allOf whose cause lists no branches":  {SchemaField: "allOf", Origin: errors.New("x")},
		"allOf the library rendered otherwise": {SchemaField: "allOf", Origin: openapi3.MultiError{errors.New("x")}},
	}

	for name, schemaErr := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Empty(t, legacyAllOfMessage(schemaErr))
		})
	}
}
