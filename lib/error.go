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
		Code:       0001,
		CodeString: "INTERNAL_SERVER_ERROR",
		HTTPCode:   http.StatusInternalServerError,
	}
)
