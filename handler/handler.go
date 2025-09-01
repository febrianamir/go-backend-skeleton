package handler

import (
	"app"
	"app/lib"
	"app/lib/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	App *app.App
}

func NewHandler(a *app.App) Handler {
	return Handler{
		App: a,
	}
}

func (handler *Handler) Healthz(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(rw, "server is ok")
}

type SuccessBody struct {
	Data    any          `json:"data,omitempty"`
	Message string       `json:"message,omitempty"`
	Meta    ResponseMeta `json:"meta"`
}

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

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	var resp any
	code := http.StatusInternalServerError

	switch errOrig := err.(type) {
	case lib.CustomError:
		errInfo := ErrorInfo{
			Message:    errOrig.Message,
			Code:       errOrig.Code,
			CodeString: errOrig.CodeString,
		}
		if len(errOrig.ErrDetails) > 0 {
			errInfo.ErrDetails = errOrig.ErrDetails
		}

		resp = ErrorBody{
			Error: errInfo,
			Meta: ResponseMeta{
				HTTPStatus: errOrig.HTTPCode,
			},
		}
		code = errOrig.HTTPCode
		logger.LogError(ctx, "error response", []zap.Field{
			zap.Error(err),
			zap.Int("code", errOrig.Code),
			zap.String("code_string", errOrig.CodeString),
			zap.Any("err_details", errOrig.ErrDetails),
			zap.Int("http_code", errOrig.Code),
		}...)
	default:
		resp = ErrorBody{
			Error: ErrorInfo{
				Message:    lib.ErrorInternalServer.Message,
				Code:       lib.ErrorInternalServer.Code,
				CodeString: lib.ErrorInternalServer.CodeString,
			},
			Meta: ResponseMeta{
				HTTPStatus: lib.ErrorInternalServer.HTTPCode,
			},
		}
		logger.LogError(ctx, "internal server error response", []zap.Field{
			zap.Error(err),
			zap.Int("code", lib.ErrorInternalServer.Code),
			zap.String("code_string", lib.ErrorInternalServer.CodeString),
			zap.Int("http_code", lib.ErrorInternalServer.HTTPCode),
		}...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func WriteSuccess(ctx context.Context, w http.ResponseWriter, data any, message string, meta ResponseMeta) {
	resp := SuccessBody{
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(meta.HTTPStatus)
	json.NewEncoder(w).Encode(resp)
}
