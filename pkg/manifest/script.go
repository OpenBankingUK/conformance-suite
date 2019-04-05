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
	Body                string            `json:"body,omitempty"`
	Permissions         []string          `json:"permissions,omitemtpy"`
	PermissionsExcluded []string          `json:"permissions-excluded,omitemtpy"`
	Resource            string            `json:"resource,omitempty"`
	Asserts             []string          `json:"asserts,omitempty"`
	Method              string            `json:"method,omitempty"`
	URI                 string            `json:"uri,omitempty"`
	URIImplemenation    string            `json:"uri_implemenation,omitempty"`
	SchemaCheck         bool              `json:"schemaCheck,omitempty"`
	ContextPut          map[string]string `json:"keepContextOnSuccess,omitempty"`
}

// References - reference collection
type References struct {
	References map[string]Reference `json:"references,omitempty"`
}

// Reference is an item referred to by the test script list an assert of token reqirement
type Reference struct {
	Expect      model.Expect `json:"expect,omitempty"`
	Permissions []string     `json:"permissions,omitempty"`
	Body        interface{}  `json:"body,omitempty"`
	BodyData    string       `json:"bodyData"`
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

// ConsentJobs Holds jobs required only to provide consent so should not show on the ui
type ConsentJobs struct {
	jobs map[string]model.TestCase
}

var cj *ConsentJobs

// GetConsentJobs - makes a structure to hold a list of payment consent jobs than need to be run before the main tests
// and so aren't included in the main test list
func GetConsentJobs() *ConsentJobs {
	if cj == nil {
		jobs := make(map[string]model.TestCase)
		cj = &ConsentJobs{jobs: jobs}
		return cj
	}
	return cj
}

// Add a consent Job
func (cj *ConsentJobs) Add(tc model.TestCase) {
	cj.jobs[tc.ID] = tc
}

// Get a consentJob
func (cj *ConsentJobs) Get(testid string) (model.TestCase, bool) {
	value, exist := cj.jobs[testid]
	return value, exist

}

// add/get ....

// GenerateTestCases examines a manifest file, asserts file and resources definition, then builds the associated test cases
func GenerateTestCases(spec string, baseurl string, ctx *model.Context) ([]model.TestCase, error) {
	logger := logrus.WithFields(logrus.Fields{
		"function": "GenerateTestCases",
	})

	specType, err := GetSpecType(spec)
	if err != nil {
		return nil, errors.New("unknown specification " + spec)
	}
	logrus.Debug("GenerateManifestTestCases for spec type:" + specType)
	scripts, refs, err := loadGenerationResources(specType)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error on loadGenerationResources")
		return nil, err
	}

	// accumulate context data from accountsData ...
	// accountCtx := model.Context{}
	// for k, v := range resources.Ais { //TODO:Get Account info from config file
	// 	accountCtx.PutString(k, v)
	// }

	ctx.DumpContext("Incoming Ctx")

	tests := []model.TestCase{}
	for _, script := range scripts.Scripts {
		localCtx, err := script.processParameters(&refs, ctx)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error on processParameters")
			return nil, err
		}

		consents := []string{}
		tc, err := testCaseBuilder(script, refs.References, localCtx, consents, baseurl)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("Error on testCaseBuilder")
		}

		localCtx.PutContext(ctx)
		showReplacementErrors := true
		tc.ProcessReplacementFields(localCtx, showReplacementErrors)

		tests = append(tests, tc)
	}
	return tests, nil
}

