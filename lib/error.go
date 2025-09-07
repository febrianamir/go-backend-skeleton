package lib

import (
	"net/http"
)

type CustomError struct {
	Message    string
	Code       int    // Code should be unique
	CodeString string // CodeString should be unique
	ErrDetails map[string]any
	HTTPCode   int
}

func (err CustomError) Error() string {
	return err.Message
}

var (
	ErrorInternalServer = CustomError{
		Message:    "Internal Server Error",
		Code:       1000,
		CodeString: "INTERNAL_SERVER_ERROR",
		HTTPCode:   http.StatusInternalServerError,
	}
	ErrorValidation = CustomError{
		Message:    "Error Validation",
		Code:       1001,
		CodeString: "ERROR_VALIDATION",
		HTTPCode:   http.StatusBadRequest,
	}
	ErrorParseQuery = CustomError{
		Message:    "Error Parse Query",
		Code:       1002,
		CodeString: "ERROR_PARSE_QUERY",
		HTTPCode:   http.StatusBadRequest,
	}
	ErrorParseParam = CustomError{
		Message:    "Error Parse Param",
		Code:       1003,
		CodeString: "ERROR_PARSE_PARAM",
		HTTPCode:   http.StatusBadRequest,
	}
	ErrorNotFound = CustomError{
		Message:    "Error Not Found",
		Code:       1004,
		CodeString: "ERROR_NOT_FOUND",
		HTTPCode:   http.StatusNotFound,
	}
	ErrorParseRequest = CustomError{
		Message:    "Error Parse Request",
		Code:       1005,
		CodeString: "ERROR_PARSE_REQUEST",
		HTTPCode:   http.StatusBadRequest,
	}
	ErrorVerificationDelay = CustomError{
		Message:    "Error Verification Delay",
		Code:       1006,
		CodeString: "ERROR_VERIFICATION_DELAY",
		HTTPCode:   http.StatusUnprocessableEntity,
	}
	ErrorWrongCredential = CustomError{
		Message:    "Error Wrong Credential",
		Code:       1007,
		CodeString: "ERROR_WRONG_CREDENTIAL",
		HTTPCode:   http.StatusUnauthorized,
	}
	ErrorVerificationInactive = CustomError{
		Message:    "Error Verification Inactive",
		Code:       1008,
		CodeString: "ERROR_VERIFICATION_INACTIVE",
		HTTPCode:   http.StatusBadRequest,
	}
	ErrorUnauthorized = CustomError{
		Message:    "Error Unauthorized",
		Code:       1009,
		CodeString: "ERROR_UNAUTHORIZED",
		HTTPCode:   http.StatusUnauthorized,
	}
	ErrorOtpRateLimit = CustomError{
		Message:    "Error Otp Rate Limit Exceeded",
		Code:       1010,
		CodeString: "ERROR_OTP_RATE_LIMIT",
		HTTPCode:   http.StatusUnprocessableEntity,
	}
	ErrorOtpDelay = CustomError{
		Message:    "Error Otp Delay",
		Code:       1011,
		CodeString: "ERROR_OTP_DELAY",
		HTTPCode:   http.StatusUnprocessableEntity,
	}
)
