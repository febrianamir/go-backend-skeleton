package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type GetUsers struct {
	BasePaginateRequest
	Preloads []string
}

func (query *GetUsers) GetOrderQuery() string {
	fieldMap := map[string]string{
		"name":       "name",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	return buildOrderQuery(query.Sort, fieldMap)
}

type GetUser struct {
	ID       uint
	Name     string
	Email    string
	Preloads []string
}

type CreateUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (r *CreateUser) Validate() error {
	validationErrDetails := map[string]any{}

	validateField(r.Name, "name", validationErrDetails, validation.Required)
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	validateField(r.PhoneNumber, "phone_number", validationErrDetails, validation.Required)
	validateField(r.Password, "password", validationErrDetails, IsPassword...)

	return buildValidationError(validationErrDetails)
}

type UpdateUser struct {
	ID          uint
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (r *UpdateUser) Validate() error {
	validationErrDetails := map[string]any{}

	validateField(r.Name, "name", validationErrDetails, validation.Required)
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	validateField(r.PhoneNumber, "phone_number", validationErrDetails, validation.Required)

	return buildValidationError(validationErrDetails)
}

type DeleteUser struct {
	ID uint
}
