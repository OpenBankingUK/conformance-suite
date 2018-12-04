package generation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"errors"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/utils"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumerateOpenApiTestcases(t *testing.T) {
	dmodel, err := loadModelOBv3Ozone()
	require.NoError(t, err)
	for _, dItem := range dmodel.DiscoveryModel.DiscoveryItems {
		fmt.Printf("\n=========================================\n%s\n=========================================", dItem.APISpecification.Name)
		fmt.Printf("\n%s\n--------------\n", dItem.APISpecification.Version)
		doc, err := loadSpec(dItem.APISpecification.SchemaVersion, false)
		require.NoError(t, err)
		printSpec(doc, dItem.ResourceBaseURI, dItem.APISpecification.Version) // print the endpoints in the spec
		fmt.Printf("\nResourceIds\n-----------\n")
		printResourceIds(&dItem)
		fmt.Printf("\nImplemented\n--------------\n")
		printImplemented(dItem, dItem.Endpoints, dItem.APISpecification.Version) // print what this org has implemeneted
		_, _ = doc, err
		break
	}
}

func TestGenerateTestCases(t *testing.T) {
	results := []model.TestCase{}
	disco, _ := loadModelOBv3Ozone()

	for _, v := range disco.DiscoveryModel.DiscoveryItems {
		result := getImplementedTestCases(&v)
		results = append(results, result...)
	}

	fmt.Println("Dumping test cases")
	for _, tc := range results {
		fmt.Println(string(pkgutils.DumpJSON(&tc)))
	}
}

// func printSpec(doc *loads.Document, base, spec string) {
// for path, props := range doc.Spec().Paths.Paths {

func getImplementedTestCases(disco *discovery.ModelDiscoveryItem) []model.TestCase {
	var testcases []model.TestCase
	var testNo = 1000
	endpoints := disco.Endpoints
	doc, err := loadSpec(disco.APISpecification.SchemaVersion, false)	
	if err != nil {
		fmt.Println(err)
		return testcases
	}

	// sp := doc.Spec()
	// fmt.Printf("%s\n============================================>>>>>>>>>>>>>\n",sp.Info.Title)
	// fmt.Println(sp.Info.Description)
	// fmt.Printf("%#v\n",sp.Parameters)
	// if true {
	// 	return testcases
	// }
	

	for _, v := range endpoints {
		var responseCodes []int
		var goodResponseCode int
		condition := getConditionality(v.Method, v.Path, disco.APISpecification.SchemaVersion)
		newpath := ReplaceResourceIds(disco, v.Path)
		fmt.Printf("[%s] %s %s\n", condition, v.Method, newpath)

		 for path, props := range doc.Spec().Paths.Paths {
			for meth,op := range getOperations(&props) {
				if (meth == v.Method) && (v.Path == path) {
					responseCodes = getResponseCodes(op)					
					fmt.Printf("Response Codes, %v\n", responseCodes)
					goodResponseCode,_ = getGoodResponseCode(responseCodes)
					fmt.Printf("Good Response Code %d",goodResponseCode)

					opProps := op.OperationProps
					fmt.Printf("\nOpProps: %#v\n",opProps)
					fmt.Printf("\nID: %#v\n",opProps.ID)
					fmt.Printf("Summary: %s\n",opProps.Summary)
					break
				}
			}			  
		 }

		// doc, err := loadSpec(disco.APISpecification.SchemaVersion, false)
		// for meth, op := range getOperations(&props) {
		// 	successStatus := 0
		// 	testNo++

		// 	for i := range op.OperationProps.Responses.ResponsesProps.StatusCodeResponses {
		// 		if i > 199 && i < 300 {
		// 			successStatus = i
		// 		}
		// 	}

		// 	input := model.Input{Method: meth, Endpoint: newpath}
		// 	expect := model.Expect{StatusCode: successStatus, SchemaValidation: true}
		// 	testcase := model.TestCase{ID: fmt.Sprintf("#t%4.4d", testNo), Input: input, Expect: expect, Name: op.Description}
		// 	testcases = append(testcases, testcase)
		// }
	}

	_ = testNo

	return testcases
}

func getGoodResponseCode(codes []int) (int, error) {
	for _,i := range codes {
		if i > 199 && i < 300 {
			return i,nil
		}
	}
	return 0, errors.New("Cannot find good response code between 200 and 299")
}


// given an operation specification, return all resultcodes for that operation
func getResponseCodes(op *spec.Operation) []int {
	var result []int
	for i := range op.OperationProps.Responses.ResponsesProps.StatusCodeResponses {
		result = append(result,i)
	}
	return result
}



func getConditionality(method, path, specification string) (string) {
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

func printImplemented(ditem discovery.ModelDiscoveryItem, endpoints []discovery.ModelEndpoint, spec string) {
	for _, v := range endpoints {
		condition := getConditionality(v.Method, v.Path, spec)
		newpath := ReplaceResourceIds(&ditem, v.Path)
		fmt.Printf("[%s] %s %s\n", condition, v.Method, newpath)
	}
}

func TestIt(t *testing.T) {
	dmodel, _ := loadModelOBv3Ozone()
	fmt.Println(string(pkgutils.DumpJSON(&dmodel)))
	for id, dItem := range dmodel.DiscoveryModel.DiscoveryItems {
		fmt.Println("dim")
		for _, ep := range dItem.Endpoints {
			fmt.Println("ep: ", ep)
			ReplaceResourceIds(&dItem, ep.Path)
			_ = id
		}
	}
	assert.True(t, false)
}

func ReplaceResourceIds(item *discovery.ModelDiscoveryItem, path string) string {
	newstr := path
	for k, v := range item.ResourceIds {
		key := strings.Join([]string{"{", k, "}"}, "")
		newstr = strings.Replace(newstr, key, v, 1)
	}
	return newstr
}

func printResourceIds(item *discovery.ModelDiscoveryItem) {
	for k, v := range item.ResourceIds {
		fmt.Printf("%-20.20s %s\n", k, v)
	}
}

// match disco endpoints to spec
// expect disco to be subset of spec
// figure out the differences

func printSpec(doc *loads.Document, base, spec string) {
	for path, props := range doc.Spec().Paths.Paths {
		for method := range getOperations(&props) {
			newPath := base + path
			condition := getConditionality(method, path, spec)
			fmt.Printf("[%s] %s %s\n", condition, method, newPath) // give to testcase along with any conditionality?
			// map disco
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
