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
		Code:       1004,
		CodeString: "ERROR_PARSE_REQUEST",
		HTTPCode:   http.StatusBadRequest,
	}
)
