package generation

import (
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/permissions"
)

func setHeader(consentRequirements []model.SpecConsentRequirements, tc model.TestCase) model.TestCase {
	if isAccountAccessConsentEndpoint(tc.Input.Endpoint) {
		// do nothing it's a special case
		return tc
	}
	if tc.Input.Headers == nil {
		tc.Input.Headers = map[string]string{}
	}
	nameSet, ok := authorizationNamedSet(consentRequirements, tc.ID)
	if ok {
		tc.Input.Headers["Authorization"] = "Bearer $" + nameSet
	}
	return tc
}

// authorizationNamedSet find named set in consent requirements for a testId
func authorizationNamedSet(consentRequirements []model.SpecConsentRequirements, testID string) (string, bool) {
	for _, consentRequirement := range consentRequirements {
		for _, namedPermissions := range consentRequirement.NamedPermissions {
			for _, namedTestID := range namedPermissions.CodeSet.TestIds {
				if permissions.TestId(testID) == namedTestID {
					return namedPermissions.Name, true
				}
			}
		}
	}
	return "", false
}

func isAccountAccessConsentEndpoint(path string) bool {
	return path == "/account-access-consents"
}
