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

type BlockHash struct {
	*uniqueHandle
}

func newBlockHash(ptr *C.btck_BlockHash) *BlockHash {
	h := newUniqueHandle(unsafe.Pointer(ptr), blockHashCFuncs{})
	return &BlockHash{uniqueHandle: h}
}

// Bytes returns the raw hash bytes
func (bh *BlockHash) Bytes() []byte {
	// BlockHash is a 32-byte array in the C struct
	return C.GoBytes(unsafe.Pointer(&(*C.btck_BlockHash)(bh.ptr).hash[0]), 32)
}
