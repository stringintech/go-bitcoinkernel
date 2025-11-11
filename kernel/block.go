package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
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
	ptr := C.btck_block_create(unsafe.Pointer(&rawBlock[0]), C.size_t(len(rawBlock)))
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
