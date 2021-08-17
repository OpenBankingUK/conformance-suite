package schema

// NullValidator - type
type NullValidator struct{}

// NewNullValidator -
func NewNullValidator() NullValidator {
	return NullValidator{}
}

// Validate - nop
func (v NullValidator) Validate(r HTTPResponse) ([]Failure, error) {
	return nil, nil
}

// IsRequestProperty - nop
func (v NullValidator) IsRequestProperty(method, path, propertpath string) (bool, string, error) {
	return false, "", nil
}
