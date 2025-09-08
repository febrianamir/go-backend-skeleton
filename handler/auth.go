package handler

import (
	"app/lib"
	"app/lib/auth"
	"app/request"
	"net/http"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	ctx := r.Context()

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
	ctx := r.Context()

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
	ctx := r.Context()

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

func (handler *Handler) SendMfaOtp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
	ctx := r.Context()

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
	ctx := r.Context()

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
	ctx := r.Context()

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