func (s *Script) processParameters(refs *References, resources *model.Context) (*model.Context, error) {
	localCtx := model.Context{}

	for k, value := range s.Parameters {
		if k == "consentId" {
			localCtx.PutString("consentId", value)
			continue
		}
		if strings.Contains(value, "$") {
			str := value[1:]
			//lookup parameter in resources - accountids
			value, _ = resources.GetString(str)
			//lookup parameter in reference data
			ref := refs.References[str]
			val := ref.getValue()
			if len(val) != 0 {
				value = val
			}
			if len(value) == 0 {
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

func (r *Reference) getValue() string {
	return r.BodyData
}

// sets testCase Bearer Header to match requested consent token - for non-consent tests
func updateTestAuthenticationFromToken(tcs []model.TestCase, rts []RequiredTokens) []model.TestCase {
	for _, rt := range rts {
		for x, tc := range tcs {
			for _, id := range rt.IDs {
				if id == tc.ID {
					reqConsent, err := tc.Context.GetString("requestConsent")
					if err == nil && len(reqConsent) > 0 {
						continue
					}

					tc.InjectBearerToken("$" + rt.Name)
					tcs[x] = tc
				}
			}
		}
	}
	return tcs
}

func testCaseBuilder(s Script, refs map[string]Reference, ctx *model.Context, consents []string, baseurl string) (model.TestCase, error) {
	tc := model.MakeTestCase()
	tc.ID = s.ID
	tc.Name = s.Description

	//TODO: make these more configurable - header also get set in buildInput Section
	tc.Input.Headers["x-fapi-financial-id"] = "$x-fapi-financial-id"
	tc.Input.Headers["x-fapi-interaction-id"] = "b4405450-febe-11e8-80a5-0fcebb1574e1"
	tc.Input.Headers["x-fcs-testcase-id"] = tc.ID
	buildInputSection(s, &tc.Input)

	tc.Purpose = s.Detail
	tc.Context = model.Context{}

	tc.Context.PutContext(ctx)
	tc.Context.PutString("x-fapi-financial-id", "$x-fapi-financial-id")
	tc.Context.PutString("baseurl", baseurl)

	for _, a := range s.Asserts {
		ref, exists := refs[a]
		if !exists {
			msg := fmt.Sprintf("assertion %s do not exist in reference data", a)
			logrus.Error(msg)
			return tc, errors.New(msg)
		}
		clone := ref.Expect.Clone()
		if ref.Expect.StatusCode != 0 {
			tc.Expect.StatusCode = clone.StatusCode
		}
		tc.Expect.Matches = append(tc.Expect.Matches, clone.Matches...)
		tc.Expect.SchemaValidation = s.SchemaCheck

	}

	// Handled PutContext parameters
	putMatches := processPutContext(&s)
	if len(putMatches) > 0 {
		tc.Expect.ContextPut.Matches = putMatches
	}

	ctx.PutContext(&tc.Context)
	tc.ProcessReplacementFields(ctx, false)

	_, exists := tc.Context.GetString("postData")
	if exists == nil {
		tc.Context.Delete("postData") // tidy context as bodydata potentially large
	}

	return tc, nil
}

func processPutContext(s *Script) []model.Match {
	m := []model.Match{}
	name, exists := s.ContextPut["name"]
	if !exists {
		return m
	}
	value, exists := s.ContextPut["value"]
	if !exists {
		return m
	}
	mx := model.Match{ContextName: name, JSON: value}
	m = append(m, mx)
	return m

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
	i.RequestBody = s.Body
}

func loadGenerationResources(specType string) (Scripts, References, error) {
	assertions, err := loadAssertions()
	if err != nil {
		return Scripts{}, References{}, err
	}
	switch specType {
	case "accounts":
		sc, err := loadTransactions31()
		return sc, assertions, err
	case "payments":
		pay, err := loadPayments31()
		return pay, assertions, err
	case "cbpii":
	case "notifications":
	}
	return Scripts{}, References{}, errors.New("loadGenerationResources: invalid spec type")
}

func loadPayments31() (Scripts, error) {
	sc, err := loadScripts("manifests/ob_3.1_payment_fca.json")
	if err != nil {
		sc, err = loadScripts("../../manifests/ob_3.1_payment_fca.json")
	}
	return sc, err
}

func loadTransactions31() (Scripts, error) {
	sc, err := loadScripts("manifests/ob_3.1_accounts_transactions_fca.json")
	if err != nil {
		sc, err = loadScripts("../../manifests/ob_3.1_accounts_transactions_fca.json")
	}
	return sc, err
}

func loadAssertions() (References, error) {
	refs, err := loadReferences("manifests/assertions.json")
	if err != nil {
		refs, err = loadReferences("../../manifests/assertions.json")
		if err != nil {
			return References{}, err
		}
	}
	refs2, err := loadReferences("manifests/data.json")

	if err != nil {
		refs2, err = loadReferences("../../manifests/data.json")
		if err != nil {
			return References{}, err
		}
	}
	for k, v := range refs2.References { // read in data references with body payloads
		body := jsonString(v.Body)
		l := len(body)
		if l > 0 {
			v.BodyData = body
			v.Body = ""
			refs2.References[k] = v
		}
		refs.References[k] = refs2.References[k]
	}

	return refs, err
}

func jsonString(i interface{}) string {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	return string(model)
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

// ScriptPermission -
type ScriptPermission struct {
	ID          string
	Permissions []string
	Path        string
}

// GetPermissions -
func getAccountPermissions(tests []model.TestCase) ([]ScriptPermission, error) {
	permCollector := []ScriptPermission{}

	for _, test := range tests {
		ctx := test.Context
		perms, err := ctx.GetStringSlice("permissions")
		if err != nil {
			continue
		}

		sp := ScriptPermission{ID: test.ID, Permissions: perms, Path: test.Input.Method + " " + test.Input.Endpoint}
		permCollector = append(permCollector, sp)
	}

	return permCollector, nil
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
}
