package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ cManagedResource = &Transaction{}

// Transaction wraps the C kernel_Transaction
type Transaction struct {
	ptr *C.kernel_Transaction
}

// NewTransactionFromRaw creates a new transaction from raw serialized data
func NewTransactionFromRaw(rawTransaction []byte) (*Transaction, error) {
	if len(rawTransaction) == 0 {
		return nil, ErrEmptyTransactionData
	}
	ptr := C.kernel_transaction_create((*C.uchar)(unsafe.Pointer(&rawTransaction[0])), C.size_t(len(rawTransaction)))
	if ptr == nil {
		return nil, ErrKernelTransactionCreate
	}

	transaction := &Transaction{ptr: ptr}
	runtime.SetFinalizer(transaction, (*Transaction).destroy)
	return transaction, nil
}

func (t *Transaction) destroy() {
	if t.ptr != nil {
		C.kernel_transaction_destroy(t.ptr)
		t.ptr = nil
	}
}

func (t *Transaction) Destroy() {
	runtime.SetFinalizer(t, nil)
	t.destroy()
}

func (t *Transaction) isReady() bool {
	return t != nil && t.ptr != nil
}

func (t *Transaction) uninitializedError() error {
	return ErrTransactionUninitialized
}
