package generation

// The purpose of the generation package is to initially to explore the options around test case generation.
// As such the code is necessarily 'experimental' and subject to change.

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
func GetImplementedTestCases(disco *discovery.ModelDiscoveryItem, print bool, beginTestNo int) []model.TestCase {
	var testcases []model.TestCase
	endpoints := disco.Endpoints
	testNo := beginTestNo
	doc, err := loadSpec(disco.APISpecification.SchemaVersion, false)
	if err != nil {
		fmt.Println(err)
		return testcases
	}

	for _, v := range endpoints {
		var responseCodes []int
		var goodResponseCode int
		condition := getConditionality(v.Method, v.Path, disco.APISpecification.SchemaVersion)
		newpath := replaceResourceIds(disco, v.Path)
		if print {
			fmt.Printf("[%s] %s %s\n", condition, v.Method, newpath)
		}

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
					input := model.Input{Method: meth, Endpoint: newpath}
					expect := model.Expect{StatusCode: goodResponseCode, SchemaValidation: true}
					testcase := model.TestCase{ID: fmt.Sprintf("#t%4.4d", testNo), Input: input, Expect: expect, Name: op.Summary}
					testcases = append(testcases, testcase)
					testNo++
					break
				}
			}
		}
	}

	return testcases
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
func getResponseCodes(op *spec.Operation) []int {
	var result []int
	for i := range op.OperationProps.Responses.ResponsesProps.StatusCodeResponses {
		result = append(result, i)
	}
	return result
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

// helper to walk through and display information from discovery
func printImplemented(ditem discovery.ModelDiscoveryItem, endpoints []discovery.ModelEndpoint, spec string) {
	for _, v := range endpoints {
		condition := getConditionality(v.Method, v.Path, spec)
		newpath := replaceResourceIds(&ditem, v.Path)
		fmt.Printf("[%s] %s %s\n", condition, v.Method, newpath)
	}
}

// helper to replace path name resource ids specificed between brackets e.g. `{AccountId}`
// with the values "ResourceIds" section of the discovery model
func replaceResourceIds(item *discovery.ModelDiscoveryItem, path string) string {
	newstr := path
	for k, v := range item.ResourceIds {
		key := strings.Join([]string{"{", k, "}"}, "")
		newstr = strings.Replace(newstr, key, v, 1)
	}
	return newstr
}

// help to dump out resourceIds to console
func printResourceIds(item *discovery.ModelDiscoveryItem) {
	for k, v := range item.ResourceIds {
		fmt.Printf("%-30.30s %s\n", k, v)
	}
}

func printSpec(doc *loads.Document, base, spec string) {
	for path, props := range doc.Spec().Paths.Paths {
		for method := range getOperations(&props) {
			newPath := base + path
			condition := getConditionality(method, path, spec)
			fmt.Printf("[%s] %s %s\n", condition, method, newPath) // give to testcase along with any conditionality?
		}
	}
}

// loads an openapi specification via http or file
func loadSpec(spec string, print bool) (*loads.Document, error) {
	doc, err := loads.Spec(spec)
	if print {
		var jsondoc []byte
		jsondoc, _ = json.MarshalIndent(doc.Spec(), "", "    ")
		fmt.Println(string(jsondoc))
	}
	return doc, err
}

// Utility to load Manifest Data Model containing all Rules, Tests and Conditions
func loadModelOBv3Ozone() (discovery.Model, error) {
	filedata, _ := ioutil.ReadFile("testdata/disco.json")
	var d discovery.Model
	err := json.Unmarshal(filedata, &d)
	if err != nil {
		return discovery.Model{}, err
	}
	return d, nil
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
