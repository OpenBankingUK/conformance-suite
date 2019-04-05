package manifest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

func TestGenerateTestCases(t *testing.T) {
	tests, err := GenerateTestCases(accountSwaggerLocation31, "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)

	perms, err := getAccountPermissions(tests)
	assert.Nil(t, err)
	m := map[string]string{}
	for _, v := range perms {
		t.Logf("perms: %s %-50.50s %s\n", v.ID, v.Path, v.Permissions)
		m[v.Path] = v.ID
	}
	requiredTokens, err := GetRequiredTokensFromTests(tests, "accounts")
	for _, v := range requiredTokens {
		fmt.Println(v)
	}
}

func TestPaymentPermissions(t *testing.T) {
	tests, err := GenerateTestCases(paymentsSwaggerLocation30, "http://mybaseurl", &model.Context{})
	fmt.Printf("we have %d tests\n", len(tests))
	for _, v := range tests {
		dumpJSON(v)
	}

	requiredTokens, err := GetPaymentPermissions(tests)
	assert.Nil(t, err)

	for _, v := range requiredTokens {
		fmt.Printf("%#v\n", v)
	}

	fmt.Println("where are my tests?")
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
