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
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func DataOK(data interface{}) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

func SendDataOK(w http.ResponseWriter, log *slog.Logger, r *http.Request, op string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(DataOK(data))

	log.Info("operation successful", slog.String("op", op))
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func Err(log *slog.Logger, w http.ResponseWriter, r *http.Request, op string, err error, msg string, httpStatus int) {
	log.Error(msg, slog.String("op", op), slog.String("err", err.Error()))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(Error(msg))
}
