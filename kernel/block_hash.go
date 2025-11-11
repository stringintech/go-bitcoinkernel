package kernel

/*
#include "bitcoinkernel.h"
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
	blockHashApi
}

func newBlockHash(ptr *C.btck_BlockHash, fromOwned bool) *BlockHash {
	h := newHandle(unsafe.Pointer(ptr), blockHashCFuncs{}, fromOwned)
	return &BlockHash{handle: h, blockHashApi: blockHashApi{(*C.btck_BlockHash)(h.ptr)}}
}

// BlockHashView is a type-safe identifier for a block.
type BlockHashView struct {
	blockHashApi
	ptr *C.btck_BlockHash
}

func newBlockHashView(ptr *C.btck_BlockHash) *BlockHashView {
	return &BlockHashView{
		blockHashApi: blockHashApi{ptr},
		ptr:          ptr,
	}
}

type blockHashApi struct {
	ptr *C.btck_BlockHash
}

func (bh *blockHashApi) blockHashPtr() *C.btck_BlockHash {
	return bh.ptr
}

// BlockHashLike is an interface for types that can provide a block hash pointer.
type BlockHashLike interface {
	blockHashPtr() *C.btck_BlockHash
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
func (bh *blockHashApi) Bytes() [32]byte {
	var output [32]C.uchar
	C.btck_block_hash_to_bytes(bh.ptr, &output[0])
	return *(*[32]byte)(unsafe.Pointer(&output[0]))
}

// Copy creates a copy of the block hash.
func (bh *blockHashApi) Copy() *BlockHash {
	return newBlockHash(bh.ptr, false)
}

// Equals checks if two block hashes are equal.
//
// Parameters:
//   - other: Block hash to compare against (can be *BlockHash or *BlockHashView)
//
// Returns true if the block hashes are equal.
func (bh *blockHashApi) Equals(other BlockHashLike) bool {
	return C.btck_block_hash_equals(bh.ptr, other.blockHashPtr()) != 0
}
