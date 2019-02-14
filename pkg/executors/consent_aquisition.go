package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

// InitiationConsentAcquisition - get required tokens
func InitiationConsentAcquisition(consentRequirements []model.SpecConsentRequirements, definition RunDefinition) {
	tokenParameters := make(map[string][]string)

	for _, v := range consentRequirements {
		for _, namedPermission := range v.NamedPermissions {
			codeset := namedPermission.CodeSet
			for _, b := range codeset.CodeSet {
				mystring := string(b)
				set := tokenParameters[namedPermission.Name]
				set = append(set, mystring)
				tokenParameters[namedPermission.Name] = set
			}
		}
	}
	logrus.Debugf("required tokens: %#v", tokenParameters)

	runner := NewConsentAcquisitionRunner(definition, NewBufferedDaemonController())

	for tokenName, permissionList := range tokenParameters {
		logrus.Debugf("processing: token: %s permissionList %v", tokenName, permissionList)
		runner.RunConsentAcquisition(tokenName, buildPermissionString(permissionList), definition.Context, definition.DiscoModel.DiscoveryModel.TokenAcquisition)
	}
}

func buildPermissionString(permissionSlice []string) string {
	var permissions string
	first := true
	for _, perms := range permissionSlice {
		if !first {
			permissions += ","
		} else {
			first = !first
		}
		permissions += "\"" + perms + "\""
	}
	return permissions
}
