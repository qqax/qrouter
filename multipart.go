package qrouter

import (
	"errors"
	"fmt"
	"mime/multipart"
)

type MultipartReader interface {
	ReadFromMultipart(fileHeader *multipart.FileHeader, allowedFileTypes []string) error
}

func ReadFile(fileHeader *multipart.FileHeader, ir MultipartReader, MaxSize int64, allowedFileTypes []string) error {
	// Restrict the size of each uploaded file to 1MB.
	// To prevent the aggregate size from exceeding
	// a specified value, use the http.MaxBytesReader() method
	// before calling ParseMultipartForm()
	if fileHeader.Size > MaxSize {
		return errors.New(fmt.Sprintf("the uploaded file is too big: %s. Please use an file less than %vMB in size", fileHeader.Filename, MaxSize/(1024*1024)))
	}

	err := ir.ReadFromMultipart(fileHeader, allowedFileTypes)
	if err != nil {
		return err
	}

	return nil
}
