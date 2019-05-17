package schema

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
)

// Response represents a response object from a HTTP Call
type Response struct {
	Method     string
	Path       string
	Header     http.Header
	Body       io.Reader
	StatusCode int
}

// Failure represents a validation failure
type Failure struct {
	Message string
}

func newFailure(message string) Failure {
	return Failure{
		Message: message,
	}
}

// Validator validates a HTTP response object against a schema
type Validator interface {
	Validate(Response) ([]Failure, error)
}

func NewSwaggerOBSpecValidator(specName, version string) (Validator, error) {
	const dirname = "pkg/schema/spec/v3.1.0"
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		wd, errGetwd := os.Getwd()
		if errGetwd != nil {
			return nil, errors.Wrapf(errGetwd, "schema: opening spec folder failed in os.Getwd, dirname=%q", dirname)
		}

		return nil, errors.Wrapf(err, "schema: opening spec folder failed, dirname=%q, wd=%q", dirname, wd)
	}

	for _, f := range files {
		filename := dirname + "/" + f.Name()
		doc, err := loads.Spec(filename)
		if err != nil {
			return nil, errors.Wrapf(err, "schema: opening spec file, filename=%q", filename)
		}
		if doc.Spec().Info.Version == version && doc.Spec().Info.Title == specName {
			return NewSwaggerValidator(filename)
		}
	}

	return nil, fmt.Errorf("schema: could not find spec file for spec %s version %s", specName, version)
}

// NewSwaggerValidator returns a swagger validator implementation
// Takes a schema file path as source, can be remote http(s) or local
func NewSwaggerValidator(schemaPath string) (Validator, error) {
	doc, err := loads.Spec(schemaPath)
	if err != nil {
		return nil, err
	}
	return newValidator(doc)
}

type validators struct {
	validators []Validator
}

func newValidator(doc *loads.Document) (Validator, error) {
	finder := newFinder(doc)

	if doc.Version() != "2.0" {
		return nil, errors.New("unsupported swagger version")
	}

	specVersion := doc.Spec().Info.Version
	switch specVersion {
	case "v3.0.0":
		fallthrough
	case "v3.1.0":
		return validators{
			validators: []Validator{
				newContentTypeValidator(finder),
				newStatusCodeValidator(finder),
				newBodyValidator(finder),
			},
		}, nil
	}

	return nil, errors.New("unsupported spec version")
}

func (v validators) Validate(r Response) ([]Failure, error) {
	var allFailures []Failure
	for _, validator := range v.validators {
		failures, err := validator.Validate(r)
		if err != nil {
			return nil, err
		}
		allFailures = append(allFailures, failures...)
	}
	return allFailures, nil
}
