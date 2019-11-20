package schema

import (
	"fmt"
)

// statusCodeValidator implements a validator that checks that
// the status code response exists in the swagger spec definition
type statusCodeValidator struct {
	finder finder
}

func newStatusCodeValidator(finder finder) Validator {
	return statusCodeValidator{
		finder: finder,
	}
}

func (v statusCodeValidator) Validate(r Response) ([]Failure, error) {
	_, err := v.finder.Response(r.Method, r.Path, r.StatusCode)
	if err == ErrNotFound {
		message := fmt.Sprintf("server Status %d not defined by the spec", r.StatusCode)
		failure := []Failure{newFailure(message)}
		return failure, nil
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

func (v statusCodeValidator) IsRequestProperty(method, path, propertpath string) (bool, string, error) {
	return false, "", nil
}
