package schema

import (
	"fmt"
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

	contentType := r.Header.Get("Content-type")
	if contentType != expectedContentType {
		message := fmt.Sprintf("Content-Type Error: Should produce '%s', but got: '%s'", expectedContentType, contentType)
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

	if len(spec.Produces) > 0 {
		expectedContentType = spec.Produces[0]
	}

	if len(operation.Produces) > 0 {
		expectedContentType = operation.Produces[0]
	}

	return expectedContentType, nil
}
