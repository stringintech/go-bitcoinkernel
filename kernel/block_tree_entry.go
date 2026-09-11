package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import "runtime"

// BlockTreeEntry represents a pointer to an element in the block index currently
// in memory of the chainstate manager.
//
// It is valid for the lifetime of the chainstate manager it was retrieved from.
// The entry is part of a tree-like structure that is maintained internally. Every
// entry, besides the genesis, points to a single parent. Multiple entries may share
// a parent, thus forming a tree. Each entry corresponds to a single block and may
// be used to retrieve its data and validation status.
type BlockTreeEntry struct {
	ptr   *C.btck_BlockTreeEntry
	owner cResourceHandle
}

// Height returns the height of this block in the block tree.
func (bi *BlockTreeEntry) Height() int32 {
	height := int32(C.btck_block_tree_entry_get_height(bi.ptr))
	runtime.KeepAlive(bi)
	return height
}

// Hash returns the block hash associated with this block tree entry.
func (bi *BlockTreeEntry) Hash() *BlockHashView {
	ptr := C.btck_block_tree_entry_get_block_hash(bi.ptr)
	runtime.KeepAlive(bi)
	return newBlockHashView(check(ptr), bi.owner)
}

// Previous returns the previous block tree entry in the chain.
//
// Returns nil if this is the genesis block. The returned entry is a non-owned
// pointer valid for the lifetime of the chainstate manager.
func (bi *BlockTreeEntry) Previous() *BlockTreeEntry {
	ptr := C.btck_block_tree_entry_get_previous(bi.ptr)
	runtime.KeepAlive(bi)
	if ptr == nil {
		return nil
	}
	prevIndex := &BlockTreeEntry{ptr: ptr, owner: bi.owner}
	return prevIndex
}

// Equals compares two block tree entries for equality.
// Returns true if both entries point to the same block in the tree.
func (bi *BlockTreeEntry) Equals(other *BlockTreeEntry) bool {
	if other == nil {
		return false
	}
	equal := C.btck_block_tree_entry_equals(bi.ptr, other.ptr) != 0
	runtime.KeepAlive(bi)
	runtime.KeepAlive(other)
	return equal
}

// Ancestor returns the ancestor block tree entry at the given height.
//
// Returns nil if height is negative or greater than this entry's height.
func (bi *BlockTreeEntry) Ancestor(height int32) *BlockTreeEntry {
	if height < 0 || height > bi.Height() {
		return nil
	}
	ptr := C.btck_block_tree_entry_get_ancestor(bi.ptr, C.int32_t(height))
	runtime.KeepAlive(bi)
	return &BlockTreeEntry{ptr: check(ptr), owner: bi.owner}
}

// GetHeader returns the block header associated with this block tree entry.
func (bi *BlockTreeEntry) GetHeader() *BlockHeader {
	header := newBlockHeader(C.btck_block_tree_entry_get_block_header(bi.ptr), true)
	runtime.KeepAlive(bi)
	return header
}
