package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &BlockIndex{}

// BlockIndex wraps the C kernel_BlockIndex
type BlockIndex struct {
	ptr *C.kernel_BlockIndex
}

func (bi *BlockIndex) Height() int32 {
	checkReady(bi)
	return int32(C.kernel_block_index_get_height(bi.ptr))
}

func (bi *BlockIndex) Hash() (*BlockHash, error) {
	checkReady(bi)

	ptr := C.kernel_block_index_get_block_hash(bi.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockGetHash
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

func (bi *BlockIndex) Previous() *BlockIndex {
	checkReady(bi)

	ptr := C.kernel_block_index_get_previous(bi.ptr)
	if ptr == nil {
		return nil
	}

	prevIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(prevIndex, (*BlockIndex).destroy)
	return prevIndex
}

func (bi *BlockIndex) destroy() {
	if bi.ptr != nil {
		C.kernel_block_index_destroy(bi.ptr)
		bi.ptr = nil
	}
}

func (bi *BlockIndex) Destroy() {
	runtime.SetFinalizer(bi, nil)
	bi.destroy()
}

func (bi *BlockIndex) isReady() bool {
	return bi != nil && bi.ptr != nil
}

func (bi *BlockIndex) uninitializedError() error {
	return ErrBlockIndexUninitialized
}
