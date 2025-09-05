package handler

import (
	"app/request"
	"net/http"
)

func (handler *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.GetUsers{}
	extractor := URLQueryExtractor{Request: r}
	mapDataFunc := map[string]func(string) (any, error){
		"limit":  extractor.ExtractNumber,
		"page":   extractor.ExtractNumber,
		"search": extractor.ExtractString,
		"sort":   extractor.ExtractSliceStringWithComma,
	}

	err := extractor.ExtractData(mapDataFunc, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	res, err := handler.App.Usecase.GetUsers(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	meta := ResponseMeta{HTTPStatus: http.StatusOK}
	meta.SerializeFromResponse(res.BasePaginateResponse)
	WriteSuccess(ctx, w, res.Data, "success", meta)
}

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.GetUser{}
	id, err := getParamUint(r, "ID")
	if err != nil {
		WriteError(ctx, w, err)
		return
	}
	req.ID = id

	res, err := handler.App.Usecase.GetUser(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.CreateUser{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	err = handler.App.Usecase.CreateUser(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.UpdateUser{}
	err := decodeAndValidateRequest(r, &req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	id, err := getParamUint(r, "ID")
	if err != nil {
		WriteError(ctx, w, err)
		return
	}
	req.ID = id

	err = handler.App.Usecase.UpdateUser(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}

func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.DeleteUser{}
	id, err := getParamUint(r, "ID")
	if err != nil {
		WriteError(ctx, w, err)
		return
	}
	req.ID = id

	err = handler.App.Usecase.DeleteUser(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, nil, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}
