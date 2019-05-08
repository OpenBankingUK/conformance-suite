package schema

type nullValidator struct{}

func NewNullValidator() nullValidator {
	return nullValidator{}
}

func (v nullValidator) Validate(r Response) ([]Failure, error) {
	return nil, nil
}
