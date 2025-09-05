package handler

import (
	"app"
	"app/lib"
	"app/lib/logger"
	"app/request"
	"app/response"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
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

func (handler *ResponseMeta) SerializeFromResponse(resp response.BasePaginateResponse) {
	handler.Total = &resp.Total
	if resp.Limit == 0 {
		return
	}
	handler.Page = &resp.Page
	handler.Limit = &resp.Limit

	lastPage := resp.Total / resp.Limit
	if (lastPage * resp.Limit) < resp.Total {
		lastPage++
	}
	if lastPage == 0 {
		lastPage++
	}

	handler.LastPage = &lastPage
}

type URLQueryExtractor struct {
	Request *http.Request
}

// ExtractDate parse time with layout 2006-01-02, this function will return YYYY-MM-DD.
func (q *URLQueryExtractor) ExtractDate(s string) (any, error) {
	return time.Parse("2006-01-02", s)
}

// ExtractBool will parse data with function strconv.ParseBool().
func (q *URLQueryExtractor) ExtractBool(s string) (any, error) {
	return strconv.ParseBool(s)
}

// ExtractNumber parse data with function strconv.ParseInt(s, 10, 64).
func (q *URLQueryExtractor) ExtractNumber(s string) (any, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ExtractFloat parse data with function strconv.ParseFloat(s, 64).
func (q *URLQueryExtractor) ExtractFloat(s string) (any, error) {
	return strconv.ParseFloat(s, 64)
}

// ExtractString return string directly.
func (q *URLQueryExtractor) ExtractString(s string) (any, error) {
	return s, nil
}

// ExtractSliceStringWithComma parse data to be []string{}, the separator on this function is comma(,).
func (q *URLQueryExtractor) ExtractSliceStringWithComma(s string) (any, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ','
	records, err := r.Read()
	if err != nil {
		return nil, err
	}

	res := []string{}
	res = append(res, records...)
	return res, nil
}

// ExtractSliceStringWithPipe parse data to be []string{}, the separator on this function is pipe(|).
func (q *URLQueryExtractor) ExtractSliceStringWithPipe(s string) (any, error) {
	list := strings.Split(s, "|")
	res := []string{}
	for _, l := range list {
		if l != "" {
			res = append(res, l)
		}
	}
	return res, nil
}

// ExtractSliceNumberWithComma parse data to be []int64{}, the separator on this function is comma(,).
func (q *URLQueryExtractor) ExtractSliceNumberWithComma(s string) (any, error) {
	list := strings.Split(s, ",")
	res := []int64{}
	for _, l := range list {
		if l != "" {
			n, _ := strconv.ParseInt(l, 10, 64)
			if n > 0 {
				res = append(res, n)
			}
		}
	}
	return res, nil
}

// Extract data to the receiver, the receiver must be pointer. Basically this use json.Marshal and json.Unmarshal
func (q *URLQueryExtractor) ExtractData(fn map[string]func(string) (val any, err error), receiver any) error {
	urlQuery := q.Request.URL.Query()
	mapData := map[string]any{}

	parseQueryError := lib.ErrorParseQuery
	parseQueryError.ErrDetails = map[string]any{}
	for k, f := range fn {
		if urlQuery.Get(k) != "" {
			v, err := f(urlQuery.Get(k))
			if err != nil {
				parseQueryError.ErrDetails[k] = err.Error()
				continue
			}
			mapData[k] = v
		}
	}
	if len(parseQueryError.ErrDetails) > 0 {
		return parseQueryError
	}

	bytes, err := json.Marshal(mapData)
	if err != nil {
		parseQueryError.Message = "Failed Marshal Data"
		return parseQueryError
	}

	err = json.Unmarshal(bytes, &receiver)
	if err != nil {
		parseQueryError.Message = "Failed Unmarshal Data"
		return parseQueryError
	}

	return nil
}

func getParamUint(r *http.Request, key string) (uint, error) {
	valueString := chi.URLParam(r, key)
	value, err := strconv.ParseUint(valueString, 10, 64)
	if err != nil {
		parseParamError := lib.ErrorParseParam
		parseParamError.ErrDetails = map[string]any{
			key: err.Error(),
		}
		return 0, parseParamError
	}
	return uint(value), nil
}

func decodeAndValidateRequest[T request.Validator](r *http.Request, req T) error {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	err = req.Validate()
	if err != nil {
		return err
	}

	return nil
}
