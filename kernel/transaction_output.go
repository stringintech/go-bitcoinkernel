package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &TransactionOutput{}

// TransactionOutput wraps the C kernel_TransactionOutput
type TransactionOutput struct {
	ptr *C.kernel_TransactionOutput
}

func NewTransactionOutput(scriptPubkey *ScriptPubkey, amount int64) (*TransactionOutput, error) {
	if err := validateReady(scriptPubkey); err != nil {
		return nil, err
	}

	ptr := C.kernel_transaction_output_create(scriptPubkey.ptr, C.int64_t(amount))
	if ptr == nil {
		return nil, ErrKernelTransactionOutputCreate
	}

	output := &TransactionOutput{ptr: ptr}
	runtime.SetFinalizer(output, (*TransactionOutput).destroy)
	return output, nil
}

// ScriptPubkey returns a copy of the script pubkey from this transaction output
func (t *TransactionOutput) ScriptPubkey() (*ScriptPubkey, error) {
	checkReady(t)

	ptr := C.kernel_transaction_output_copy_script_pubkey(t.ptr)
	if ptr == nil {
		return nil, ErrKernelCopyScriptPubkeyFromOutput
	}

	scriptPubkey := &ScriptPubkey{ptr: ptr}
	runtime.SetFinalizer(scriptPubkey, (*ScriptPubkey).destroy)
	return scriptPubkey, nil
}

func (t *TransactionOutput) Amount() int64 {
	checkReady(t)
	return int64(C.kernel_transaction_output_get_amount(t.ptr))
}

func (t *TransactionOutput) destroy() {
	if t.ptr != nil {
		C.kernel_transaction_output_destroy(t.ptr)
		t.ptr = nil
	}
}

func (t *TransactionOutput) Destroy() {
	runtime.SetFinalizer(t, nil)
	t.destroy()
}

func (t *TransactionOutput) isReady() bool {
	return t != nil && t.ptr != nil
}

func (t *TransactionOutput) uninitializedError() error {
	return ErrTransactionOutputUninitialized
}
