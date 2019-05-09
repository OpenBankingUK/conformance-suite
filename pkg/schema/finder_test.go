package schema

import (
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFinder_Spec(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)

	spec := finder.Spec()

	assert.Equal(t, "v3.1.0", spec.Info.Version)
	assert.Equal(t, "Swagger for Account and Transaction API Specification", spec.Info.Description)
}

func TestFinder_Operation(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)

	operation, err := finder.Operation("post", "/account-access-consents")

	require.NoError(t, err)
	assert.Len(t, operation.Responses.StatusCodeResponses, 10)
	assert.Equal(t, "Create Account Access Consents", operation.Summary)
}

func TestFinder_Response(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/account-info-swagger.flattened.json")
	require.NoError(t, err)
	finder := newFinder(doc)

	response, err := finder.Response("post", "/account-access-consents", 201)

	require.NoError(t, err)
	assert.NotNil(t, response.Schema)
}
