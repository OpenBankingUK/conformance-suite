package schema

import (
	"errors"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Compatibility shims. The OpenAPI 3 validator was originally written against
// kin-openapi v0.128.0; these keep responses being accepted and rejected exactly
// as they were after upgrading to v0.149.0.

// kin-openapi v0.149.0 tightened its built-in "date" and "date-time" checks: a
// timezone became mandatory and month/day/hour ranges are now enforced. These
// are the patterns shipped with v0.128.0.
const (
	legacyFormatOfStringDate     = `^[0-9]{4}-(0[0-9]|10|11|12)-([0-2][0-9]|30|31)$`
	legacyFormatOfStringDateTime = `^[0-9]{4}-(0[0-9]|10|11|12)-([0-2][0-9]|30|31)T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|(\+|-)[0-9]{2}:[0-9]{2})?$`
)

func init() {
	openapi3.DefineStringFormatValidator("date", openapi3.NewRegexpFormatValidator(legacyFormatOfStringDate))
	openapi3.DefineStringFormatValidator("date-time", openapi3.NewRegexpFormatValidator(legacyFormatOfStringDateTime))
}

// emptyHeaderValuesAsNull returns header with every empty-valued entry replaced by
// an entry that has no values. kin-openapi v0.128.0 decoded an empty response header
// as null, so a declared string header sent empty failed validation ("Value is not
// nullable"); v0.149.0 decodes it as "" and accepts it. Both versions decode a header
// with no values as null, so this restores the original result and message.
func emptyHeaderValuesAsNull(header http.Header) http.Header {
	var out http.Header
	for name, values := range header {
		if len(values) == 0 || values[0] != "" {
			continue
		}
		if out == nil {
			out = header.Clone()
		}
		out[name] = []string{}
	}
	if out == nil {
		return header
	}
	return out
}

// legacyAllOfMessage restores the wording kin-openapi v0.128.0 used for a failing
// "allOf": the first failing branch only. v0.149.0 prefixes "doesn't match schema
// due to:" and appends every failing branch. It is installed as the schema error
// message customizer; returning "" keeps the library's own message, which happens
// for every other error and whenever the old text cannot be reproduced exactly.
func legacyAllOfMessage(err *openapi3.SchemaError) string {
	const libraryPrefix = "doesn't match schema due to: "

	if err.SchemaField != "allOf" || err.Origin == nil {
		return ""
	}

	var branches openapi3.MultiError
	for cause := err.Origin; cause != nil && branches == nil; cause = errors.Unwrap(cause) {
		branches, _ = cause.(openapi3.MultiError)
	}
	if len(branches) == 0 {
		return ""
	}

	// The library renders its text when the allOf error is created, but adds the
	// parents' keys to the branch errors' paths afterwards, so a branch must not
	// be rendered again as it is now. Rebuild its original text and check it
	// against the library's own before using it.
	path := err.JSONPointer()
	first := branchTextAtCreation(branches[0], len(path))
	rendered := err.Origin.Error()
	if rendered != libraryPrefix+first && !strings.HasPrefix(rendered, libraryPrefix+first+" And ") {
		return ""
	}
	return errorAt(path) + first
}

// branchTextAtCreation returns the text of an allOf branch error before the
// parentKeys outermost path elements were added to it.
func branchTextAtCreation(branch error, parentKeys int) string {
	schemaErr, ok := branch.(*openapi3.SchemaError)
	if !ok {
		return branch.Error()
	}
	path := schemaErr.JSONPointer()
	text := strings.TrimPrefix(schemaErr.Error(), errorAt(path))
	return errorAt(path[min(parentKeys, len(path)):]) + text
}

// errorAt renders a schema error path the way kin-openapi prefixes its messages.
func errorAt(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return `Error at "/` + strings.Join(path, "/") + `": `
}
