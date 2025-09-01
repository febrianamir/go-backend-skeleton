package middleware

type ErrorBody struct {
	Error ErrorInfo `json:"error"`
	Meta  any       `json:"meta"`
}

type ResponseMeta struct {
	HTTPStatus int   `json:"http_status"`
	Total      *uint `json:"total,omitempty"`
	Offset     *uint `json:"offset,omitempty"`
	Limit      *uint `json:"limit,omitempty"`
	Page       *uint `json:"page,omitempty"`
	LastPage   *uint `json:"last_page,omitempty"`
}

type ErrorInfo struct {
	Message    string         `json:"message"`
	Code       int            `json:"code,omitempty"`
	CodeString string         `json:"code_string,omitempty"`
	ErrDetails map[string]any `json:"err_details,omitempty"`
}
