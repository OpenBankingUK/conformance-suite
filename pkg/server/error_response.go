package server

// ErrorResponse wraps `error` into a JSON object.
type ErrorResponse struct {
	Error interface{} `json:"error"`
}

// NewErrorResponse creates a ErrorResponse object from a go standard error
func NewErrorResponse(err error) *ErrorResponse {
	return NewErrorMessageResponse(err.Error())
}

func NewErrorMessageResponse(msg string) *ErrorResponse {
	return &ErrorResponse{msg}
}
