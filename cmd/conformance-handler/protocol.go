package main

import (
	"encoding/json"
)

type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	ID      string      `json:"id"`
	Success interface{} `json:"success,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Type    string `json:"type"`
	Variant string `json:"variant,omitempty"`
}

// NewErrorResponse creates an error response
func NewErrorResponse(id, errorType, variant string) Response {
	return Response{
		ID: id,
		Error: &Error{
			Type:    errorType,
			Variant: variant,
		},
	}
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(id string, result interface{}) Response {
	return Response{
		ID:      id,
		Success: result,
	}
}

// NewEmptySuccessResponse creates an empty success response
func NewEmptySuccessResponse(id string) Response {
	return NewSuccessResponse(id, map[string]interface{}{})
}
