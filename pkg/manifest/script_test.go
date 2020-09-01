package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

func TestGenerateTestCases(t *testing.T) {
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}
	specType, err := GetSpecType(apiSpec.SchemaVersion)
	fmt.Println(specType)
	assert.Nil(t, err)

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1")
	context := model.Context{"apiversions": values}

	scripts, _, err := LoadGenerationResources(specType, manifestPath, &context)

	params := GenerationParameters{Scripts: scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &context,
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.Nil(t, err)

	perms := getAccountPermissions(tests)
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
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}
	specType, err := GetSpecType(apiSpec.SchemaVersion)
	assert.Nil(t, err)
	scripts, _, err := LoadGenerationResources(specType, manifestPath, nil)
	if err != nil {
		fmt.Println("Error on loadGenerationResources")
		return
	}

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &model.Context{},
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.NoError(t, err)

	fmt.Printf("we have %d tests\n", len(tests))
	for _, v := range tests {
		dumpJSON(v)
	}

	requiredTokens, err := GetPaymentPermissions(tests)
	assert.NoError(t, err)

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

func TestPermissionFiteringAccounts(t *testing.T) {

	ctx := model.Context{
		"accountId":           "123123123",
		"client_access_token": "abc-defg-hijk-lmno-pqrs",
	}

	endpoints := readDiscovery()
	apiSpec := discovery.ModelAPISpecification{
		SchemaVersion: accountSwaggerLocation31,
	}
	scripts, _, err := LoadGenerationResources("accounts", manifestPath, nil)
	if err != nil {
		fmt.Println("Error on loadGenerationResources")
		return
	}
	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          &ctx,
		Endpoints:    readDiscovery(),
		ManifestPath: manifestPath,
		Validator:    schema.NewNullValidator(),
	}
	tests, _, err := GenerateTestCases(&params)
	assert.NoError(t, err)

	fmt.Printf("%d tests loaded", len(tests))

	filteredScripts, err := FilterTestsBasedOnDiscoveryEndpointsPlayground(scripts, endpoints)
	if err != nil {

	}
	for _, v := range filteredScripts.Scripts {
		dumpJSON(v)
	}
}

func readDiscovery() []discovery.ModelEndpoint {
	discoveryJSON, err := ioutil.ReadFile("../discovery/templates/ob-v3.1-ozone.json")
	if err != nil {
		fmt.Println("discovery read failed")
		return nil
	}

	disco := &discovery.Model{}

	err = json.Unmarshal(discoveryJSON, &disco)

	return disco.DiscoveryModel.DiscoveryItems[0].Endpoints

}

func FilterTestsBasedOnDiscoveryEndpointsPlayground(scripts Scripts, endpoints []discovery.ModelEndpoint) (Scripts, error) {

	lookupMap := make(map[string]bool)
	_ = lookupMap
	filteredScripts := []Script{}
	fmt.Println("***Discovery Endpoint URLs")

	for _, ep := range endpoints {
		for _, regpath := range accountsRegex {
			matched, err := regexp.MatchString(regpath.Regex, ep.Path)
			if err != nil {
				continue
			}
			if matched {
				lookupMap[regpath.Regex] = true
				fmt.Printf("endpoint %40.40s matched by regex %42.42s: %s\n", ep.Path, regpath.Regex, regpath.Name)
			}
		}
	}
	fmt.Println("***Script URLs")
	for _, scr := range scripts.Scripts {
		for _, regpath := range accountsRegex {
			stripped := strings.Replace(scr.URI, "$", "", -1) // only works with a single character
			matched, err := regexp.MatchString(regpath.Regex, stripped)
			if err != nil {
				fmt.Printf("matching err " + err.Error())
				continue
			}
			if matched {
				fmt.Printf("%40.40s matched by regex %42.42s: %s\n", scr.URI, regpath.Regex, regpath.Name)
			} else {
				//fmt.Printf("No match %s\n", scr.URI)
			}
		}
	}

	fmt.Println("***dmp")
	for k := range lookupMap {
		fmt.Printf("lookup map %s\n", k)

	}

	fmt.Println("***Lookup Map")
	for k := range lookupMap {
		for i, scr := range scripts.Scripts {
			stripped := strings.Replace(scr.URI, "$", "", -1) // only works with a single character
			matched, err := regexp.MatchString(k, stripped)
			if err != nil {
				continue
			}
			if matched {
				fmt.Printf("endpoint %40.40s matched by regex %42.42s\n", scr.URI, k)
				filteredScripts = append(filteredScripts, scripts.Scripts[i])
			}
		}
	}
	myscripts := Scripts{Scripts: filteredScripts}

	return myscripts, nil
}

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
	scripts, _, err := LoadGenerationResources(specType, manifestPath, ctx)
	assert.Nil(t, err)

	params := GenerationParameters{
		Scripts:      scripts,
		Spec:         apiSpec,
		Baseurl:      "http://mybaseurl",
		Ctx:          ctx,
		Endpoints:    readDiscovery(),
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
