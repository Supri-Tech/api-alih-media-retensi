package pkg

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempy"`
}

func JSON(w http.ResponseWriter, code int, status, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, "success", message, data)
}

func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, "error", message, nil)
}
