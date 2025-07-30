// Copyright (c) 2025 LullNil. All rights reserved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file.

package apperr

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func New(code int, msg string) HTTPError {
	return HTTPError{Code: code, Message: msg}
}
