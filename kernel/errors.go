package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrBlockUninitialized        = &UninitializedError{ObjectName: "block"}
	ErrKernelBlockCreate         = &KernelError{Operation: "btck_block_create"}
	ErrKernelBlockCopy           = &KernelError{Operation: "btck_block_copy"}
	ErrKernelBlockGetHash        = &KernelError{Operation: "btck_block_get_hash"}
	ErrKernelBlockGetTransaction = &KernelError{Operation: "btck_block_get_transaction_at"}
	ErrEmptyBlockData            = errors.New("empty block data") // constructor error

	ErrBlockSpentOutputsUninitialized                      = &UninitializedError{ObjectName: "blockSpentOutputs"}
	ErrKernelBlockSpentOutputsCopy                         = &KernelError{Operation: "btck_block_spent_outputs_copy"}
	ErrKernelBlockSpentOutputsGetTransactionSpentOutputsAt = &KernelError{Operation: "btck_block_spent_outputs_get_transaction_spent_outputs_at"}

	ErrBlockTreeEntryUninitialized = &UninitializedError{ObjectName: "blockTreeEntry"}

	ErrBlockHashUninitialized = &UninitializedError{ObjectName: "blockHash"}

	ErrBlockValidationStateUninitialized = &UninitializedError{ObjectName: "blockValidationState"}

	ErrTransactionUninitialized   = &UninitializedError{ObjectName: "transaction"}
	ErrKernelTransactionCreate    = &KernelError{Operation: "btck_transaction_create"}
	ErrKernelTransactionCopy      = &KernelError{Operation: "btck_transaction_copy"}
	ErrKernelTransactionGetOutput = &KernelError{Operation: "btck_transaction_get_output_at"}
	ErrEmptyTransactionData       = errors.New("empty transaction data") // constructor error

	ErrTransactionOutputUninitialized   = &UninitializedError{ObjectName: "transactionOutput"}
	ErrKernelTransactionOutputCreate    = &KernelError{Operation: "btck_transaction_output_create"}
	ErrKernelTransactionOutputCopy      = &KernelError{Operation: "btck_transaction_output_copy"}
	ErrKernelCopyScriptPubkeyFromOutput = &KernelError{Operation: "btck_transaction_output_copy_script_pubkey"}

	ErrTransactionSpentOutputsUninitialized   = &UninitializedError{ObjectName: "transactionSpentOutputs"}
	ErrKernelTransactionSpentOutputsCopy      = &KernelError{Operation: "btck_transaction_spent_outputs_copy"}
	ErrKernelTransactionSpentOutputsGetCoinAt = &KernelError{Operation: "btck_transaction_spent_outputs_get_coin_at"}

	ErrScriptPubkeyUninitialized                 = &UninitializedError{ObjectName: "scriptPubkey"}
	ErrKernelScriptPubkeyCreate                  = &KernelError{Operation: "btck_script_pubkey_create"}
	ErrKernelScriptPubkeyCopy                    = &KernelError{Operation: "btck_script_pubkey_copy"}
	ErrKernelScriptVerify                        = &KernelError{Operation: "btck_verify_script"}
	ErrKernelScriptVerifyInvalidFlags            = &KernelError{Operation: "btck_verify_script", Detail: "the provided bitfield for the flags was invalid"}
	ErrKernelScriptVerifyInvalidFlagsCombination = &KernelError{Operation: "btck_verify_script", Detail: "the flags were combined in an invalid way"}
	ErrKernelScriptVerifySpentOutputsRequired    = &KernelError{Operation: "btck_verify_script", Detail: "the taproot flag was set, so valid spent_outputs have to be provided"}
	ErrEmptyScriptPubkeyData                     = errors.New("empty script pubkey data") // constructor error

	ErrContextUninitialized = &UninitializedError{ObjectName: "context"}
	ErrKernelContextCreate  = &KernelError{Operation: "btck_context_create"}

	ErrContextOptionsUninitialized = &UninitializedError{ObjectName: "contextOptions"}
	ErrKernelContextOptionsCreate  = &KernelError{Operation: "btck_context_options_create"}

	ErrChainstateManagerUninitialized       = &UninitializedError{ObjectName: "chainstateManager"}
	ErrKernelChainstateManagerCreate        = &KernelError{Operation: "btck_chainstate_manager_create"}
	ErrKernelChainstateManagerProcessBlock  = &KernelError{Operation: "btck_chainstate_manager_process_block"}
	ErrKernelImportBlocks                   = &KernelError{Operation: "btck_chainstate_manager_import_blocks"}
	ErrKernelChainstateManagerReadBlock     = &KernelError{Operation: "btck_block_read"}
	ErrKernelChainstateManagerReadBlockUndo = &KernelError{Operation: "btck_block_undo_read"}

	ErrChainstateManagerOptionsUninitialized = &UninitializedError{ObjectName: "chainstateManagerOptions"}
	ErrKernelChainstateManagerOptionsCreate  = &KernelError{Operation: "btck_chainstate_manager_options_create"}

	ErrChainParametersUninitialized = &UninitializedError{ObjectName: "chainParameters"}
	ErrKernelChainParametersCreate  = &KernelError{Operation: "btck_chain_parameters_create"}

	ErrChainUninitialized = &UninitializedError{ObjectName: "chain"}

	ErrCoinUninitialized = &UninitializedError{ObjectName: "coin"}
	ErrKernelCoinCopy    = &KernelError{Operation: "btck_coin_copy"}

	ErrLoggingConnectionUninitialized = &UninitializedError{ObjectName: "loggingConnection"}
	ErrKernelLoggingConnectionCreate  = &KernelError{Operation: "btck_logging_connection_create"}

	ErrInvalidChainType                = errors.New("invalid chain type")
	ErrInvalidLogLevel                 = errors.New("invalid log level")
	ErrInvalidLogCategory              = errors.New("invalid log category")
	ErrInvalidScriptVerifyStatus       = errors.New("invalid script verify status")
	ErrNilNotificationCallbacks        = errors.New("nil notification callbacks")
	ErrNilValidationInterfaceCallbacks = errors.New("nil validation interface callbacks")
)

// UninitializedError is returned when an operation is attempted on a
// Go wrapper struct that has not been properly initialized (i.e., its C pointer is nil).
// This typically indicates a programmer error.
type UninitializedError struct {
	// ObjectName is the type of object that was uninitialized
	ObjectName string
}

func (e *UninitializedError) Error() string {
	return fmt.Sprintf("%s is not initialized", e.ObjectName)
}

// KernelError is returned when a call to the underlying library fails.
type KernelError struct {
	// Operation describes the C-level action that failed
	Operation string
	Detail    string
}

func (e *KernelError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("kernel error during '%s': %s", e.Operation, e.Detail)
	}
	return fmt.Sprintf("kernel error during '%s'", e.Operation)
}
