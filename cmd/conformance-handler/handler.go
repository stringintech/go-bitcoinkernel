package main

import "fmt"

// handleRequest dispatches a request to the appropriate handler
func handleRequest(req Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			resp = NewErrorResponse(req.ID, "InternalError", fmt.Sprintf("%v", r))
		}
	}()

	switch req.Method {
	case "script_pubkey.verify":
		return handleScriptPubkeyVerify(req)
	default:
		return NewErrorResponse(req.ID, "MethodNotFound", req.Method)
	}
}
