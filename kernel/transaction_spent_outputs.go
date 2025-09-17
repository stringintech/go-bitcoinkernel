package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type transactionSpentOutputsCFuncs struct{}

func (transactionSpentOutputsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_spent_outputs_destroy((*C.btck_TransactionSpentOutputs)(ptr))
}

func (transactionSpentOutputsCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_spent_outputs_copy((*C.btck_TransactionSpentOutputs)(ptr)))
}

type TransactionSpentOutputs struct {
	*handle
	transactionSpentOutputsApi
}

func newTransactionSpentOutputs(ptr *C.btck_TransactionSpentOutputs, fromOwned bool) *TransactionSpentOutputs {
	h := newHandle(unsafe.Pointer(ptr), transactionSpentOutputsCFuncs{}, fromOwned)
	return &TransactionSpentOutputs{handle: h, transactionSpentOutputsApi: transactionSpentOutputsApi{(*C.btck_TransactionSpentOutputs)(h.ptr)}}
}

type TransactionSpentOutputsView struct {
	transactionSpentOutputsApi
	ptr *C.btck_TransactionSpentOutputs
}

func newTransactionSpentOutputsView(ptr *C.btck_TransactionSpentOutputs) *TransactionSpentOutputsView {
	return &TransactionSpentOutputsView{
		transactionSpentOutputsApi: transactionSpentOutputsApi{ptr},
		ptr:                        ptr,
	}
}

type transactionSpentOutputsApi struct {
	ptr *C.btck_TransactionSpentOutputs
}

func (t *transactionSpentOutputsApi) Copy() *TransactionSpentOutputs {
	return newTransactionSpentOutputs(t.ptr, false)
}

func (t *transactionSpentOutputsApi) Count() uint64 {
	return uint64(C.btck_transaction_spent_outputs_count(t.ptr))
}

func (t *transactionSpentOutputsApi) GetCoinAt(index uint64) (*CoinView, error) {
	if index >= t.Count() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_transaction_spent_outputs_get_coin_at(t.ptr, C.size_t(index))
	return newCoinView(check(ptr)), nil
}
