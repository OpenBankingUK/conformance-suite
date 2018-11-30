package generation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/stretchr/testify/require"
)

// Enumerates all OpenAPI calls from swagger file
func TestEnumerateOpenApiTestcases(t *testing.T) {
	dmodel, err := loadModelOBv3Ozone()
	require.NoError(t, err)
	for _, dItem := range dmodel.DiscoveryModel.DiscoveryItems {
		fmt.Printf("\n=========================================\n%s\n=========================================",dItem.APISpecification.Name)
		fmt.Printf("\n%s\n--------------\n", dItem.APISpecification.Version)
		doc, err := loadSpec(dItem.APISpecification.SchemaVersion, false)
		require.NoError(t, err)
		printSpec(doc, dItem.ResourceBaseURI,dItem.APISpecification.Version) // print the endpoints in the spec
		fmt.Printf("\nImplemented\n--------------\n")
		printImplemented(dItem.Endpoints,dItem.APISpecification.Version) // print what this org has implemeneted
	}
}

func getConditionality(method, path, specification string) (string, error) {
	condition, err := model.GetConditionality(method, path, specification)
	if err != nil {
		return "", err
	}
	switch condition {
	case model.Mandatory:
		return "M", nil
	case model.Conditional:
		return "C", nil
	case model.Optional:
		return "O", nil
	default:
		return "U", nil
	}
}

func printImplemented(endpoints []discovery.ModelEndpoint, spec string) {
	for _, v := range endpoints {
		condition, err := getConditionality(v.Method,v.Path,spec)
		if err != nil {
			fmt.Printf("%s",err)
		}
		fmt.Printf("[%s] %s %s\n", condition, v.Method, v.Path)
	}
}

// match disco endpoints to spec
// expect disco to be subset of spec
// figure out the differences

func printSpec(doc *loads.Document, base , spec string) {
	for path, props := range doc.Spec().Paths.Paths {
		for method := range getOperations(&props) {
			newPath := base + path
			condition, err := getConditionality(method,path,spec)
			if err != nil {
				fmt.Printf("%s",err)
			}			
			fmt.Printf("[%s] %s %s\n", condition, method, newPath) // give to testcase along with any conditionality?
			// map disco
		}
	}
}

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
