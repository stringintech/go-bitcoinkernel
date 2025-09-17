package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type transactionCFuncs struct{}

func (transactionCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_destroy((*C.btck_Transaction)(ptr))
}

func (transactionCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_copy((*C.btck_Transaction)(ptr)))
}

type Transaction struct {
	*handle
	transactionApi
}

func newTransaction(ptr *C.btck_Transaction, fromOwned bool) *Transaction {
	h := newHandle(unsafe.Pointer(ptr), transactionCFuncs{}, fromOwned)
	return &Transaction{handle: h, transactionApi: transactionApi{(*C.btck_Transaction)(h.ptr)}}
}

// NewTransaction creates a new transaction from raw serialized data
func NewTransaction(rawTransaction []byte) (*Transaction, error) {
	ptr := C.btck_transaction_create(unsafe.Pointer(&rawTransaction[0]), C.size_t(len(rawTransaction)))
	if ptr == nil {
		return nil, &InternalError{"Failed to create transaction from bytes"}
	}
	return newTransaction(ptr, true), nil
}

type TransactionView struct {
	transactionApi
	ptr *C.btck_Transaction
}

func newTransactionView(ptr *C.btck_Transaction) *TransactionView {
	return &TransactionView{
		transactionApi: transactionApi{ptr},
		ptr:            ptr,
	}
}

type transactionApi struct {
	ptr *C.btck_Transaction
}

func (t *transactionApi) Copy() *Transaction {
	return newTransaction(t.ptr, false)
}

func (t *transactionApi) CountInputs() uint64 {
	return uint64(C.btck_transaction_count_inputs(t.ptr))
}

func (t *transactionApi) CountOutputs() uint64 {
	return uint64(C.btck_transaction_count_outputs(t.ptr))
}

func (t *transactionApi) GetOutput(index uint64) (*TransactionOutputView, error) {
	if index >= t.CountOutputs() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_transaction_get_output_at(t.ptr, C.size_t(index))
	return newTransactionOutputView(check(ptr)), nil
}

// Bytes returns the consensus serialized transaction
func (t *transactionApi) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, userData unsafe.Pointer) C.int {
		return C.btck_transaction_to_bytes(t.ptr, writer, userData)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize transaction"}
	}
	return bytes, nil
}
