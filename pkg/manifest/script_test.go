package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
	"github.com/blang/semver/v4"

	"github.com/stretchr/testify/assert"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
)

func TestGetConsentJobs(t *testing.T) {
	assert := assert.New(t)

	consentJobs := GetConsentJobs()
	assert.NotNil(consentJobs)
	assert.NotNil(consentJobs.jobs)
	assert.Empty(consentJobs.jobs)
}

var (
	consentJobsTcs = map[string]struct {
		testCase model.TestCase
	}{
		"EmptyID": {
			testCase: model.TestCase{},
		},
		"VRPID": {
			testCase: model.TestCase{ID: "OB-301-VRP-10670"},
		},
	}

	generationTests = map[string]struct {
		discoveryPath string
		specType      string
	}{
		"AISP": {
			discoveryPath: "../discovery/templates/ob-v3.1.11-ozone-AISP.json",
			specType:      "accounts",
		},
		"PISP": {
			discoveryPath: "../discovery/templates/ob-v3.1.11-ozone-PISP.json",
			specType:      "payments",
		},
		"CBPII": {
			discoveryPath: "../discovery/templates/ob-v3.1.11-ozone-CBPII.json",
			specType:      "cbpii"},
		"VRP": {
			discoveryPath: "../discovery/templates/ob-v3.1.11-ozone-VRP.json",
			specType:      "vrps",
		},
		"VRP with conditionalProperties": {
			discoveryPath: "../discovery/templates/ob-v3.1.11-ozone-VRP_conditionalProperties.json",
			specType:      "vrps",
		},
	}
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	consentJobs := GetConsentJobs()

	for name, tc := range consentJobsTcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			consentJobs.Add(tc.testCase)

			value, ok := consentJobs.jobs[tc.testCase.ID]
			assert.True(ok)
			assert.Equal(tc.testCase, value)
		})
	}

	assert.Len(consentJobs.jobs, 2)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	consentJobs := GetConsentJobs()

	for name, tc := range consentJobsTcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			consentJobs.Add(tc.testCase)
			value, ok := consentJobs.jobs[tc.testCase.ID]
			assert.True(ok)
			assert.Equal(tc.testCase, value)

			concentJobTc, ok := consentJobs.Get(tc.testCase.ID)
			assert.True(ok)
			assert.Equal(tc.testCase, concentJobTc)
		})
	}

	assert.Len(consentJobs.jobs, 2)
}

func TestLoadGenerationResources(t *testing.T) {
	for name, tc := range generationTests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assert.NoDirExists(tc.discoveryPath)
			discoFile, err := os.ReadFile(tc.discoveryPath)
			assert.NoError(err)
			assert.NotNil(discoFile)

			disco, err := discovery.UnmarshalDiscoveryJSON(string(discoFile))
			assert.NoError(err)
			assert.NotNil(disco)

			apiSpec := disco.DiscoveryModel.DiscoveryItems[0].APISpecification

			schemaVersion := apiSpec.SchemaVersion
			specType, err := GetSpecType(schemaVersion)
			assert.NoError(err)
			assert.Equal(tc.specType, specType)

			manifestPath := apiSpec.Manifest
			apiVersions := DetermineAPIVersions(disco.DiscoveryModel.DiscoveryItems)
			assert.NotEmpty(apiVersions)

			context := &model.Context{}
			context.PutStringSlice("apiversions", apiVersions)

			scripts, references, err := LoadGenerationResources(specType, manifestPath, context)
			assert.NoError(err)
			assert.NotEmpty(scripts)
			assert.NotEmpty(references)
		})
	}

	t.Run("No apiversions in the context", func(t *testing.T) {
		assert := assert.New(t)
		context := &model.Context{}
		var specType string

		scripts, references, err := LoadGenerationResources(specType, "", context)
		assert.EqualError(err, "loadGenerationResources: apiversions - context variable not found")
		assert.Empty(references)
		assert.Empty(scripts)
	})

	t.Run("Invalid specType (notifications)", func(t *testing.T) {
		assert := assert.New(t)
		context := &model.Context{}
		context.PutStringSlice("apiversions", []string{"apiVersions"})
		specType := "notifications"

		scripts, references, err := LoadGenerationResources(specType, "", context)
		assert.EqualError(err, "loadGenerationResources: invalid spec type")
		assert.Empty(references)
		assert.Empty(scripts)
	})

	t.Run("Cannot get specVersion", func(t *testing.T) {
		assert := assert.New(t)
		context := &model.Context{}
		apiVersions := []string{"payments_vX"}
		context.PutStringSlice("apiversions", apiVersions)
		specType := "payments"

		scripts, references, err := LoadGenerationResources(specType, "", context)
		assert.EqualErrorf(err, "loadGenerationResources: cannot get spec version from spec type payments:[payments_vX]", "loadGenerationResources: cannot get spec version from spec type %s:%v", specType, apiVersions)
		assert.Empty(references)
		assert.Empty(scripts)
	})

	t.Run("Cannot get specVersion from specType", func(t *testing.T) {
		assert := assert.New(t)
		context := &model.Context{}
		apiVersions := []string{"payments_vX"}
		context.PutStringSlice("apiversions", apiVersions)
		specType := "payments"

		scripts, references, err := LoadGenerationResources(specType, "", context)
		assert.EqualErrorf(err, "loadGenerationResources: cannot get spec version from spec type payments:[payments_vX]", "loadGenerationResources: cannot get spec version from spec type %s:%v", specType, apiVersions)
		assert.Empty(references)
		assert.Empty(scripts)
	})
}

