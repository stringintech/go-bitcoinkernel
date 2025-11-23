package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"iter"
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

// CountInputs returns the number of inputs in the transaction.
func (t *transactionApi) CountInputs() uint64 {
	return uint64(C.btck_transaction_count_inputs(t.ptr))
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

// Inputs returns an iterator over all inputs in the transaction.
//
// The returned inputs are non-owned views that depend on the lifetime of this transaction.
//
// Example usage:
//
//	for input := range tx.Inputs() {
//	    // Process input
//	}
func (t *transactionApi) Inputs() iter.Seq[*TransactionInputView] {
	return func(yield func(*TransactionInputView) bool) {
		t.iterInputs(0, t.CountInputs(), yield)
	}
}

// InputsRange returns an iterator over a range of inputs in the transaction.
//
// Parameters:
//   - from: Starting index (inclusive)
//   - to: Ending index (exclusive)
//
// The returned inputs are non-owned views that depend on the lifetime of this transaction.
// Safe for out-of-bounds arguments: 'to' is clamped to the count,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for input := range tx.InputsRange(0, 5) {
//	    // Process inputs 0-4
//	}
func (t *transactionApi) InputsRange(from, to uint64) iter.Seq[*TransactionInputView] {
	return func(yield func(*TransactionInputView) bool) {
		if count := t.CountInputs(); to > count {
			to = count
		}
		t.iterInputs(from, to, yield)
	}
}

// InputsFrom returns an iterator over inputs starting from the given index.
//
// Parameters:
//   - from: Starting index (inclusive)
//
// The returned inputs are non-owned views that depend on the lifetime of this transaction.
// If from is beyond the input count, returns an empty iterator.
//
// Example usage:
//
//	for input := range tx.InputsFrom(5) {
//	    // Process inputs from index 5 to the end
//	}
func (t *transactionApi) InputsFrom(from uint64) iter.Seq[*TransactionInputView] {
	return func(yield func(*TransactionInputView) bool) {
		t.iterInputs(from, t.CountInputs(), yield)
	}
}

// iterInputs is a helper that iterates over inputs in [from, to).
func (t *transactionApi) iterInputs(from, to uint64, yield func(*TransactionInputView) bool) {
	for i := from; i < to; i++ {
		input, err := t.GetInput(i)
		if err != nil {
			panic(err)
		}
		if !yield(input) {
			return
		}
	}
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

// Outputs returns an iterator over all outputs in the transaction.
//
// The returned outputs are non-owned views that depend on the lifetime of this transaction.
//
// Example usage:
//
//	for output := range tx.Outputs() {
//	    // Process output
//	}
func (t *transactionApi) Outputs() iter.Seq[*TransactionOutputView] {
	return func(yield func(*TransactionOutputView) bool) {
		t.iterOutputs(0, t.CountOutputs(), yield)
	}
}

// OutputsRange returns an iterator over a range of outputs in the transaction.
//
// Parameters:
//   - from: Starting index (inclusive)
//   - to: Ending index (exclusive)
//
// The returned outputs are non-owned views that depend on the lifetime of this transaction.
// Safe for out-of-bounds arguments: 'to' is clamped to the count,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for output := range tx.OutputsRange(0, 5) {
//	    // Process outputs 0-4
//	}
func (t *transactionApi) OutputsRange(from, to uint64) iter.Seq[*TransactionOutputView] {
	return func(yield func(*TransactionOutputView) bool) {
		if count := t.CountOutputs(); to > count {
			to = count
		}
		t.iterOutputs(from, to, yield)
	}
}

// OutputsFrom returns an iterator over outputs starting from the given index.
//
// Parameters:
//   - from: Starting index (inclusive)
//
// The returned outputs are non-owned views that depend on the lifetime of this transaction.
// If from is beyond the output count, returns an empty iterator.
//
// Example usage:
//
//	for output := range tx.OutputsFrom(5) {
//	    // Process outputs from index 5 to the end
//	}
func (t *transactionApi) OutputsFrom(from uint64) iter.Seq[*TransactionOutputView] {
	return func(yield func(*TransactionOutputView) bool) {
		t.iterOutputs(from, t.CountOutputs(), yield)
	}
}

// iterOutputs is a helper that iterates over outputs in [from, to).
func (t *transactionApi) iterOutputs(from, to uint64, yield func(*TransactionOutputView) bool) {
	for i := from; i < to; i++ {
		output, err := t.GetOutput(i)
		if err != nil {
			panic(err)
		}
		if !yield(output) {
			return
		}
	}
}
