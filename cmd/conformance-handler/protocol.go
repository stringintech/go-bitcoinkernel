package main

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	Ref    string          `json:"ref,omitempty"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error,omitempty"`
}

type Error struct {
	Code *ErrorCode `json:"code,omitempty"`
}

type ErrorCode struct {
	Type   string `json:"type"`
	Member string `json:"member"`
}

type RefObject struct {
	Ref string `json:"ref"`
}

// NewErrorResponse creates an error response with the given code type and member.
// Use directly for C API error codes (e.g., "btck_ScriptVerifyStatus").
func NewErrorResponse(codeType, codeMember string) Response {
	return Response{
		Error: &Error{
			Code: &ErrorCode{
				Type:   codeType,
				Member: codeMember,
			},
		},
	}
}

// NewEmptyErrorResponse creates an error response with an empty error object {}.
// Use when an operation fails but no specific error code applies (e.g., C API returned null).
func NewEmptyErrorResponse() Response {
	return Response{Error: &Error{}}
}

// NewSuccessResponse creates a success response with a result value.
// Use when an operation succeeds and returns data.
func NewSuccessResponse(result interface{}) Response {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal result for request: %v", err))
	}
	return Response{
		Result: resultJSON,
	}
}

// NewSuccessResponseWithRef creates a success response returning a reference object.
// Use for methods that create objects and store them in the registry.
func NewSuccessResponseWithRef(ref string) Response {
	return NewSuccessResponse(RefObject{Ref: ref})
}

// NewEmptySuccessResponse creates a success response with no result.
// Use for void/nullptr operations that succeed but return no data.
func NewEmptySuccessResponse() Response {
	return Response{}
}
