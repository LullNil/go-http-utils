// Copyright (c) 2025 LullNil. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package apperr

type HTTPError struct {
	Code    int
	Message string
	Data    any
}

// Error implements the error interface, returning the message string.
func (e HTTPError) Error() string {
	return e.Message
}

// New returns a new HTTPError with the given code and message.
// The Data field is set to nil.
func New(code int, msg string) *HTTPError {
	return &HTTPError{Code: code, Message: msg}
}

// NewWithData returns a new HTTPError with the given code, message, and data.
// The Data field is used to provide additional context or information related to the error.
func NewWithData(code int, msg string, data any) *HTTPError {
	return &HTTPError{Code: code, Message: msg, Data: data}
}
