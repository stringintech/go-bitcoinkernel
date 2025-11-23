package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"iter"
	"unsafe"
)

type blockSpentOutputsCFuncs struct{}

func (blockSpentOutputsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_block_spent_outputs_destroy((*C.btck_BlockSpentOutputs)(ptr))
}

func (blockSpentOutputsCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_block_spent_outputs_copy((*C.btck_BlockSpentOutputs)(ptr)))
}

// BlockSpentOutputs holds all the previous outputs consumed by all transactions
// in a specific block.
//
// Internally it holds a nested vector. The top level vector has an entry for each
// transaction in a block (in order of the actual transactions of the block and
// without the coinbase transaction). This is exposed through TransactionSpentOutputs.
// Each TransactionSpentOutputs is in turn a vector of all the previous outputs of a
// transaction (in order of their corresponding inputs).
type BlockSpentOutputs struct {
	*handle
}

func newBlockSpentOutputs(ptr *C.btck_BlockSpentOutputs, fromOwned bool) *BlockSpentOutputs {
	h := newHandle(unsafe.Pointer(ptr), blockSpentOutputsCFuncs{}, fromOwned)
	return &BlockSpentOutputs{handle: h}
}

// Count returns the number of transaction spent outputs contained in this block's spent outputs.
func (bso *BlockSpentOutputs) Count() uint64 {
	return uint64(C.btck_block_spent_outputs_count((*C.btck_BlockSpentOutputs)(bso.ptr)))
}

// GetTransactionSpentOutputsAt retrieves the spent outputs for a specific transaction.
//
// The returned TransactionSpentOutputsView is a non-owned pointer that depends on
// the lifetime of this BlockSpentOutputs.
//
// Parameters:
//   - index: Index of the transaction spent outputs to retrieve
//
// Returns an error if the index is out of bounds.
func (bso *BlockSpentOutputs) GetTransactionSpentOutputsAt(index uint64) (*TransactionSpentOutputsView, error) {
	if index >= bso.Count() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_block_spent_outputs_get_transaction_spent_outputs_at((*C.btck_BlockSpentOutputs)(bso.ptr), C.size_t(index))
	return newTransactionSpentOutputsView(check(ptr)), nil
}

// Copy creates a shallow copy of the block spent outputs by incrementing its reference count.
//
// The block spent outputs is reference-counted internally, so this operation is efficient
// and does not duplicate the underlying data.
func (bso *BlockSpentOutputs) Copy() *BlockSpentOutputs {
	return newBlockSpentOutputs((*C.btck_BlockSpentOutputs)(bso.ptr), false)
}

// TransactionsSpentOutputs returns an iterator over all transaction spent outputs in the block.
//
// The returned transaction spent outputs are non-owned views that depend on the lifetime of this BlockSpentOutputs.
//
// Example usage:
//
//	for txSpentOutputs := range blockSpentOutputs.TransactionsSpentOutputs() {
//	    // Process transaction spent outputs
//	}
func (bso *BlockSpentOutputs) TransactionsSpentOutputs() iter.Seq[*TransactionSpentOutputsView] {
	return func(yield func(*TransactionSpentOutputsView) bool) {
		bso.iterTransactionsSpentOutputs(0, bso.Count(), yield)
	}
}

// TransactionsSpentOutputsRange returns an iterator over a range of transaction spent outputs in the block.
//
// Parameters:
//   - from: Starting index (inclusive)
//   - to: Ending index (exclusive)
//
// The returned transaction spent outputs are non-owned views that depend on the lifetime of this BlockSpentOutputs.
// Safe for out-of-bounds arguments: 'to' is clamped to the count,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for txSpentOutputs := range blockSpentOutputs.TransactionsSpentOutputsRange(0, 5) {
//	    // Process transaction spent outputs 0-4
//	}
func (bso *BlockSpentOutputs) TransactionsSpentOutputsRange(from, to uint64) iter.Seq[*TransactionSpentOutputsView] {
	return func(yield func(*TransactionSpentOutputsView) bool) {
		if count := bso.Count(); to > count {
			to = count
		}
		bso.iterTransactionsSpentOutputs(from, to, yield)
	}
}

// TransactionsSpentOutputsFrom returns an iterator over transaction spent outputs starting from the given index.
//
// Parameters:
//   - from: Starting index (inclusive)
//
// The returned transaction spent outputs are non-owned views that depend on the lifetime of this BlockSpentOutputs.
// If from is beyond the count, returns an empty iterator.
//
// Example usage:
//
//	for txSpentOutputs := range blockSpentOutputs.TransactionsSpentOutputsFrom(5) {
//	    // Process transaction spent outputs from index 5 to the end
//	}
func (bso *BlockSpentOutputs) TransactionsSpentOutputsFrom(from uint64) iter.Seq[*TransactionSpentOutputsView] {
	return func(yield func(*TransactionSpentOutputsView) bool) {
		bso.iterTransactionsSpentOutputs(from, bso.Count(), yield)
	}
}

// iterTransactionsSpentOutputs is a helper that iterates over transaction spent outputs in [from, to).
func (bso *BlockSpentOutputs) iterTransactionsSpentOutputs(from, to uint64, yield func(*TransactionSpentOutputsView) bool) {
	for i := from; i < to; i++ {
		txSpentOutputs, err := bso.GetTransactionSpentOutputsAt(i)
		if err != nil {
			panic(err)
		}
		if !yield(txSpentOutputs) {
			return
		}
	}
}
