package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestCheckRequestSchema(t *testing.T) {
	doc, err := loads.Spec("spec/v3.1.0/payment-initiation-swagger.flattened.json")
	require.NoError(t, err)

	spec := doc.Spec()

	for path, props := range spec.Paths.Paths {
		for meth, op := range getOperations(&props) {
			_, _, _ = path, meth, op
			if path == "/domestic-standing-order-consents" {
				for _, param := range op.Parameters {
					if param.ParamProps.In == "body" {
						schema := param.ParamProps.Schema
						recurseProps(schema, 0)
					}
				}
			}
		}
	}

	t.Fail()
}

func recurseProps(schema *spec.Schema, level int) {
	level++
	for k, j := range schema.SchemaProps.Properties {
		fmt.Printf("%s%s\n", strings.Repeat("   ", level), k)
		recurseProps(&j, level)
	}
}

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
