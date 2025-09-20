package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"app/lib"
	"app/lib/auth"
	"app/lib/logger"
	"app/lib/signoz"

	"github.com/felixge/httpsnoop"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
			body, _ := io.ReadAll(request.Body)
			request.Body = io.NopCloser(bytes.NewBuffer(body))
			reqBody := string(body)
			reqBodyIsForm := strings.Contains(reqBody, "=")

			if reqBodyIsForm {
				const maxUploadSize = 20 * 1024 * 1024

				err := request.ParseMultipartForm(maxUploadSize)
				if err == nil {
					reqBodyFormMap := map[string]any{}
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

		ctx, span := signoz.StartSpan(request.Context(), fmt.Sprintf("[%s] %s", request.Method, generateTransactionNameFromURLPath(request.URL.Path)))
		defer span.Finish()

		reqID := request.Header.Get(string(logger.CtxRequestID))
		if reqID == "" {
			reqID = span.TraceID()
		}
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx = context.WithValue(request.Context(), logger.CtxRequestID, reqID)
		m := httpsnoop.CaptureMetrics(handler.PanicMiddleware(next), writer, request.WithContext(ctx))

		var signozSpan trace.Span = *span.SignozSpan

		if signozSpan != nil {
			signozSpan.SetAttributes(attribute.String("path", request.URL.Path))
			signozSpan.SetAttributes(attribute.String("duration", fmt.Sprintf("%d ms", m.Duration.Milliseconds())))
			signozSpan.SetAttributes(attribute.String("method", request.Method))
			signozSpan.SetAttributes(attribute.Int("status", m.Code))
			signozSpan.SetAttributes(attribute.String("user_agent", request.UserAgent()))
			signozSpan.SetAttributes(attribute.String("request_body", reqBodyJson))
			signozSpan.SetAttributes(attribute.String("request_form", reqBodyForm))
			signozSpan.SetAttributes(attribute.String("host", request.Host))

			var code codes.Code = codes.Unset

			if m.Code == 500 {
				code = codes.Error
			} else {
				code = codes.Ok
			}

			signozSpan.SetStatus(code, fmt.Sprintf("http handler with status: %d", m.Code))
		}

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

func (handler *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		idTokenClaim, err := handler.getAndValidateIDToken(ctx, request)
		if err != nil {
			WriteError(ctx, writer, err)
			return
		}

		if idTokenClaim.IsMfaToken {
			WriteError(ctx, writer, lib.ErrorUnauthorized)
			return
		}

		ctx = auth.NewFromCtx(ctx, idTokenClaim)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (handler *Handler) AuthMfaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		idTokenClaim, err := handler.getAndValidateIDToken(ctx, request)
		if err != nil {
			WriteError(ctx, writer, err)
			return
		}

		ctx = auth.NewFromCtx(ctx, idTokenClaim)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (handler *Handler) WebSocketAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		// Two authentication method for websocket:
		// - API Key - for server to server connection
		// - Access Token - for frontend to server connection
		apiKey := request.URL.Query().Get("api_key")
		accessToken := request.URL.Query().Get("token")

		if apiKey == "" && accessToken == "" {
			WriteError(ctx, writer, lib.ErrorUnauthorized)
			return
		}

		if apiKey != "" {
			if !handler.App.Usecase.ValidateWebsocketAPIKey(ctx, apiKey) {
				WriteError(ctx, writer, lib.ErrorUnauthorized)
				return
			}

			idTokenClaim := &auth.IDTokenClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "server",
				},
			}
			ctx = auth.NewFromCtx(ctx, idTokenClaim)
		}

		if accessToken != "" {
			idTokenClaim, err := handler.validateIDToken(ctx, accessToken)
			if err != nil {
				WriteError(ctx, writer, err)
				return
			}

			if idTokenClaim.IsMfaToken {
				WriteError(ctx, writer, lib.ErrorUnauthorized)
				return
			}

			ctx = auth.NewFromCtx(ctx, idTokenClaim)
		}

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (handler *Handler) getAndValidateIDToken(ctx context.Context, request *http.Request) (*auth.IDTokenClaims, error) {
	headerAuthorization := request.Header.Get("Authorization")
	if headerAuthorization == "" {
		return nil, lib.ErrorUnauthorized
	}

	splitToken := strings.Split(headerAuthorization, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return nil, lib.ErrorUnauthorized
	}

	// Exchange access_token -> id_token
	accessToken := splitToken[1]
	idTokenClaim, err := handler.validateIDToken(ctx, accessToken)
	if err != nil {
		return nil, lib.ErrorUnauthorized
	}

	return idTokenClaim, nil
}

func (handler *Handler) validateIDToken(ctx context.Context, accessToken string) (*auth.IDTokenClaims, error) {
	idToken, err := handler.App.Usecase.GetIDToken(ctx, accessToken)
	if err != nil || idToken == "" {
		return nil, lib.ErrorUnauthorized
	}

	idTokenClaim, err := handler.App.Usecase.ParseIDToken(ctx, idToken)
	if err != nil {
		return nil, lib.ErrorUnauthorized
	}

	return idTokenClaim, nil
}

func generateTransactionNameFromURLPath(s string) string {
	parts := strings.Split(s, "/")
	result := "home"

	if len(parts) > 1 {
		result = parts[1]
	}

	if len(parts) > 2 {
		for i, part := range parts {
			var onlyText = regexp.MustCompile(`^[A-Za-z\s]+$`)
			if i > 1 && onlyText.MatchString(part) {
				result = result + "/" + part
			}
		}
	}

	return result
}
