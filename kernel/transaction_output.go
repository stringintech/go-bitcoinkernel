package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type transactionOutputCFuncs struct{}

func (transactionOutputCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_output_destroy((*C.btck_TransactionOutput)(ptr))
}

func (transactionOutputCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_output_copy((*C.btck_TransactionOutput)(ptr)))
}

type TransactionOutput struct {
	*handle
	transactionOutputApi
}

func newTransactionOutput(ptr *C.btck_TransactionOutput, fromOwned bool) *TransactionOutput {
	h := newHandle(unsafe.Pointer(ptr), transactionOutputCFuncs{}, fromOwned)
	return &TransactionOutput{handle: h, transactionOutputApi: transactionOutputApi{(*C.btck_TransactionOutput)(h.ptr)}}
}

// NewTransactionOutput creates a transaction output from a script pubkey and an amount.
//
// Parameters:
//   - scriptPubkey: ScriptPubkey defining the conditions to spend this output
//   - amount: The amount associated with the script pubkey for this output
func NewTransactionOutput(scriptPubkey *ScriptPubkey, amount int64) *TransactionOutput {
	ptr := C.btck_transaction_output_create((*C.btck_ScriptPubkey)(scriptPubkey.handle.ptr), C.int64_t(amount))
	return newTransactionOutput(check(ptr), true)
}

type TransactionOutputView struct {
	transactionOutputApi
	ptr *C.btck_TransactionOutput
}

func newTransactionOutputView(ptr *C.btck_TransactionOutput) *TransactionOutputView {
	return &TransactionOutputView{
		transactionOutputApi: transactionOutputApi{ptr},
		ptr:                  ptr,
	}
}

type transactionOutputApi struct {
	ptr *C.btck_TransactionOutput
}

// Copy creates a copy of the transaction output.
func (t *transactionOutputApi) Copy() *TransactionOutput {
	return newTransactionOutput(t.ptr, false)
}

// ScriptPubkey returns the script pubkey of this output.
//
// The returned ScriptPubkeyView is a non-owned pointer valid for the lifetime of
// this transaction output.
func (t *transactionOutputApi) ScriptPubkey() *ScriptPubkeyView {
	ptr := C.btck_transaction_output_get_script_pubkey(t.ptr)
	return newScriptPubkeyView(check(ptr))
}

// Amount returns the amount in the output
func (t *transactionOutputApi) Amount() int64 {
	return int64(C.btck_transaction_output_get_amount(t.ptr))
}
