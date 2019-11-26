package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// bodyValidator implements a schema body validator
// validates body schema using swagger spec document
type bodyValidator struct {
	finder finder
}

func newBodyValidator(finder finder) Validator {
	return bodyValidator{
		finder: finder,
	}
}

func (v bodyValidator) Validate(r Response) ([]Failure, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return v.validate(r, body)
}

func (v bodyValidator) IsRequestProperty(method, path, propertpath string) (bool, string, error) {
	return false, "", nil
}

func (v bodyValidator) validate(r Response, body []byte) ([]Failure, error) {
	response, err := v.finder.Response(r.Method, r.Path, r.StatusCode)
	if err == ErrNotFound {
		message := fmt.Sprintf("could't find a schema to validate for status code %d", r.StatusCode)
		return []Failure{newFailure(message)}, nil
	} else if err != nil {
		return nil, err
	}

	if r.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		message := fmt.Sprintf("could not unmarshal request body %s", err.Error())
		return []Failure{newFailure(message)}, nil
	}

	// swagger API call
	val := validate.NewSchemaValidator(response.Schema, v.finder.doc, "", strfmt.Default)
	result := val.Validate(data)
	if result.HasErrors() {
		return mapToFailures(result), nil
	}

	return nil, nil
}

// mapToFailures maps between swagger error and this package Failure object
func mapToFailures(result *validate.Result) []Failure {
	failures := []Failure{}
	for _, err := range result.Errors {
		failures = append(failures, newFailure(err.Error()))
	}
	return failures
}
