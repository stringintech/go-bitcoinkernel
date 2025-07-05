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

	ErrKernelBlockCreate                            = &KernelError{Operation: "kernel_block_create"}
	ErrKernelBlockGetHash                           = &KernelError{Operation: "kernel_block_get_hash"}
	ErrKernelCopyBlockData                          = &KernelError{Operation: "kernel_copy_block_data"}
	ErrKernelScriptPubkeyCreate                     = &KernelError{Operation: "kernel_script_pubkey_create"}
	ErrKernelCopyScriptPubkeyData                   = &KernelError{Operation: "kernel_copy_script_pubkey_data"}
	ErrKernelTransactionOutputCreate                = &KernelError{Operation: "kernel_transaction_output_create"}
	ErrKernelCopyScriptPubkeyFromOutput             = &KernelError{Operation: "kernel_copy_script_pubkey_from_output"}
	ErrKernelTransactionCreate                      = &KernelError{Operation: "kernel_transaction_create"}
	ErrKernelLoggingConnectionCreate                = &KernelError{Operation: "kernel_logging_connection_create"}
	ErrKernelGetUndoOutputByIndex                   = &KernelError{Operation: "kernel_get_undo_output_by_index"}
	ErrKernelChainstateManagerCreate                = &KernelError{Operation: "kernel_chainstate_manager_create"}
	ErrKernelChainstateManagerOptionsCreate         = &KernelError{Operation: "kernel_chainstate_manager_options_create"}
	ErrKernelContextCreate                          = &KernelError{Operation: "kernel_context_create"}
	ErrKernelChainstateManagerReadBlockFromDisk     = &KernelError{Operation: "kernel_read_block_from_disk"}
	ErrKernelChainstateManagerReadBlockUndoFromDisk = &KernelError{Operation: "kernel_read_block_undo_from_disk"}
	ErrKernelChainstateManagerProcessBlock          = &KernelError{Operation: "kernel_chainstate_manager_process_block"}
	ErrKernelImportBlocks                           = &KernelError{Operation: "kernel_import_blocks"}
	ErrKernelChainParametersCreate                  = &KernelError{Operation: "kernel_chain_parameters_create"}
	ErrKernelContextOptionsCreate                   = &KernelError{Operation: "kernel_context_options_create"}

	ErrInvalidChainType      = errors.New("invalid chain type")
	ErrInvalidLogLevel       = errors.New("invalid log level")
	ErrInvalidLogCategory    = errors.New("invalid log category")
	ErrEmptyBlockData        = errors.New("empty block data")
	ErrEmptyScriptPubkeyData = errors.New("empty script pubkey data")
	ErrEmptyTransactionData  = errors.New("empty transaction data")
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
}

func (e *KernelError) Error() string {
	return fmt.Sprintf("kernel error during '%s'", e.Operation)
}
