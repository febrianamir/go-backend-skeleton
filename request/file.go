package request

import (
	"app/lib"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type UploadFile struct {
	FileData
	TargetPath string
}

type FileData struct {
	File            multipart.File
	FileHeader      *multipart.FileHeader
	FileExtension   string
	Filename        string
	FileContentType string
}

// ParseFile response for FileExtension is set to lower case and the dot of extension also was removed. ex: pdf, xlsx
func ParseFile(r *http.Request, field string, maxFileSize int64, allowedExtensions []string) (FileData, error) {
	var (
		fd  FileData
		f   multipart.File
		fh  *multipart.FileHeader
		err error
	)

	err = r.ParseMultipartForm(int64((maxFileSize + 1) * 1024 * 1024))
	if err != nil {
		return fd, err
	}

	f, fh, err = r.FormFile(field)
	if err != nil {
		return fd, err
	}

	// Validate file size
	size := fh.Size / 1024 / 1024 // In MB
	if size > maxFileSize {
		customErr := lib.ErrorValidation
		customErr.ErrDetails = map[string]any{
			"file": fmt.Sprintf("File too large. Max %d MB.", maxFileSize),
		}
		return fd, customErr
	}

	// Validate extension
	fExt := filepath.Ext(fh.Filename)
	fExt = strings.ReplaceAll(fExt, ".", "")
	fExt = strings.ToLower(fExt)
	if len(allowedExtensions) > 0 {
		if !validFileExtension(fExt, allowedExtensions) && fExt != "" {
			customErr := lib.ErrorValidation
			customErr.ErrDetails = map[string]any{
				"file": fmt.Sprintf("Invalid file extension. Allowed extensions: %s.", strings.Join(allowedExtensions, ", ")),
			}
			return fd, customErr
		}
	}

	fd.File = f
	fd.FileHeader = fh
	fd.FileExtension = fExt
	fd.Filename = fh.Filename
	fd.FileContentType = fh.Header.Get("Content-Type")
	return fd, nil
}

// ParseFiles response for FileExtension is set to lower case and the dot of extension also was removed. ex: pdf, xlsx
func ParseFiles(r *http.Request, field string, maxFileSize int64, allowedExtensions []string) ([]FileData, error) {
	var files []FileData

	maxSize := int64(maxFileSize * 1024 * 1024) // In MB
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		return nil, err
	}

	form := r.MultipartForm
	fileHeaders := form.File[field]
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("no files uploaded with field name: %s", field)
	}

	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}

		size := float64(fh.Size) / 1024 / 1024 // In MB
		if size > float64(maxFileSize) {
			f.Close()
			customErr := lib.ErrorValidation
			customErr.ErrDetails = map[string]any{
				"file": fmt.Sprintf("File too large. Max %d MB.", maxFileSize),
			}
			return nil, customErr
		}

		fExt := filepath.Ext(fh.Filename)
		fExt = strings.ToLower(strings.TrimPrefix(fExt, "."))
		if len(allowedExtensions) > 0 {
			if !validFileExtension(fExt, allowedExtensions) && fExt != "" {
				f.Close()
				customErr := lib.ErrorValidation
				customErr.ErrDetails = map[string]any{
					"file": fmt.Sprintf("Invalid file extension. Allowed extensions: %s.", strings.Join(allowedExtensions, ", ")),
				}
				return nil, customErr
			}
		}

		fd := FileData{
			File:            f,
			FileHeader:      fh,
			FileExtension:   fExt,
			Filename:        fh.Filename,
			FileContentType: fh.Header.Get("Content-Type"),
		}
		files = append(files, fd)
	}

	return files, nil
}

func validFileExtension(ext string, extensions []string) bool {
	return slices.Contains(extensions, ext)
}
