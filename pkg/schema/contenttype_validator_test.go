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
	doc, err := loads.Spec("spec/v3.1.0/event-notifications-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)
	validator := newContentTypeValidator(finder)
	body := strings.NewReader(getAccountsResponse)
	r := Response{
		Method:     "POST",
		Path:       "/event-notifications",
		StatusCode: http.StatusAccepted,
		Body:       body,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 0)
}

func TestContentTypeValidator_Validate_WrongContentType(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/confirmation-funds-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)
	validator := newContentTypeValidator(finder)
	body := strings.NewReader(getAccountsResponse)
	header := &http.Header{}
	header.Add("Content-type", "application/klingon")
	r := Response{
		Method:     "POST",
		Path:       "/funds-confirmation-consents",
		StatusCode: http.StatusOK,
		Body:       body,
		Header:     *header,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 1)
	expected := []Failure{
		{Message: "Content-Type Error: Should produce 'application/json; charset=utf-8', but got: 'application/klingon'"},
	}
	assert.Equal(t, expected, failures)
}
