package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

type Chain struct {
	ptr *C.btck_Chain
}

// GetTip returns the block tree entry of the current tip, or nil if chain is empty
func (c *Chain) GetTip() *BlockTreeEntry {
	ptr := C.btck_chain_get_tip(c.ptr)
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// GetGenesis returns the block tree entry of the genesis block, or nil if chain is empty
func (c *Chain) GetGenesis() *BlockTreeEntry {
	ptr := C.btck_chain_get_genesis(c.ptr)
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// GetByHeight returns the block tree entry for the given height in the chain, or nil if it does not exist
func (c *Chain) GetByHeight(height int32) *BlockTreeEntry {
	ptr := C.btck_chain_get_by_height(c.ptr, C.int(height))
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr}
}

func (c *Chain) Contains(blockTreeEntry *BlockTreeEntry) bool {
	return C.btck_chain_contains(c.ptr, blockTreeEntry.ptr) != 0
}

// GetHeight returns the height of the tip of the chain
func (c *Chain) GetHeight() int32 {
	return int32(C.btck_chain_get_height(c.ptr))
}
