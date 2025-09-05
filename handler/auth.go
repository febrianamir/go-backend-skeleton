package handler

import (
	"app/request"
	"encoding/json"
	"net/http"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := request.Register{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = req.Validate()
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
