package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

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
