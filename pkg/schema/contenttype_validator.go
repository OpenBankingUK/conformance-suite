package schema

import (
	"fmt"
	"mime"
	"strings"

	"github.com/pkg/errors"
)

// contentTypeValidator implements a validator for content type check on header
// matches what the swagger spec is expecting with the content type found in
// the response object
type contentTypeValidator struct {
	finder finder
}

func newContentTypeValidator(finder finder) Validator {
	return contentTypeValidator{
		finder: finder,
	}
}

func (v contentTypeValidator) Validate(r Response) ([]Failure, error) {
	expectedContentType, err := v.expectedContentType(r)
	if err == ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	mediaExpected, paramsExpected, err := mime.ParseMediaType(expectedContentType)
	if err != nil {
		return nil, errors.Wrap(err, "parse expected content type validator")
	}

	contentTypeRequest := r.Header.Get("Content-type")
	mediaRequest, paramsRequest, err := mime.ParseMediaType(contentTypeRequest)
	if err != nil {
		return nil, errors.Wrap(err, "parse request content type validator")
	}

	if mediaRequest != mediaExpected {
		message := fmt.Sprintf("Content-Type Error: Should produce '%s', but got: '%s'", mediaExpected, contentTypeRequest)
		return []Failure{newFailure(message)}, nil
	}

	if !sameParams(paramsExpected, paramsRequest) {
		message := fmt.Sprintf("Content-Type Error: Should produce params '%s', but got: '%s'", mapToString(paramsExpected), mapToString(paramsRequest))
		return []Failure{newFailure(message)}, nil
	}

	return nil, nil
}

func (v contentTypeValidator) expectedContentType(r Response) (string, error) {
	spec := v.finder.Spec()

	operation, err := v.finder.Operation(r.Method, r.Path)
	if err != nil {
		return "", err
	}

	if len(spec.Produces) == 0 && len(operation.Produces) == 0 {
		return "", ErrNotFound
	}

	var expectedContentType string

	if spec != nil && len(spec.Produces) > 0 {
		expectedContentType = spec.Produces[0]
	}

	if operation != nil && len(operation.Produces) > 0 {
		expectedContentType = operation.Produces[0]
	}

	return expectedContentType, nil
}

func sameParams(params1, params2 map[string]string) bool {
	if len(params1) != len(params2) {
		return false
	}

	for key, value := range params1 {
		otherValue, ok := params2[key]
		if !ok {
			return false
		}
		if strings.ToUpper(value) != strings.ToUpper(otherValue) {
			return false
		}
	}
	return true
}

func mapToString(params map[string]string) string {
	result := []string{}
	for key, value := range params {
		result = append(result, key+"="+value)
	}
	return strings.Join(result, ";")
}
