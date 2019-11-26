package schema

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCodeValidator_Validate(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/event-notifications-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)
	validator := newStatusCodeValidator(f)
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

func TestStatusCodeValidator_Validate_UnexpectedStatusCode(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/event-notifications-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)
	validator := newStatusCodeValidator(f)
	body := strings.NewReader(getAccountsResponse)
	r := Response{
		Method:     "POST",
		Path:       "/event-notifications",
		StatusCode: http.StatusOK,
		Body:       body,
	}

	failures, err := validator.Validate(r)

	require.NoError(t, err)
	assert.Len(t, failures, 1)
	expected := []Failure{
		{Message: "server Status 200 not defined by the spec"},
	}
	assert.Equal(t, expected, failures)
}
