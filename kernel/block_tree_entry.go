package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"

// BlockTreeEntry represents a pointer to an element in the block index currently
// in memory of the chainstate manager.
//
// It is valid for the lifetime of the chainstate manager it was retrieved from.
// The entry is part of a tree-like structure that is maintained internally. Every
// entry, besides the genesis, points to a single parent. Multiple entries may share
// a parent, thus forming a tree. Each entry corresponds to a single block and may
// be used to retrieve its data and validation status.
type BlockTreeEntry struct {
	ptr *C.btck_BlockTreeEntry
}

// Height returns the height of this block in the block tree.
func (bi *BlockTreeEntry) Height() int32 {
	return int32(C.btck_block_tree_entry_get_height(bi.ptr))
}

// Hash returns the block hash associated with this block tree entry.
func (bi *BlockTreeEntry) Hash() *BlockHashView {
	ptr := C.btck_block_tree_entry_get_block_hash(bi.ptr)
	return newBlockHashView(check(ptr))
}

// Previous returns the previous block tree entry in the chain.
//
// Returns nil if this is the genesis block. The returned entry is a non-owned
// pointer valid for the lifetime of the chainstate manager.
func (bi *BlockTreeEntry) Previous() *BlockTreeEntry {
	ptr := C.btck_block_tree_entry_get_previous(bi.ptr)
	if ptr == nil {
		return nil
	}
	prevIndex := &BlockTreeEntry{ptr: ptr}
	return prevIndex
}

// Equals compares two block tree entries for equality.
// Returns true if both entries point to the same block in the tree.
func (bi *BlockTreeEntry) Equals(other *BlockTreeEntry) bool {
	if other == nil {
		return false
	}
	return C.btck_block_tree_entry_equals(bi.ptr, other.ptr) != 0
}
