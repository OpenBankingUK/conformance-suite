package manifest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadScriptsFile(t *testing.T) {
	// read scripts file
	tp, err := loadTestPlan("testdata/get-accounts.json")
	assert.Nil(t, err)
	dumpJSON(tp)
	introspect(&tp)
}

func introspect(tp *TestPlan) {
	scr := tp.Scripts
	fmt.Println("----------------------------------------")
	for k, v := range scr.Scripts {
		_, _ = k, v
		fmt.Printf("\n%s\n%s\n resource:%20s\n", v.ID, v.Description, v.Resource)
		fmt.Println("Parameters:")
		for key, value := range v.Parameters {
			fmt.Printf("\t%-25s:%s\n", key, value)
		}
		fmt.Println("Asserts")
		for _, value := range v.Asserts {
			fmt.Printf("\t%s\n", value)
		}
	}
	fmt.Println("----------------------------------------")
	ref := tp.References
	for k, v := range ref.References {
		fmt.Printf("%s\n", k)
		if v.Expect.StatusCode != 0 {
			fmt.Printf("expected:\n %#v\n", v.Expect)
		}
		if len(v.Permissions) > 0 {
			fmt.Printf("permission:\n %#v\n", v.Permissions)
		}
	}
}

func loadTestPlan2() (*TestPlan, error) {
	//sc, err := loadScripts("testdata/oneAccountScript.json")
	sc, err := loadScripts("../../manifests/ob_3.1_accounts_transactions_fca.json")
	if err != nil {
		return nil, err
	}
	as, err := loadReferences("../../manifests/assertions.json")
	if err != nil {
		return nil, err
	}
	tp := TestPlan{Scripts: sc, References: as}
	return &tp, nil
}

func TestGenerateTestCases(t *testing.T) {
	tests, err := GenerateTestCases("TestSpec", "http://mybaseurl")
	assert.Nil(t, err)

	perms, err := GetPermissions(tests)
	assert.Nil(t, err)
	m := make(map[string]string, 0)
	for _, v := range perms {
		fmt.Printf("perms: %s %-50.50s %s\n", v.ID, v.Path, v.Permissions)
		m[v.Path] = v.ID
	}
	fmt.Println("----------------------==")
	for k := range m {
		fmt.Println(k)
	}

	for _, v := range tests {
		dumpJSON(v)
	}

}
