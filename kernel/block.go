package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"iter"
	"unsafe"
)

type blockCFuncs struct{}

func (blockCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_block_destroy((*C.btck_Block)(ptr))
}

func (blockCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_block_copy((*C.btck_Block)(ptr)))
}

type Block struct {
	*handle
}

func newBlock(ptr *C.btck_Block, fromOwned bool) *Block {
	h := newHandle(unsafe.Pointer(ptr), blockCFuncs{}, fromOwned)
	return &Block{handle: h}
}

// NewBlock creates a new block from raw serialized consensus block data.
//
// Parameters:
//   - rawBlock: Serialized block data in Bitcoin's consensus format
//
// Returns an error if the block data is malformed or cannot be parsed.
func NewBlock(rawBlock []byte) (*Block, error) {
	ptr := C.btck_block_create(unsafe.Pointer(unsafe.SliceData(rawBlock)), C.size_t(len(rawBlock)))
	if ptr == nil {
		return nil, &InternalError{"Failed to create block from bytes"}
	}
	return newBlock(ptr, true), nil
}

// Hash calculates and returns the hash of this block.
func (b *Block) Hash() *BlockHash {
	return newBlockHash(C.btck_block_get_hash((*C.btck_Block)(b.ptr)), true)
}

// Bytes returns the consensus serialized representation of the block.
//
// Returns an error if the serialization fails.
func (b *Block) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, userData unsafe.Pointer) C.int {
		return C.btck_block_to_bytes((*C.btck_Block)(b.ptr), writer, userData)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize block"}
	}
	return bytes, nil
}

// Copy creates a shallow copy of the block by incrementing its reference count.
//
// Blocks are reference-counted internally,
// so this operation is efficient and does not duplicate the underlying data.
func (b *Block) Copy() *Block {
	return newBlock((*C.btck_Block)(b.ptr), false)
}

// CountTransactions returns the number of transactions contained in the block.
func (b *Block) CountTransactions() uint64 {
	return uint64(C.btck_block_count_transactions((*C.btck_Block)(b.ptr)))
}

// GetTransactionAt retrieves the transaction at the specified index.
//
// The returned transaction is a non-owned view that depends on the lifetime of this Block.
//
// Parameters:
//   - index: Index of the transaction to retrieve
//
// Returns an error if the index is out of bounds.
func (b *Block) GetTransactionAt(index uint64) (*TransactionView, error) {
	if index >= b.CountTransactions() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_block_get_transaction_at((*C.btck_Block)(b.ptr), C.size_t(index))
	return newTransactionView(check(ptr)), nil
}

// Transactions returns an iterator over all transactions in the block.
//
// The returned transactions are non-owned views that depend on the lifetime of this Block.
//
// Example usage:
//
//	for tx := range block.Transactions() {
//	    // Process transaction
//	}
func (b *Block) Transactions() iter.Seq[*TransactionView] {
	return func(yield func(*TransactionView) bool) {
		b.iterTransactions(0, b.CountTransactions(), yield)
	}
}

// TransactionsRange returns an iterator over a range of transactions in the block.
//
// Parameters:
//   - from: Starting index (inclusive)
//   - to: Ending index (exclusive)
//
// The returned transactions are non-owned views that depend on the lifetime of this Block.
// Safe for out-of-bounds arguments: 'to' is clamped to the count,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for tx := range block.TransactionsRange(0, 5) {
//	    // Process transactions 0-4
//	}
func (b *Block) TransactionsRange(from, to uint64) iter.Seq[*TransactionView] {
	return func(yield func(*TransactionView) bool) {
		if count := b.CountTransactions(); to > count {
			to = count
		}
		b.iterTransactions(from, to, yield)
	}
}

// TransactionsFrom returns an iterator over transactions starting from the given index.
//
// Parameters:
//   - from: Starting index (inclusive)
//
// The returned transactions are non-owned views that depend on the lifetime of this Block.
// If from is beyond the transaction count, returns an empty iterator.
//
// Example usage:
//
//	for tx := range block.TransactionsFrom(5) {
//	    // Process transactions from index 5 to the end
//	}
func (b *Block) TransactionsFrom(from uint64) iter.Seq[*TransactionView] {
	return func(yield func(*TransactionView) bool) {
		b.iterTransactions(from, b.CountTransactions(), yield)
	}
}

// iterTransactions is a helper that iterates over transactions in [from, to).
func (b *Block) iterTransactions(from, to uint64, yield func(*TransactionView) bool) {
	for i := from; i < to; i++ {
		tx, err := b.GetTransactionAt(i)
		if err != nil {
			panic(err)
		}
		if !yield(tx) {
			return
		}
	}
}
