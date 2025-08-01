// Copyright (c) 2025 LullNil. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func DataOK(data any) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

func DataWithError(msg string, data any) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
		Data:   data,
	}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func Err(w http.ResponseWriter, r *http.Request, log *slog.Logger, op string, err error, msg string, httpStatus int) {
	log.Error(msg, slog.String("op", op), slog.String("err", err.Error()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(Error(msg))
}
