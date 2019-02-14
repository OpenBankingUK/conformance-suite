package generation

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

func TestSetHeaderDoesNothingOnAccountAccessConsent(t *testing.T) {
	assert := test.NewAssert(t)

	consentRequirements := []model.SpecConsentRequirements{}
	tc := model.TestCase{}
	tc.Input.Endpoint = "/account-access-consents"

	updatedTc := setHeader(consentRequirements, tc)

	assert.Nil(updatedTc.Input.Headers)
}

func TestSetHeader(t *testing.T) {
	assert := test.NewAssert(t)

	consentRequirements := []model.SpecConsentRequirements{
		{
			NamedPermissions: []model.NamedPermission{
				{
					Name: "permission-set-1",
					CodeSet: permissions.CodeSetResult{
						TestIds: []permissions.TestId{"1"},
					},
				},
			},
		},
	}
	tc := model.TestCase{}

	// Creates headers with authorization
	tc.ID = "1"
	updatedTc := setHeader(consentRequirements, tc)
	assert.NotNil(updatedTc.Input.Headers)
	assert.NotNil(updatedTc.Input.Headers["authorization"])
	assert.Equal("Bearer $permission-set-1", updatedTc.Input.Headers["authorization"])

	// replaces header with authorization
	tc.ID = "1"
	updatedTc.Input.Headers["authorization"] = "123"
	updatedTc = setHeader(consentRequirements, tc)
	assert.NotNil(updatedTc.Input.Headers)
	assert.Equal("Bearer $permission-set-1", updatedTc.Input.Headers["authorization"])
}

func TestAuthorizationNamedSet(t *testing.T) {
	assert := test.NewAssert(t)

	consentRequirements := []model.SpecConsentRequirements{
		{
			NamedPermissions: []model.NamedPermission{
				{
					Name: "permission set 1",
					CodeSet: permissions.CodeSetResult{
						TestIds: []permissions.TestId{"1"},
					},
				},
				{
					Name: "permission set 2",
					CodeSet: permissions.CodeSetResult{
						TestIds: []permissions.TestId{"2", "3"},
					},
				},
			},
		},
	}

	name, found := authorizationNamedSet(consentRequirements, "1")
	assert.True(found)
	assert.Equal("permission set 1", name)

	name, found = authorizationNamedSet(consentRequirements, "3")
	assert.True(found)
	assert.Equal("permission set 2", name)

	name, found = authorizationNamedSet(consentRequirements, "4")
	assert.False(found)
	assert.Equal("", name)
}

func TestIsAccountAccessConsentEndpoint(t *testing.T) {
	assert := test.NewAssert(t)

	assert.True(isAccountAccessConsentEndpoint("/account-access-consents"))
	assert.False(isAccountAccessConsentEndpoint("account-access-consents"))
}