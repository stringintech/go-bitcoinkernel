package kernel

/*
#include "bitcoinkernel.h"
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
	return &TransactionOutput{
		handle: h,
		transactionOutputApi: transactionOutputApi{
			ptr: func() *C.btck_TransactionOutput {
				return (*C.btck_TransactionOutput)(h.ptr)
			},
		},
	}
}

// NewTransactionOutput creates a transaction output from a script pubkey and an amount.
//
// Parameters:
//   - scriptPubkey: ScriptPubkey defining the conditions to spend this output
//   - amount: The amount associated with the script pubkey for this output
func NewTransactionOutput(scriptPubkey ScriptPubkeyLike, amount int64) *TransactionOutput {
	ptr := C.btck_transaction_output_create(scriptPubkey.cPtr(), C.int64_t(amount))
	return newTransactionOutput(check(ptr), true)
}

type TransactionOutputView struct {
	transactionOutputApi
}

func newTransactionOutputView(ptr *C.btck_TransactionOutput) *TransactionOutputView {
	return &TransactionOutputView{
		transactionOutputApi: transactionOutputApi{
			ptr: func() *C.btck_TransactionOutput {
				return ptr
			},
		},
	}
}

type transactionOutputApi struct {
	ptr func() *C.btck_TransactionOutput
}

func (t *transactionOutputApi) cPtr() *C.btck_TransactionOutput {
	return t.ptr()
}

// TransactionOutputLike is implemented by *TransactionOutput and *TransactionOutputView.
type TransactionOutputLike interface {
	cPtr() *C.btck_TransactionOutput
	Copy() *TransactionOutput
	ScriptPubkey() *ScriptPubkeyView
	Amount() int64
}

var _ TransactionOutputLike = (*TransactionOutput)(nil)
var _ TransactionOutputLike = (*TransactionOutputView)(nil)

// Copy creates a copy of the transaction output.
func (t *transactionOutputApi) Copy() *TransactionOutput {
	return newTransactionOutput(t.ptr(), false)
}

// ScriptPubkey returns the script pubkey of this output.
//
// The returned ScriptPubkeyView is a non-owned pointer valid for the lifetime of
// this transaction output.
func (t *transactionOutputApi) ScriptPubkey() *ScriptPubkeyView {
	ptr := C.btck_transaction_output_get_script_pubkey(t.ptr())
	return newScriptPubkeyView(check(ptr))
}

// Amount returns the amount in the output
func (t *transactionOutputApi) Amount() int64 {
	return int64(C.btck_transaction_output_get_amount(t.ptr()))
}
