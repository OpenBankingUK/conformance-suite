package models

import validation "github.com/go-ozzo/ozzo-validation"

// ImportRequest - Request to `/api/import/review` or `/api/import/rerun` POST.
// TODO(mbana): Needs more work.
type ImportRequest struct {
	Report string `json:"report" form:"report"` // The exported report ZIP archive.
}

// Validate - used by github.com/go-ozzo/ozzo-validation to validate struct.
func (r ImportRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Report, validation.Required),
	)
}

// ImportReviewResponse - Response to `/api/import/review` POST.
type ImportReviewResponse struct {
}

// ImportRerunResponse - Response to `/api/import/rerun` POST.
type ImportRerunResponse struct {
}
