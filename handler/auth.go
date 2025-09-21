package handler

import (
	"app/lib"
	"app/lib/auth"
	"app/lib/signoz"
	"app/request"
	"net/http"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.Register")
	defer span.Finish()

	req := request.Register{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.Register(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) RegisterResendVerification(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.RegisterResendVerification")
	defer span.Finish()

	req := request.RegisterResendVerification{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.RegisterResendVerification(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.VerifyAccount")
	defer span.Finish()

	req := request.VerifyAccount{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.VerifyAccount(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.Login")
	defer span.Finish()

	req := request.Login{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.Login(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.RefreshSession")
	defer span.Finish()

	req := request.RefreshSession{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.RefreshSession(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) SendMfaOtp(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.SendMfaOtp")
	defer span.Finish()

	idTokenClaims := auth.GetAuthFromCtx(ctx)
	if idTokenClaims.UserID == 0 {
		WriteError(ctx, w, lib.ErrorUnauthorized)
		return
	}

	req := request.SendMfaOtp{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.SendOtp(ctx, request.SendOtp{
		Channel: req.Channel,
		UserId:  idTokenClaims.UserID,
	})
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) ValidateMfaOtp(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.ValidateMfaOtp")
	defer span.Finish()

	idTokenClaims := auth.GetAuthFromCtx(ctx)
	if idTokenClaims.UserID == 0 {
		WriteError(ctx, w, lib.ErrorUnauthorized)
		return
	}

	req := request.ValidateMfaOtp{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.ValidateOtp(ctx, request.ValidateMfaOtp{
		OtpCode: req.OtpCode,
		UserId:  idTokenClaims.UserID,
	})
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.ForgotPassword")
	defer span.Finish()

	req := request.ForgotPassword{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.ForgotPassword(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.ResetPassword")
	defer span.Finish()

	req := request.ResetPassword{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.ResetPassword(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) SsoGoogle(w http.ResponseWriter, r *http.Request) {
	ctx, span := signoz.StartSpan(r.Context(), "handler.SsoGoogle")
	defer span.Finish()

	req := request.SsoGoogle{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.SsoGoogle(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}
