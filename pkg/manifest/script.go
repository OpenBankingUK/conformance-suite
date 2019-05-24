package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"github.com/pkg/errors"

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
	APIName             string            `json:"apiName"`
	APIVersion          string            `json:"apiVersion"`
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
	UseCCGToken         bool              `json:"useCCGToken,omitempty"`
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

// GenerateTestCases examines a manifest file, asserts file and resources definition, then builds the associated test cases
func GenerateTestCases(scripts Scripts, spec discovery.ModelAPISpecification, baseurl string, ctx *model.Context, endpoints []discovery.ModelEndpoint, manifestPath string, validator schema.Validator) ([]model.TestCase, Scripts, error) {
	logger := logrus.WithFields(logrus.Fields{
		"function": "GenerateTestCases",
	})

	specType, err := GetSpecType(spec.SchemaVersion)
	if err != nil {
		return nil, Scripts{}, errors.New("unknown specification " + spec.SchemaVersion)

	}
	logrus.Debug("GenerateManifestTestCases for spec type:" + specType)
	scripts, refs, err := LoadGenerationResources(specType, manifestPath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("Error on loadGenerationResources")
		return nil, Scripts{}, err
	}
	var filteredScripts Scripts
	if specType == "accounts" {
		filteredScripts, err = FilterTestsBasedOnDiscoveryEndpoints(scripts, endpoints, accountsRegex)
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("error filter scripts based on accounts discovery")
		}
	} else if specType == "payments" {
		filteredScripts, err = FilterTestsBasedOnDiscoveryEndpoints(scripts, endpoints, paymentsRegex)
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("error filter scripts based on payments discovery")
		}
	} else {
		filteredScripts = scripts // normal processing
	}

	ctx.DumpContext("Incoming Ctx")

	tests := []model.TestCase{}

	for _, script := range filteredScripts.Scripts {
		localCtx, err := script.processParameters(&refs, ctx)
		if err != nil {
			logger.WithError(err).Error("Error on processParameters")
			return nil, Scripts{}, err
		}

		consents := []string{}
		tc, err := testCaseBuilder(script, refs.References, localCtx, consents, baseurl, specType, validator, spec)
		if err != nil {
			logger.WithError(err).Error("Error on testCaseBuilder")
		}

		localCtx.PutContext(ctx)
		showReplacementErrors := true
		tc.ProcessReplacementFields(localCtx, showReplacementErrors)

		tests = append(tests, tc)
	}

	return tests, filteredScripts, nil
}

