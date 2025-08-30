package handler

import (
	"app/request"
	"context"
	"encoding/json"
	"net/http"
)

func (handler *Handler) TestSendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req := request.TestSendEmail{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.TestSendEmail(context.Background(), req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}
