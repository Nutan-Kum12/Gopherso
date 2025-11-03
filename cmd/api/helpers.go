package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// writeJSONResponse writes a JSON response with the given status code
func (app *application) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// writeErrorResponse writes an error response with the given status code and message
func (app *application) writeErrorResponse(w http.ResponseWriter, status int, message string) {
	errorResp := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}
	app.writeJSONResponse(w, status, errorResp)
}

// writeSuccessResponse writes a success response with optional data
func (app *application) writeSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	app.writeJSONResponse(w, status, response)
}
