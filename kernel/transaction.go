package kernel

/*
#include "bitcoinkernel.h"
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

// NewTransaction creates a new transaction from raw serialized transaction data.
//
// Parameters:
//   - rawTransaction: Serialized transaction data in Bitcoin's consensus format
//
// Returns an error if the transaction data is malformed or cannot be parsed.
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

// Copy creates a shallow copy of the transaction by incrementing its reference count.
//
// Transactions are reference-counted internally, so this operation is efficient and does
// not duplicate the underlying data.
func (t *transactionApi) Copy() *Transaction {
	return newTransaction(t.ptr, false)
}

// CountInputs returns the number of inputs in the transaction.
func (t *transactionApi) CountInputs() uint64 {
	return uint64(C.btck_transaction_count_inputs(t.ptr))
}

// CountOutputs returns the number of outputs in the transaction.
func (t *transactionApi) CountOutputs() uint64 {
	return uint64(C.btck_transaction_count_outputs(t.ptr))
}

// GetOutput retrieves the output at the specified index.
//
// The returned output is a non-owned view that depends on the lifetime of this transaction.
//
// Parameters:
//   - index: Index of the output to retrieve
//
// Returns an error if the index is out of bounds.
func (t *transactionApi) GetOutput(index uint64) (*TransactionOutputView, error) {
	if index >= t.CountOutputs() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_transaction_get_output_at(t.ptr, C.size_t(index))
	return newTransactionOutputView(check(ptr)), nil
}

// GetInput retrieves the input at the specified index.
//
// The returned input is a non-owned view that depends on the lifetime of this transaction.
//
// Parameters:
//   - index: Index of the input to retrieve
//
// Returns an error if the index is out of bounds.
func (t *transactionApi) GetInput(index uint64) (*TransactionInputView, error) {
	if index >= t.CountInputs() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_transaction_get_input_at(t.ptr, C.size_t(index))
	return newTransactionInputView(check(ptr)), nil
}

// Bytes returns the consensus serialized representation of the transaction.
//
// Returns an error if the serialization fails.
func (t *transactionApi) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, userData unsafe.Pointer) C.int {
		return C.btck_transaction_to_bytes(t.ptr, writer, userData)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize transaction"}
	}
	return bytes, nil
}

// GetTxid returns the txid for this transaction.
func (t *transactionApi) GetTxid() *TxidView {
	ptr := C.btck_transaction_get_txid(t.ptr)
	return newTxidView(check(ptr))
}
