package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type chainstateManagerCFuncs struct{}

func (chainstateManagerCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_chainstate_manager_destroy((*C.btck_ChainstateManager)(ptr))
}

type ChainstateManager struct {
	*uniqueHandle
}

func newChainstateManager(ptr *C.btck_ChainstateManager) *ChainstateManager {
	h := newUniqueHandle(unsafe.Pointer(ptr), chainstateManagerCFuncs{})
	return &ChainstateManager{uniqueHandle: h}
}

// NewChainstateManager creates a new chainstate manager.
// Kernel copies all necessary data from the options during construction,
// so the caller can safely free the options object after this call returns successfully.
func NewChainstateManager(options *ChainstateManagerOptions) (*ChainstateManager, error) {
	ptr := C.btck_chainstate_manager_create((*C.btck_ChainstateManagerOptions)(options.ptr))
	if ptr == nil {
		return nil, &InternalError{"Failed to create chainstate manager"}
	}
	return newChainstateManager(ptr), nil
}

func (cm *ChainstateManager) ReadBlock(blockTreeEntry *BlockTreeEntry) (*Block, error) {
	ptr := C.btck_block_read((*C.btck_ChainstateManager)(cm.ptr), blockTreeEntry.ptr)
	if ptr == nil {
		return nil, &InternalError{"Failed to read block"}
	}
	return newBlock(ptr, true), nil
}

func (cm *ChainstateManager) ReadBlockSpentOutputs(blockTreeEntry *BlockTreeEntry) (*BlockSpentOutputs, error) {
	ptr := C.btck_block_spent_outputs_read((*C.btck_ChainstateManager)(cm.ptr), blockTreeEntry.ptr)
	if ptr == nil {
		return nil, &InternalError{"Failed to read block spent outputs"}
	}
	return newBlockSpentOutputs(ptr, true), nil
}

// ProcessBlock processes and validates the given block with the chainstate manager.
// It returns ok=true if processing was successful (including for duplicate blocks)
// and duplicate=true if this block was processed before.
func (cm *ChainstateManager) ProcessBlock(block *Block) (ok bool, duplicate bool) {
	var newBlock C.int
	result := C.btck_chainstate_manager_process_block((*C.btck_ChainstateManager)(cm.ptr), (*C.btck_Block)(block.ptr), &newBlock)
	ok = result == 0
	duplicate = newBlock == 0
	return
}

func (cm *ChainstateManager) GetActiveChain() *Chain {
	return &Chain{C.btck_chainstate_manager_get_active_chain((*C.btck_ChainstateManager)(cm.ptr))}
}

// GetBlockTreeEntryByHash returns the block tree entry for a given block hash, or null if the hash is not found
func (cm *ChainstateManager) GetBlockTreeEntryByHash(blockHash *BlockHash) *BlockTreeEntry {
	ptr := C.btck_chainstate_manager_get_block_tree_entry_by_hash((*C.btck_ChainstateManager)(cm.ptr), (*C.btck_BlockHash)(blockHash.ptr))
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// ImportBlocks triggers a reindex if the option was previously set and can also import
// existing block files from the specified filesystem paths.
func (cm *ChainstateManager) ImportBlocks(blockFilePaths []string) error {
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
		(*C.btck_ChainstateManager)(cm.ptr),
		cPathsPtr,
		cLensPtr,
		C.size_t(len(blockFilePaths)),
	)
	if success != 0 {
		return &InternalError{"Failed to import blocks"}
	}
	return nil
}
