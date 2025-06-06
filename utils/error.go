package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents a standard JSON error message.
type ErrorResponse struct {
	Errors Errors `json:"errors"`
}

type Errors struct {
	Message string `json:"message"`
}

// HTTPResponseError is used to signal an HTTP error in systems like Lura.
type HTTPResponseError struct {
	Code         int    `json:"http_status_code"`
	Msg          string `json:"http_body,omitempty"`
	HTTPEncoding string `json:"http_encoding"`
}

// Error implements the error interface.
func (e HTTPResponseError) Error() string {
	return e.Msg
}

// StatusCode returns the HTTP status code.
func (e HTTPResponseError) StatusCode() int {
	return e.Code
}

// Encoding returns the Content-Type for the response.
func (e HTTPResponseError) Encoding() string {
	return e.HTTPEncoding
}

// WriteJSONError writes a simple JSON error response to an http.ResponseWriter.
func WriteJSONError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Errors: Errors{
			Message: http.StatusText(statusCode),
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write JSON error: %v", err)
	}
}

func NewHTTPResponseError(statusCode int) HTTPResponseError {
	body := toJSON(ErrorResponse{
		Errors: Errors{
			Message: http.StatusText(statusCode),
		},
	})
	return HTTPResponseError{
		Code:         statusCode,
		Msg:          body,
		HTTPEncoding: "application/json",
	}
}

func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		return `{}`
	}
	return string(b)
}
