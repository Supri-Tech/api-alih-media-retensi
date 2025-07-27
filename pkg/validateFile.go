package pkg

import (
	"errors"
	"mime/multipart"
)

var MaxFileSize = 1 << 20

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

func ValidateFile(header *multipart.FileHeader) error {
	if header.Size > int64(MaxFileSize) {
		return errors.New("File is too large; must be less than 1mb")
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return errors.New("File is not allowed")
	}

	return nil
}
