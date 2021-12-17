package generation

// The purpose of the generation package is to initially to explore the options around test case generation.
// As such the code is necessarily 'experimental' and subject to change.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

var httpUserAgent string

func init() {
	humanVersion := version.NewBitBucket("").GetHumanVersion()
	httpUserAgent = fmt.Sprintf("Open Banking Conformance Suite %s", humanVersion)
}

// GetImplementedTestCases takes a discovery Model and determines the implemented endpoints.
// Currently this function is experimental - meaning it contains fmt.Printlns as an aid to understanding
// and conceptualisation
func GetImplementedTestCases(disco *discovery.ModelDiscoveryItem, nameGenerator names.Generator, ctx *model.Context, headlessTokenAcquisition bool, genConfig GeneratorConfig) ([]model.TestCase, map[string]string) {
	logger := logrus.StandardLogger()
	originalEndpoints := make(map[string]string)
	var testcases []model.TestCase
	endpoints := disco.Endpoints
	doc, err := loadSpec(disco.APISpecification.SchemaVersion, false)
	if err != nil {
		logger.Errorln(err)
		return nil, nil
	}

	for _, v := range endpoints {
		var responseCodes []int
		var goodResponseCode int
		newpath := getResourceIds(disco, v.Path, genConfig)

		for path, props := range doc.Spec().Paths.Paths {
			for meth, op := range getOperations(&props) {
				if (meth == v.Method) && (v.Path == path) {
					responseCodes = getResponseCodes(op)
					goodResponseCode, err = getGoodResponseCode(responseCodes)
					if err != nil {
						logger.WithFields(logrus.Fields{
							"testcase": op.Summary,
							"method":   meth,
							"endpoint": newpath,
							"err":      err,
						}).Error("Cannot get good response code")
						return nil, nil
					}

					headers := map[string]string{
						"Authorization":         "Bearer $access_token",
						"X-Fapi-Financial-Id":   "$fapi_financial_id",
						"X-Fapi-Interaction-Id": "b4405450-febe-11e8-80a5-0fcebb1574e1",
						"Content-Type":          "application/json",
						"User-Agent":            httpUserAgent,
						"Accept":                "*/*",
					}

					if strings.Contains(newpath, "account-access-consents") { // consent endpoints require a different access_token + custom chain
						headers["Authorization"] = "Bearer $client_access_token"
						customTestCases, err := getTemplatedTestCases(newpath)
						if err != nil {
							logger.WithFields(logrus.Fields{
								"testcase": op.Summary,
								"method":   meth,
								"endpoint": newpath,
								"err":      err,
							}).Warn("error getting Templated TestCase")
							return nil, nil
						}
						for i := range customTestCases {
							showReplacementErrors := false
							customTestCases[i].ProcessReplacementFields(ctx, showReplacementErrors)
						}
						if customTestCases != nil {
							testcases = append(testcases, customTestCases...)
						}
						continue
					}

					input := model.Input{Method: meth, Endpoint: newpath, Headers: headers}
					expect := model.Expect{StatusCode: goodResponseCode, SchemaValidation: true}
					context := model.Context{"baseurl": disco.ResourceBaseURI}
					testcase := model.TestCase{ID: nameGenerator.Generate(), Input: input, Context: context, Expect: expect, Name: op.Summary}
					if !headlessTokenAcquisition {
						showReplacementErrors := false
						testcase.ProcessReplacementFields(ctx, showReplacementErrors)
					}
					originalEndpoints[testcase.ID] = v.Path // capture original spec paths
					testcases = append(testcases, testcase)
					break
				}
			}
		}
	}
	return testcases, originalEndpoints
}

func getTemplatedTestCases(path string) (tc []model.TestCase, err error) {
	if !isAccountAccessConsentEndpoint(path) {
		return tc, nil
	}

	filedata, err := ioutil.ReadFile("components/account_consent.json")
	if err != nil {
		filedata, err = ioutil.ReadFile("../../components/account_consent.json") // handle testing
		if err != nil {
			logrus.StandardLogger().Error("Cannot read: components/account_consent " + err.Error())
			return nil, err
		}
	}
	testcases := []model.TestCase{}
	err = json.Unmarshal(filedata, &testcases)
	if err != nil {
		return testcases, err
	}
	return testcases, nil
}

// GetCustomTestCases retrieves custom tests from the discovery file
func GetCustomTestCases(discoReader *discovery.CustomTest, ctx *model.Context, headlessTokenAcquisition bool) SpecificationTestCases {
	spec := discovery.ModelAPISpecification{Name: discoReader.Name}
	specTestCases := SpecificationTestCases{Specification: spec}
	testcases := []model.TestCase{}
	for _, testcase := range discoReader.Sequence {
		if !headlessTokenAcquisition {
			showReplacementErrors := false
			testcase.ProcessReplacementFields(ctx, showReplacementErrors)
		}
		testcases = append(testcases, testcase)
	}
	specTestCases.TestCases = testcases
	return specTestCases
}

// check if a response code is in the range 200-299 - therefore a 'good' response code
func getGoodResponseCode(codes []int) (int, error) {
	for _, i := range codes {
		if i > 199 && i < 300 {
			return i, nil
		}
	}
	return 0, errors.New("Cannot find good response code between 200 and 299")
}

// given an operation specification, return all resultcodes for that operation
func getResponseCodes(op *spec.Operation) (result []int) {
	for i := range op.OperationProps.Responses.ResponsesProps.StatusCodeResponses {
		result = append(result, i)
	}
	return
}

// helper to replace path name resource ids specified between brackets e.g. `{AccountId}`
// with the values "ResourceIds" section of the discovery model
func getResourceIds(item *discovery.ModelDiscoveryItem, path string, genConfig GeneratorConfig) string {
	result := path
	for k, v := range item.ResourceIds {
		key := strings.Join([]string{"{", k, "}"}, "")
		result = strings.Replace(result, key, v, 1)
	}

	// Update the account ids in based on the discovery configuration
	if len(genConfig.ResourceIDs.AccountIDs) > 0 {
		// At the moment, according to requirements, we only need support the first ID.
		logrus.StandardLogger().Warn("Using the {AccountId} value at index 0 - ignoring others")

		v := genConfig.ResourceIDs.AccountIDs[0]
		result = strings.Replace(result, "{AccountId}", v.AccountID, 1)
	}
	// Update the statement ids in based on the discovery configuration
	if len(genConfig.ResourceIDs.StatementIDs) > 0 {
		// At the moment, according to requirements, we only need support the first ID.
		logrus.StandardLogger().Warn("Using the {StatementId} value at index 0 - ignoring others")

		v := genConfig.ResourceIDs.StatementIDs[0]
		result = strings.Replace(result, "{StatementId}", v.StatementID, 1)
	}

	return result
}

// loads an openapi specification via http or file
func loadSpec(spec string, print bool) (*loads.Document, error) {
	doc, err := loads.Spec(spec)
	if err != nil {
		return nil, err
	}
	if print {
		var jsondoc []byte
		jsondoc, err = json.MarshalIndent(doc.Spec(), "", "    ")
		if err != nil {
			return nil, err
		}

		fmt.Println(string(jsondoc))
	}
	return doc, err
}

// Utilities to walk the swagger tree
// getOperations returns a mapping of HTTP Verb name to "spec operation name"
func getOperations(props *spec.PathItem) map[string]*spec.Operation {
	ops := map[string]*spec.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}

	// Keep those != nil
	for key, op := range ops {
		if op == nil {
			delete(ops, key)
		}
	}
	return ops
}