func TestAddQueryParametersToRequest(t *testing.T) {
	assert := assert.New(t)

	tc := model.MakeTestCase()
	queryParameters := map[string]string{
		"firstPaymentDateTime": "2024-03-17T00:00:00+01:00",
		"symbols":              "[]{}/\\=+-_*&^%$#@!",
	}

	addQueryParametersToRequest(&tc, queryParameters)
	assert.NotEmpty(tc.Input.QueryParameters)

	for k, v := range queryParameters {
		inputQueryParameter := tc.Input.QueryParameters[k]
		assert.Equal(v, inputQueryParameter)
	}
}

func TestBuildTestCase(t *testing.T) {
	tests := map[string]struct {
		script        Script
		references    map[string]Reference
		ctx           *model.Context
		baseurl       string
		specType      string
		validator     schema.Validator
		apiSpec       discovery.ModelAPISpecification
		interactionId string
		err           string
	}{
		"Empty references": {
			script: Script{
				Asserts: []string{"OB3GLOAssertOn404"},
			},
			references: map[string]Reference{},
			ctx:        &model.Context{},
			err:        "assertion OB3GLOAssertOn404 do not exist in reference data",
		},
		"Script's asserts match references": {
			script: Script{
				Asserts: []string{"OB3GLOAssertOn404"},
			},
			references: map[string]Reference{
				"OB3GLOAssertOn404": {Expect: model.Expect{StatusCode: 404}},
			},
			ctx: &model.Context{},
		},
		"Setting CGG token": {
			script: Script{
				Asserts:     []string{"OB3GLOAssertOn404"},
				UseCCGToken: true,
			},
			references: map[string]Reference{
				"OB3GLOAssertOn404": {Expect: model.Expect{StatusCode: 404}},
			},
			ctx: &model.Context{},
		},
		"postData in Context": {
			script: Script{
				Asserts: []string{"OB3GLOAssertOn404"},
			},
			references: map[string]Reference{
				"OB3GLOAssertOn404": {Expect: model.Expect{StatusCode: 404}},
			},
			ctx: &model.Context{"postData": "{\"Data\": \"test\"}"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			tc, err := buildTestCase(tt.script, tt.references, tt.ctx, tt.baseurl, tt.specType, tt.validator, tt.apiSpec, tt.interactionId)

			if tt.err != "" {
				assert.EqualError(err, tt.err)
			} else {
				assert.NoError(err)
			}

			assert.NotEmpty(t, tc)

			useCCGTokenValue, err := tc.Context.GetString("useCCGToken")
			if tt.script.UseCCGToken {
				assert.NoError(err)
				assert.Equal("yes", useCCGTokenValue)
			} else {
				assert.Error(err)
				assert.Empty(useCCGTokenValue)
			}

			if _, present := tt.ctx.Get("postData"); present {
				_, tcPresent := tc.Context.Get("postData")
				assert.False(tcPresent)
			}
		})
	}
}

