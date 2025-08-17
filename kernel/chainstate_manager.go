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

var _ cManagedResource = &ChainstateManager{}

// ChainstateManager wraps the C btck_ChainstateManager
type ChainstateManager struct {
	ptr     *C.btck_ChainstateManager
	context *Context
}

// NewChainstateManager creates a new chainstate manager.
// Kernel copies all necessary data from the options during construction,
// so the caller can safely free the options object after this call returns successfully.
// However, the context must remain valid for the entire lifetime of the returned ChainstateManager.
func NewChainstateManager(context *Context, options *ChainstateManagerOptions) (*ChainstateManager, error) {
	if err := validateReady(context); err != nil {
		return nil, err
	}
	if err := validateReady(options); err != nil {
		return nil, err
	}

	ptr := C.btck_chainstate_manager_create(options.ptr)
	if ptr == nil {
		return nil, ErrKernelChainstateManagerCreate
	}

	manager := &ChainstateManager{
		ptr:     ptr,
		context: context,
	}
	runtime.SetFinalizer(manager, (*ChainstateManager).destroy)
	return manager, nil
}

// ReadBlock reads a block using the provided block tree entry
func (cm *ChainstateManager) ReadBlock(blockTreeEntry *BlockTreeEntry) (*Block, error) {
	checkReady(cm)
	if err := validateReady(blockTreeEntry); err != nil {
		return nil, err
	}

	ptr := C.btck_block_read(cm.ptr, blockTreeEntry.ptr)
	if ptr == nil {
		return nil, ErrKernelChainstateManagerReadBlock
	}
	return newBlockFromPtr(ptr), nil
}

// ReadBlockSpentOutputs reads block spent outputs data for a given block tree entry
func (cm *ChainstateManager) ReadBlockSpentOutputs(blockTreeEntry *BlockTreeEntry) (*BlockSpentOutputs, error) {
	checkReady(cm)
	if err := validateReady(blockTreeEntry); err != nil {
		return nil, err
	}

	ptr := C.btck_block_spent_outputs_read(cm.ptr, blockTreeEntry.ptr)
	if ptr == nil {
		return nil, ErrKernelChainstateManagerReadBlockUndo
	}

	blockSpentOutputs := &BlockSpentOutputs{ptr: ptr}
	runtime.SetFinalizer(blockSpentOutputs, (*BlockSpentOutputs).destroy)
	return blockSpentOutputs, nil
}

// ProcessBlock processes and validates a block
func (cm *ChainstateManager) ProcessBlock(block *Block) (bool, bool, error) {
	checkReady(cm)
	if err := validateReady(block); err != nil {
		return false, false, err
	}

	var newBlock C.int
	result := C.btck_chainstate_manager_process_block(
		cm.ptr,
		block.ptr,
		&newBlock,
	)
	if result != 0 {
		return false, false, ErrKernelChainstateManagerProcessBlock
	}
	return true, newBlock != 0, nil
}

// GetActiveChain returns the currently active chain
func (cm *ChainstateManager) GetActiveChain() (*Chain, error) {
	checkReady(cm)

	ptr := C.btck_chainstate_manager_get_active_chain(cm.ptr)
	if ptr == nil {
		return nil, ErrChainUninitialized
	}

	chain := &Chain{ptr: ptr}
	runtime.SetFinalizer(chain, (*Chain).destroy)
	return chain, nil
}

// GetBlockTreeEntryByHash returns the block tree entry for a given block hash, or null if the hash is not found
func (cm *ChainstateManager) GetBlockTreeEntryByHash(blockHash *BlockHash) (*BlockTreeEntry, error) {
	checkReady(cm)
	if err := validateReady(blockHash); err != nil {
		return nil, err
	}

	ptr := C.btck_chainstate_manager_get_block_tree_entry_by_hash(cm.ptr, blockHash.ptr)
	if ptr == nil {
		return nil, nil
	}

	blockTreeEntry := &BlockTreeEntry{ptr: ptr}
	runtime.SetFinalizer(blockTreeEntry, (*BlockTreeEntry).destroy)
	return blockTreeEntry, nil
}

// ImportBlocks imports blocks from the specified file paths
func (cm *ChainstateManager) ImportBlocks(blockFilePaths []string) error {
	checkReady(cm)

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
			if cPaths[i] != nil {
				C.free(unsafe.Pointer(cPaths[i]))
			}
		}
	}()

	var cPathsPtr **C.char
	var cLensPtr *C.size_t
	if len(cPaths) > 0 {
		cPathsPtr = &cPaths[0]
		cLensPtr = &cLens[0]
	}

	success := C.btck_chainstate_manager_import_blocks(
		cm.ptr,
		cPathsPtr,
		cLensPtr,
		C.size_t(len(blockFilePaths)),
	)

	if success != 0 {
		return ErrKernelImportBlocks
	}
	return nil
}

func (cm *ChainstateManager) destroy() {
	if cm.isReady() {
		C.btck_chainstate_manager_destroy(cm.ptr)
		cm.ptr = nil
		cm.context = nil
	}
}

func (cm *ChainstateManager) Destroy() {
	runtime.SetFinalizer(cm, nil)
	cm.destroy()
}

func (cm *ChainstateManager) isReady() bool {
	return cm != nil && cm.ptr != nil && cm.context.isReady()
}

func (cm *ChainstateManager) uninitializedError() error {
	return ErrChainstateManagerUninitialized
}
