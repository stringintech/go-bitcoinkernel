package kernel

import (
	"errors"
	"fmt"
)

var (
	ErrBlockUninitialized                    = &UninitializedError{ObjectName: "block"}
	ErrBlockHashUninitialized                = &UninitializedError{ObjectName: "blockHash"}
	ErrBlockTreeEntryUninitialized           = &UninitializedError{ObjectName: "blockTreeEntry"}
	ErrScriptPubkeyUninitialized             = &UninitializedError{ObjectName: "scriptPubkey"}
	ErrBlockSpentOutputsUninitialized        = &UninitializedError{ObjectName: "blockSpentOutputs"}
	ErrTransactionSpentOutputsUninitialized  = &UninitializedError{ObjectName: "transactionSpentOutputs"}
	ErrCoinUninitialized                     = &UninitializedError{ObjectName: "coin"}
	ErrTransactionOutputUninitialized        = &UninitializedError{ObjectName: "transactionOutput"}
	ErrChainUninitialized                    = &UninitializedError{ObjectName: "chain"}
	ErrChainstateManagerUninitialized        = &UninitializedError{ObjectName: "chainstateManager"}
	ErrContextUninitialized                  = &UninitializedError{ObjectName: "context"}
	ErrChainParametersUninitialized          = &UninitializedError{ObjectName: "chainParameters"}
	ErrContextOptionsUninitialized           = &UninitializedError{ObjectName: "contextOptions"}
	ErrChainstateManagerOptionsUninitialized = &UninitializedError{ObjectName: "chainstateManagerOptions"}
	ErrBlockValidationStateUninitialized     = &UninitializedError{ObjectName: "blockValidationState"}
	ErrTransactionUninitialized              = &UninitializedError{ObjectName: "transaction"}
	ErrLoggingConnectionUninitialized        = &UninitializedError{ObjectName: "loggingConnection"}

	ErrKernelBlockCreate                             = &KernelError{Operation: "btck_block_create"}
	ErrKernelBlockCopy                               = &KernelError{Operation: "btck_block_copy"}
	ErrKernelBlockGetHash                            = &KernelError{Operation: "btck_block_get_hash"}
	ErrKernelBlockGetTransaction                     = &KernelError{Operation: "btck_block_get_transaction_at"}
	ErrKernelScriptPubkeyCreate                      = &KernelError{Operation: "btck_script_pubkey_create"}
	ErrKernelTransactionOutputCreate                 = &KernelError{Operation: "btck_transaction_output_create"}
	ErrKernelCopyScriptPubkeyFromOutput              = &KernelError{Operation: "btck_transaction_output_copy_script_pubkey"}
	ErrKernelTransactionCreate                       = &KernelError{Operation: "btck_transaction_create"}
	ErrBlockSpentOutputsGetTransactionSpentOutputsAt = &KernelError{Operation: "btck_block_spent_outputs_get_transaction_spent_outputs_at"}
	ErrKernelTransactionCopy                         = &KernelError{Operation: "btck_transaction_copy"}
	ErrKernelTransactionGetOutput                    = &KernelError{Operation: "btck_transaction_get_output_at"}
	ErrKernelScriptPubkeyCopy                        = &KernelError{Operation: "btck_script_pubkey_copy"}
	ErrKernelTransactionOutputCopy                   = &KernelError{Operation: "btck_transaction_output_copy"}
	ErrKernelBlockSpentOutputsCopy                   = &KernelError{Operation: "btck_block_spent_outputs_copy"}
	ErrKernelTransactionSpentOutputsCopy             = &KernelError{Operation: "btck_transaction_spent_outputs_copy"}
	ErrKernelTransactionSpentOutputsGetCoinAt        = &KernelError{Operation: "btck_transaction_spent_outputs_get_coin_at"}
	ErrKernelCoinCopy                                = &KernelError{Operation: "btck_coin_copy"}
	ErrKernelLoggingConnectionCreate                 = &KernelError{Operation: "btck_logging_connection_create"}
	ErrKernelChainstateManagerCreate                 = &KernelError{Operation: "btck_chainstate_manager_create"}
	ErrKernelChainstateManagerOptionsCreate          = &KernelError{Operation: "btck_chainstate_manager_options_create"}
	ErrKernelContextCreate                           = &KernelError{Operation: "btck_context_create"}
	ErrKernelChainstateManagerReadBlock              = &KernelError{Operation: "btck_block_read"}
	ErrKernelChainstateManagerReadBlockUndo          = &KernelError{Operation: "btck_block_undo_read"}
	ErrKernelChainstateManagerProcessBlock           = &KernelError{Operation: "btck_chainstate_manager_process_block"}
	ErrKernelImportBlocks                            = &KernelError{Operation: "btck_chainstate_manager_import_blocks"}
	ErrKernelChainParametersCreate                   = &KernelError{Operation: "btck_chain_parameters_create"}
	ErrKernelContextOptionsCreate                    = &KernelError{Operation: "btck_context_options_create"}

	ErrScriptVerify                        = &KernelError{Operation: "btck_verify_script"}
	ErrScriptVerifyInvalidFlags            = &KernelError{Operation: "btck_verify_script", Detail: "the provided bitfield for the flags was invalid"}
	ErrScriptVerifyInvalidFlagsCombination = &KernelError{Operation: "btck_verify_script", Detail: "the flags were combined in an invalid way"}
	ErrScriptVerifySpentOutputsRequired    = &KernelError{Operation: "btck_verify_script", Detail: "the taproot flag was set, so valid spent_outputs have to be provided"}

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
