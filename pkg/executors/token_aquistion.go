package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// InitiationTokenAcquisition - get required tokens
func InitiationTokenAcquisition(definition RunDefinition, consentRequirements []model.SpecConsentRequirements) {
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
	fmt.Printf("%#v\n", tokenParameters)
}

func stuff(definition RunDefinition, consentReqs []model.SpecConsentRequirements, tokens map[string][]string) {
	comp, err := model.LoadComponent("../../templates/tokenProviderComponent.json")
	_, _ = comp, err

	for k, v := range tokens {
		ctx := duplicateContext(definition.TestCaseRun.GlobalContext)
		params := make(map[string]interface{})
		ctx.PutString("token_names", k)
		ctx.PutStringSlice("permissions", v)
		params["permissions"] = v
	}
}

func duplicateContext(context map[string]string) model.Context {

	return model.Context{}

}
