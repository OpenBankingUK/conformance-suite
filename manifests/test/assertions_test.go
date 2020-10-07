package assertionstest

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

var (
	accountSpecPath               = flag.String("acc_spec", "../../pkg/schema/spec/v3.1.6/account-info-swagger-flattened.json", "Path to the accounts specification swagger file.")
	paymentSpecPath               = flag.String("pay_spec", "../../pkg/schema/spec/v3.1.6/payment-initiation-swagger-flattened.json", "Path to the payments specification swagger file.")
	cbpiiSpecPath                 = flag.String("cbpii_spec", "../../pkg/schema/spec/v3.1.6/confirmation-funds-flattened.json", "Path to the funds confirmations specification swagger file.")
	assertionsPath                = flag.String("assertions", "../assertions.json", "Path to the JSON file containing the assertion rules.")
	accountsManifestPath          = flag.String("acc_man", "../ob_3.1_accounts_transactions_fca.json", "Path to accounts tests json file.")
	paymentsManifestPath          = flag.String("pay_man", "../ob_3.1_payment_fca.json", "Path to payments tests json file.")
	fundsConfirmationManifestPath = flag.String("cbpii_man", "../ob_3.1_cbpii_fca.json", "Path to funds confirmations tests json file.")

	// Load scripts from all the paths above. They contain the assertion 'sets' tested here.
	scripts = func() []manifest.Script {
		s := []manifest.Script{}
		for _, path := range []string{*accountsManifestPath, *paymentsManifestPath, *fundsConfirmationManifestPath} {
			scripts := &manifest.Scripts{}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(b, scripts)
			if err != nil {
				log.Fatal(err)
			}
			s = append(s, scripts.Scripts...)
		}
		return s
	}()

	refs = func() map[string]manifest.Reference {
		b, err := ioutil.ReadFile(*assertionsPath)
		if err != nil {
			log.Fatal(err)
		}
		refs := &manifest.References{}
		err = json.Unmarshal(b, refs)
		if err != nil {
			log.Fatal(err)
		}
		return refs.References
	}()
)

