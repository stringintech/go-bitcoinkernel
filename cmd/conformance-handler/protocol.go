package main

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	ID     string          `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

type Error struct {
	Code ErrorCode `json:"code"`
}

type ErrorCode struct {
	Type   string `json:"type"`
	Member string `json:"member"`
}

// NewErrorResponse creates an error response with the given code type and member.
// Use directly for C API error codes (e.g., "btck_ScriptVerifyStatus").
// For handler errors, use NewHandlerErrorResponse.
func NewErrorResponse(id, codeType, codeMember string) Response {
	return Response{
		ID: id,
		Error: &Error{
			Code: ErrorCode{
				Type:   codeType,
				Member: codeMember,
			},
		},
	}
}

// NewHandlerErrorResponse creates an error response for handler layer errors.
// Use for request validation, method routing, and parameter parsing errors.
// Optional detail parameter adds context to the error (e.g., "INVALID_PARAMS (missing field 'foo')").
func NewHandlerErrorResponse(id, codeMember, detail string) Response {
	member := codeMember
	if detail != "" {
		member += fmt.Sprintf(" (%s)", detail)
	}
	return NewErrorResponse(id, "Handler", member)
}

// NewInvalidParamsResponse creates an INVALID_PARAMS error with optional detail.
// Use when request parameters are malformed or missing. Detail provides context about the issue.
func NewInvalidParamsResponse(id, detail string) Response {
	return NewHandlerErrorResponse(id, "INVALID_PARAMS", detail)
}
