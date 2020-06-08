package schema

import (
	"fmt"
	"mime"
	"net/http"
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
	expectedContentTypes, err := v.expectedContentTypes(r)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
	}

	contentTypeRequest := r.Header.Get("Content-type")
	mediaRequest, paramsRequest, err := mime.ParseMediaType(contentTypeRequest)
	if err != nil {
		return nil, errors.Wrap(err, "parse request content type validator")
	}

	for _, contentType := range expectedContentTypes {
		mediaType, parameters, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, errors.Wrap(err, "parse expected content type validator")
		}

		if mediaRequest == mediaType {
			if len(parameters) == 0 && len(paramsRequest) == 0 {
				return nil, nil
			}
			if sameParams(parameters, paramsRequest) {
				return nil, nil
			}
		}
	}
	message := fmt.Sprintf("Content-Type Error: acceptable content types: '%s', : actual content type is '%s'", strings.Join(expectedContentTypes, "','"), contentTypeRequest)
	return []Failure{newFailure(message)}, nil

}

// here to satisfy Validator interface
func (v contentTypeValidator) IsRequestProperty(method, path, propertpath string) (bool, string, error) {
	return false, "", nil
}

func (v contentTypeValidator) expectedContentTypes(r Response) ([]string, error) {
	if r.StatusCode == http.StatusNoContent {
		return nil, ErrNotFound
	}
	spec := v.finder.Spec()

	operation, err := v.finder.Operation(r.Method, r.Path)
	if err != nil {
		return nil, err
	}

	if len(spec.Produces) == 0 && len(operation.Produces) == 0 {
		return nil, ErrNotFound
	}

	var expectedContentTypes []string

	if spec != nil && len(spec.Produces) > 0 {
		expectedContentTypes = make([]string, len(spec.Produces))
		copy(expectedContentTypes, spec.Produces)
	}

	if operation != nil && len(operation.Produces) > 0 {
		expectedContentTypes = make([]string, len(spec.Produces))
		copy(expectedContentTypes, operation.Produces)
	}

	return expectedContentTypes, nil

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
