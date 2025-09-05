package handler

import (
	"app/lib/constant"
	"app/request"
	"net/http"
	"strings"
)

func (handler *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := request.UploadFile{}
	targetPath := strings.ToLower(r.FormValue("path"))

	f, err := request.ParseFile(
		r,
		constant.MapUploadFileProps[targetPath].FormValue,
		constant.MapUploadFileProps[targetPath].MaxFilesize,
		constant.MapUploadFileProps[targetPath].AllowedExtensions,
	)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}
	req.TargetPath = targetPath
	req.FileData = f

	res, err := handler.App.Usecase.UploadFile(ctx, req)
	if err != nil {
		WriteError(ctx, w, err)
		return
	}

	WriteSuccess(ctx, w, res, "success", ResponseMeta{HTTPStatus: http.StatusOK})
}
