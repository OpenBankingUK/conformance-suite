package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
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
func GenerateTestCases(spec string, baseurl string, ctx *model.Context, endpoints []discovery.ModelEndpoint, manifestPath string) ([]model.TestCase, error) {
	logger := logrus.WithFields(logrus.Fields{
		"function": "GenerateTestCases",
	})

	specType, err := GetSpecType(spec)
	if err != nil {
		return nil, errors.New("unknown specification " + spec)
	}
	logrus.Debug("GenerateManifestTestCases for spec type:" + specType)
	scripts, refs, err := loadGenerationResources(specType, manifestPath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error on loadGenerationResources")
		return nil, err
	}
	var filteredScripts Scripts
	if specType == "accounts" { //TODO: Complete so it makes sense for payments
		filteredScripts, err = filterTestsBasedOnDiscoveryEndpoints(scripts, endpoints)
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("error filter scripts based on discovery")
		}
	} else {
		filteredScripts = scripts // normal processing
	}

	ctx.DumpContext("Incoming Ctx")

	tests := []model.TestCase{}
	//for _, script := range scripts.Scripts {
	for _, script := range filteredScripts.Scripts {
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

func loadGenerationResources(specType, manifestPath string) (Scripts, References, error) {
	assertions, err := loadAssertions()
	if err != nil {
		return Scripts{}, References{}, err
	}
	switch specType {
	case "accounts":
		sc, err := loadScripts(manifestPath)
		return sc, assertions, err
	case "payments":
		pay, err := loadScripts(manifestPath)
		return pay, assertions, err
	case "cbpii":
	case "notifications":
	}
	return Scripts{}, References{}, errors.New("loadGenerationResources: invalid spec type")
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
	const schemeHttps = "https://"
	const schemeFile = "file://"

	var scrBytes []byte
	if strings.HasPrefix(strings.ToLower(filename), schemeHttps) {
		return Scripts{}, errors.New("https:// download of scripts not yet supported")
	} else if strings.HasPrefix(strings.ToLower(filename), schemeFile) {
		f := strings.TrimPrefix(filename, schemeFile)
		sb, err := ioutil.ReadFile(f)
		if err != nil && os.IsNotExist(err) {
			sb, err = ioutil.ReadFile(f)
		}
		scrBytes = sb
		if err != nil {
			return Scripts{}, err
		}

	} else {
		return Scripts{}, errors.New("unable to load scripts")
	}

	var m Scripts
	err := json.Unmarshal(scrBytes, &m)
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

func filterTestsBasedOnDiscoveryEndpoints(scripts Scripts, endpoints []discovery.ModelEndpoint) (Scripts, error) {
	lookupMap := make(map[string]bool)
	filteredScripts := []Script{}

	for _, ep := range endpoints {
		for _, regpath := range accountsRegex {
			matched, err := regexp.MatchString(regpath.Regex, ep.Path)
			if err != nil {
				continue
			}
			if matched {
				lookupMap[regpath.Regex] = true
				logrus.Tracef("endpoint %40.40s matched by regex %42.42s: %s", ep.Path, regpath.Regex, regpath.Name)
			}
		}
	}

	for k := range lookupMap {
		for i, scr := range scripts.Scripts {
			stripped := strings.Replace(scr.URI, "$", "", -1) // only works with a single character
			if strings.Contains(stripped, "foobar") {         //exceptions
				nofoobar := strings.Replace(stripped, "/foobar", "", -1) // only works with a single character
				matched, err := regexp.MatchString(k, nofoobar)
				if err != nil {
					continue
				}
				if matched {
					if !contains(filteredScripts, scripts.Scripts[i]) {
						logrus.Tracef("endpoint %40.40s matched by regex %42.42s", scr.URI, k)
						filteredScripts = append(filteredScripts, scripts.Scripts[i])
					}
				}

				if scr.URI == "/foobar" {
					if !contains(filteredScripts, scripts.Scripts[i]) {
						filteredScripts = append(filteredScripts, scripts.Scripts[i])
					}
					continue
				}
			}

			matched, err := regexp.MatchString(k, stripped)
			if err != nil {
				continue
			}
			if matched {
				if !contains(filteredScripts, scripts.Scripts[i]) {
					logrus.Tracef("endpoint %40.40s matched by regex %42.42s", scr.URI, k)
					filteredScripts = append(filteredScripts, scripts.Scripts[i])
				}
			}
		}
	}
	resultscripts := Scripts{Scripts: filteredScripts}
	sort.Slice(resultscripts.Scripts, func(i, j int) bool { return resultscripts.Scripts[i].ID < resultscripts.Scripts[j].ID })

	return resultscripts, nil
}

func contains(s []Script, e Script) bool {
	for _, a := range s {
		if a.ID == e.ID {
			return true
		}
	}
	return false
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
}

var subPathx = "[a-zA-Z0-9_{}-]+" // url sub path regex

type pathRegex struct {
	Regex string
	Name  string
}

var accountsRegex = []pathRegex{
	{"^/accounts$", "Get Accounts"},
	{"^/accounts/" + subPathx + "$", "Get Accounts Resource"},
	{"^/accounts/" + subPathx + "/balances$", "Get Balances Resource"},
	{"^/accounts/" + subPathx + "/beneficiaries$", "Get Beneficiaries Resource"},
	{"^/accounts/" + subPathx + "/direct-debits$", "Get Direct Debits Resource"},
	{"^/accounts/" + subPathx + "/offers$", "Get Offers Resource"},
	{"^/accounts/" + subPathx + "/party$", "Get Party Rsource"},
	{"^/accounts/" + subPathx + "/product$", "Get Product Resource"},
	{"^/accounts/" + subPathx + "/scheduled-payments$", "Get Schedulated Payment resource"},
	{"^/accounts/" + subPathx + "/standing-orders$", "Get Standing Orders resource"},
	{"^/accounts/" + subPathx + "/statements$", "Get Statements Resource"},
	{"^/accounts/" + subPathx + "/statements/" + subPathx + "/file$", "Get statement files resource"},
	{"^/accounts/" + subPathx + "/statements/" + subPathx + "/transactions$", "Get statement transactions resource"},
	{"^/accounts/" + subPathx + "/transactions$", "Get transactions resource"},
	{"^/balances$", "Get Balances"},
	{"^/beneficiaries$", "Get Beneficiaries"},
	{"^/direct-debits$", "Get directory debits"},
	{"^/offers$", "Get Offers"},
	{"^/party$", "Get party"},
	{"^/products$", "Get Products"},
	{"^/scheduled-payments$", "Get Payments"},
	{"^/standing-orders$", "Get Orders"},
	{"^/statements$", "Get Statements"},
	{"^/transactions$", "Get Transactions"},
}
