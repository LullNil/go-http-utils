package httputils

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/LullNil/go-http-utils/apperr"
	"github.com/LullNil/go-http-utils/response"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeRequest parses and validates JSON body from the request into the given generic struct.
// Returns the struct and a boolean indicating success or failure.
func DecodeRequest[T any](r *http.Request, log *slog.Logger, op string, w http.ResponseWriter) (T, bool) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.String("op", op), slog.String("err", err.Error()))
		response.Err(log, w, r, op, err, "invalid request payload", http.StatusBadRequest)
		return req, false
	}

	return req, true
}

// ValidateRequest checks the struct against validation tags.
// Returns true if valid, otherwise sends error response and returns false.
func ValidateRequest[T any](req T, log *slog.Logger, op string, w http.ResponseWriter, r *http.Request) bool {
	if err := validate.Struct(req); err != nil {
		log.Error("validation failed",
			slog.String("op", op),
			slog.String("err", err.Error()),
			slog.Any("validation_errors", err.Error()),
		)
		response.Err(log, w, r, op, err, "invalid input data", http.StatusBadRequest)
		return false
	}
	return true
}

// SendOK sends a standard JSON success response with HTTP 200.
func SendOK(w http.ResponseWriter, log *slog.Logger, op string, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.OK())
	log.Info("operation successful", slog.String("op", op))
}

// WriteHTTPError writes an HTTP error response to w based on the given error.
// If the error is an apperr.HTTPError, it will be used directly.
// Otherwise, it will be logged and an internal server error will be written.
func WriteHTTPError(w http.ResponseWriter, log *slog.Logger, op string, err error) {
	var httpErr apperr.HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	log.Error("internal error", slog.String("op", op), slog.String("err", err.Error()))
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
