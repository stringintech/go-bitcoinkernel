package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &TransactionSpentOutputs{}

// TransactionSpentOutputs wraps the C btck_TransactionSpentOutputs
type TransactionSpentOutputs struct {
	ptr *C.btck_TransactionSpentOutputs
}

// Size returns the number of spent transaction outputs for the transaction
func (tso *TransactionSpentOutputs) Size() uint64 {
	checkReady(tso)
	return uint64(C.btck_transaction_spent_outputs_size(tso.ptr))
}

// GetCoinAt returns a coin contained in the transaction spent outputs at the specified index
func (tso *TransactionSpentOutputs) GetCoinAt(index uint64) (*Coin, error) {
	checkReady(tso)
	ptr := C.btck_transaction_spent_outputs_get_coin_at(tso.ptr, C.uint64_t(index))
	if ptr == nil {
		return nil, ErrKernelTransactionSpentOutputsGetCoinAt
	}

	coin := &Coin{ptr: ptr}
	runtime.SetFinalizer(coin, (*Coin).destroy)
	return coin, nil
}

// Copy creates a copy of the transaction spent outputs
func (tso *TransactionSpentOutputs) Copy() (*TransactionSpentOutputs, error) {
	checkReady(tso)

	ptr := C.btck_transaction_spent_outputs_copy(tso.ptr)
	if ptr == nil {
		return nil, ErrKernelTransactionSpentOutputsCopy
	}

	outputs := &TransactionSpentOutputs{ptr: ptr}
	runtime.SetFinalizer(outputs, (*TransactionSpentOutputs).destroy)
	return outputs, nil
}

func (tso *TransactionSpentOutputs) destroy() {
	if tso.ptr != nil {
		C.btck_transaction_spent_outputs_destroy(tso.ptr)
		tso.ptr = nil
	}
}

func (tso *TransactionSpentOutputs) Destroy() {
	runtime.SetFinalizer(tso, nil)
	tso.destroy()
}

func (tso *TransactionSpentOutputs) isReady() bool {
	return tso != nil && tso.ptr != nil
}

func (tso *TransactionSpentOutputs) uninitializedError() error {
	return ErrTransactionSpentOutputsUninitialized
}
