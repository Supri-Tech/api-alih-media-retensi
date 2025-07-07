package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
)

const (
	RequestTimeout = 5 * time.Minute
)

var SignatureSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifyHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestampStr := r.Header.Get("X-Timestamp")
		if timestampStr == "" {
			pkg.Error(w, http.StatusBadRequest, "X-Timestamp header required")
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			pkg.Error(w, http.StatusBadRequest, "Invalid X-Timestamp format")
			return
		}

		requestTime := time.Unix(timestamp, 0)
		if time.Since(requestTime) > RequestTimeout {
			pkg.Error(w, http.StatusBadRequest, "Request expired")
			return
		}

		signature := r.Header.Get("X-Signature")
		if signature == "" {
			pkg.Error(w, http.StatusBadRequest, "X-Signature header required")
			return
		}

		expectedSignature := calculateSignature(r, timestampStr)
		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			pkg.Error(w, http.StatusUnauthorized, "Invalid signature")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func calculateSignature(r *http.Request, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(SignatureSecret))

	mac.Write([]byte(r.Method))
	mac.Write([]byte(r.URL.Path))
	mac.Write([]byte(timestamp))

	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		mac.Write(body)
	}

	return hex.EncodeToString(mac.Sum(nil))
}
