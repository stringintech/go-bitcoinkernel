package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrBlockUninitialized                    = &UninitializedError{ObjectName: "block"}
	ErrBlockHashUninitialized                = &UninitializedError{ObjectName: "blockHash"}
	ErrBlockIndexUninitialized               = &UninitializedError{ObjectName: "blockIndex"}
	ErrScriptPubkeyUninitialized             = &UninitializedError{ObjectName: "scriptPubkey"}
	ErrTransactionOutputUninitialized        = &UninitializedError{ObjectName: "transactionOutput"}
	ErrBlockUndoUninitialized                = &UninitializedError{ObjectName: "blockUndo"}
	ErrChainstateManagerUninitialized        = &UninitializedError{ObjectName: "chainstateManager"}
	ErrContextUninitialized                  = &UninitializedError{ObjectName: "context"}
	ErrChainParametersUninitialized          = &UninitializedError{ObjectName: "chainParameters"}
	ErrContextOptionsUninitialized           = &UninitializedError{ObjectName: "contextOptions"}
	ErrChainstateManagerOptionsUninitialized = &UninitializedError{ObjectName: "chainstateManagerOptions"}
	ErrBlockValidationStateUninitialized     = &UninitializedError{ObjectName: "blockValidationState"}
	ErrTransactionUninitialized              = &UninitializedError{ObjectName: "transaction"}
	ErrLoggingConnectionUninitialized        = &UninitializedError{ObjectName: "loggingConnection"}

	ErrKernelBlockCreate                    = &KernelError{Operation: "kernel_block_create"}
	ErrKernelBlockGetHash                   = &KernelError{Operation: "kernel_block_get_hash"}
	ErrKernelCopyBlockData                  = &KernelError{Operation: "kernel_block_copy_data"}
	ErrKernelScriptPubkeyCreate             = &KernelError{Operation: "kernel_script_pubkey_create"}
	ErrKernelCopyScriptPubkeyData           = &KernelError{Operation: "kernel_script_pubkey_copy_data"}
	ErrKernelTransactionOutputCreate        = &KernelError{Operation: "kernel_transaction_output_create"}
	ErrKernelCopyScriptPubkeyFromOutput     = &KernelError{Operation: "kernel_transaction_output_copy_script_pubkey"}
	ErrKernelTransactionCreate              = &KernelError{Operation: "kernel_transaction_create"}
	ErrKernelLoggingConnectionCreate        = &KernelError{Operation: "kernel_logging_connection_create"}
	ErrKernelGetUndoOutputByIndex           = &KernelError{Operation: "kernel_block_undo_copy_transaction_output_by_index"}
	ErrKernelChainstateManagerCreate        = &KernelError{Operation: "kernel_chainstate_manager_create"}
	ErrKernelChainstateManagerOptionsCreate = &KernelError{Operation: "kernel_chainstate_manager_options_create"}
	ErrKernelContextCreate                  = &KernelError{Operation: "kernel_context_create"}
	ErrKernelChainstateManagerReadBlock     = &KernelError{Operation: "kernel_block_read"}
	ErrKernelChainstateManagerReadBlockUndo = &KernelError{Operation: "kernel_block_undo_read"}
	ErrKernelChainstateManagerProcessBlock  = &KernelError{Operation: "kernel_chainstate_manager_process_block"}
	ErrKernelImportBlocks                   = &KernelError{Operation: "kernel_chainstate_manager_import_blocks"}
	ErrKernelChainParametersCreate          = &KernelError{Operation: "kernel_chain_parameters_create"}
	ErrKernelContextOptionsCreate           = &KernelError{Operation: "kernel_context_options_create"}

	ErrScriptVerify                        = &KernelError{Operation: "kernel_verify_script"}
	ErrScriptVerifyTxInputIndex            = &KernelError{Operation: "kernel_verify_script", Detail: "the provided input index is out of range of the actual number of inputs of the transaction"}
	ErrScriptVerifyInvalidFlags            = &KernelError{Operation: "kernel_verify_script", Detail: "the provided bitfield for the flags was invalid"}
	ErrScriptVerifyInvalidFlagsCombination = &KernelError{Operation: "kernel_verify_script", Detail: "the flags were combined in an invalid way"}
	ErrScriptVerifySpentOutputsRequired    = &KernelError{Operation: "kernel_verify_script", Detail: "the taproot flag was set, so valid spent_outputs have to be provided"}
	ErrScriptVerifySpentOutputsMismatch    = &KernelError{Operation: "kernel_verify_script", Detail: "the number of spent outputs does not match the number of inputs of the transaction"}

	ErrInvalidChainType                = errors.New("invalid chain type")
	ErrInvalidLogLevel                 = errors.New("invalid log level")
	ErrInvalidLogCategory              = errors.New("invalid log category")
	ErrInvalidScriptVerifyStatus       = errors.New("invalid script verify status")
	ErrEmptyBlockData                  = errors.New("empty block data")
	ErrEmptyScriptPubkeyData           = errors.New("empty script pubkey data")
	ErrEmptyTransactionData            = errors.New("empty transaction data")
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
