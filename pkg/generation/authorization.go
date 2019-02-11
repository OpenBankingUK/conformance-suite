package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

func updateSpecsBearer(consentRequirements []model.SpecConsentRequirements, specs []SpecificationTestCases) []SpecificationTestCases {
	updatedTests := make([]SpecificationTestCases, len(specs))
	for key, spec := range specs {
		updatedTests[key] = updateSpecBearer(consentRequirements, spec)
	}
	return updatedTests
}

func updateSpecBearer(consentRequirements []model.SpecConsentRequirements, spec SpecificationTestCases) SpecificationTestCases {
	tcs := make([]model.TestCase, len(spec.TestCases))
	for key, tc := range spec.TestCases {
		updatedTc := setHeader(consentRequirements, tc)
		tcs[key] = updatedTc
	}
	updateSpec := spec
	updateSpec.TestCases = tcs
	return updateSpec
}

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
		tc.Input.Headers["authorization"] = "Bearer $" + nameSet
	}
	return tc
}

// authorizationNamedSet find named set in consent requirements for a testId
func authorizationNamedSet(consentRequirements []model.SpecConsentRequirements, testId string) (string, bool) {
	for _, consentRequirement := range consentRequirements {
		for _, namedPermissions := range consentRequirement.NamedPermissions {
			for _, namedTestId := range namedPermissions.CodeSet.TestIds {
				if permissions.TestId(testId) == namedTestId {
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
