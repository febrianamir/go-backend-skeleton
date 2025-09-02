package handler

import (
	"app/lib"
	"app/lib/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (handler *Handler) PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				var errInfo ErrorInfo

				switch err := r.(type) {
				case error:
					errInfo.Message = fmt.Sprintf("PANIC: %s", err.Error())
				default:
					errInfo.Message = fmt.Sprintf("PANIC: unknown error: %v", err)
				}

				errInfo.Code = lib.ErrorInternalServer.Code
				errInfo.CodeString = lib.ErrorInternalServer.CodeString
				logger.LogError(request.Context(), errInfo.Message, []zap.Field{}...)

				res := ErrorBody{
					Error: errInfo,
					Meta: ResponseMeta{
						HTTPStatus: lib.ErrorInternalServer.HTTPCode,
					},
				}

				writer.WriteHeader(lib.ErrorInternalServer.HTTPCode)
				responseBody, _ := json.Marshal(res)
				_, _ = writer.Write(responseBody)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func (handler *Handler) InstrumentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var reqBodyJson string
		var reqBodyForm string

		if request.Method != http.MethodGet {
			var body, _ = io.ReadAll(request.Body)
			request.Body = io.NopCloser(bytes.NewBuffer(body))
			reqBody := string(body)
			reqBodyIsForm := strings.Contains(reqBody, "=")

			if reqBodyIsForm {
				const maxUploadSize = 20 * 1024 * 1024

				err := request.ParseMultipartForm(maxUploadSize)
				if err == nil {
					var reqBodyFormMap = map[string]any{}
					for key, values := range request.MultipartForm.Value {
						for _, value := range values {
							reqBodyFormMap[key] = value
						}
					}

					reqBodyFormJson, _ := json.Marshal(reqBodyFormMap)
					reqBodyForm = string(reqBodyFormJson)
				}
			} else {
				reqBodyJson = reqBody
			}
		}

		reqID := request.Header.Get(string(logger.CtxRequestID))
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := context.WithValue(request.Context(), logger.CtxRequestID, reqID)
		m := httpsnoop.CaptureMetrics(handler.PanicMiddleware(next), writer, request.WithContext(ctx))

		logger.LogInfo(ctx, fmt.Sprintf("http handler ([%s] - %s) completed", request.Method, request.URL.Path), []zap.Field{
			zap.Int("status_code", m.Code),
			zap.String("duration", fmt.Sprintf("%d ms", m.Duration.Milliseconds())),
			zap.String("method", request.Method),
			zap.String("path", request.URL.Path),
		}...)

		logger.TrafficLogInfo(ctx, fmt.Sprintf("Traffic log: [%s] - %s", request.Method, request.URL.Path), []zap.Field{
			zap.String("path", request.URL.Path),
			zap.String("host", request.Host),
			zap.String("method", request.Method),
			zap.String("duration", fmt.Sprintf("%d ms", m.Duration.Milliseconds())),
			zap.Any("user_agent", request.UserAgent()),
			zap.String("request_body", reqBodyJson),
			zap.String("request_form", reqBodyForm),
			zap.Int("status_code", m.Code),
		}...)
	})
}
