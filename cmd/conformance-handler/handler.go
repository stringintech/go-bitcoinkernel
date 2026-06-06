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
	case "btck_script_pubkey_to_bytes":
		return handleScriptPubkeyToBytes(registry, req)
	case "btck_script_pubkey_copy":
		return handleScriptPubkeyCopy(registry, req)
	case "btck_script_pubkey_verify":
		return handleScriptPubkeyVerify(registry, req)

	// Transaction operations
	case "btck_transaction_create":
		return handleTransactionCreate(registry, req)
	case "btck_transaction_destroy":
		return handleTransactionDestroy(registry, req)
	case "btck_transaction_get_txid":
		return handleTransactionGetTxid(registry, req)
	case "btck_txid_to_bytes":
		return handleTxidToBytes(registry, req)
	case "btck_txid_equals":
		return handleTxidEquals(registry, req)
	case "btck_txid_copy":
		return handleTxidCopy(registry, req)
	case "btck_txid_destroy":
		return handleTxidDestroy(registry, req)
	case "btck_transaction_count_inputs":
		return handleTransactionCountInputs(registry, req)
	case "btck_transaction_count_outputs":
		return handleTransactionCountOutputs(registry, req)
	case "btck_transaction_to_bytes":
		return handleTransactionToBytes(registry, req)
	case "btck_transaction_get_output_at":
		return handleTransactionGetOutputAt(registry, req)
	case "btck_transaction_get_input_at":
		return handleTransactionGetInputAt(registry, req)
	case "btck_transaction_copy":
		return handleTransactionCopy(registry, req)
	case "btck_transaction_input_destroy":
		return handleTransactionInputDestroy(registry, req)
	case "btck_transaction_input_get_out_point":
		return handleTransactionInputGetOutPoint(registry, req)
	case "btck_transaction_input_copy":
		return handleTransactionInputCopy(registry, req)
	case "btck_transaction_out_point_get_txid":
		return handleTransactionOutPointGetTxid(registry, req)
	case "btck_transaction_out_point_get_index":
		return handleTransactionOutPointGetIndex(registry, req)
	case "btck_transaction_out_point_copy":
		return handleTransactionOutPointCopy(registry, req)
	case "btck_transaction_out_point_destroy":
		return handleTransactionOutPointDestroy(registry, req)
	case "btck_transaction_output_create":
		return handleTransactionOutputCreate(registry, req)
	case "btck_transaction_output_destroy":
		return handleTransactionOutputDestroy(registry, req)
	case "btck_transaction_output_copy":
		return handleTransactionOutputCopy(registry, req)
	case "btck_transaction_output_get_amount":
		return handleTransactionOutputGetAmount(registry, req)
	case "btck_transaction_output_get_script_pubkey":
		return handleTransactionOutputGetScriptPubkey(registry, req)

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
	case "btck_block_destroy":
		return handleBlockDestroy(registry, req)
	case "btck_block_create":
		return handleBlockCreate(registry, req)
	case "btck_block_get_hash":
		return handleBlockGetHash(registry, req)
	case "btck_block_get_header":
		return handleBlockGetHeader(registry, req)
	case "btck_block_count_transactions":
		return handleBlockCountTransactions(registry, req)
	case "btck_block_get_transaction_at":
		return handleBlockGetTransactionAt(registry, req)
	case "btck_block_to_bytes":
		return handleBlockToBytes(registry, req)
	case "btck_block_copy":
		return handleBlockCopy(registry, req)

	// BlockTreeEntry operations
	case "btck_block_tree_entry_get_block_hash":
		return handleBlockTreeEntryGetBlockHash(registry, req)
	case "btck_block_tree_entry_get_height":
		return handleBlockTreeEntryGetHeight(registry, req)

	// Notification callbacks operations
	case "notification_callbacks_create":
		return handleNotificationCallbacksCreate(registry, req)
	case "notification_callbacks_drain":
		return handleNotificationCallbacksDrain(registry, req)

	// Validation interface callbacks operations
	case "validation_interface_callbacks_create":
		return handleValidationInterfaceCallbacksCreate(registry, req)
	case "validation_callbacks_drain":
		return handleValidationCallbacksDrain(registry, req)

	// Block validation state operations
	case "btck_block_validation_state_get_validation_mode":
		return handleBlockValidationStateGetValidationMode(registry, req)
	case "btck_block_validation_state_destroy":
		return handleBlockValidationStateDestroy(registry, req)

	// Block hash operations
	case "btck_block_hash_create":
		return handleBlockHashCreate(registry, req)
	case "btck_block_hash_to_bytes":
		return handleBlockHashToBytes(registry, req)
	case "btck_block_hash_equals":
		return handleBlockHashEquals(registry, req)
	case "btck_block_hash_copy":
		return handleBlockHashCopy(registry, req)
	case "btck_block_hash_destroy":
		return handleBlockHashDestroy(registry, req)

	// Block header operations
	case "btck_block_header_create":
		return handleBlockHeaderCreate(registry, req)
	case "btck_block_header_get_hash":
		return handleBlockHeaderGetHash(registry, req)
	case "btck_block_header_get_prev_hash":
		return handleBlockHeaderGetPrevHash(registry, req)
	case "btck_block_header_get_timestamp":
		return handleBlockHeaderGetTimestamp(registry, req)
	case "btck_block_header_get_bits":
		return handleBlockHeaderGetBits(registry, req)
	case "btck_block_header_get_version":
		return handleBlockHeaderGetVersion(registry, req)
	case "btck_block_header_get_nonce":
		return handleBlockHeaderGetNonce(registry, req)
	case "btck_block_header_to_bytes":
		return handleBlockHeaderToBytes(registry, req)
	case "btck_block_header_copy":
		return handleBlockHeaderCopy(registry, req)
	case "btck_block_header_destroy":
		return handleBlockHeaderDestroy(registry, req)

	default:
		return Response{}, fmt.Errorf("unknown method: %s", req.Method)
	}
}
