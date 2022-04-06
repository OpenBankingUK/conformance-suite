package manifest

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/schema"
	"github.com/stretchr/testify/assert"
)

type spec struct {
	Name          string
	Version       string
	SchemaVersion string
	Manifest      string
	URL           string
}

var (
	specs = []spec{
		// VRP
		{
			Name:          "OBIE VRP Profile",
			Version:       "v3.1.9",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.9/dist/openapi/payment-initiation-openapi.json",
			Manifest:      "file://manifests/ob_3.1_variable_recurring_payments.json",
			URL:           "https://openbankinguk.github.io/read-write-api-site3/v3.1.9/profiles/payment-initiation-api-profile.html",
		},
		// AIS
		{
			Name:          "Account and Transaction API Specification",
			Version:       "v3.1.9",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.9/dist/openapi/account-info-openapi.json",
			Manifest:      "file://manifests/ob_3.1_accounts_transactions_fca.json",
			URL:           "https://openbankinguk.github.io/read-write-api-site3/v3.1.9/profiles/account-and-transaction-api-profile.html",
		},
		// PIS
		{
			Name:          "Payment Initiation API",
			Version:       "v3.1.9",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.9/dist/swagger/payment-initiation-swagger.json",
			Manifest:      "file://manifests/ob_3.1_payment_fca.json",
			URL:           "https://openbankinguk.github.io/read-write-api-site3/v3.1.9/profiles/payment-initiation-api-profile.html",
		},
		// CBPII
		{
			Name:          "Confirmation of Funds API Specification",
			Version:       "v3.1.9",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.9/dist/swagger/confirmation-funds-swagger.json",
			Manifest:      "file://manifests/ob_3.1_cbpii_fca.json",
			URL:           "https://openbankinguk.github.io/read-write-api-site3/v3.1.9/profiles/confirmation-of-funds-api-profile.html",
		},
	}

	fieldReplacements = map[string]string{
		"$transactionFromDate":        "2017-06-05T15:15:13+00:00",
		"$transactionToDate":          "2020-06-05T15:15:13+00:00",
		"$instructedAmountCurrency":   "GBP",
		"$instructedAmountValue":      "100.00",
		"$requestedExecutionDateTime": "2017-06-05T15:15:22+00:00",
		"$firstPaymentDateTime":       "2017-06-05T15:15:22+00:00",
		"$currencyOfTransfer":         "GBP",
		"$frequency":                  "EvryDay",
		"$expirationDateTime":         "2017-06-05T15:15:22+00:00",
	}

	baseurl = "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/pisp"
)

func completeReferences(refs References) References {
	fixedRefs := References{
		References: map[string]Reference{},
	}
	for i, v := range refs.References {
		v.BodyData = completeBodyFields(v.BodyData)
		fixedRefs.References[i] = v
	}
	return fixedRefs
}

func completeBodyFields(body string) string {
	for oldField, newField := range fieldReplacements {
		body = strings.ReplaceAll(body, oldField, newField)
	}
	return body
}

// TestValidateRequest - Test manfest files against sqaggers(validators)
func TestValidateRequest(t *testing.T) {
	for _, manifest := range specs {
		validator, err := schema.NewRawOpenAPI3Validator(manifest.Name, manifest.Version)
		assert.NoError(t, err)

		specType, err := GetSpecType(manifest.SchemaVersion)

		var values []interface{}
		values = append(values, "accounts_v3.1.1", "payments_v3.1.1", "cbpii_v3.1.1", "vrps_v3.1.1")
		ctx := model.Context{"apiversions": values}

		manifestScripts, refs, err := LoadGenerationResources(specType, manifest.Manifest, &ctx)
		refs = completeReferences(refs)
		assert.NoError(t, err)

		spec := discovery.ModelAPISpecification{
			Name:          manifest.Name,
			URL:           manifest.URL,
			Version:       manifest.Version,
			SchemaVersion: manifest.SchemaVersion,
			Manifest:      manifest.Manifest,
			SpecType:      specType,
		}

		for _, script := range manifestScripts.Scripts {
			localCtx, err := script.processParameters(&refs, &ctx)
			assert.NoError(t, err)

			tc, err := buildTestCase(script, refs.References, localCtx, baseurl, specType, validator, spec)
			assert.NoError(t, err)
			localCtx.PutContext(&ctx)

			// skip fails and StatusCode == 0 (means no expected value)
			if tc.Expect.StatusCode >= 400 || tc.Expect.StatusCode == 0 {
				continue
			}

			tc.Header = http.Header{}
			for headerKey, header_value := range tc.Input.Headers {
				tc.Header[headerKey] = []string{header_value}
			}

			req := schema.HTTPRequest{
				Method:  tc.Input.Method,
				Path:    tc.Input.Endpoint,
				Header:  tc.Header,
				StrBody: tc.Input.RequestBody,
			}

			err = validator.ValidateRequest(req)
			if err != nil {
				err = fmt.Errorf("%s, %s: %s", manifest.Name, script.ID, err.Error())
			}
			assert.NoError(t, err)
		}
	}
}
