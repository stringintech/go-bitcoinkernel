package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

// TransactionOutput wraps the C kernel_TransactionOutput
type TransactionOutput struct {
	ptr *C.kernel_TransactionOutput
}

func NewTransactionOutput(scriptPubkey *ScriptPubkey, amount int64) (*TransactionOutput, error) {
	if scriptPubkey == nil || scriptPubkey.ptr == nil {
		return nil, ErrInvalidScriptPubkey
	}

	ptr := C.kernel_transaction_output_create(scriptPubkey.ptr, C.int64_t(amount))
	if ptr == nil {
		return nil, ErrTransactionOutputCreation
	}

	output := &TransactionOutput{ptr: ptr}
	runtime.SetFinalizer(output, (*TransactionOutput).destroy)
	return output, nil
}

// ScriptPubkey returns a copy of the script pubkey from this transaction output
func (t *TransactionOutput) ScriptPubkey() (*ScriptPubkey, error) {
	if t.ptr == nil {
		return nil, ErrInvalidTransactionOutput
	}

	ptr := C.kernel_copy_script_pubkey_from_output(t.ptr)
	if ptr == nil {
		return nil, ErrScriptPubkeyCopyFromOutput
	}

	scriptPubkey := &ScriptPubkey{ptr: ptr}
	runtime.SetFinalizer(scriptPubkey, (*ScriptPubkey).destroy)
	return scriptPubkey, nil
}

func (t *TransactionOutput) Amount() int64 {
	if t.ptr == nil {
		return 0
	}

	return int64(C.kernel_get_transaction_output_amount(t.ptr))
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
