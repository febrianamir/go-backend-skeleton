package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Register struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (r *Register) Validate() error {
	validationErrDetails := map[string]any{}

	validateField(r.Name, "name", validationErrDetails, validation.Required)
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	validateField(r.PhoneNumber, "phone_number", validationErrDetails, validation.Required)
	validateField(r.Password, "password", validationErrDetails, IsPassword...)

	return buildValidationError(validationErrDetails)
}
