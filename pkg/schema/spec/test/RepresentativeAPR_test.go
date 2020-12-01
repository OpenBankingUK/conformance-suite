package assertionstest

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/schema"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestRepresentativeAPR(t *testing.T) {
	var err error
	productsResponse200 := `
	{
		"Data": {
			"Product": [
				{
					"AccountId": "22392",
					"ProductType": "BusinessCurrentAccount",
					"ProductName": "Barclays Business Current Account",
					"BCA": {
						"Overdraft": {
							"OverdraftTierBandSet": [
								{
									"TierBandMethod": "Tiered",
									"OverdraftTierBand": [
										{
											"TierValueMin": "1.0000",
											"RepresentativeAPR": "%v"
										}
									]
								}
							]
						}
					}
				}
			]
		}
	}`

	t.Run("Schema validation should PASS using response with correctly formatted RepresentativeAPR value.", func(t *testing.T) {
		testCase := model.TestCase{
			Input:  model.Input{Method: "GET", Endpoint: "/open-banking/v3.1/aisp/products"},
			Expect: model.Expect{SchemaValidation: true}}

		emptyContext := &model.Context{}
		response := fmt.Sprintf(productsResponse200, 42.42)
		for _, specPath := range []string{*accountSpecPath313, *accountSpecPath314, *accountSpecPath315, *accountSpecPath316} {
			testCase.Validator, err = schema.NewSwaggerValidator(specPath)
			if err != nil {
				t.Fatal(err)
			}

			resp := test.CreateHTTPResponse(200, "OK", response, "Content-Type", "application/json")
			valisationDone, errs := testCase.Validate(resp, emptyContext)
			if len(errs) != 0 {
				t.Fatal(errs)
			}
			assert.True(t, valisationDone, "expected: validated=%v actual: validated=%v", true, valisationDone)
		}
	})

	t.Run("Schema validation should FAIL using response with incorrectly formatted RepresentativeAPR value.", func(t *testing.T) {
		testCase := model.TestCase{
			Input:  model.Input{Method: "GET", Endpoint: "/open-banking/v3.1/aisp/products"},
			Expect: model.Expect{SchemaValidation: true}}

		emptyContext := &model.Context{}
		invalidValue := 4242.4242 // added an extra digit: must be +/-ddd.dddd
		response := fmt.Sprintf(productsResponse200, invalidValue)
		for _, specPath := range []string{*accountSpecPath313, *accountSpecPath314, *accountSpecPath315, *accountSpecPath316} {
			testCase.Validator, err = schema.NewSwaggerValidator(specPath)
			if err != nil {
				t.Fatal(err)
			}

			resp := test.CreateHTTPResponse(200, "OK", response, "Content-Type", "application/json")
			validationDone, errs := testCase.Validate(resp, emptyContext)
			assert.Len(t, errs, 1) // expect one error
			assert.True(t, validationDone, "expected: validated=%v actual: validated=%v", true, validationDone)
		}
	})
}
