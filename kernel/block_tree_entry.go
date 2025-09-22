package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

type BlockTreeEntry struct {
	ptr *C.btck_BlockTreeEntry
}

func (bi *BlockTreeEntry) Height() int32 {
	return int32(C.btck_block_tree_entry_get_height(bi.ptr))
}

func (bi *BlockTreeEntry) Hash() *BlockHash {
	ptr := C.btck_block_tree_entry_get_block_hash(bi.ptr)
	return newBlockHash(check(ptr), true)
}

func (bi *BlockTreeEntry) Previous() *BlockTreeEntry {
	ptr := C.btck_block_tree_entry_get_previous(bi.ptr)
	if ptr == nil {
		return nil
	}
	prevIndex := &BlockTreeEntry{ptr: ptr}
	return prevIndex
}
