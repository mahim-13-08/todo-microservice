package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success writes a successful JSON response.
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes an error JSON response.
func Error(w http.ResponseWriter, statusCode int, errMsg string) {
	JSON(w, statusCode, APIResponse{
		Success: false,
		Error:   errMsg,
	})
}
