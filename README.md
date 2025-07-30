[![Go Report Card](https://goreportcard.com/badge/github.com/LullNil/go-http-utils)](https://goreportcard.com/report/github.com/LullNil/go-http-utils)

# HTTPUtils - Simplified Go HTTP Handler Library

A lightweight Go library that provides utilities for building consistent, clean, and maintainable HTTP handlers with built-in request validation, error handling, and response formatting.

## Features

- **Generic Request Decoding**: Type-safe JSON request parsing with compile-time validation
- **Built-in Validation**: Automatic struct validation using `validator` tags
- **Consistent Error Handling**: Standardized error responses with proper HTTP status codes
- **Structured Logging**: Integration with `slog` for consistent operation logging
- **Clean Response Format**: Unified JSON response structure across all endpoints

## Installation

```bash
go get github.com/LullNil/httputils
```

## Quick Start

### Basic Handler Example

```go
package main

import (
    "net/http"
    "log/slog"
    "github.com/LullNil/httputils"
)

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    const op = "handler.CreateUser"
    
    // Decode and validate request in one step
    req, ok := httputils.DecodeRequest[CreateUserRequest](r, h.log, op, w)
    if !ok {
        return
    }
    
    // Validate request
    if ok := httputils.ValidateRequest(req, h.log, op, w, r); !ok {
        return
    }
    
    // Your business logic here
    err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }
    
    // Send success response
    httputils.SendOK(w, h.log, op, r)
}
```

## Comparison: Before vs After

### Standard Go HTTP Handler (Before)

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
    
    // Email validation would require additional logic...
    
    // Business logic
    err := service.CreateUser(r.Context(), req)
    if err != nil {
        // Manual error type checking and response
        log.Printf("service error: %v", err)
        if isValidationError(err) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{
                "error": err.Error(),
            })
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "internal server error",
        })
        return
    }
    
    // Success response
    log.Printf("user created successfully")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "OK",
    })
}
```

### With HTTPUtils (After)

```go
func CreateUserWithUtils(w http.ResponseWriter, r *http.Request) {
    const op = "handler.CreateUser"
    
    req, ok := httputils.DecodeRequest[CreateUserRequest](r, h.log, op, w)
    if !ok {
        return
    }
    
    if ok := httputils.ValidateRequest(req, h.log, op, w, r); !ok {
        return
    }
    
    err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }
    
    httputils.SendOK(w, h.log, op, r)
}
```

**Lines of code: 67 → 18 (73% reduction)**

## API Reference

### Core Functions

#### `DecodeRequest[T any](r *http.Request, log *slog.Logger, op string, w http.ResponseWriter) (T, bool)`

Decodes JSON request body into a generic struct type.

**Parameters:**
- `r`: HTTP request
- `log`: Structured logger
- `op`: Operation name for logging
- `w`: HTTP response writer

**Returns:**
- Decoded struct of type T
- Boolean indicating success (false means error response already sent)

**Example:**
```go
req, ok := httputils.DecodeRequest[UserRequest](r, log, "CreateUser", w)
if !ok {
    return // Error already handled and response sent
}
```

#### `ValidateRequest[T any](req T, log *slog.Logger, op string, w http.ResponseWriter, r *http.Request) bool`

Validates struct using `validator` tags.

**Parameters:**
- `req`: Struct to validate
- `log`: Structured logger
- `op`: Operation name
- `w`: HTTP response writer
- `r`: HTTP request

**Returns:**
- `true` if valid, `false` if validation failed (error response sent)

**Example:**
```go
type UserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=18,max=120"`
}

if ok := httputils.ValidateRequest(req, log, op, w, r); !ok {
    return
}
```

#### `SendOK(w http.ResponseWriter, log *slog.Logger, op string, r *http.Request)`

Sends a standard success response.

**Response Format:**
```json
{
    "status": "OK"
}
```

#### `WriteHTTPError(w http.ResponseWriter, log *slog.Logger, op string, err error)`

Handles error responses based on error type.

- If error implements `HTTPError` interface: uses custom status code and message
- Otherwise: logs error and returns 500 Internal Server Error

### Response Utilities

#### `SendDataOK(w http.ResponseWriter, log *slog.Logger, r *http.Request, op string, data interface{})`

Sends success response with data payload.

**Response Format:**
```json
{
    "status": "OK",
    "data": { /* your data here */ }
}
```

#### `Err(log *slog.Logger, w http.ResponseWriter, r *http.Request, op string, err error, msg string, httpStatus int)`

Sends error response with custom message and status code.

**Response Format:**
```json
{
    "status": "Error",
    "error": "error message"
}
```

### Error Handling

#### `HTTPError` Type

```go
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
```

**Usage in Service Layer:**
```go
func (s *Service) GetUser(id string) (*User, error) {
    if id == "" {
        return nil, apperr.New(http.StatusBadRequest, "user ID is required")
    }
    
    user, err := s.repo.GetUser(id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return nil, apperr.New(http.StatusNotFound, "user not found")
        }
        return nil, err // Will be handled as 500 Internal Server Error
    }
    
    return user, nil
}
```

## Advanced Examples

### Handler with Data Response

```go
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    const op = "handler.GetUsers"
    
    req, ok := httputils.DecodeRequest[GetUsersRequest](r, h.log, op, w)
    if !ok {
        return
    }
    
    if ok := httputils.ValidateRequest(req, h.log, op, w, r); !ok {
        return
    }
    
    users, err := h.service.GetUsers(r.Context(), req)
    if err != nil {
        httputils.WriteHTTPError(w, h.log, op, err)
        return
    }
    
    response.SendDataOK(w, h.log, r, op, users)
}
```

### Request Struct with Validation

```go
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=2,max=100"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    CategoryID  int     `json:"category_id" validate:"required,gt=0"`
    Description string  `json:"description" validate:"max=500"`
    Tags        []string`json:"tags" validate:"dive,min=1,max=50"`
}

type GetProductsRequest struct {
    Page     int    `json:"page" validate:"min=1"`
    Limit    int    `json:"limit" validate:"min=1,max=100"`
    Category string `json:"category" validate:"omitempty,min=2"`
}
```

## Benefits

### 🔧 **Consistency**
- All handlers follow the same pattern
- Uniform error responses across the application
- Standardized logging format

### 🚀 **Productivity**
- Reduces boilerplate code by ~70%
- Generic functions work with any struct type
- Built-in validation eliminates manual checks

### 🛡️ **Reliability**
- Type-safe request handling
- Comprehensive error handling
- Structured logging for better debugging

### 🧹 **Maintainability**
- Clean, readable handler code
- Single responsibility functions
- Easy to test and mock

### 📊 **Observability**
- Consistent operation logging
- Error tracking with context
- Request/response lifecycle visibility

## Best Practices

1. **Use Consistent Operation Names**: Follow a pattern like `"handler.package.Method"`
2. **Validate at Struct Level**: Use comprehensive `validate` tags
3. **Handle Errors Gracefully**: Use `HTTPError` for business logic errors
4. **Log Contextually**: Include operation names in all log entries
5. **Keep Handlers Thin**: Move business logic to service layer

## Migration Guide

To migrate existing handlers:

1. Replace manual JSON decoding with `DecodeRequest`
2. Replace manual validation with `ValidateRequest` 
3. Use `WriteHTTPError` for error responses
4. Use `SendOK` or `SendDataOK` for success responses
5. Add structured logging with operation names

Your handlers will become more maintainable, consistent, and significantly shorter!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
