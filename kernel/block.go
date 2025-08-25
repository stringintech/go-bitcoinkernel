package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ cManagedResource = &Block{}

// Block wraps the C btck_Block
type Block struct {
	ptr *C.btck_Block
}

// NewBlockFromRaw creates a new block from raw serialized data
func NewBlockFromRaw(rawBlock []byte) (*Block, error) {
	if len(rawBlock) == 0 {
		return nil, ErrEmptyBlockData
	}
	ptr := C.btck_block_create(unsafe.Pointer(&rawBlock[0]), C.size_t(len(rawBlock)))
	if ptr == nil {
		return nil, ErrKernelBlockCreate
	}
	return newBlockFromPtr(ptr), nil
}

func newBlockFromPtr(ptr *C.btck_Block) *Block {
	block := &Block{ptr: ptr}
	runtime.SetFinalizer(block, (*Block).destroy)
	return block
}

func (b *Block) Hash() (*BlockHash, error) {
	checkReady(b)

	ptr := C.btck_block_get_hash(b.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockGetHash
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

// Bytes returns the serialized block
func (b *Block) Bytes() ([]byte, error) {
	checkReady(b)

	// Use the callback helper to collect bytes from btck_block_to_bytes
	return writeToBytes(func(writer C.btck_WriteBytes, userData unsafe.Pointer) C.int {
		return C.btck_block_to_bytes(b.ptr, writer, userData)
	})
}

// Copy creates a copy of the block
func (b *Block) Copy() (*Block, error) {
	checkReady(b)

	ptr := C.btck_block_copy(b.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockCopy
	}

	return newBlockFromPtr(ptr), nil
}

// CountTransactions returns the number of transactions in the block
func (b *Block) CountTransactions() (uint64, error) {
	checkReady(b)

	count := C.btck_block_count_transactions(b.ptr)
	return uint64(count), nil
}

// GetTransactionAt returns the transaction at the specified index.
func (b *Block) GetTransactionAt(index uint64) (*Transaction, error) {
	checkReady(b)

	ptr := C.btck_block_get_transaction_at(b.ptr, C.size_t(index))
	if ptr == nil {
		return nil, ErrKernelBlockGetTransaction
	}

	transaction := &Transaction{ptr: ptr}
	runtime.SetFinalizer(transaction, (*Transaction).destroy)
	return transaction, nil
}

func (b *Block) destroy() {
	if b.ptr != nil {
		C.btck_block_destroy(b.ptr)
		b.ptr = nil
	}
}

func (b *Block) Destroy() {
	runtime.SetFinalizer(b, nil)
	b.destroy()
}

func (b *Block) isReady() bool {
	return b != nil && b.ptr != nil
}

func (b *Block) uninitializedError() error {
	return ErrBlockUninitialized
}
