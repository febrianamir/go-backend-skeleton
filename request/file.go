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

	if err := validateFileSize(f, fh.Size, maxFileSize); err != nil {
		return fd, err
	}

	if err := validateFileExtension(f, fh.Filename, allowedExtensions); err != nil {
		return fd, err
	}

	fd.File = f
	fd.FileHeader = fh
	fd.FileExtension = extractFileExt(fh.Filename)
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

		if err := validateFileSize(f, fh.Size, maxFileSize); err != nil {
			return nil, err
		}

		if err := validateFileExtension(f, fh.Filename, allowedExtensions); err != nil {
			return nil, err
		}

		fd := FileData{
			File:            f,
			FileHeader:      fh,
			FileExtension:   extractFileExt(fh.Filename),
			Filename:        fh.Filename,
			FileContentType: fh.Header.Get("Content-Type"),
		}
		files = append(files, fd)
	}

	return files, nil
}

func validateFileExtension(f multipart.File, filename string, allowedExtensions []string) error {
	fExt := extractFileExt(filename)
	if len(allowedExtensions) > 0 {
		if !slices.Contains(allowedExtensions, fExt) && fExt != "" {
			f.Close()
			customErr := lib.ErrorValidation
			customErr.ErrDetails = map[string]any{
				"file": fmt.Sprintf("Invalid file extension. Allowed extensions: %s.", strings.Join(allowedExtensions, ", ")),
			}
			return customErr
		}
	}
	return nil
}

func validateFileSize(f multipart.File, sizeMB, maxFileSize int64) error {
	size := sizeMB / 1024 / 1024 // In MB
	if size > maxFileSize {
		f.Close()
		customErr := lib.ErrorValidation
		customErr.ErrDetails = map[string]any{
			"file": fmt.Sprintf("File too large. Max %d MB.", maxFileSize),
		}
		return customErr
	}
	return nil
}

func extractFileExt(filename string) string {
	fExt := filepath.Ext(filename)
	return strings.ToLower(strings.TrimPrefix(fExt, "."))
}
