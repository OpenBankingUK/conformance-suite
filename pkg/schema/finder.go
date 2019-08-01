package schema

import (
	"errors"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

var ErrNotFound = errors.New("operation for method/path not found")

// finder is a helper to find schema and operation in a depp nested
// swagger spec document
type finder struct {
	doc     *loads.Document
	matcher Matcher
}

func newFinder(doc *loads.Document) finder {
	return finder{
		doc:     doc,
		matcher: NewMatcher(),
	}
}

func (f finder) Spec() *spec.Swagger {
	return f.doc.Spec()
}

// Operation returns a Operation object from the spec relative to a method and path
func (f finder) Operation(method, path string) (*spec.Operation, error) {
	for specPath, props := range f.doc.Spec().Paths.Paths {
		if f.matcher.Match(specPath, path) {
			var operation *spec.Operation

			switch strings.ToUpper(method) {
			case "DELETE":
				operation = props.Delete
			case "GET":
				operation = props.Get
			case "HEAD":
				operation = props.Head
			case "OPTIONS":
				operation = props.Options
			case "PATCH":
				operation = props.Patch
			case "POST":
				operation = props.Post
			case "PUT":
				operation = props.Put
			}

			if operation != nil {
				return operation, nil
			}
		}
	}
	return nil, ErrNotFound
}

// Response returns a Response object from the spec relative to a method, path and a
// specific response code
func (f finder) Response(method, path string, statusCode int) (spec.Response, error) {
	operation, err := f.Operation(method, path)
	if err != nil {
		return spec.Response{}, err
	}

	response, ok := operation.Responses.StatusCodeResponses[statusCode]
	if !ok {
		return spec.Response{}, ErrNotFound
	}

	return response, nil
}