func (s *Script) processParameters(refs *References, resources *model.Context) (*model.Context, error) {
	localCtx := model.Context{}

	for k, value := range s.Parameters {
		contextValue := value
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
				contextValue = val
			}
			if len(value) == 0 {
				localCtx.PutString(k, contextValue)
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

func testCaseBuilder(s Script, refs map[string]Reference, ctx *model.Context, consents []string, baseurl string, specType string, validator schema.Validator, apiSpec discovery.ModelAPISpecification) (model.TestCase, error) {
	tc := model.MakeTestCase()
	tc.ID = s.ID
	tc.Name = s.Description
	tc.Detail = s.Detail
	tc.RefURI = s.RefURI
	tc.APIName = apiSpec.Name
	tc.APIVersion = apiSpec.Version
	tc.Validator = validator

	//TODO: make these more configurable - header also get set in buildInput Section
	tc.Input.Headers["x-fapi-financial-id"] = "$x-fapi-financial-id"
	// TODO: use automated interaction-id generation - one id per run - injected into context at journey
	tc.Input.Headers["x-fapi-interaction-id"] = "c4405450-febe-11e8-80a5-0fcebb157400"
	tc.Input.Headers["x-fcs-testcase-id"] = tc.ID
	buildInputSection(s, &tc.Input)

	tc.Purpose = s.Detail
	tc.Context = model.Context{}

	tc.Context.PutContext(ctx)
	tc.Context.PutString("x-fapi-financial-id", "$x-fapi-financial-id")
	tc.Context.PutString("baseurl", baseurl)
	if s.UseCCGToken {
		tc.Context.PutString("useCCGToken", "yes") // used for payment posts
	}

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

	if specType == "payments" && tc.Input.Method == "POST" {
		tc.Input.JwsSig = true
		tc.Input.IdempotencyKey = true
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

func LoadGenerationResources(specType, manifestPath string) (Scripts, References, error) {
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
	const schemeHTTPS = "https://"
	const schemeHTTP = "http://"
	const schemeFile = "file://"

	var scriptBytes []byte
	var err error
	if strings.HasPrefix(strings.ToLower(filename), schemeHTTPS) || strings.HasPrefix(strings.ToLower(filename), schemeHTTP) {
		return Scripts{}, errors.New("loadscripts: https:// and http:// download of scripts not implemented")
	} else if strings.HasPrefix(strings.ToLower(filename), schemeFile) {
		fp := strings.TrimPrefix(filename, schemeFile)
		scriptBytes, err = ioutil.ReadFile(fp)
		if err != nil && os.IsNotExist(err) {
			scriptBytes, err = ioutil.ReadFile(fmt.Sprintf("../../%s", fp))
			if err != nil {
				return Scripts{}, errors.Wrap(err, "loadScripts ioutil.ReadFile()")
			}
		} else if err != nil {
			return Scripts{}, errors.Wrap(err, "loadScripts ioutil.ReadFile()")
		}

	} else {
		return Scripts{}, errors.New("loadScripts - no scheme present: (file://)")
	}

	var m Scripts
	err = json.Unmarshal(scriptBytes, &m)
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

// FilterTestsBasedOnDiscoveryEndpoints returns a subset of the first `scripts` parameter, thus filtering `scripts`.
// Filtering is performed by matching (via `regPaths` regex's) the provided `endpoints` against the provided `scripts`.
// The result is: For each path in the collection of scripts returned, there is at least one matching path in the `endpoint`
// list.
func FilterTestsBasedOnDiscoveryEndpoints(scripts Scripts, endpoints []discovery.ModelEndpoint, regPaths []PathRegex) (Scripts, error) {
	lookupMap := make(map[string]bool)
	var filteredScripts []Script

	for _, ep := range endpoints {
		for _, regPath := range regPaths {
			matched, err := regexp.MatchString(regPath.Regex, ep.Path)
			if err != nil {
				continue
			}
			if matched {
				lookupMap[regPath.Regex] = true
				logrus.Tracef("endpoint %40.40s matched by regex %42.42s: %s", ep.Path, regPath.Regex, regPath.Name)
			}
		}
	}

	for k := range lookupMap {
		for _, scr := range scripts.Scripts {
			stripped := strings.Replace(scr.URI, "$", "", -1) // only works with a single character
			if strings.Contains(stripped, "foobar") {         //exceptions
				noFoobar := strings.Replace(stripped, "/foobar", "", -1) // only works with a single character
				matched, err := regexp.MatchString(k, noFoobar)
				if err != nil {
					continue
				}
				if matched {
					if !contains(filteredScripts, scr) {
						logrus.Tracef("endpoint %40.40s matched by regex %42.42s", scr.URI, k)
						filteredScripts = append(filteredScripts, scr)
					}
				}

				if scr.URI == "/foobar" {
					if !contains(filteredScripts, scr) {
						filteredScripts = append(filteredScripts, scr)
					}
					continue
				}
			}

			matched, err := regexp.MatchString(k, stripped)
			if err != nil {
				continue
			}
			if matched {
				if !contains(filteredScripts, scr) {
					logrus.Tracef("endpoint %40.40s matched by regex %42.42s", scr.URI, k)
					filteredScripts = append(filteredScripts, scr)
				}
			}
		}
	}
	result := Scripts{Scripts: filteredScripts}
	sort.Slice(result.Scripts, func(i, j int) bool { return result.Scripts[i].ID < result.Scripts[j].ID })

	return result, nil
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

var subPathx = "[a-zA-Z0-9_{}-]+" // url s	ub path regex

type PathRegex struct {
	Regex  string
	Method string
	Name   string
}

var accountsRegex = []PathRegex{
	{
		Regex: "^/accounts$",
		Name:  "Get Accounts",
	},
	{
		Regex: "^/accounts/" + subPathx + "$",
		Name:  "Get Accounts Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/balances$",
		Name:  "Get Balances Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/beneficiaries$",
		Name:  "Get Beneficiaries Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/direct-debits$",
		Name:  "Get Direct Debits Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/offers$",
		Name:  "Get Offers Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/party$",
		Name:  "Get Party Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/product$",
		Name:  "Get Product Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/scheduled-payments$",
		Name:  "Get Scheduled Payment resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/standing-orders$",
		Name:  "Get Standing Orders resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/statements$",
		Name:  "Get Statements Resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/statements/" + subPathx + "/file$",
		Name:  "Get statement files resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/statements/" + subPathx + "/transactions$",
		Name:  "Get statement transactions resource",
	},
	{
		Regex: "^/accounts/" + subPathx + "/transactions$",
		Name:  "Get transactions resource",
	},
	{
		Regex: "^/balances$",
		Name:  "Get Balances",
	},
	{
		Regex: "^/beneficiaries$",
		Name:  "Get Beneficiaries",
	},
	{
		Regex: "^/direct-debits$",
		Name:  "Get directory debits",
	},
	{
		Regex: "^/offers$",
		Name:  "Get Offers",
	},
	{
		Regex: "^/party$",
		Name:  "Get party",
	},
	{
		Regex: "^/products$",
		Name:  "Get Products",
	},

	{
		Regex: "^/scheduled-payments$",
		Name:  "Get Payments",
	},
	{
		Regex: "^/standing-orders$",
		Name:  "Get Orders",
	},
	{
		Regex: "^/statements$",
		Name:  "Get Statements",
	},
	{
		Regex: "^/transactions$",
		Name:  "Get Transactions",
	},
}

var paymentsRegex = []PathRegex{
	{
		Regex:  "^/domestic-payment-consents$",
		Method: "POST",
		Name:   "Create a domestic payment consent",
	},
	{
		Regex:  "^/domestic-payment-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic payment consent by by consent ID",
	},
	{
		Regex:  "^/domestic-payment-consents/" + subPathx + "/funds-confirmation$",
		Method: "GET",
		Name:   "Get domestic payment consents funds confirmation, by consentID",
	},
	{
		Regex:  "^/domestic-payments$",
		Method: "POST",
		Name:   "Create a domestic payment",
	},
	{
		Regex:  "^/domestic-payments/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic payment by domesticPaymentID",
	},
	{
		Regex:  "^/domestic-scheduled-payment-consents$",
		Method: "POST",
		Name:   "Create a domestic scheduled payment consent",
	},
	{
		Regex:  "^/domestic-scheduled-payment-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic scheduled payment consent by consentID",
	},
	{
		Regex:  "^/domestic-scheduled-payments$",
		Method: "POST",
		Name:   "Create a domestic scheduled payment",
	},
	{
		Regex:  "^/domestic-scheduled-payment/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic scheduled payments by consentID",
	},
	{
		Regex:  "^/domestic-standing-order-consents$",
		Method: "POST",
		Name:   "Create a domestic standing order consent",
	},
	{
		Regex:  "^/domestic-standing-order-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic standing order consent by consentID",
	},
	{
		Regex:  "^/domestic-standing-orders$",
		Method: "POST",
		Name:   "Create a domestic standing order",
	},
	{
		Regex:  "^/domestic-standing-orders/" + subPathx + "$",
		Method: "GET",
		Name:   "Get domestic standing order by domesticStandingOrderID",
	},
	{
		Regex:  "^/international-payment-consents$",
		Method: "POST",
		Name:   "Create an international payment consent",
	},
	{
		Regex:  "^/international-payment-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get international payment consent by consentID",
	},
	{
		Regex:  "^/international-payment-consents/" + subPathx + "/funds-confirmation$",
		Method: "GET",
		Name:   "Get international payment consent funds confirmation by consentID",
	},
	{
		Regex:  "^/international-payments$",
		Method: "POST",
		Name:   "Create an international payment",
	},
	{
		Regex:  "^/international-payments/" + subPathx + "$",
		Method: "GET",
		Name:   "Get international payment by internationalPaymentID",
	},
	{
		Regex:  "^/international-scheduled-payment-consents$",
		Method: "POST",
		Name:   "Create an international scheduled payment consent",
	},
	{
		Regex:  "^/international-scheduled-payment-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get international scheduled payment consents by consentID",
	},
	{
		Regex:  "^/international-scheduled-payments/" + subPathx + "/funds-confirmation$",
		Method: "GET",
		Name:   "Get international scheduled payment funds confirmation by consentID",
	},
	{
		Regex:  "^/international-scheduled-payments$",
		Method: "POST",
		Name:   "Create an international scheduled payment",
	},
	{
		Regex:  "^/international-scheduled-payments/" + subPathx + "$",
		Method: "GET",
		Name:   "Create an international scheduled payment by internationalScheduledPaymentID",
	},
	{
		Regex:  "^/international-standing-order-consents$",
		Method: "POST",
		Name:   "Create international standing order consent",
	},
	{
		Regex:  "^/international-standing-order-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get international standing order consent by consentID",
	},
	{
		Regex:  "^/international-standing-orders$",
		Method: "POST",
		Name:   "Create international standing order",
	},
	{
		Regex:  "^/international-standing-orders/" + subPathx + "$",
		Method: "GET",
		Name:   "Get an international standing order by internationalStandingOrderID",
	},
	{
		Regex:  "^/file-payment-consents$",
		Method: "POST",
		Name:   "Create a file payment consent",
	},
	{
		Regex:  "^/file-payment-consents/" + subPathx + "$",
		Method: "GET",
		Name:   "Get a file payment consent by consentID",
	},
	{
		Regex:  "^/file-payment-consents/" + subPathx + "/file$",
		Method: "POST",
		Name:   "Create a file payment consent file by consentID",
	},
	{
		Regex:  "^/file-payment-consents/" + subPathx + "/file$",
		Method: "GET",
		Name:   "Get a file payment consents file by consentID",
	},
	{
		Regex:  "^/file-payments$",
		Method: "POST",
		Name:   "Create a file payment",
	},
	{
		Regex:  "^/file-payments/" + subPathx + "$",
		Method: "GET",
		Name:   "Get a file payment by filePaymentID",
	},
	{
		Regex:  "^/file-payments/" + subPathx + "/report-file$",
		Method: "GET",
		Name:   "Get a file payment report file by filePaymentID",
	},
}

func timeNowMillis() string {
	tm := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", tm)
}
