package request

import (
	"fmt"
	"strings"
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
