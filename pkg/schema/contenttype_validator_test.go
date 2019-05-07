package schema

import (
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestContentTypeValidator_Validate(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/confirmation-funds-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)
	validator := newContentTypeValidator(finder)
	body := strings.NewReader(getAccountsResponse)

	var testCases = []struct {
		name                string
		responseContentType string
		failures            []Failure
	}{
		{
			name:                "expected usage",
			responseContentType: "application/json;charset=utf-8",
			failures:            nil,
		},
		{
			name:                "expected usage uppercase param value",
			responseContentType: "application/JSON;Charset=UTF-8",
			failures:            nil,
		},
		{
			name:                "expected usage quoted param value",
			responseContentType: `application/json;charset="utf-8"`,
			failures:            nil,
		},
		{
			name:                "expected usage extra space between media and params",
			responseContentType: `application/json; charset="utf-8"`,
			failures:            nil,
		},
		{
			name:                "wrong media type",
			responseContentType: "text/html",
			failures:            []Failure{{Message: "Content-Type Error: Should produce 'application/json', but got: 'text/html'"}},
		},
		{
			name:                "wrong param expected",
			responseContentType: "application/json;charset=klingon",
			failures:            []Failure{{Message: "Content-Type Error: Should produce params 'charset=utf-8', but got: 'charset=klingon'"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := &http.Header{}
			header.Add("Content-type", tc.responseContentType)
			r := Response{
				Method:     "POST",
				Path:       "/funds-confirmation-consents",
				StatusCode: http.StatusOK,
				Body:       body,
				Header:     *header,
			}

			failures, err := validator.Validate(r)

			assert.NoError(t, err)
			assert.Equal(t, tc.failures, failures)
		})
	}
}
