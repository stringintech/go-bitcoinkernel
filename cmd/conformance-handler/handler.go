package main

import "fmt"

// handleRequest dispatches a request to the appropriate handler
func handleRequest(registry *Registry, req Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			resp = NewHandlerErrorResponse(req.ID, "INTERNAL_ERROR", fmt.Sprintf("%v", r))
		}
	}()

	switch req.Method {
	// Script verification
	case "btck_script_pubkey_verify":
		return handleScriptPubkeyVerify(req)

	// Context management
	case "btck_context_create":
		return handleContextCreate(registry, req)
	case "btck_context_destroy":
		return handleContextDestroy(registry, req)

	// Chainstate manager operations
	case "btck_chainstate_manager_create":
		return handleChainstateManagerCreate(registry, req)
	case "btck_chainstate_manager_get_active_chain":
		return handleChainstateManagerGetActiveChain(registry, req)
	case "btck_chainstate_manager_process_block":
		return handleChainstateManagerProcessBlock(registry, req)
	case "btck_chainstate_manager_destroy":
		return handleChainstateManagerDestroy(registry, req)

	// Chain operations
	case "btck_chain_get_height":
		return handleChainGetHeight(registry, req)
	case "btck_chain_get_by_height":
		return handleChainGetByHeight(registry, req)
	case "btck_chain_contains":
		return handleChainContains(registry, req)

	// Block operations
	case "btck_block_create":
		return handleBlockCreate(registry, req)
	case "btck_block_tree_entry_get_block_hash":
		return handleBlockTreeEntryGetBlockHash(registry, req)

	default:
		return NewHandlerErrorResponse(req.ID, "METHOD_NOT_FOUND", "")
	}
}
