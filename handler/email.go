package handler

import (
	"app/request"
	"net/http"
)

func (handler *Handler) TestSendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.TestSendEmail{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.TestSendEmail(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) TestSendNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.TestSendNotification{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.TestSendNotification(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}
