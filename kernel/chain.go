package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &Chain{}

// Chain wraps the C btck_Chain
type Chain struct {
	ptr *C.btck_Chain
}

// GetTip returns the block tree entry of the current tip, or nil if chain is empty
func (c *Chain) GetTip() (*BlockTreeEntry, error) {
	checkReady(c)
	ptr := C.btck_chain_get_tip(c.ptr)
	if ptr == nil {
		return nil, nil
	}

	entry := &BlockTreeEntry{ptr: ptr}
	runtime.SetFinalizer(entry, (*BlockTreeEntry).destroy)
	return entry, nil
}

// GetGenesis returns the block tree entry of the genesis block, or nil if chain is empty
func (c *Chain) GetGenesis() (*BlockTreeEntry, error) {
	checkReady(c)
	ptr := C.btck_chain_get_genesis(c.ptr)
	if ptr == nil {
		return nil, nil
	}

	entry := &BlockTreeEntry{ptr: ptr}
	runtime.SetFinalizer(entry, (*BlockTreeEntry).destroy)
	return entry, nil
}

// GetByHeight returns the block tree entry at the specified height
func (c *Chain) GetByHeight(height int) *BlockTreeEntry {
	checkReady(c)
	ptr := C.btck_chain_get_by_height(c.ptr, C.int(height))
	if ptr == nil {
		panic(ErrKernelChainGetByHeight)
	}

	entry := &BlockTreeEntry{ptr: ptr}
	runtime.SetFinalizer(entry, (*BlockTreeEntry).destroy)
	return entry
}

// Contains returns true if the chain contains the block tree entry
func (c *Chain) Contains(blockTreeEntry *BlockTreeEntry) bool {
	checkReady(c)
	if err := validateReady(blockTreeEntry); err != nil {
		return false
	}

	return C.btck_chain_contains(c.ptr, blockTreeEntry.ptr) != 0
}

func (c *Chain) destroy() {
	if c.ptr != nil {
		C.btck_chain_destroy(c.ptr)
		c.ptr = nil
	}
}

func (c *Chain) Destroy() {
	runtime.SetFinalizer(c, nil)
	c.destroy()
}

func (c *Chain) isReady() bool {
	return c != nil && c.ptr != nil
}

func (c *Chain) uninitializedError() error {
	return ErrChainUninitialized
}