func TestProcessPutContext(t *testing.T) {
	tests := map[string]struct {
		script          Script
		matches         []model.Match
		isEmptyExpected bool
	}{
		"Empty ContextPut": {
			script:          Script{ContextPut: map[string]string{}},
			isEmptyExpected: true,
		},
		"ContextPut with the name": {
			script:          Script{ContextPut: map[string]string{"name": "name"}},
			isEmptyExpected: true,
		},
		"ContextPut with the name, value": {
			script:          Script{ContextPut: map[string]string{"name": "name", "value": "value"}},
			matches:         []model.Match{{ContextName: "name", JSON: "value"}},
			isEmptyExpected: false,
		},
		"ContextPut with the value": {
			script:          Script{ContextPut: map[string]string{"value": "value"}},
			isEmptyExpected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			matches := processPutContext(&tt.script)

			if tt.isEmptyExpected {
				assert.Empty(matches)
			} else {
				assert.NotEmpty(matches)
				assert.Equal(tt.matches, matches)
			}

		})
	}
}

// TODO: add failing tests
func TestGenerateTestCases(t *testing.T) {
	for name, tt := range generationTests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assert.NoDirExists(tt.discoveryPath)
			discoFile, err := os.ReadFile(tt.discoveryPath)
			assert.NoError(err)
			assert.NotNil(discoFile)

			disco, err := discovery.UnmarshalDiscoveryJSON(string(discoFile))
			assert.NoError(err)
			assert.NotNil(disco)

			apiSpec := disco.DiscoveryModel.DiscoveryItems[0].APISpecification

			manifestPath := apiSpec.Manifest
			apiVersions := DetermineAPIVersions(disco.DiscoveryModel.DiscoveryItems)
			assert.NotEmpty(apiVersions)

			context := &model.Context{}
			context.PutStringSlice("apiversions", apiVersions)

			validator, err := schema.NewOpenAPI3Validator(apiSpec.Name, apiSpec.Version)
			assert.NoError(err)
			assert.NotNil(validator)

			endpoints := disco.DiscoveryModel.DiscoveryItems[0].Endpoints
			assert.NotEmpty(endpoints)

			conditionalAPIProperties, _, err := discovery.GetConditionalProperties(disco)
			assert.NoError(err)

			params := &GenerationParameters{
				Spec:         apiSpec,
				Baseurl:      "http://mybaseurl",
				Ctx:          context,
				Endpoints:    endpoints,
				ManifestPath: manifestPath,
				Validator:    validator,
				Conditional:  conditionalAPIProperties,
			}

			tests, scripts, err := GenerateTestCases(params)
			assert.NoError(err)
			assert.NotEmpty(tests)
			assert.NotEmpty(scripts)
		})
	}
}

// permx related
func TestGetRequiredTokensFromTests(t *testing.T) {
	for name, tc := range generationTests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assert.NoDirExists(tc.discoveryPath)
			discoFile, err := os.ReadFile(tc.discoveryPath)
			assert.NoError(err)
			assert.NotNil(discoFile)

			disco, err := discovery.UnmarshalDiscoveryJSON(string(discoFile))
			assert.NoError(err)
			assert.NotNil(disco)

			apiSpec := disco.DiscoveryModel.DiscoveryItems[0].APISpecification

			manifestPath := apiSpec.Manifest
			apiVersions := DetermineAPIVersions(disco.DiscoveryModel.DiscoveryItems)
			assert.NotEmpty(apiVersions)

			context := &model.Context{}
			context.PutStringSlice("apiversions", apiVersions)

			validator, err := schema.NewOpenAPI3Validator(apiSpec.Name, apiSpec.Version)
			assert.NoError(err)
			assert.NotNil(validator)

			endpoints := disco.DiscoveryModel.DiscoveryItems[0].Endpoints
			assert.NotEmpty(endpoints)

			conditionalAPIProperties, _, err := discovery.GetConditionalProperties(disco)
			assert.NoError(err)

			params := &GenerationParameters{
				Spec:         apiSpec,
				Baseurl:      "http://mybaseurl",
				Ctx:          context,
				Endpoints:    endpoints,
				ManifestPath: manifestPath,
				Validator:    validator,
				Conditional:  conditionalAPIProperties,
			}

			tests, _, err := GenerateTestCases(params)
			assert.NoError(err)
			assert.NotEmpty(tests)

			schemaVersion := apiSpec.SchemaVersion
			specType, err := GetSpecType(schemaVersion)
			assert.NoError(err)
			assert.Equal(tc.specType, specType)

			requiredTokens, err := GetRequiredTokensFromTests(tests, specType)
			assert.NoError(err)
			assert.NotEmpty(requiredTokens)
		})
	}
}

