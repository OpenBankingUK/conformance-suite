package generation

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestSetHeaderDoesNothingOnAccountAccessConsent(t *testing.T) {
	assert := test.NewAssert(t)

	consentRequirements := []model.SpecConsentRequirements{}
	tc := model.TestCase{}
	tc.Input.Endpoint = "/account-access-consents"

	updatedTc := setHeader(consentRequirements, tc)

	assert.Nil(updatedTc.Input.Headers)
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
