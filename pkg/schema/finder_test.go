package schema

import (
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinder_Spec(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)

	spec := f.Spec()

	assert.Equal(t, "v3.1.0", spec.Info.Version)
	assert.Equal(t, "Swagger for Account and Transaction API Specification", spec.Info.Description)
}

func TestFinder_Operation(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)

	operation, err := f.Operation("post", "/account-access-consents")

	require.NoError(t, err)
	assert.Len(t, operation.Responses.StatusCodeResponses, 10)
	assert.Equal(t, "Create Account Access Consents", operation.Summary)
}

func TestFinder_Response(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	f := newFinder(doc)

	response, err := f.Response("post", "/account-access-consents", 201)

	require.NoError(t, err)
	assert.NotNil(t, response.Schema)
}
