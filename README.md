# HTTPUtils - Go HTTP Handler Library

[![Go Report Card](https://goreportcard.com/badge/github.com/LullNil/go-http-utils)](https://goreportcard.com/report/github.com/LullNil/go-http-utils)

A lightweight Go library that simplifies HTTP handler development with automatic request validation, error handling, and consistent JSON responses.

## Features

* **Type-safe JSON decoding** with generics
* **Built-in validation** using struct tags
* **Consistent error handling** with proper HTTP status codes and optional data payloads
* **Structured logging** integration
* **Unified response format** across all endpoints

## Installation

```bash
go get github.com/LullNil/go-http-utils
```

## Quick Start

```go
package main

import (
    "net/http"
    "log/slog"
    "github.com/LullNil/go-http-utils/httputils"
)

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    const op = "handler.CreateUser"

    // Decode and validate request
    req, ok := httputils.DecodeRequest[CreateUserRequest](w, r, h.log, op)
    if !ok {
        return
    }

    if !httputils.ValidateRequest(w, r, h.log, op, req) {
        return
    }

    // Business logic
    err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }

    httputils.SendOK(w, r, h.log, op)
}
```

## Code Comparison

**Standard Go (67 lines):**

```go
func CreateUserStandard(w http.ResponseWriter, r *http.Request) {
    // Manual JSON decoding with error handling
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("failed to decode request: %v", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "invalid request payload",
        })
        return
    }
    
    // Manual validation
    if req.Name == "" {
        log.Printf("validation failed: name is required")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "name is required",
        })
        return
    }
    
    if len(req.Name) < 2 {
        log.Printf("validation failed: name too short")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "name must be at least 2 characters",
        })
        return
    }
    
    // Email validation (simplified)
    if req.Email == "" {
        log.Printf("validation failed: email is required")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "email is required",
        })
        return
    }
    
    if !strings.Contains(req.Email, "@") {
        log.Printf("validation failed: invalid email format")
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "invalid email format",
        })
        return
    }
    
    // Business logic
    err := service.CreateUser(r.Context(), req)
    if err != nil {
        log.Printf("service error: %v", err)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "internal server error",
        })
        return
    }

    log.Printf("user created successfully")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "OK",
    })
}
```

**With HTTPUtils (18 lines):**

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    const op = "handler.CreateUser"

    req, ok := httputils.DecodeRequest[CreateUserRequest](w, r, h.log, op)
    if !ok {
        return
    }

    if !httputils.ValidateRequest(w, r, h.log, op, req) {
        return
    }

    err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }

    httputils.SendOK(w, r, h.log, op)
}
```

**Result: 73% less code**

## API Reference

### Request Handling

#### `DecodeRequest[T](w, r, log, op) (T, bool)`

Decodes JSON request body into generic struct type.

```go
req, ok := httputils.DecodeRequest[UserRequest](w, r, log, "CreateUser")
if !ok {
    return // Error already handled
}
```

#### `ValidateRequest[T](w, r, log, op, req) bool`

Validates struct using `validator` tags.

```go
type UserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}

if !httputils.ValidateRequest(w, r, log, op, req) {
    return
}
```

### Response Handling

#### `SendOK(w, r, log, op)`

Sends standard success response:

```json
{"status": "OK"}
```

#### `SendDataOK(w, r, log, op, data)`

Sends success response with data:

```json
{
    "status": "OK",
    "data": { /* your data */ }
}
```

#### `WriteHTTPError(w, log, op, err)`

Handles errors automatically:

* If the error is a `HTTPError` (including `WithData`), it extracts the message, status code, and optional `data`
* Logs structured error and responds accordingly

### Error Handling

#### HTTPError Type

```go
type HTTPError struct {
    Code    int
    Message string
    Data    any // optional
}

func New(code int, msg string) HTTPError {
    return HTTPError{Code: code, Message: msg}
}

func NewWithData(code int, msg string, data any) HTTPError {
    return HTTPError{Code: code, Message: msg, Data: data}
}
```

**Usage in services:**

```go
func (s *Service) GetUser(id string) (*User, error) {
    if id == "" {
        return nil, apperr.New(http.StatusBadRequest, "user ID required")
    }

    user, err := s.repo.GetUser(id)
    if errors.Is(err, ErrUserNotFound) {
        return nil, apperr.New(http.StatusNotFound, "user not found")
    }

    return user, err
}
```

## Advanced Examples

### Handler with Data Response

```go
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    const op = "handler.GetUsers"

    req, ok := httputils.DecodeRequest[GetUsersRequest](w, r, h.log, op)
    if !ok {
        return
    }

    if !httputils.ValidateRequest(w, r, h.log, op, req) {
        return
    }

    users, err := h.service.GetUsers(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }

    httputils.SendDataOK(w, r, h.log, op, users)
}
```

### Error Response with Data (Partial Success)

```go
func (h *Handler) ProcessBatch(w http.ResponseWriter, r *http.Request) {
    const op = "handler.ProcessBatch"

    req, ok := httputils.DecodeRequest[BatchRequest](w, r, h.log, op)
    if !ok {
        return
    }

    results, err := h.service.ProcessBatch(r.Context(), req)
    if err != nil {
        // Use NewWithData and standard WriteHTTPError
        err = apperr.NewWithData(http.StatusConflict, "partial failure", results)
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }

    httputils.SendDataOK(w, r, h.log, op, results)
}
```

### Complex Validation Example

```go
type CreateProductRequest struct {
    Name        string   `json:"name" validate:"required,min=2,max=100"`
    Price       float64  `json:"price" validate:"required,gt=0"`
    CategoryID  int      `json:"category_id" validate:"required,gt=0"`
    Description string   `json:"description" validate:"max=500"`
    Tags        []string `json:"tags" validate:"dive,min=1,max=50"`
}
```

## Key Benefits

* **🔧 Consistency**: Uniform error responses and logging
* **🚀 Productivity**: 73% less boilerplate code
* **🛡️ Reliability**: Type-safe request handling
* **🧹 Maintainability**: Clean, testable handler code
* **📊 Observability**: Structured logging with context

## Best Practices

1. **Operation Names**: Use pattern `"handler.package.Method"`
2. **Validation**: Leverage comprehensive `validate` tags
3. **Error Handling**: Use `HTTPError` or `NewWithData` for rich error information
4. **Logging**: Include operation context in all logs
5. **Clean Architecture**: Keep handlers thin, logic in services

## Migration Guide

1. Replace `json.NewDecoder(r.Body).Decode()` → `DecodeRequest`
2. Replace manual validation → `ValidateRequest`
3. Replace custom error responses → `WriteHTTPError`
4. Replace manual success responses → `SendOK`/`SendDataOK`
5. Remove `WriteHTTPErrorWithData` and use `NewWithData` instead
6. Add structured logging with operation names

## License

MIT License - see [LICENSE](LICENSE) file for details.
