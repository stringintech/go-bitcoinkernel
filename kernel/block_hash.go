package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type blockHashCFuncs struct{}

func (blockHashCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_block_hash_destroy((*C.btck_BlockHash)(ptr))
}

func (blockHashCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_block_hash_copy((*C.btck_BlockHash)(ptr)))
}

// BlockHash is a type-safe identifier for a block.
type BlockHash struct {
	*handle
}

func newBlockHash(ptr *C.btck_BlockHash, fromOwned bool) *BlockHash {
	h := newHandle(unsafe.Pointer(ptr), blockHashCFuncs{}, fromOwned)
	return &BlockHash{handle: h}
}

// NewBlockHash creates a new BlockHash from a 32-byte hash value.
//
// Parameters:
//   - hashBytes: 32-byte array containing the block hash
func NewBlockHash(hashBytes [32]byte) *BlockHash {
	ptr := C.btck_block_hash_create((*C.uchar)(unsafe.Pointer(&hashBytes[0])))
	return newBlockHash(ptr, true)
}

// Bytes returns the 32-byte representation of the block hash.
func (bh *BlockHash) Bytes() [32]byte {
	var output [32]C.uchar
	C.btck_block_hash_to_bytes((*C.btck_BlockHash)(bh.ptr), &output[0])
	return *(*[32]byte)(unsafe.Pointer(&output[0]))
}

// Copy creates a copy of the block hash.
func (bh *BlockHash) Copy() *BlockHash {
	return newBlockHash((*C.btck_BlockHash)(bh.ptr), false)
}