// GetPaymentPermissions - permx related
func TestPaymentPermissions(t *testing.T) {
	assert := assert.New(t)

	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: paymentsSwaggerLocation31,
	}

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1")
	context := model.Context{"apiversions": values}

	specType, err := GetSpecType(apiSpec.SchemaVersion)
	assert.Equal("payments", specType)
	assert.NoError(err)

	manifestPath := getManifestPaths()[specType]
	scripts, _, err := LoadGenerationResources(specType, manifestPath, &context)
	assert.NoError(err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &context,
		Endpoints:    readDiscovery2(specType),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, scripts, err := GenerateTestCases(&params)
	assert.NoError(err)
	assert.NotEmpty(tests)
	assert.NotEmpty(scripts)

	requiredTokens, err := GetPaymentPermissions(tests)
	assert.NoError(err)
	assert.NotEmpty(requiredTokens)
}

// replace it with loadAssertations test
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

func TestPermissionFiteringAccounts(t *testing.T) {

	ctx := model.Context{
		"accountId":           "123123123",
		"client_access_token": "abc-defg-hijk-lmno-pqrs",
		"apiversions":         []string{"accounts_v3.1.1", "payments_v3.1.1"},
	}

	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}

	specType, err := GetSpecType(apiSpec.SchemaVersion)
	assert.Nil(t, err)
	manifestPath := getManifestPaths()[specType]

	scripts, _, err := LoadGenerationResources(specType, manifestPath, &ctx)
	if err != nil {
		fmt.Println("Error on loadGenerationResources")
		return
	}
	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &ctx,
		Endpoints:    readDiscovery2(specType),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.NoError(t, err)

	fmt.Printf("%d tests loaded", len(tests))

	// filteredScripts, err := FilterTestsBasedOnDiscoveryEndpointsPlayground(scripts, endpoints)
	// if err != nil {

	// }
	// for _, v := range filteredScripts.Scripts {
	// 	dumpJSON(v)
	// }
}

func readDiscovery2(apiSpec string) []discovery.ModelEndpoint {
	discoveryJSON, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-ozone.json")
	if err != nil {
		fmt.Println("discovery read failed")
		return nil
	}

	disco := &discovery.Model{}

	if err = json.Unmarshal(discoveryJSON, &disco); err != nil {
		return nil
	}

	switch apiSpec {
	case "accounts":
		return disco.DiscoveryModel.DiscoveryItems[0].Endpoints
	case "payments":
		return disco.DiscoveryModel.DiscoveryItems[1].Endpoints
	case "cbpii":
		return disco.DiscoveryModel.DiscoveryItems[2].Endpoints
	// TODO: Add VRP support
	// case "vrp":
	// 	return disco.DiscoveryModel.DiscoveryItems[2].Endpoints
	default:
		return nil
	}
}

// look at context
func TestPaymentTestCaseCreation(t *testing.T) {
	ctx := &model.Context{
		"consent_id":                          "aac-fee2b8eb-ce1b-48f1-af7f-dc8f576d53dc",
		"xchange_code":                        "10e9d80b-10d4-4abd-9fe0-15789cc512b5",
		"baseurl":                             "https://matls-sso.openbankingtest.org.uk",
		"access_token":                        "18d5a754-0b76-4a8f-9c68-dc5caaf812e2",
		"client_id":                           "12312",
		"scope":                               "AuthoritiesReadAccess ASPSPReadAccess TPPReadAll",
		"authorisation_endpoint":              "https://example.com/authorisation",
		"OB-301-DOP-100300-ConsentId":         "100100-ConsentId",
		"OB-301-DOP-100600-DomesticPaymentId": "100600-DomesticPaymentId-PaymentId",
		"OB-301-DOP-100100-ConsentId":         "100100-ConsentId",
		"OB-301-DOP-100800-ConsentId":         "100800-Consentid",
		"creditorIdentification":              "1231231231",
		"thisCurrency":                        "GBP",
		"creditorScheme":                      "default",
	}

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1")
	ctx.Put("apiversions", values)

	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: paymentsSwaggerLocation31,
	}

	specType, err := GetSpecType(apiSpec.SchemaVersion)
	assert.Nil(t, err)
	manifestPath := getManifestPaths()[specType]
	scripts, _, err := LoadGenerationResources(specType, manifestPath, ctx)
	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          ctx,
		Endpoints:    readDiscovery2(specType),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.Nil(t, err)

	fmt.Printf("we have %d tests\n", len(tests))
	for _, v := range tests {
		dumpJSON(v)
	}

}

