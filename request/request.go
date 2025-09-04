package request

import (
	"app/lib"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BasePaginateRequest struct {
	Sort   []string `json:"sort"`
	Search string   `json:"search"`
	Page   uint     `json:"page"`
	Limit  uint     `json:"limit"`
}

func (query *BasePaginateRequest) GetOffset() uint {
	return (query.Page - 1) * query.Limit
}

func buildOrderQuery(querySort []string, fieldMap map[string]string) string {
	result := []string{}
	for _, s := range querySort {
		if len(s) == 0 {
			continue
		}

		order, key := "ASC", s
		if s[len(s)-1:] == "-" {
			order, key = "DESC", s[:len(s)-1]
		}

		fieldName, ok := fieldMap[key]
		if !ok {
			continue
		}

		result = append(result, fmt.Sprintf("%s %s", fieldName, order))
	}

	return strings.Join(result, ",")
}

// validateField validates a single field and adds any error to the provided map
func validateField(field, key string, errDetails map[string]any, rules ...validation.Rule) {
	if err := validation.Validate(field, rules...); err != nil {
		errDetails[key] = err.Error()
	}
}

// buildValidationError creates a customer error validation from error details map
func buildValidationError(errDetails map[string]any) error {
	if len(errDetails) == 0 {
		return nil
	}
	validationError := lib.ErrorValidation
	validationError.ErrDetails = errDetails
	return validationError
}

var IsPassword = []validation.Rule{
	validation.Required,
	validation.Length(6, 0),
	validation.Match(regexp.MustCompile(`[0-9]`)).Error("at least one number"),
	validation.Match(regexp.MustCompile(`[a-zA-Z]`)).Error("at least one letter"),
	validation.Match(regexp.MustCompile(`[!@#$%^&*()]`)).Error("at least one special character"),
	validation.Match(regexp.MustCompile(`^[ a-zA-Z0-9!@#$%^&*()?]+$`)).Error("use only allowed special characters: !@#$%^&*()"),
}
