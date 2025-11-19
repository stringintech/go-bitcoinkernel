package main

import "fmt"

// handleRequest dispatches a request to the appropriate handler
func handleRequest(registry *Registry, req Request) (resp Response, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	switch req.Method {
	// Script pubkey operations
	case "btck_script_pubkey_create":
		return handleScriptPubkeyCreate(registry, req)
	case "btck_script_pubkey_destroy":
		return handleScriptPubkeyDestroy(registry, req)
	case "btck_script_pubkey_verify":
		return handleScriptPubkeyVerify(registry, req)

	// Transaction operations
	case "btck_transaction_create":
		return handleTransactionCreate(registry, req)
	case "btck_transaction_destroy":
		return handleTransactionDestroy(registry, req)

	// Transaction output operations
	case "btck_transaction_output_create":
		return handleTransactionOutputCreate(registry, req)
	case "btck_transaction_output_destroy":
		return handleTransactionOutputDestroy(registry, req)

	// Precomputed transaction data operations
	case "btck_precomputed_transaction_data_create":
		return handlePrecomputedTransactionDataCreate(registry, req)
	case "btck_precomputed_transaction_data_destroy":
		return handlePrecomputedTransactionDataDestroy(registry, req)

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
		return Response{}, fmt.Errorf("unknown method: %s", req.Method)
	}
}
