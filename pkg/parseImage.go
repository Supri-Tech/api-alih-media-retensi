package pkg

import (
	"encoding/base64"
	"io"
	"net/http"
)

func ParseImage(r *http.Request, fieldName string) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	if err := ValidateFile(header); err != nil {
		return "", err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileBytes), nil
}
