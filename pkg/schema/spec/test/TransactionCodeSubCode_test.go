package assertionstest

import (
	"flag"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/assert"
)

var (
	accountSpecPath316 = flag.String("spec316", "../v3.1.6/account-info-swagger-flattened.json", "Path to the specification swagger file.")
	accountSpecPath315 = flag.String("spec315", "../v3.1.5/account-info-swagger-flattened.json", "Path to the specification swagger file.")
	accountSpecPath314 = flag.String("spec314", "../v3.1.4/account-info-swagger-flattened.json", "Path to the specification swagger file.")
	accountSpecPath313 = flag.String("spec313", "../v3.1.5/account-info-swagger-flattened.json", "Path to the specification swagger file.")
)

// Tests to check for issues with Code and Subcode swagger field lengths
// Fields should be unrestricted, not limited to 4 Characters
// See REFAPP-1083
func TestBankTransactionCodeSubCode(t *testing.T) {
	var err error
	emptyContext := &model.Context{}

	testCase := model.TestCase{
		Input:  model.Input{Method: "GET", Endpoint: "/open-banking/v3.1/aisp/accounts/00000005/transactions"},
		Expect: model.Expect{SchemaValidation: true}}

	t.Run("Transaction with Codes and Subcodes more than 4 chars should PASS 3.1.6 accounts", func(t *testing.T) {
		testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath316)
		if err != nil {
			t.Fatal(err)
		}

		resp := test.CreateHTTPResponse(200, "OK", string(transactionData), "Content-Type", "application/json")
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Transaction with Codes and Subcodes more than 4 chars should PASS 3.1.5 accounts", func(t *testing.T) {
		testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath315)
		if err != nil {
			t.Fatal(err)
		}

		resp := test.CreateHTTPResponse(200, "OK", string(transactionData), "Content-Type", "application/json")
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Transaction with Codes and Subcodes more than 4 chars should PASS 3.1.4 accounts", func(t *testing.T) {
		testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath314)
		if err != nil {
			t.Fatal(err)
		}

		resp := test.CreateHTTPResponse(200, "OK", string(transactionData), "Content-Type", "application/json")
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

	t.Run("Transaction with Codes and Subcodes more than 4 chars should PASS 3.1.3 accounts", func(t *testing.T) {
		testCase.Validator, err = schema.NewSwaggerValidator(*accountSpecPath313)
		if err != nil {
			t.Fatal(err)
		}

		resp := test.CreateHTTPResponse(200, "OK", string(transactionData), "Content-Type", "application/json")
		result, err := testCase.Validate(resp, emptyContext)
		if len(err) != 0 {
			t.Fatal(err)
		}
		assert.True(t, result, "expected: %v actual: %v", true, result)
	})

}

var transactionData = []byte(`{
	"Data": {
	   "Transaction": [
		  {
			 "AccountId": "700004000000000000000003",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "19.88",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "e11bc082-f680-4101-8509-70a0254e811a",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-24T22:43:48.901Z",
			 "ValueDateTime": "2019-09-24T22:43:48.901Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.5",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  }
	   ]
	},
	"Links": {
	   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/transactions"
	},
	"Meta": {
	   "TotalPages": 1,
	   "FirstAvailableDateTime": "2016-01-01T08:40:00.000Z",
	   "LastAvailableDateTime": "2025-12-31T08:40:00.000Z"
	}
 }`)
