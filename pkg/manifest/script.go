package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// TestPlan species a list of scripts, asserts and other entities required to run a set of test
type TestPlan struct {
	Scripts    Scripts
	References References
}

// Scripts -
type Scripts struct {
	Scripts []Script `json:"scripts,omitempty"`
}

// Script represents a highlevel test definition
type Script struct {
	Description         string            `json:"description,omitempty"`
	Detail              string            `json:"detail,omitempty"`
	ID                  string            `json:"id,omitempty"`
	RefURI              string            `json:"refURI,omitempty"`
	Parameters          map[string]string `json:"parameters,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	Permissions         []string          `json:"permissions,omitemtpy"`
	PermissionsExcluded []string          `json:"permissions-excluded,omitemtpy"`
	Resource            string            `json:"resource,omitempty"`
	Asserts             []string          `json:"asserts,omitempty"`
	Method              string            `json:"method,omitempty"`
	URI                 string            `json:"uri,omitempty"`
	URIImplemenation    string            `json:"uri_implemenation,omitempty"`
	SchemaCheck         bool              `json:"schemaCheck,omitempty"`
}

// References - reference collection
type References struct {
	References map[string]Reference `json:"references,omitempty"`
}

// Reference is an item referred to by the test script list an assert of token reqirement
type Reference struct {
	Expect      model.Expect `json:"expect,omitempty"`
	Permissions []string     `json:"permissions,omitempty"`
}

// AccountData stores account number to be used in the test scripts
type AccountData struct {
	Ais           map[string]string `json:"ais,omitempty"`
	AisConsentIds []string          `json:"ais.ConsetnAccoutId,omitempty"`
	Pis           PisData           `json:"pis,omitempty"`
}

// PisData contains information about PIS accounts required for the test scrips
type PisData struct {
	Currency        string            `json:"Currency,omitempty"`
	DebtorAccount   map[string]string `json:"DebtorAccount,omitempty"`
	MADebtorAccount map[string]string `json:"MADebtorAccount,omitempty"`
}

// GenerateTestCases examines a manifest file, asserts file and resources definition, then builds the associated test cases
func GenerateTestCases(spec string, baseurl string) ([]model.TestCase, error) {
	scripts, refs, resources, err := loadGenerationResources()
	if err != nil {
		return nil, err
	}

	// accumulate context data from accountsData ...
	accountCtx := model.Context{}
	for k, v := range resources.Ais {
		accountCtx.PutString(k, v)
	}

	tests := []model.TestCase{}
	for _, script := range scripts.Scripts {

		localCtx, err := script.processParameters(&refs, &accountCtx)
		if err != nil {
			return nil, err
		}
		consents := []string{}
		tc, _ := testCaseBuilder(script, refs.References, localCtx, consents, baseurl)
		tc.ProcessReplacementFields(localCtx, false)
		tests = append(tests, tc)
	}

	return tests, nil
}

func (s *Script) processParameters(refs *References, resources *model.Context) (*model.Context, error) {
	localCtx := model.Context{}
	var err error

	for k, value := range s.Parameters {
		if strings.Contains(value, "$") {
			str := value[1:]
			value, err = resources.GetString(str)
			if err != nil {
				continue
			}
		}
		switch k {
		case "tokenRequestScope":
			localCtx.PutString("tokenScope", value)
		default:
			localCtx.PutString(k, value)
		}
	}
	if len(s.Permissions) > 0 {
		localCtx.PutStringSlice("permissions", s.Permissions)
	}
	if len(s.PermissionsExcluded) > 0 {
		localCtx.PutStringSlice("permissions-excluded", s.PermissionsExcluded)
	}

	return &localCtx, nil

}

func testCaseBuilder(s Script, refs map[string]Reference, ctx *model.Context, consents []string, baseurl string) (model.TestCase, error) {
	tc := model.MakeTestCase()
	tc.ID = s.ID
	tc.Name = s.Description

	//TODO: make these more configurable - header also get set in buildInput Section
	tc.Input.Headers["x-fapi-financial-id"] = "$fapi_financial_id"
	tc.Input.Headers["x-fapi-interaction-id"] = "b4405450-febe-11e8-80a5-0fcebb1574e1"
	buildInputSection(s, &tc.Input)

	tc.Purpose = s.Detail
	tc.Context = model.Context{}

	tc.Context.PutContext(ctx)
	tc.Context.PutString("baseurl", baseurl)
	tc.InjectBearerToken("$access_token")

	for _, a := range s.Asserts {
		ref, exists := refs[a]
		if !exists {
			msg := fmt.Sprintf("assertion %s do not exist in reference data", a)
			logrus.Error(msg)
			return tc, errors.New(msg)
		}
		tc.Expect = ref.Expect.Clone()
		tc.Expect.SchemaValidation = s.SchemaCheck

	}
	tc.ProcessReplacementFields(ctx, false)
	tc.ProcessReplacementFields(&tc.Context, false)
	return tc, nil
}

func getAccountConsent(refs *References, vx string) []string {
	ref := refs.References[vx]
	return ref.Permissions
}

func buildInputSection(s Script, i *model.Input) {
	i.Method = strings.ToUpper(s.Method)
	i.Endpoint = s.URI
	for k, v := range s.Headers {
		i.Headers[k] = v
	}
}

func loadGenerationResources() (Scripts, References, AccountData, error) {
	return loadScriptFiles()
}

func loadScriptFiles() (Scripts, References, AccountData, error) {
	sc, err := loadScripts("../../manifests/ob_3.1_accounts_transactions_fca.json")
	if err != nil {
		sc, err = loadScripts("manifests/ob_3.1_accounts_transactions_fca.json")
		if err != nil {
			return Scripts{}, References{}, AccountData{}, err
		}
	}

	// sc, err = loadScripts("testdata/oneAccountScript.json")
	// if err != nil {
	// 	sc, err = loadScripts("pkg/manifest/testdata/oneAccountScript.json")
	// 	if err != nil {
	// 		return Scripts{}, References{}, AccountData{}, err
	// 	}
	// }

	refs, err := loadReferences("../../manifests/assertions.json")
	if err != nil {
		refs, err = loadReferences("manifests/assertions.json")
		if err != nil {
			return Scripts{}, References{}, AccountData{}, err
		}
	}

	ad, err := loadAccountData("testdata/resources.json") // temp integration shiv
	if err != nil {
		ad, err = loadAccountData("pkg/manifest/testdata/resources.json")
		if err != nil {
			return Scripts{}, References{}, AccountData{}, err
		}
	}

	return sc, refs, ad, nil
}

func loadAccountData(filename string) (AccountData, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return AccountData{}, err
	}
	var m AccountData
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return AccountData{}, err
	}
	return m, nil
}

func loadScripts(filename string) (Scripts, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return Scripts{}, err
	}
	var m Scripts
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return Scripts{}, err
	}
	return m, nil
}

func loadReferences(filename string) (References, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return References{}, err
	}
	var m References
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return References{}, err
	}
	return m, nil
}

func loadTestPlan(filename string) (TestPlan, error) {
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		return TestPlan{}, err
	}
	var m TestPlan
	err = json.Unmarshal(plan, &m)
	if err != nil {
		return TestPlan{}, err
	}
	return m, nil
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
}
