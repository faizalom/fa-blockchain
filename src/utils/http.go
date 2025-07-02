package utils

import (
	"encoding/json"
	"net/http"
)

func Success(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func Message(message string) map[string]any {
	return map[string]any{"message": message}
}

func Error(w http.ResponseWriter, data any, responseCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(data)
}
