package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ cManagedResource = &BlockHash{}

// BlockHash wraps the C kernel_BlockHash
type BlockHash struct {
	ptr *C.kernel_BlockHash
}

// Bytes returns the raw hash bytes
func (bh *BlockHash) Bytes() []byte {
	checkReady(bh)
	// BlockHash is a 32-byte array in the C struct
	return C.GoBytes(unsafe.Pointer(&bh.ptr.hash[0]), 32)
}

func (bh *BlockHash) destroy() {
	if bh.ptr != nil {
		C.kernel_block_hash_destroy(bh.ptr)
		bh.ptr = nil
	}
}

func (bh *BlockHash) Destroy() {
	runtime.SetFinalizer(bh, nil)
	bh.destroy()
}

func (bh *BlockHash) isReady() bool {
	return bh != nil && bh.ptr != nil
}

func (bh *BlockHash) uninitializedError() error {
	return ErrBlockHashUninitialized
}