// TestFilterTestsBasedOnDiscoveryEndpoints with this test we want to test filtering of Scripts.
// Given a collection of `Scripts` and a collection of `endpoints`, we want the tested function return
// a subset of `Scripts`, where the URI of each returned script matches an endpoint (via regex) of at least one of
// the paths in the collection of `endpoints`.
func TestFilterTestsBasedOnDiscoveryEndpoints(t *testing.T) {
	scripts := Scripts{
		Scripts: []Script{
			{
				ID:  "0000",
				URI: "/domestic-payment-consents/ConsentID-Here1234",
			},
			{
				ID:  "1000",
				URI: "/domestic-payment-consents",
			},
			{
				ID:  "2000",
				URI: "/domestic-payment-consents/ConsentID-Here1234/funds-confirmation",
			},
			{
				ID:  "3000",
				URI: "/domestic-payments",
			},
			{
				ID:  "4000",
				URI: "/domestic-payments/ConsentID-Here1234",
			},
			{
				ID:  "5000",
				URI: "/domestic-scheduled-payment-consents",
			},
			{
				ID:  "6000",
				URI: "/domestic-scheduled-payment-consents/ConsentID-Here1234",
			},
			{
				ID:  "7000",
				URI: "/domestic-scheduled-payment-consents/ConsentID-Here1234",
			},
			{
				ID:  "8000",
				URI: "/domestic-scheduled-payments/DomesticSceduledPaymentID-Here1234",
			},
			{
				ID:  "90000",
				URI: "/domestic-standing-order-consents",
			},
			{
				ID:  "10000",
				URI: "/domestic-standing-order-consents/ConsentID-Here1234",
			},
			{
				ID:  "11000",
				URI: "/domestic-standing-orders/DomesticStandingOrderID-Here1234",
			},
			{
				ID:  "12000",
				URI: "/international-payment-consents",
			},
			{
				ID:  "13000",
				URI: "/international-payment-consents/ConsentID-Here1234",
			},
			{
				ID:  "14000",
				URI: "/international-payments",
			},
			{
				ID:  "15000",
				URI: "/international-payments/ConsentID-Here1234",
			},
			{
				ID:  "16000",
				URI: "/international-scheduled-payment-consents",
			},
			{
				ID:  "17000",
				URI: "/international-scheduled-payment-consents/ConsentID-Here1234",
			},
			{
				ID:  "18000",
				URI: "/international-scheduled-payments",
			},
			{
				ID:  "19000",
				URI: "/international-scheduled-payments/InternationalScheduledPaymentID-Here1234",
			},
		},
	}
	endpoints := []discovery.ModelEndpoint{
		{
			Path: "/domestic-payment-consents/1234",
		},
		{
			Path: "/domestic-payment-consents",
		},
		{
			Path: "/domestic-payment-consents/2345678987DFGHJGH/funds-confirmation",
		},
		{
			Path: "/international-payment-consents",
		},
		{
			Path: "/international-payments/INT-PAY-1234-ID",
		},
		{
			Path: "/international-scheduled-payments/InternationalScheduledPaymentID-Here1234",
		},
	}
	filtered, err := FilterTestsBasedOnDiscoveryEndpoints(scripts, endpoints, paymentsRegex)
	assert.NoError(t, err)

	// As a simple check validate the lengths match
	assert.Equal(t, len(endpoints), len(filtered.Scripts))

	// Now, we need to check that the scripts we expect are actually in the result set
	// We know what to check for by manually matching the paths in `endpoints` to paths in `scripts`
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[0]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[1]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[2]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[12]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[15]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[19]))
}

