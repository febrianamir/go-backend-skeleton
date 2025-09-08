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

type RegisterResendVerification struct {
	Email string `json:"email"`
}

func (r *RegisterResendVerification) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	return buildValidationError(validationErrDetails)
}

type VerifyAccount struct {
	Code string `json:"code"`
}

func (r *VerifyAccount) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Code, "code", validationErrDetails, validation.Required)
	return buildValidationError(validationErrDetails)
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *Login) Validate() error {
	validationErrDetails := map[string]any{}

	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	validateField(r.Password, "password", validationErrDetails, IsPassword...)

	return buildValidationError(validationErrDetails)
}

type SendMfaOtp struct {
	Channel string `json:"channel"`
}

func (r *SendMfaOtp) Validate() error {
	return nil
}

type SendOtp struct {
	Channel string `json:"channel"`
	UserId  uint   `json:"user_id"`
}

type ValidateMfaOtp struct {
	OtpCode string `json:"otp_code"`
	UserId  uint   `json:"user_id"`
}

func (r *ValidateMfaOtp) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.OtpCode, "otp_code", validationErrDetails, validation.Required)
	return buildValidationError(validationErrDetails)
}

type ForgotPassword struct {
	Email string `json:"email"`
}

func (r *ForgotPassword) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Email, "email", validationErrDetails, validation.Required, is.EmailFormat)
	return buildValidationError(validationErrDetails)
}

type ResetPassword struct {
	Code               string `json:"code"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

func (r *ResetPassword) Validate() error {
	validationErrDetails := map[string]any{}
	validateField(r.Code, "code", validationErrDetails, validation.Required)
	validateField(r.NewPassword, "new_password", validationErrDetails, IsPassword...)
	validateField(r.ConfirmNewPassword, "confirm_new_password", validationErrDetails, validation.By(isEqual(r.NewPassword, "password")))
	return buildValidationError(validationErrDetails)
}
