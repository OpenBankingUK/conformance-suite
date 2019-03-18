package manifest

import (
	"encoding/json"
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestReadScriptsFile(t *testing.T) {
	// read scripts file
	tp, err := loadTestPlan("testdata/get-accounts.json")
	assert.Nil(t, err)
	dumpJSON(tp)
	introspect(&tp)
}

func TestLoadStuff(t *testing.T) {
	sc, err := loadScripts("testdata/ob31_testscript.json")
	assert.Nil(t, err)
	dumpJSON(sc)
	as, err := loadReferences("testdata/assertions.json")
	assert.Nil(t, err)
	dumpJSON(as)
	tp := TestPlan{Scripts: sc, References: as}
	introspect(&tp)
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
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

	perms, err := getPermissions(tests)
	assert.Nil(t, err)
	m := make(map[string]string, 0)
	for _, v := range perms {
		fmt.Printf("perms: %#v\n", v)
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

type ScriptPermission struct {
	ID          string
	Permissions []string
	Path        string
}

func getPermissions(tests []model.TestCase) ([]ScriptPermission, error) {
	permCollector := []ScriptPermission{}

	for _, test := range tests {
		ctx := test.Context
		// for k, v := range ctx {
		// 	fmt.Printf("[Context] %s:%v\n", k, v)
		// }

		perms, err := ctx.GetStringSlice("permissions")
		if err != nil {
			return nil, err
		}

		sp := ScriptPermission{ID: test.ID, Permissions: perms, Path: test.Input.Method + " " + test.Input.Endpoint}
		permCollector = append(permCollector, sp)
	}

	return permCollector, nil
}
