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

	"github.com/OpenBankingUK/conformance-suite/pkg/manifest"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
	"github.com/OpenBankingUK/conformance-suite/pkg/test"
)

var (
	accountSpecPath               = flag.String("acc_spec", "../../pkg/schema/spec/v3.1.6/account-info-swagger-flattened.json", "Path to the accounts specification swagger file.")
	paymentSpecPath               = flag.String("pay_spec", "../../pkg/schema/spec/v3.1.6/payment-initiation-swagger-flattened.json", "Path to the payments specification swagger file.")
	cbpiiSpecPath                 = flag.String("cbpii_spec", "../../pkg/schema/spec/v3.1.6/confirmation-funds-flattened.json", "Path to the funds confirmations specification swagger file.")
	assertionsPath                = flag.String("assertions", "../assertions.json", "Path to the JSON file containing the assertion rules.")
	accountsManifestPath          = flag.String("acc_man", "../ob_3.1_accounts_transactions_fca.json", "Path to accounts tests json file.")
	paymentsManifestPath          = flag.String("pay_man", "../ob_3.1_payment_fca.json", "Path to payments tests json file.")
	fundsConfirmationManifestPath = flag.String("cbpii_man", "../ob_3.1_cbpii_fca.json", "Path to funds confirmations tests json file.")
	accountsManifestPathV4        = flag.String("acc_man_v4", "../ob_4.0_accounts_transactions_fca.json", "Path to v4 accounts tests json file.")

	// Load scripts from all the paths above. They contain the assertion 'sets' tested here.
	scripts = func() []manifest.Script {
		s := []manifest.Script{}
		for _, path := range []string{*accountsManifestPath, *paymentsManifestPath, *fundsConfirmationManifestPath, *accountsManifestPathV4} {
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
		//{
		//	name:       "OB-xxx-yyy-zzzzzz pass if ASPSP returns correct error",
		//	manifestID: "OB-xxx-yyy-zzzzzz",
		//	response: mockResponse{
		//		400,
		//		map[string]string{},
		//		`{"Errors":[{"ErrorCode":"UK.OBIE.??????????????"}]}`,
		//	},
		//	schemaSpec:                    *paymentSpecPath,
		//	ExpectValidationPass:          true,
		//	ExpectValidationErrorContains: "????????????????",
		//},
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
		{
			name:       "OB-400-TRA-101200 fails if ASPSP doesn't return BankTransactionCode or ProprietaryBankTransactionCode",
			manifestID: "OB-400-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-101200 passes if ASPSP returns at least BankTransactionCode",
			manifestID: "OB-400-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"BankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-TRA-101200 passes if ASPSP returns at least ProprietaryBankTransactionCode",
			manifestID: "OB-400-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"ProprietaryBankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-TRA-101200 passes if ASPSP returns both BankTransactionCode and ProprietaryBankTransactionCode",
			manifestID: "OB-400-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"BankTransactionCode": {}, "ProprietaryBankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-TRA-101200 fails if ASPSP doesn't return BankTransactionCode or ProprietaryBankTransactionCode",
			manifestID: "OB-301-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-101200 passes if ASPSP returns at least BankTransactionCode",
			manifestID: "OB-301-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"BankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-TRA-101200 passes if ASPSP returns at least ProprietaryBankTransactionCode",
			manifestID: "OB-301-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"ProprietaryBankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-TRA-101200 passes if ASPSP returns both BankTransactionCode and ProprietaryBankTransactionCode",
			manifestID: "OB-301-TRA-101200",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Transaction": [{"BankTransactionCode": {}, "ProprietaryBankTransactionCode": {}, "SomeRandomThing": {}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.StatementFrequencyAndFormat Tests 100000
		{
			name:       "OB-400-ACC-100000 should fail when StatementFrequencyAndFormat is present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"StatementFrequencyAndFormat":[{"Frequency": "YEAR"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100000 should pass when StatementFrequencyAndFormat is not present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.Servicer Tests 100000
		{
			name:       "OB-400-ACC-100000 should fail when Servicer is present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"Servicer":{"SchemeName":"UK.OBIE.BICFI","Identification":"8020441910203345","Name":"ServicerName"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100000 should pass when Servicer is not present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.Account Tests 100000
		{
			name:       "OB-400-ACC-100000 should fail when Account.Account is present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"Account":[{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200110203345","Name":"Mr Kevin","SecondaryIdentification":"00021","LEI":"9193001QZMP2PQT4AK86"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100000 should pass when Account.Account is not present",
			manifestID: "OB-400-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.StatementFrequencyAndFormat Tests 100300
		{
			name:       "OB-400-ACC-100300 should fail when StatementFrequencyAndFormat is present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"StatementFrequencyAndFormat":[{"Frequency": "YEAR"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100300 should pass when StatementFrequencyAndFormat is not present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.Servicer Tests 100300
		{
			name:       "OB-400-ACC-100300 should fail when Servicer is present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"Servicer":{"SchemeName":"UK.OBIE.BICFI","Identification":"8020441910203345","Name":"ServicerName"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100300 should pass when Servicer is not present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		// Account.Account Tests 100300
		{
			name:       "OB-400-ACC-100300 should fail when Account.Account is present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{},
				`{"Data":{"Account":[{"Account":[{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200110203345","Name":"Mr Kevin","SecondaryIdentification":"00021","LEI":"9193001QZMP2PQT4AK86"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-ACC-100300 should pass when Account.Account is not present",
			manifestID: "OB-400-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-TRA-105000 should pass when no Transaction Details are present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-TRA-105000 should fail when TransactionInformation is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"TransactionInformation": "Cash from Aubrey"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when Balance is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"Balance":{"Amount":{"Amount":"230.00","Currency":"GBP"},"CreditDebitIndicator":"Credit","Type":"ITBD"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when MerchantDetails is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"MerchantDetails":{"MerchantName":"Merchant's Name","MerchantCategoryCode":"5874"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when CreditorAgent is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200112345678","Name":"Mrs Juniper","SecondaryIdentification":"80200112374165","Proxy":{"Identification":"2360549017905188","Code":"TELE"}}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when CreditorAgent is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAgent":{"LEI":"IZ9Q00LZEVUKWCQY6X15","SchemeName":"UK.OBIE.BICFI","Identification":"80200112344562","Name":"The Credit Agent"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when CreditorAccount is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200112345678","Name":"Mrs Juniper","SecondaryIdentification":"80200112374165","Proxy":{"Identification":"2360549017905188","Code":"TELE"}}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when UltimateCreditor is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"UltimateCreditor":{"SchemeName":"UK.OBIE.BICFI","Identification":"2360549017905161589","Name":"Ultimate Creditor","LEI":"60450004FECVJV7YN339"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when DebtorAgent is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAgent":{"LEI":"IZ9Q00LZEVUKWCQY8i14","SchemeName":"UK.OBIE.BICFI","Identification":"8020011234487","Name":"The Debtor Agent"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when DebtorAccount is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200112345784","Name":"Mr Juniper","SecondaryIdentification":"80200112378745"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105000 should fail when UltimateDebtor is present",
			manifestID: "OB-400-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"UltimateDebtor":{"SchemeName":"UK.OBIE.BICFI","Identification":"2360549017905161589","Name":"Ultimate Debtor","LEI":"8200007YHFDMEODY1965"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		// Transaction Details Tests OB-400-TRA-105200
		{
			name:       "OB-400-TRA-105200 should pass when no Transaction Details are present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-TRA-105200 should fail when TransactionInformation is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"TransactionInformation": "Cash from Aubrey"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when Balance is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"Balance":{"Amount":{"Amount":"230.00","Currency":"GBP"},"CreditDebitIndicator":"Credit","Type":"ITBD"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when MerchantDetails is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"MerchantDetails":{"MerchantName":"Merchant's Name","MerchantCategoryCode":"5874"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when CreditorAgent is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAgent":{"LEI":"IZ9Q00LZEVUKWCQY6X15","SchemeName":"UK.OBIE.BICFI","Identification":"80200112344562","Name":"The Credit Agent"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when CreditorAccount is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200112345678","Name":"Mrs Juniper","SecondaryIdentification":"80200112374165","Proxy":{"Identification":"2360549017905188","Code":"TELE"}}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when UltimateCreditor is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"UltimateCreditor":{"SchemeName":"UK.OBIE.BICFI","Identification":"2360549017905161589","Name":"Ultimate Creditor","LEI":"60450004FECVJV7YN339"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when DebtorAgent is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAgent":{"LEI":"IZ9Q00LZEVUKWCQY8i14","SchemeName":"UK.OBIE.BICFI","Identification":"8020011234487","Name":"The Debtor Agent"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when DebtorAccount is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAccount":{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200112345784","Name":"Mr Juniper","SecondaryIdentification":"80200112378745"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-TRA-105200 should fail when UltimateDebtor is present",
			manifestID: "OB-400-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"UltimateDebtor":{"SchemeName":"UK.OBIE.BICFI","Identification":"2360549017905161589","Name":"Ultimate Debtor","LEI":"8200007YHFDMEODY1965"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		// v3 AccountDetails Tests 100000
		{
			name:       "OB-301-ACC-100000 should pass when no AccountDetails are present",
			manifestID: "OB-301-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-ACC-100000 should fail when Servicer is present",
			manifestID: "OB-301-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"Servicer":{"SchemeName":"UK.OBIE.BICFI","Identification":"8020441910203345","Name":"ServicerName"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-ACC-100000 should fail when Account is present",
			manifestID: "OB-301-ACC-100000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"Account":[{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200110203345","Name":"Mr Kevin","SecondaryIdentification":"00021","LEI":"9193001QZMP2PQT4AK86"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		// v3 AccountDetails Tests 100300
		{
			name:       "OB-301-ACC-100300 should pass when no AccountDetails are present",
			manifestID: "OB-301-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-ACC-100300 should fail when Servicer is present",
			manifestID: "OB-301-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"Servicer":{"SchemeName":"UK.OBIE.BICFI","Identification":"8020441910203345","Name":"ServicerName"}}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-ACC-100300 should fail when Account is present",
			manifestID: "OB-301-ACC-100300",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Account":[{"Account":[{"SchemeName":"UK.OBIE.SortCodeAccountNumber","Identification":"80200110203345","Name":"Mr Kevin","SecondaryIdentification":"00021","LEI":"9193001QZMP2PQT4AK86"}]}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		// Test v3 Transactions Permissions tests (105000 & 105200)
		{
			name:       "OB-301-TRA-105000 should pass when no TransactionDetails are present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-TRA-105200 should pass when no TransactionDetails are present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"AccountId": "22289"}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-TRA-105000 should fail when TransactionInformation is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"TransactionInformation": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when Balance is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"Balance": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when MerchantDetails is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"MerchantDetails": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when CreditorAgent is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when CreditorAccount is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when DebtorAgent is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105000 should fail when DebtorAccount is present",
			manifestID: "OB-301-TRA-105000",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when TransactionInformation is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"TransactionInformation": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when Balance is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"Balance": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when MerchantDetails is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"MerchantDetails": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when CreditorAgent is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when CreditorAccount is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when DebtorAgent is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-TRA-105200 should fail when DebtorAccount is present",
			manifestID: "OB-301-TRA-105200",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Transaction":[{"DebtorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-BEN-101800 should pass when no Beneficiary Details are present",
			manifestID: "OB-301-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"AccountId": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-BEN-101800 should fail when CreditorAgent is present",
			manifestID: "OB-301-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-BEN-101800 should fail when CreditorAccount is present",
			manifestID: "OB-301-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-BEN-101900 should pass when no Beneficiary Details are present",
			manifestID: "OB-301-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"AccountId": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-BEN-101900 should fail when CreditorAgent is present",
			manifestID: "OB-301-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-BEN-101900 should fail when CreditorAccount is present",
			manifestID: "OB-301-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-BEN-101800 should pass when no Beneficiary Details are present",
			manifestID: "OB-400-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"AccountId": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-BEN-101800 should fail when CreditorAgent is present",
			manifestID: "OB-400-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-BEN-101800 should fail when CreditorAccount is present",
			manifestID: "OB-400-BEN-101800",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-BEN-101900 should pass when no Beneficiary Details are present",
			manifestID: "OB-400-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"AccountId": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-BEN-101900 should fail when CreditorAgent is present",
			manifestID: "OB-400-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-BEN-101900 should fail when CreditorAccount is present",
			manifestID: "OB-400-BEN-101900",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"Beneficiary":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STO-103801 should pass when no Standing Order Details are present",
			manifestID: "OB-301-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-STO-103801 should fail when CreditorAgent is present",
			manifestID: "OB-301-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STO-103801 should fail when CreditorAccount is present",
			manifestID: "OB-301-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STO-103801 should pass when no Standing Order Details are present",
			manifestID: "OB-400-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-STO-103801 should fail when CreditorAgent is present",
			manifestID: "OB-400-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STO-103801 should fail when CreditorAccount is present",
			manifestID: "OB-400-STO-103801",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STO-103902 should pass when no Standing Order Details are present",
			manifestID: "OB-301-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-STO-103902 should fail when CreditorAgent is present",
			manifestID: "OB-301-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STO-103902 should fail when CreditorAccount is present",
			manifestID: "OB-301-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STO-103902 should pass when no Standing Order Details are present",
			manifestID: "OB-400-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-STO-103902 should fail when CreditorAgent is present",
			manifestID: "OB-400-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STO-103902 should fail when CreditorAccount is present",
			manifestID: "OB-400-STO-103902",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": ""},
				`{"Data":{"StandingOrder":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STA-106400 should pass when no Standing Order Details are present",
			manifestID: "OB-301-STA-106400",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-STA-106400 should fail when CreditorAgent is present",
			manifestID: "OB-301-STA-106400",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"StatementAmount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-STA-106500 should pass when no Standing Order Details are present",
			manifestID: "OB-301-STA-106500",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-STA-106500 should fail when CreditorAgent is present",
			manifestID: "OB-301-STA-106500",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"StatementAmount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STA-106400 should pass when no Standing Order Details are present",
			manifestID: "OB-400-STA-106400",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-STA-106400 should fail when CreditorAgent is present",
			manifestID: "OB-400-STA-106400",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"StatementAmount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-STA-106500 should pass when no Standing Order Details are present",
			manifestID: "OB-400-STA-106500",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-STA-106500 should fail when CreditorAgent is present",
			manifestID: "OB-400-STA-106500",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"Statement":[{"StatementAmount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-SCP-103501 should pass when no Standing Order Details are present",
			manifestID: "OB-301-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-SCP-103501 should fail when CreditorAgent is present",
			manifestID: "OB-301-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-SCP-103501 should fail when CreditorAgent is present",
			manifestID: "OB-301-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-SCP-103501 should pass when no Standing Order Details are present",
			manifestID: "OB-301-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-301-SCP-103601 should fail when CreditorAgent is present",
			manifestID: "OB-301-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-301-SCP-103601 should fail when CreditorAgent is present",
			manifestID: "OB-301-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-SCP-103501 should pass when no Standing Order Details are present",
			manifestID: "OB-400-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-SCP-103501 should fail when CreditorAgent is present",
			manifestID: "OB-400-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-SCP-103501 should fail when CreditorAgent is present",
			manifestID: "OB-400-SCP-103501",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-SCP-103501 should pass when no Standing Order Details are present",
			manifestID: "OB-400-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"AccountID": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: true,
		},
		{
			name:       "OB-400-SCP-103601 should fail when CreditorAgent is present",
			manifestID: "OB-400-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAgent": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
		},
		{
			name:       "OB-400-SCP-103601 should fail when CreditorAgent is present",
			manifestID: "OB-400-SCP-103601",
			response: mockResponse{
				200,
				map[string]string{"x-fapi-interaction-id": "$x-fapi-interaction-id"},
				`{"Data":{"ScheduledPayment":[{"CreditorAccount": ""}]}}`,
			},
			schemaSpec:           *accountSpecPath,
			ExpectValidationPass: false,
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
		errorMessages := &strings.Builder{}
		if result != true {
			for _, e := range errors {
				errorMessages.WriteString(fmt.Sprintf("\t- %s\n", e))
			}
			assert.True(t, errorsContain(errors, test.ExpectValidationErrorContains),
				"%s: validation errors should contain an error with text: '%s'\nERRORS:\n%v", test.name, test.ExpectValidationErrorContains, errorMessages)
		}

		assert.Equal(t, test.ExpectValidationPass, result, "%s result - expected: %v actual: %v\nCaught error(s):\n%s", test.name, test.ExpectValidationPass, result, errorMessages)
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