// TestFilterTestsBasedOnDiscoveryEndpoints with this test we want to test filtering of Scripts.
// Given a collection of `Scripts` and a collection of `endpoints`, we want the tested function return
// a subset of `Scripts`, where the URI of each returned script matches an endpoint (via regex) of at least one of
// the paths in the collection of `endpoints`.
func TestVrpFilterTestsBasedOnDiscoveryEndpoints(t *testing.T) {
	scripts := Scripts{
		Scripts: []Script{
			{
				ID:  "0000",
				URI: "/domestic-vrp-consents/ConsentID-Here1234",
			},
			{
				ID:  "1000",
				URI: "/domestic-vrp-consents",
			},
			{
				ID:  "2000",
				URI: "/domestic-vrp-consents/ConsentID-Here1234/funds-confirmation",
			},
			{
				ID:  "3000",
				URI: "/domestic-vrps/VRPId-Here1234/payment-details",
			},
		},
	}
	endpoints := []discovery.ModelEndpoint{
		{
			Path: "/domestic-vrp-consents/1234",
		},
		{
			Path: "/domestic-vrp-consents",
		},
		{
			Path: "/domestic-vrp-consents/2345678987DFGHJGH/funds-confirmation",
		},
	}
	filtered, err := FilterTestsBasedOnDiscoveryEndpoints(scripts, endpoints, vrpRegex)
	assert.NoError(t, err)

	// As a simple check validate the lengths match
	assert.Equal(t, len(endpoints), len(filtered.Scripts))

	// Now, we need to check that the scripts we expect are actually in the result set
	// We know what to check for by manually matching the paths in `endpoints` to paths in `scripts`
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[0]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[1]))
	assert.True(t, contains(filtered.Scripts, scripts.Scripts[2]))
}

func TestContains(t *testing.T) {
	collection := []Script{
		{
			ID: "123",
		},
		{
			ID: "456",
		},
	}

	subjectExists := Script{
		ID: "123",
	}
	subjectNotExists := Script{
		ID: "789",
	}

	assert.True(t, contains(collection, subjectExists))
	assert.False(t, contains(collection, subjectNotExists))
}

func TestGetSpecVersion(t *testing.T) {
	assert := assert.New(t)

	expectedSpecVersion := semver.Version{Major: 3, Minor: 1, Patch: 11}

	tests := map[string]struct {
		specType    string
		apiVersions []string
		specVersion semver.Version
		err         string
	}{
		"OK - VRP": {
			specType:    "vrps",
			apiVersions: []string{"vrps_v3.1.11"},
			specVersion: expectedSpecVersion,
		},
		// if specType matches the first value from apiVersions it takes it
		"OK - VRP - multiple VRP apiVersions": {
			specType:    "vrps",
			apiVersions: []string{"vrps_v3.1.11", "vrps_v3.1.12"},
			specVersion: expectedSpecVersion,
		},
		"OK - VRP - multiple apiVersions (one standard version)": {
			specType:    "vrps",
			apiVersions: []string{"accounts_v3.1.11", "vrps_v3.1.11", "payments_v3.1.11", "cbpii_v3.1.11"},
			specVersion: expectedSpecVersion,
		},
		"OK - VRP - multiple apiVersions (multiple standard version)": {
			specType:    "vrps",
			apiVersions: []string{"accounts_v3.1.10", "vrps_v3.1.11"},
			specVersion: expectedSpecVersion,
		},
		"Bad - empty": {
			err: "getSpecVersion: cannot parse versions []",
		},
		"Bad - empty specType": {
			apiVersions: []string{"vrps_v3.1.11"},
			err:         "getSpecVersion: cannot parse versions [vrps_v3.1.11]",
		},
		// TODO:
		"Bad - empty apiVersions": {
			specType: "vrps",
			err:      "getSpecVersion: cannot parse versions []",
		},
		// TODO: specVersion should be empty, but is 3.1.11
		"Bad - apiVersions with incorrect formatting (missed API group)": {
			apiVersions: []string{"_v3.1.11"},
			specVersion: expectedSpecVersion,
			// err:         "getSpecVersion: cannot parse versions []",
		},
		"Bad - apiVersions with incorrect formatting (missed standard version)": {
			apiVersions: []string{"vrps_v"},
			err:         "getSpecVersion: cannot parse versions [vrps_v]",
		},
		"Bad - correct specType, apiVersions with incorrect formatting (missed API group)": {
			specType:    "vrps",
			apiVersions: []string{"_v3.1.11"},
			err:         "getSpecVersion: cannot parse versions [_v3.1.11]",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			specVersion, err := getSpecVersion(tt.specType, tt.apiVersions)

			if tt.err != "" {
				assert.EqualError(err, tt.err)
			} else {
				assert.NoError(err)
			}

			assert.Equal(tt.specVersion, specVersion)
		})
	}
}

func TestLoadAssert(t *testing.T) {
	assert := assert.New(t)

	references, err := loadAssert()
	assert.NoError(err)
	assert.NotEmpty(references)
}

func TestLoadAssertions(t *testing.T) {
	assert := assert.New(t)

	references, err := loadAssertions()
	assert.NoError(err)
	assert.NotEmpty(references)
}
