package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"iter"
	"unsafe"
)

type transactionSpentOutputsCFuncs struct{}

func (transactionSpentOutputsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_spent_outputs_destroy((*C.btck_TransactionSpentOutputs)(ptr))
}

func (transactionSpentOutputsCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_spent_outputs_copy((*C.btck_TransactionSpentOutputs)(ptr)))
}

// TransactionSpentOutputs holds the coins consumed by a transaction.
//
// Retrieved through BlockSpentOutputs. The coins are in the same order as the
// transaction's inputs consuming them.
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

// Copy creates a copy of the transaction spent outputs.
func (t *transactionSpentOutputsApi) Copy() *TransactionSpentOutputs {
	return newTransactionSpentOutputs(t.ptr, false)
}

// Count returns the number of previous transaction outputs contained in the transaction spent outputs data.
func (t *transactionSpentOutputsApi) Count() uint64 {
	return uint64(C.btck_transaction_spent_outputs_count(t.ptr))
}

// GetCoinAt returns a coin contained in the transaction spent outputs at a
// certain index. The returned CoinView is unowned and only valid for the
// lifetime of transaction_spent_outputs.
//
// Parameters:
//   - index: The index of the to be retrieved coin within the transaction spent outputs
//
// Returns an error if the index is out of bounds.
func (t *transactionSpentOutputsApi) GetCoinAt(index uint64) (*CoinView, error) {
	if index >= t.Count() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_transaction_spent_outputs_get_coin_at(t.ptr, C.size_t(index))
	return newCoinView(check(ptr)), nil
}

// Coins returns an iterator over all coins in the transaction spent outputs.
//
// The returned coins are non-owned views that depend on the lifetime of this transaction spent outputs.
//
// Example usage:
//
//	for coin := range tso.Coins() {
//	    // Process coin
//	}
func (t *transactionSpentOutputsApi) Coins() iter.Seq[*CoinView] {
	return func(yield func(*CoinView) bool) {
		t.iterCoins(0, t.Count(), yield)
	}
}

// CoinsRange returns an iterator over a range of coins in the transaction spent outputs.
//
// Parameters:
//   - from: Starting index (inclusive)
//   - to: Ending index (exclusive)
//
// The returned coins are non-owned views that depend on the lifetime of this transaction spent outputs.
// Safe for out-of-bounds arguments: 'to' is clamped to the count,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for coin := range tso.CoinsRange(0, 5) {
//	    // Process coins 0-4
//	}
func (t *transactionSpentOutputsApi) CoinsRange(from, to uint64) iter.Seq[*CoinView] {
	return func(yield func(*CoinView) bool) {
		if count := t.Count(); to > count {
			to = count
		}
		t.iterCoins(from, to, yield)
	}
}

// CoinsFrom returns an iterator over coins starting from the given index.
//
// Parameters:
//   - from: Starting index (inclusive)
//
// The returned coins are non-owned views that depend on the lifetime of this transaction spent outputs.
// If from is beyond the coin count, returns an empty iterator.
//
// Example usage:
//
//	for coin := range tso.CoinsFrom(5) {
//	    // Process coins from index 5 to the end
//	}
func (t *transactionSpentOutputsApi) CoinsFrom(from uint64) iter.Seq[*CoinView] {
	return func(yield func(*CoinView) bool) {
		t.iterCoins(from, t.Count(), yield)
	}
}

// iterCoins is a helper that iterates over coins in [from, to).
func (t *transactionSpentOutputsApi) iterCoins(from, to uint64, yield func(*CoinView) bool) {
	for i := from; i < to; i++ {
		coin, err := t.GetCoinAt(i)
		if err != nil {
			panic(err)
		}
		if !yield(coin) {
			return
		}
	}
}
