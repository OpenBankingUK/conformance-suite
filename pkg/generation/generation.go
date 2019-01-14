package generation

// The purpose of the generation package is to initially to explore the options around test case generation.
// As such the code is necessarily 'experimental' and subject to change.

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// GetImplementedTestCases takes a discovery Model and determines the implemented endpoints.
// Currently this function is experimental - meaning it contains fmt.Printlns as an aid to understanding
// and conceptualisation
func GetImplementedTestCases(disco *discovery.ModelDiscoveryItem, beginTestNo int) []model.TestCase {
	var testcases []model.TestCase
	endpoints := disco.Endpoints
	testNo := beginTestNo
	doc, err := loadSpec(disco.APISpecification.SchemaVersion, false)
	if err != nil {
		logrus.Errorln(err)
		return testcases
	}

	for _, v := range endpoints {
		var responseCodes []int
		var goodResponseCode int
		newpath := getResourceIds(disco, v.Path)

		for path, props := range doc.Spec().Paths.Paths {
			for meth, op := range getOperations(&props) {
				if (meth == v.Method) && (v.Path == path) {
					responseCodes = getResponseCodes(op)
					goodResponseCode, err = getGoodResponseCode(responseCodes)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"testcase": op.Summary,
							"method":   meth,
							"endpoint": newpath,
							"err":      err,
						}).Warn("Cannot get good response code")
						continue
					}
					headers := map[string]string{
						"authorization":         "Bearer $access_token",
						"X-Fapi-Financial-Id":   "0015800001041RHAAY",
						"X-Fapi-Interaction-Id": "b4405450-febe-11e8-80a5-0fcebb1574e1",
						"Content-Type":          "application/json",
						"User-Agent":            "Open Banking Conformance Suite v0.2.0-alpha",
						"Accept":                "*/*",
					}
					input := model.Input{Method: meth, Endpoint: newpath, Headers: headers}
					expect := model.Expect{StatusCode: goodResponseCode, SchemaValidation: true}
					context := model.Context{"baseurl": disco.ResourceBaseURI}
					testcase := model.TestCase{ID: fmt.Sprintf("#t%4.4d", testNo), Input: input, Context: context, Expect: expect, Name: op.Summary}
					testcases = append(testcases, testcase)
					testNo++
					break
				}
			}
		}
	}

	return testcases
}

// GetCustomTestCases retrieves custom tests from the discovery file
func GetCustomTestCases(discoReader *discovery.CustomTest) SpecificationTestCases {
	spec := discovery.ModelAPISpecification{Name: discoReader.Name}
	specTestCases := SpecificationTestCases{Specification: spec}
	testcases := []model.TestCase{}
	for _, testcase := range discoReader.Sequence {
		testcase.ProcessReplacementFields(discoReader.Replacements)
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

// helper to annotate generation routines with conditionality inidicator
func getConditionality(method, path, specification string) string {
	condition, err := model.GetConditionality(method, path, specification)
	if err != nil {
		return "U"
	}
	switch condition {
	case model.Mandatory:
		return "M"
	case model.Conditional:
		return "C"
	case model.Optional:
		return "O"
	default:
		return "U"
	}
}

// helper to replace path name resource ids specificed between brackets e.g. `{AccountId}`
// with the values "ResourceIds" section of the discovery model
func getResourceIds(item *discovery.ModelDiscoveryItem, path string) string {
	newstr := path
	for k, v := range item.ResourceIds {
		key := strings.Join([]string{"{", k, "}"}, "")
		newstr = strings.Replace(newstr, key, v, 1)
	}
	return newstr
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
