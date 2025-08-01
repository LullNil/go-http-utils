// Copyright (c) 2025 LullNil. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

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
func DecodeRequest[T any](w http.ResponseWriter, r *http.Request, log *slog.Logger, op string) (T, bool) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request body", slog.String("op", op), slog.String("err", err.Error()))
		response.Err(w, r, log, op, err, "invalid request payload", http.StatusBadRequest)
		return req, false
	}

	return req, true
}

// ValidateRequest checks the struct against validation tags.
// Returns true if valid, otherwise sends error response and returns false.
func ValidateRequest[T any](w http.ResponseWriter, r *http.Request, log *slog.Logger, op string, req T) bool {
	if err := validate.Struct(req); err != nil {
		log.Error("validation failed",
			slog.String("op", op),
			slog.String("err", err.Error()),
			slog.Any("validation_errors", err.Error()),
		)
		response.Err(w, r, log, op, err, "invalid input data", http.StatusBadRequest)
		return false
	}
	return true
}

// SendOK sends a standard JSON success response with HTTP 200.
func SendOK(w http.ResponseWriter, r *http.Request, log *slog.Logger, op string) {
	writeJSON(w, http.StatusOK, response.OK())
	log.Info("operation successful", slog.String("op", op))
}

// SendDataOK sends a JSON response with the given data and HTTP 200 status.
// Automatically logs the operation as successful.
func SendDataOK(w http.ResponseWriter, r *http.Request, log *slog.Logger, op string, data any) {
	writeJSON(w, http.StatusOK, response.DataOK(data))
	log.Info("operation successful", slog.String("op", op))
}

// WriteHTTPError writes an HTTP error response to w based on the given error.
// If the error is an apperr.HTTPError, it will be used directly.
// Otherwise, it will be logged and an internal server error will be written.
func WriteHTTPError(w http.ResponseWriter, log *slog.Logger, op string, err error) {
	var httpErr *apperr.HTTPError
	if errors.As(err, &httpErr) {
		log.Error("handled error",
			slog.String("op", op),
			slog.String("err", httpErr.Error()),
		)

		if httpErr.Data != nil {
			writeJSON(w, httpErr.Code, response.DataWithError(httpErr.Message, httpErr.Data))
		} else {
			writeJSON(w, httpErr.Code, response.Error(httpErr.Message))
		}
		return
	}

	log.Error("internal error",
		slog.String("op", op),
		slog.String("err", err.Error()),
	)

	writeJSON(w, http.StatusInternalServerError, response.Error("internal server error"))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
