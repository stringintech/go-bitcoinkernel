package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &BlockTreeEntry{}

// BlockTreeEntry wraps the C btck_BlockTreeEntry
type BlockTreeEntry struct {
	ptr *C.btck_BlockTreeEntry
}

func (bi *BlockTreeEntry) Height() int32 {
	checkReady(bi)
	return int32(C.btck_block_tree_entry_get_height(bi.ptr))
}

func (bi *BlockTreeEntry) Hash() (*BlockHash, error) {
	checkReady(bi)

	ptr := C.btck_block_tree_entry_get_block_hash(bi.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockGetHash
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

func (bi *BlockTreeEntry) Previous() *BlockTreeEntry {
	checkReady(bi)

	ptr := C.btck_block_tree_entry_get_previous(bi.ptr)
	if ptr == nil {
		return nil
	}

	prevIndex := &BlockTreeEntry{ptr: ptr}
	runtime.SetFinalizer(prevIndex, (*BlockTreeEntry).destroy)
	return prevIndex
}

func (bi *BlockTreeEntry) destroy() {
	if bi.ptr != nil {
		C.btck_block_tree_entry_destroy(bi.ptr)
		bi.ptr = nil
	}
}

func (bi *BlockTreeEntry) Destroy() {
	runtime.SetFinalizer(bi, nil)
	bi.destroy()
}

func (bi *BlockTreeEntry) isReady() bool {
	return bi != nil && bi.ptr != nil
}

func (bi *BlockTreeEntry) uninitializedError() error {
	return ErrBlockTreeEntryUninitialized
}
