package assertionstest

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"gopkg.in/resty.v1"
)

var (
	accountSpecPath = flag.String("spec", "../../pkg/schema/spec/v3.1.5/account-info-swagger-flattened.json", "Path to the specification swagger file.")
	assertionsPath  = flag.String("assertions", "../assertions.json", "Path to the JSON file containing the assertion rules.")
)

// Testing the scenarios in the manifests/ob_3.1_accounts_transactions_fca.json
// It is important to note that there is no automatic mechanism keeping the above file and these tests synchronised
// which means that in case the above file changes then these tests must be updated accordingly.
func TestAccountTransactions(t *testing.T) {
	const emptyBody = ""
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

	// Testing the cases where API returns either 400 or 403 status code
	// Responses with status code 400 must provide correct error message
	// but responses with 403 code may return with any body.
	testCase := model.TestCase{
		ExpectOneOf: []model.Expect{
			refs.References["OB3IPAssertResourceFieldInvalidOBErrorCode"].Expect,
			refs.References["OB3GLOAssertOn403"].Expect,
		},
	}

	testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Returning 403 without response body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := CreateHTTPResponse(403, emptyBody, headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 403 with any response body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := CreateHTTPResponse(403, "response body is not checked (TBD: does non-empty body have to follow schema?)", headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 400 with correct body should PASS", func(t *testing.T) {
		headers := map[string]string{}
		resp := CreateHTTPResponse(400, `{"Errors":[{"ErrorCode":"UK.OBIE.Field.Invalid"}]}`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Returning 400 with incorrect body should FAIL", func(t *testing.T) {
		headers := map[string]string{}
		resp := CreateHTTPResponse(400, `{"Errors":[]}`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		assert.True(t, errorsContain(err, "JSON Match Failed"))
		assert.False(t, result, "expected: %v actual: %v", false, result)
	})

	t.Run("Returning incorrect status should FAIL", func(t *testing.T) {
		headers := map[string]string{}
		resp := CreateHTTPResponse(200, `OK`, headers)
		result, err := testCase.Validate(resp, emptyContext)
		assert.True(t, errorsContain(err, "HTTP Status code does not match"))
		assert.False(t, result, "expected: %v actual: %v", true, result)
	})
}

func errorsContain(errs []error, s string) bool {
	for _, err := range errs {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}

func CreateHTTPResponse(code int, body string, headers map[string]string) *resty.Response {
	mockedServer, mockedServerURL := test.HTTPServer(code, body, headers)
	defer mockedServer.Close()
	res, _ := resty.R().Get(mockedServerURL)
	return res
}
