package handler

import (
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
