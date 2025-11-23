package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import "iter"

// Chain represents the currently known best-chain associated with a chainstate.
//
// Its lifetime depends on the chainstate manager, and state transitions within
// the manager (e.g., when processing blocks) will also change the chain. Data
// retrieved from this chain is only consistent up to the point when new data
// is processed in the chainstate manager.
type Chain struct {
	ptr *C.btck_Chain
}

// GetTip returns the block tree entry of the current chain tip.
//
// Returns nil if the chain is empty. Once returned, there is no guarantee that it
// remains in the active chain if new blocks are processed.
func (c *Chain) GetTip() *BlockTreeEntry {
	ptr := C.btck_chain_get_tip(c.ptr)
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// GetGenesis returns the block tree entry of the genesis block.
//
// Returns nil if the chain is empty.
func (c *Chain) GetGenesis() *BlockTreeEntry {
	ptr := C.btck_chain_get_genesis(c.ptr)
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// GetByHeight retrieves a block tree entry by its height in the currently active chain.
//
// Returns nil if the height is out of bounds. Once retrieved, there is no guarantee
// that it remains in the active chain if new blocks are processed.
//
// Parameters:
//   - height: Block height to retrieve
func (c *Chain) GetByHeight(height int32) *BlockTreeEntry {
	ptr := C.btck_chain_get_by_height(c.ptr, C.int(height))
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr}
}

// Contains checks whether the given block tree entry is part of this chain.
//
// Returns true if the block tree entry is in the currently active chain, false otherwise.
func (c *Chain) Contains(blockTreeEntry *BlockTreeEntry) bool {
	return C.btck_chain_contains(c.ptr, blockTreeEntry.ptr) != 0
}

// GetHeight returns the height of the chain's tip.
//
// This is the height of the most recent block in the chain.
func (c *Chain) GetHeight() int32 {
	return int32(C.btck_chain_get_height(c.ptr))
}

// Entries returns an iterator over all block tree entries in the chain.
//
// The iterator starts from height 0 (genesis) and goes up to and including the chain tip.
//
// Example usage:
//
//	for entry := range chain.Entries() {
//	    // Process block tree entry
//	}
func (c *Chain) Entries() iter.Seq[*BlockTreeEntry] {
	return func(yield func(*BlockTreeEntry) bool) {
		c.iterEntries(0, c.GetHeight()+1, yield)
	}
}

// EntriesRange returns an iterator over a range of block tree entries in the chain.
//
// Parameters:
//   - from: Starting height (inclusive)
//   - to: Ending height (exclusive)
//
// Safe for out-of-bounds arguments: 'to' is clamped to chain tip height + 1,
// and an invalid range (from >= to) yields an empty iterator.
//
// Example usage:
//
//	for entry := range chain.EntriesRange(0, 6) {
//	    // Process block tree entries at heights 0-5
//	}
func (c *Chain) EntriesRange(from, to int32) iter.Seq[*BlockTreeEntry] {
	return func(yield func(*BlockTreeEntry) bool) {
		if tipHeight := c.GetHeight(); to > tipHeight+1 {
			to = tipHeight + 1
		}
		c.iterEntries(from, to, yield)
	}
}

// EntriesFrom returns an iterator over block tree entries starting from the given height.
//
// Parameters:
//   - from: Starting height (inclusive)
//
// If from is beyond the chain height, returns an empty iterator.
//
// Example usage:
//
//	for entry := range chain.EntriesFrom(5) {
//	    // Process block tree entries from height 5 to the tip
//	}
func (c *Chain) EntriesFrom(from int32) iter.Seq[*BlockTreeEntry] {
	return func(yield func(*BlockTreeEntry) bool) {
		c.iterEntries(from, c.GetHeight()+1, yield)
	}
}

// iterEntries is a helper that iterates over block tree entries in [from, to).
func (c *Chain) iterEntries(from, to int32, yield func(*BlockTreeEntry) bool) {
	for h := from; h < to; h++ {
		entry := c.GetByHeight(h)
		if entry == nil { // Height may become out of bounds due to a reorg
			return
		}
		if !yield(entry) {
			return
		}
	}
}
