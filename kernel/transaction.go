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

// Transaction wraps the C btck_Transaction
type Transaction struct {
	ptr *C.btck_Transaction
}

// NewTransactionFromRaw creates a new transaction from raw serialized data
func NewTransactionFromRaw(rawTransaction []byte) (*Transaction, error) {
	if len(rawTransaction) == 0 {
		return nil, ErrEmptyTransactionData
	}
	ptr := C.btck_transaction_create(unsafe.Pointer(&rawTransaction[0]), C.size_t(len(rawTransaction)))
	if ptr == nil {
		return nil, ErrKernelTransactionCreate
	}

	transaction := &Transaction{ptr: ptr}
	runtime.SetFinalizer(transaction, (*Transaction).destroy)
	return transaction, nil
}

// Copy creates a copy of the transaction. Transactions are reference counted,
// so this just increments the reference count.
func (t *Transaction) Copy() (*Transaction, error) {
	checkReady(t)

	ptr := C.btck_transaction_copy(t.ptr)
	if ptr == nil {
		return nil, ErrKernelTransactionCopy
	}

	transaction := &Transaction{ptr: ptr}
	runtime.SetFinalizer(transaction, (*Transaction).destroy)
	return transaction, nil
}

// CountInputs returns the number of inputs in the transaction
func (t *Transaction) CountInputs() (uint64, error) {
	checkReady(t)

	count := C.btck_transaction_count_inputs(t.ptr)
	return uint64(count), nil
}

// CountOutputs returns the number of outputs in the transaction
func (t *Transaction) CountOutputs() (uint64, error) {
	checkReady(t)

	count := C.btck_transaction_count_outputs(t.ptr)
	return uint64(count), nil
}

// GetOutputAt returns the transaction output at the specified index.
func (t *Transaction) GetOutputAt(index uint64) (*TransactionOutput, error) {
	checkReady(t)

	ptr := C.btck_transaction_get_output_at(t.ptr, C.uint64_t(index))
	if ptr == nil {
		return nil, ErrKernelTransactionGetOutput
	}

	output := &TransactionOutput{ptr: ptr}
	runtime.SetFinalizer(output, (*TransactionOutput).destroy)
	return output, nil
}

// ToBytes serializes the transaction to bytes using consensus serialization
func (t *Transaction) ToBytes() ([]byte, error) {
	checkReady(t)

	return writeToBytes(func(writer C.btck_WriteBytes, userData unsafe.Pointer) C.int {
		return C.btck_transaction_to_bytes(t.ptr, writer, userData)
	})
}

func (t *Transaction) destroy() {
	if t.ptr != nil {
		C.btck_transaction_destroy(t.ptr)
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
