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

func getValidatorManifestScriptsSpec(manifest spec) (schema.OpenAPI3Validator, Scripts, References, discovery.ModelAPISpecification, error) {
	validator, err := schema.NewRawOpenAPI3Validator(manifest.Name, manifest.Version)
	if err != nil {
		return validator, Scripts{}, References{}, discovery.ModelAPISpecification{}, err
	}

	specType, err := GetSpecType(manifest.SchemaVersion)
	if err != nil {
		return validator, Scripts{}, References{}, discovery.ModelAPISpecification{}, err
	}

	var values []interface{}
	values = append(values, "accounts_v3.1.1", "payments_v3.1.1", "cbpii_v3.1.1", "vrps_v3.1.1")
	ctx := model.Context{"apiversions": values}

	manifestScripts, refs, err := LoadGenerationResources(specType, manifest.Manifest, &ctx)
	refs = completeReferences(refs)

	spec := discovery.ModelAPISpecification{
		Name:          manifest.Name,
		URL:           manifest.URL,
		Version:       manifest.Version,
		SchemaVersion: manifest.SchemaVersion,
		Manifest:      manifest.Manifest,
		SpecType:      specType,
	}

	return validator, manifestScripts, refs, spec, err
}

// TestValidateRequest - Test manfest files against sqaggers(validators)
func TestValidateRequest(t *testing.T) {
	for _, manifest := range specs {
		validator, manifestScripts, refs, spec, err := getValidatorManifestScriptsSpec(manifest)
		assert.NoError(t, err)

		var values []interface{}
		values = append(values, "accounts_v3.1.1", "payments_v3.1.1", "cbpii_v3.1.1", "vrps_v3.1.1")
		ctx := model.Context{"apiversions": values}

		for _, script := range manifestScripts.Scripts {
			localCtx, err := script.processParameters(&refs, &ctx)
			assert.NoError(t, err)

			tc, err := buildTestCase(script, refs.References, localCtx, baseurl, spec.SpecType, validator, spec)
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

func getManifestEndpoints(manifestScripts Scripts) map[string]map[string]bool {
	endpoints := make(map[string]map[string]bool)
	for _, s := range manifestScripts.Scripts {
		if endpoints[s.URI] == nil {
			endpoints[s.URI] = make(map[string]bool)
		}
		endpoints[s.URI][s.Method] = true
	}
	return endpoints
}

// simplifyUri - replace "{id}" and "$id" with "#"
func simplifyUri(uri, start, end string) string {
	stardIdx := strings.Index(uri, start)
	if stardIdx < 0 {
		return uri
	}
	endIdx := strings.Index(uri[stardIdx:], end)
	// if id field is the last one it doesn't contain / at the end
	if endIdx < 0 {
		endIdx = len(uri) - 1
	} else {
		if end == "/" {
			endIdx -= 1
		}
		endIdx += stardIdx
	}
	subStr := uri[stardIdx : endIdx+1]
	return simplifyUri(strings.Replace(uri, subStr, "#", 1), start, end)
}

func TestSsimplifyUri(t *testing.T) {
	simplifyUri("/accounts/$accountId/party", "$", "/")
	simplifyUri("/accounts", "$", "/")
	simplifyUri("/accounts/$accountId", "$", "/")
}

func simplifyEndpoints(endpoints map[string]map[string]bool, start, end string) map[string]map[string]bool {
	simplifiedEndpoints := make(map[string]map[string]bool)
	for uri, methods := range endpoints {
		newUri := simplifyUri(uri, start, end)
		simplifiedEndpoints[newUri] = methods
	}
	return simplifiedEndpoints
}

func findMatchingEndpoint(key string, manifestEndpoints map[string]map[string]bool) (map[string]bool, error) {
	methods := manifestEndpoints[key]
	if methods == nil {
		return nil, fmt.Errorf("%s not found in manifest", key)
	}
	return methods, nil
}

// validateMethods - validate Validator Methods against Manifest Methods. It detects missing methods in the Manifest.
func validateMethods(validatorMethods, manifestMethods map[string]bool) map[string]bool {
	missingMethods := make(map[string]bool)

	for validatorMethod, _ := range validatorMethods {
		if !manifestMethods[validatorMethod] {
			missingMethods[validatorMethod] = true
		}
	}

	return missingMethods
}

// validateEndpoints - validate Validator Endpoints against Manifest Endpoints. It detects missing tests in the Manifest.
func validateEndpoints(validatorEndpoints, manifestEndpoints map[string]map[string]bool) map[string]map[string]bool {
	simplifiedManifestEndpoints := simplifyEndpoints(manifestEndpoints, "$", "/")
	missingEndpoints := make(map[string]map[string]bool)

	for uri, methods := range validatorEndpoints {
		newUri := simplifyUri(uri, "{", "}")
		manifestMethods, err := findMatchingEndpoint(newUri, simplifiedManifestEndpoints)
		if err != nil {
			missingEndpoints[uri] = methods
		}
		missingMethods := validateMethods(methods, manifestMethods)
		if len(missingMethods) > 0 {
			missingEndpoints[uri] = missingMethods
		}
	}
	return missingEndpoints
}

func mapToSlice(boolMap map[string]bool) []string {
	slice := []string{}
	for k, _ := range boolMap {
		slice = append(slice, k)
	}
	return slice
}

func prepareMissingTestsMsg(missingTests map[string]map[string]map[string]bool) string {
	var msg strings.Builder
	for name, endpoints := range missingTests {
		if len(endpoints) > 0 {
			msg.WriteString(name)
			msg.WriteString("/n")
			for enpoint, methods := range endpoints {
				methodSlice := mapToSlice(methods)
				msg.WriteString(fmt.Sprintf("%s - %s/n", enpoint, strings.Join(methodSlice, ", ")))
			}
		}
	}
	return msg.String()
}

func TestValidateEndpointsAmount(t *testing.T) {
	missingTests := make(map[string]map[string]map[string]bool)
	for _, manifest := range specs {
		validator, manifestScripts, _, _, err := getValidatorManifestScriptsSpec(manifest)
		assert.NoError(t, err)
		validatorEndpoints := validator.GetEndpoints()
		manifestEndpoints := getManifestEndpoints(manifestScripts)
		missingTests[manifest.Name] = validateEndpoints(validatorEndpoints, manifestEndpoints)
	}
	msg := prepareMissingTestsMsg(missingTests)
	t.Log(msg)

}
