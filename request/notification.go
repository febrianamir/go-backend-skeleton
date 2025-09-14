package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TestSendNotification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (r *TestSendNotification) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Title, "title", validationErrDetails, validation.Required)
	return buildValidationError(validationErrDetails)
}
