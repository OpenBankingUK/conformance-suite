package schema

type nullValidator struct{}

func NewNullValidator() nullValidator {
	return nullValidator{}
}

func (v nullValidator) Validate(r Response) ([]Failure, error) {
	return nil, nil
}

func (v nullValidator) IsRequestProperty(method, path, propertpath string) (bool, string, error) {
	return false, "", nil
}
