package manifest

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTestCases(t *testing.T) {
	tests, err := GenerateTestCases("Payment Initiation API Specification", "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)

	perms, err := GetPermissions(tests)
	assert.Nil(t, err)
	m := make(map[string]string, 0)
	for _, v := range perms {
		fmt.Printf("perms: %s %-50.50s %s\n", v.ID, v.Path, v.Permissions)
		m[v.Path] = v.ID
	}
	for k := range m {
		fmt.Println(k)
	}

	for _, v := range tests {
		dumpJSON(v)
	}
}

func TestPaymentPermissionsCases(t *testing.T) {
	tests, err := GenerateTestCases("Payment Initiation API Specification", "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)
	fmt.Printf("we have %d tests\n", len(tests))
	requiredTokens, err := getPaymentPermissions(tests)
	requiredTokens, err = updateTokensFromConsent(requiredTokens, tests)
	assert.Nil(t, err)
	for _, v := range requiredTokens {
		fmt.Printf("perms: %s\n", v.IDs)
	}
	updateTestAuthenticationFromToken(tests, requiredTokens)

	fmt.Println("where are my tests?")
	for x, v := range tests {
		if x > 15 {
			break
		}
		dumpJSON(v)
	}
}

func TestDataReferencesAndDump(t *testing.T) {
	data, err := loadAssert()
	assert.Nil(t, err)

	for k, v := range data.References {
		body := jsonString(v.Body)
		l := len(body)
		if l > 0 {
			v.BodyData = body
			v.Body = ""
			data.References[k] = v
		}
	}
}

func loadAssert() (References, error) {
	refs, err := loadReferences("../../manifests/data.json")

	if err != nil {
		fmt.Println("what the hell is going on " + err.Error())
		refs, err = loadReferences("manifesxts/data.json")
		if err != nil {
			fmt.Println("what the hell is going on " + err.Error())
			return References{}, err
		}
	}

	for k, v := range refs.References { // read in data references with body payloads
		body := jsonString(v.Body)
		l := len(body)
		if l > 0 {
			v.BodyData = body
			v.Body = ""
			refs.References[k] = v
		}
	}
	dumpJSON(refs)
	return refs, err
}
