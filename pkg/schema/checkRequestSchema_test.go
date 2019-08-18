package schema

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/require"
)

func TestCheckRequestSchema(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/payment-initiation-swagger.flattened.json")
	require.NoError(t, err)

	spec := doc.Spec()

	//spew.Config.Indent = "\t"
	spew.Config.MaxDepth = 2
	// spew.Config.DisablePointerAddresses = true
	// spew.Config.SortKeys = true
	// spew.Config.DisablePointerMethods = true
	// spew.Config.DisableMethods = true

	// t.Logf("--> %#v\n", spec.Info)
	// paths = spec.Paths
	// for k, pathitem := range spec.Paths.Paths {
	// 	_, _ = k, pathitem
	// 	if pathitem.Post != nil {
	// 		for i, parameter := range pathitem.Post.Parameters {
	// 			_ = i
	// 			if parameter.In == "header" {
	// 				continue
	// 			}
	// 			fmt.Printf("%s in %s\n", parameter.Name, parameter.In)
	// 		}
	// 	}
	// }

	//x := spec.SwaggerProps["OBDomesticStandingOrder2"].

	x := spec.Parameters["OBDomesticStandingOrder2"]
	spew.Dump(x)

	t.Fail()
}
