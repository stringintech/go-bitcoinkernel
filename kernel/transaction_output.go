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

func (t *transactionOutputApi) Copy() *TransactionOutput {
	return newTransactionOutput(t.ptr, false)
}

func (t *transactionOutputApi) ScriptPubkey() *ScriptPubkeyView {
	ptr := C.btck_transaction_output_get_script_pubkey(t.ptr)
	return newScriptPubkeyView(check(ptr))
}

func (t *transactionOutputApi) Amount() int64 {
	return int64(C.btck_transaction_output_get_amount(t.ptr))
}
