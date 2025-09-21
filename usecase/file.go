package usecase

import (
	"app/lib"
	"app/lib/signoz"
	"app/request"
	"app/response"
	"context"
	"fmt"
)

func (usecase *Usecase) UploadFile(ctx context.Context, req request.UploadFile) (res response.UploadFile, err error) {
	ctx, span := signoz.StartSpan(ctx, "usecase.UploadFile")
	defer span.Finish()

	newFilename := fmt.Sprintf("%s.%s", lib.GenerateUUID(), req.FileExtension)
	filepath := fmt.Sprintf("%s/%s", req.TargetPath, newFilename)

	file, err := req.FileHeader.Open()
	if err != nil {
		return res, err
	}
	defer file.Close()

	bufFileHeader := make([]byte, 512)
	_, err = file.Read(bufFileHeader)
	if err != nil {
		return res, err
	}

	// Reset seek pointer to 0 after the file header being readed
	_, err = file.Seek(0, 0)
	if err != nil {
		return res, err
	}

	err = usecase.storage.UploadFile(ctx, "", filepath, req.FileContentType, file)
	if err != nil {
		return res, err
	}

	tmpURL, err := usecase.storage.GetFileTemporaryURL(ctx, "", filepath)
	if err != nil {
		return res, err
	}

	res = response.UploadFile{
		Filename: req.Filename,
		Filepath: filepath,
		URL:      tmpURL,
	}
	return res, nil
}
