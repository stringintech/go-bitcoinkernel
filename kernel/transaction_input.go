package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type transactionInputCFuncs struct{}

func (transactionInputCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_input_destroy((*C.btck_TransactionInput)(ptr))
}

func (transactionInputCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_input_copy((*C.btck_TransactionInput)(ptr)))
}

// TransactionInput holds information on the TransactionOutPoint held within.
type TransactionInput struct {
	*handle
	transactionInputApi
}

func newTransactionInput(ptr *C.btck_TransactionInput, fromOwned bool) *TransactionInput {
	h := newHandle(unsafe.Pointer(ptr), transactionInputCFuncs{}, fromOwned)
	return &TransactionInput{handle: h, transactionInputApi: transactionInputApi{(*C.btck_TransactionInput)(h.ptr)}}
}

// TransactionInputView holds information on the TransactionOutPoint held within.
type TransactionInputView struct {
	transactionInputApi
	ptr *C.btck_TransactionInput
}

func newTransactionInputView(ptr *C.btck_TransactionInput) *TransactionInputView {
	return &TransactionInputView{
		transactionInputApi: transactionInputApi{ptr},
		ptr:                 ptr,
	}
}

type transactionInputApi struct {
	ptr *C.btck_TransactionInput
}

// Copy creates a copy of the transaction input.
func (t *transactionInputApi) Copy() *TransactionInput {
	return newTransactionInput(t.ptr, false)
}

// GetOutPoint returns the transaction out point.
func (t *transactionInputApi) GetOutPoint() *TransactionOutPointView {
	ptr := C.btck_transaction_input_get_out_point(t.ptr)
	return newTransactionOutPointView(check(ptr))
}
