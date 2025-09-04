package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type TestSendEmail struct {
	Email string `json:"email"`
}

func (r *TestSendEmail) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	return buildValidationError(validationErrDetails)
}
