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

type BlockHash struct {
	*handle
}

func newBlockHash(ptr *C.btck_BlockHash, fromOwned bool) *BlockHash {
	h := newHandle(unsafe.Pointer(ptr), blockHashCFuncs{}, fromOwned)
	return &BlockHash{handle: h}
}

// NewBlockHash creates a new BlockHash from raw 32-byte hash data
func NewBlockHash(hashBytes [32]byte) *BlockHash {
	ptr := C.btck_block_hash_create((*C.uchar)(unsafe.Pointer(&hashBytes[0])))
	return newBlockHash(ptr, true)
}

// Bytes returns the raw hash bytes
func (bh *BlockHash) Bytes() [32]byte {
	var output [32]C.uchar
	C.btck_block_hash_to_bytes((*C.btck_BlockHash)(bh.ptr), &output[0])
	return *(*[32]byte)(unsafe.Pointer(&output[0]))
}

func (bh *BlockHash) Copy() *BlockHash {
	return newBlockHash((*C.btck_BlockHash)(bh.ptr), false)
}
