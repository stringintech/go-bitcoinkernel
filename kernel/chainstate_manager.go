package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ChainstateManager wraps the C kernel_ChainstateManager
type ChainstateManager struct {
	ptr     *C.kernel_ChainstateManager
	context *Context
}

// NewChainstateManager creates a new chainstate manager.
// Kernel copies all necessary data from the options during construction,
// so the caller can safely free the options object after this call returns successfully.
// However, the context must remain valid for the entire lifetime of the returned ChainstateManager.
func NewChainstateManager(context *Context, options *ChainstateManagerOptions) (*ChainstateManager, error) {
	if context == nil || context.ptr == nil {
		return nil, ErrContextCreation
	}
	if options == nil || options.ptr == nil {
		return nil, ErrChainstateManagerOptionsCreation
	}

	ptr := C.kernel_chainstate_manager_create(context.ptr, options.ptr)
	if ptr == nil {
		return nil, ErrChainstateManagerCreation
	}

	manager := &ChainstateManager{
		ptr:     ptr,
		context: context,
	}
	runtime.SetFinalizer(manager, (*ChainstateManager).destroy)
	return manager, nil
}

// ReadBlockFromDisk reads a block from disk using the provided block index
func (cm *ChainstateManager) ReadBlockFromDisk(blockIndex *BlockIndex) (*Block, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}
	if blockIndex == nil || blockIndex.ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	ptr := C.kernel_read_block_from_disk(cm.context.ptr, cm.ptr, blockIndex.ptr)
	if ptr == nil {
		return nil, ErrBlockRead
	}

	block := &Block{ptr: ptr}
	runtime.SetFinalizer(block, (*Block).destroy)
	return block, nil
}

// ProcessBlock processes and validates a block
func (cm *ChainstateManager) ProcessBlock(block *Block) (bool, bool, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return false, false, ErrChainstateManagerCreation
	}
	if block == nil || block.ptr == nil {
		return false, false, ErrInvalidBlock
	}

	var newBlock C.bool
	success := C.kernel_chainstate_manager_process_block(
		cm.context.ptr,
		cm.ptr,
		block.ptr,
		&newBlock,
	)

	return bool(success), bool(newBlock), nil
}

// GetBlockIndexFromTip returns the block index of the current chain tip
func (cm *ChainstateManager) GetBlockIndexFromTip() (*BlockIndex, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}

	ptr := C.kernel_get_block_index_from_tip(cm.context.ptr, cm.ptr)
	if ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	blockIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(blockIndex, (*BlockIndex).destroy)
	return blockIndex, nil
}

// GetBlockIndexFromGenesis returns the block index of the genesis block
func (cm *ChainstateManager) GetBlockIndexFromGenesis() (*BlockIndex, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}

	ptr := C.kernel_get_block_index_from_genesis(cm.context.ptr, cm.ptr)
	if ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	blockIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(blockIndex, (*BlockIndex).destroy)
	return blockIndex, nil
}

// GetBlockIndexFromHash returns the block index for a given block hash
func (cm *ChainstateManager) GetBlockIndexFromHash(blockHash *BlockHash) (*BlockIndex, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}
	if blockHash == nil || blockHash.ptr == nil {
		return nil, ErrHashCalculation
	}

	ptr := C.kernel_get_block_index_from_hash(cm.context.ptr, cm.ptr, blockHash.ptr)
	if ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	blockIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(blockIndex, (*BlockIndex).destroy)
	return blockIndex, nil
}

// GetBlockIndexFromHeight returns the block index for a given height in the currently active chain
func (cm *ChainstateManager) GetBlockIndexFromHeight(height int) (*BlockIndex, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}

	ptr := C.kernel_get_block_index_from_height(cm.context.ptr, cm.ptr, C.int(height))
	if ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	blockIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(blockIndex, (*BlockIndex).destroy)
	return blockIndex, nil
}

// GetNextBlockIndex returns the next block index in the active chain
func (cm *ChainstateManager) GetNextBlockIndex(blockIndex *BlockIndex) (*BlockIndex, error) {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return nil, ErrChainstateManagerCreation
	}
	if blockIndex == nil || blockIndex.ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	ptr := C.kernel_get_next_block_index(cm.context.ptr, cm.ptr, blockIndex.ptr)
	if ptr == nil {
		return nil, nil // No next block (tip or invalid)
	}

	nextIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(nextIndex, (*BlockIndex).destroy)
	return nextIndex, nil
}

// ImportBlocks imports blocks from the specified file paths
func (cm *ChainstateManager) ImportBlocks(blockFilePaths []string) error {
	if cm.ptr == nil || cm.context == nil || cm.context.ptr == nil {
		return ErrChainstateManagerCreation
	}

	if len(blockFilePaths) == 0 {
		// Import with no files triggers reindex if wipe options were set
		success := C.kernel_import_blocks(cm.context.ptr, cm.ptr, nil, nil, 0)
		if !success {
			return ErrBlockProcessing
		}
		return nil
	}

	// Convert Go strings to C strings
	cPaths := make([]*C.char, len(blockFilePaths))
	cLens := make([]C.size_t, len(blockFilePaths))

	for i, path := range blockFilePaths {
		cPaths[i] = C.CString(path)
		cLens[i] = C.size_t(len(path))
	}

	// Clean up C strings
	defer func() {
		for i := range cPaths {
			C.free(unsafe.Pointer(cPaths[i]))
		}
	}()

	success := C.kernel_import_blocks(
		cm.context.ptr,
		cm.ptr,
		(**C.char)(unsafe.Pointer(&cPaths[0])),
		(*C.size_t)(unsafe.Pointer(&cLens[0])),
		C.size_t(len(blockFilePaths)),
	)

	if !success {
		return ErrBlockProcessing
	}
	return nil
}

func (cm *ChainstateManager) destroy() {
	if cm.ptr != nil && cm.context != nil && cm.context.ptr != nil {
		C.kernel_chainstate_manager_destroy(cm.ptr, cm.context.ptr)
		cm.ptr = nil
		cm.context = nil
	}
}

func (cm *ChainstateManager) Destroy() {
	runtime.SetFinalizer(cm, nil)
	cm.destroy()
}
