package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Local struct {
	Directory string
}

func (m *Local) UploadFile(ctx context.Context, bucketName, fileName, contentType string, file io.Reader) error {
	path := filepath.Join(m.Directory, bucketName, fileName)
	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("failed open local file: ", err)
		return err
	}

	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("failed copy local file: ", err)
		return err
	}

	if fileInfo, err := os.Stat(path); fileInfo.Size() > 10000000 {
		if err != nil {
			return errors.New("file not found")
		}
		return errors.New("file too large")
	}

	return nil
}

func (m *Local) GetFileTemporaryURL(ctx context.Context, bucketName, filename string) (string, error) {
	return fmt.Sprintf("%s/%s", os.Getenv("CDN_BASE_URL"), filename), nil
}

func (m *Local) GetFilePublicURL(ctx context.Context, bucketName, filename string) (string, error) {
	return "", nil
}

func (m *Local) GetObject(_ context.Context, _, filename string) (io.Reader, error) {
	return os.Open(filename)
}

func (m *Local) FGetObject(ctx context.Context, bucketName, filename, destination string) error {
	return nil
}

func (m *Local) FPutObject(ctx context.Context, bucketName, filename, filepath string) error {
	return nil
}

func (m *Local) RemoveFile(ctx context.Context, bucketName, filename string) error {
	return nil
}

func (m *Local) IsFileExist(ctx context.Context, bucketName, fileptah string) (bool, error) {
	if _, err := os.Stat(fileptah); err == nil {
		// path/to/whatever exists
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

func (m *Local) FCopyObject(ctx context.Context, bucketName, src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed open source file: %w", err)
	}
	defer sourceFile.Close()

	err = os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return fmt.Errorf("failed create target directory: %w", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed create file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed copy file content: %w", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed sync target file: %w", err)
	}

	return nil
}
