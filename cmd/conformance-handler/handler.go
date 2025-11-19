package main

import "fmt"

// handleRequest dispatches a request to the appropriate handler
func handleRequest(req Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			resp = NewHandlerErrorResponse(req.ID, "INTERNAL_ERROR", fmt.Sprintf("%v", r))
		}
	}()

	switch req.Method {
	case "btck_script_pubkey_verify":
		return handleScriptPubkeyVerify(req)
	default:
		return NewHandlerErrorResponse(req.ID, "METHOD_NOT_FOUND", "")
	}
}