func getScript(id string) *manifest.Script {
	for _, s := range scripts {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

// Testing assertions as defined on specific scripts.
// These tests verify that the validation specified for a given test case
// passes or fails according the expectations when certain (mocked) reponses
// from the ASPSP are processed.
func TestAssertions(t *testing.T) {
	emptyContext := &model.Context{}
	_ = emptyContext

	type mockResponse struct {
		code    int
		headers map[string]string
		body    string
	}

	testCases := []struct {
		name                          string       // for our eyes - to recognise which test fails
		manifestID                    string       // the id of the test in (any of) the manifest script files
		response                      mockResponse // the mocked response from the ASPSP
		schemaSpec                    string       // path to the jsonschema spec to be used with this particular case
		ExpectValidationPass          bool         // should the scenario with the above parameters pass / fail ?
		ExpectValidationErrorContains string       // the validation step should produce an error which contains this string
	}{
		// ADD tests, for example:
		// {
		// 	name:       "OB-xxx-yyy-zzzzzz pass if ASPSP returns correct error",
		// 	manifestID: "OB-xxx-yyy-zzzzzz",
		// 	response: mockResponse{
		// 		400,
		// 		map[string]string{},
		// 		`{"Errors":[{"ErrorCode":"UK.OBIE.??????????????"}]}`,
		// 	},
		// 	schemaSpec:                    *paymentSpecPath,
		// 	ExpectValidationPass:          true,
		// 	ExpectValidationErrorContains: "????????????????",
		// },
		{
			name:       "OB-301-DOP-100110 pass if ASPSP returns missing claim error with code 400",
			manifestID: "OB-301-DOP-100110",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.MissingClaim"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-DOP-100110 pass if ASPSP returns invalid claim error with code 400",
			manifestID: "OB-301-DOP-100110",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.InvalidClaim"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-DOP-100110 pass if ASPSP returns malformed signature error with code 400",
			manifestID: "OB-301-DOP-100110",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.Malformed"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-DOP-100110 does not pass if ASPSP returns any error with code other than 400",
			manifestID: "OB-301-DOP-100110",
			response: mockResponse{
				401,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.Malformed"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-DOP-100110 does not pass if ASPSP returns incorrect error with code 400",
			manifestID: "OB-301-DOP-100110",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.Invalid"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-316-DOP-100310 pass if ASPSP returns correct error code and message",
			manifestID: "OB-316-DOP-100310",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.Missing"}]}`,
			},
			schemaSpec:           *paymentSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-316-DOP-100310 fails if ASPSP returns incorrect error code",
			manifestID: "OB-316-DOP-100310",
			response: mockResponse{
				401,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Signature.Missing"}]}`,
			},
			schemaSpec:                    *paymentSpecPath,
			ExpectValidationPass:          false,
			ExpectValidationErrorContains: "HTTP Status code does not match: expected 400 got 401",
		},
		{
			name:       "OB-316-DOP-100310 fails if ASPSP returns incorrect error message",
			manifestID: "OB-316-DOP-100310",
			response: mockResponse{
				400,
				map[string]string{},
				`{"Errors":[{"ErrorCode":"UK.OBIE.Incorrect"}]}`,
			},
			schemaSpec:                    *paymentSpecPath,
			ExpectValidationPass:          false,
			ExpectValidationErrorContains: "JSON Match Failed - expected (UK.OBIE.Signature.Missing)",
		},
	}

	for _, test := range testCases {
		manifestTC, err := makeTestCase(test.manifestID, test.schemaSpec)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		mockResp := createHTTPResponse(test.response.code, test.response.body, test.response.headers)
		result, errors := manifestTC.Validate(mockResp, emptyContext)
		if result != true {
			errorMessages := &strings.Builder{}
			for _, e := range errors {
				errorMessages.WriteString(fmt.Sprintf("\t- %s\n", e))
			}
			assert.True(t, errorsContain(errors, test.ExpectValidationErrorContains),
				"%s: validation errors should contain an error with text: '%s'\nERRORS:\n%v", test.name, test.ExpectValidationErrorContains, errorMessages)
		}
		assert.Equal(t, test.ExpectValidationPass, result, "%s result - expected: %v actual: %v", test.name, test.ExpectValidationPass, result)
	}
}

// Testing some assertions or combinations in isolation, independently from how they are used in scripts
func TestAccountTransactions(t *testing.T) {
	emptyContext := &model.Context{}

	b, err := ioutil.ReadFile(*assertionsPath)
	if err != nil {
		t.Fatal(err)
	}

	refs := manifest.References{}
	err = json.Unmarshal(b, &refs)
	if err != nil {
		t.Fatal(err)
	}

	// Testing the cases where API returns either 400 or 403 status code.
	// Responses with status code 400 must provide correct error message
	// but responses with 403 code may return with any body.
	testCase := model.TestCase{
		ExpectOneOf: []model.Expect{
			refs.References["OB3IPAssertResourceFieldInvalidOBErrorCode400"].Expect,
			refs.References["OB3GLOAssertOn403"].Expect,
		},
	}

	testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Returning 403 without response body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := createHTTPResponse(403, "", headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 403 with any response body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := createHTTPResponse(403, "response body is not checked (TBD: does non-empty body have to follow schema?)", headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 400 with correct body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := createHTTPResponse(400, `{"Errors":[{"ErrorCode":"UK.OBIE.Field.Invalid"}]}`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 400 with incorrect body should FAIL", func(t *testing.T) {
		headers := map[string]string{}
		resp := createHTTPResponse(400, `{"Errors":[]}`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		assert.True(t, errorsContain(err, "JSON Match Failed"))
		assert.False(t, result, "expected: %v actual: %v", false, result)
	})

	t.Run("Returning incorrect status should FAIL", func(t *testing.T) {
		headers := map[string]string{}
		resp := createHTTPResponse(200, `OK`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		assert.True(t, errorsContain(err, "HTTP Status code does not match"))
		assert.False(t, result, "expected: %v actual: %v", true, result)
	})
}

// duplicates some of the mechanisms used in the testcase builder
// it's somewhat brittle; consider exposing relevant bits to be imported / used here
func makeTestCase(scriptID, specPath string) (model.TestCase, error) {
	s := getScript(scriptID)

	type testEntry struct {
		// relevant bits of script
		id           string
		Asserts      []string `json:"asserts"`
		AssertsOneOf []string `json:"asserts_one_of"`
		SchemaCheck  bool
	}

	tc := model.MakeTestCase()
	for _, a := range s.Asserts {
		ref, exists := refs[a]
		if !exists {
			msg := fmt.Sprintf("assertion %s do not exist in reference data", a)
			return tc, errors.New(msg)
		}
		clone := ref.Expect.Clone()
		if ref.Expect.StatusCode != 0 {
			tc.Expect.StatusCode = clone.StatusCode
		}
		tc.Expect.Matches = append(tc.Expect.Matches, clone.Matches...)
	}

	for _, a := range s.AssertsOneOf {
		ref, exists := refs[a]
		if !exists {
			msg := fmt.Sprintf("assertion %s does not exist in reference data", a)
			return tc, errors.New(msg)
		}
		tc.ExpectOneOf = append(tc.ExpectOneOf, ref.Expect.Clone())
	}

	var err error
	tc.Validator, err = schema.NewSwaggerValidator(specPath)
	if err != nil {
		return tc, err
	}

	tc.Expect.SchemaValidation = s.SchemaCheck
	return tc, nil
}

func errorsContain(errs []error, s string) bool {
	for _, err := range errs {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}

func createHTTPResponse(code int, body string, headers map[string]string) *resty.Response {
	mockedServer, mockedServerURL := test.HTTPServer(code, body, headers)
	defer mockedServer.Close()
	res, _ := resty.R().Get(mockedServerURL)
	return res
}
