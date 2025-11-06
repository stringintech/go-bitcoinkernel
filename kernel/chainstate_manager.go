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

// ChainstateManager is the central object for doing validation tasks as well as
// retrieving data from the chain.
type ChainstateManager struct {
	*uniqueHandle
}

func newChainstateManager(ptr *C.btck_ChainstateManager) *ChainstateManager {
	h := newUniqueHandle(unsafe.Pointer(ptr), chainstateManagerCFuncs{})
	return &ChainstateManager{uniqueHandle: h}
}

// NewChainstateManager creates a new chainstate manager for validation and chain queries.
//
// This is the main object for validation tasks, retrieving data from the chain, and
// interacting with chainstate and indexes.
//
// Parameters:
//   - options: Configuration options created by NewChainstateManagerOptions
//
// Returns an error if the chainstate manager cannot be created.
func NewChainstateManager(options *ChainstateManagerOptions) (*ChainstateManager, error) {
	ptr := C.btck_chainstate_manager_create((*C.btck_ChainstateManagerOptions)(options.ptr))
	if ptr == nil {
		return nil, &InternalError{"Failed to create chainstate manager"}
	}
	return newChainstateManager(ptr), nil
}

// ReadBlock reads the block from disk that the block tree entry points to.
//
// Parameters:
//   - blockTreeEntry: Block index entry obtained from GetBlockTreeEntryByHash or chain queries
//
// Returns an error if the block cannot be read from disk.
func (cm *ChainstateManager) ReadBlock(blockTreeEntry *BlockTreeEntry) (*Block, error) {
	ptr := C.btck_block_read((*C.btck_ChainstateManager)(cm.ptr), blockTreeEntry.ptr)
	if ptr == nil {
		return nil, &InternalError{"Failed to read block"}
	}
	return newBlock(ptr, true), nil
}

// ReadBlockSpentOutputs reads the spent outputs for the block that the block tree entry points to from disk.
//
// Parameters:
//   - blockTreeEntry: Block index entry for the block whose spent outputs to read
//
// Returns an error if the undo data cannot be read from disk.
func (cm *ChainstateManager) ReadBlockSpentOutputs(blockTreeEntry *BlockTreeEntry) (*BlockSpentOutputs, error) {
	ptr := C.btck_block_spent_outputs_read((*C.btck_ChainstateManager)(cm.ptr), blockTreeEntry.ptr)
	if ptr == nil {
		return nil, &InternalError{"Failed to read block spent outputs"}
	}
	return newBlockSpentOutputs(ptr, true), nil
}

// ProcessBlock processes and validates the passed in block with the chainstate
// manager. Processing first does checks on the block, and if these passed,
// saves it to disk. It then validates the block against the utxo set. If it is
// valid, the chain is extended with it. The ok return value is not indicative of
// the block's validity. Detailed information on the validity of the block can
// be retrieved by registering the block_checked callback in the validation
// interface.
//
// Parameters:
//   - block: Block to validate and potentially add to the chain
//
// Returns ok=true if processing the block was successful (will also return true for valid,
// but duplicate blocks) and duplicate=false if this block was not processed before. Note that
// duplicate might also be false if processing was attempted before, but the block was found
// invalid before its data was persisted.
func (cm *ChainstateManager) ProcessBlock(block *Block) (ok bool, duplicate bool) {
	var newBlock C.int
	result := C.btck_chainstate_manager_process_block((*C.btck_ChainstateManager)(cm.ptr), (*C.btck_Block)(block.ptr), &newBlock)
	ok = result == 0
	duplicate = newBlock == 0
	return
}

// GetActiveChain returns the currently active best-known chain.
//
// The returned Chain can be thought of as a view on a vector of block tree entries
// that form the best chain. The chain's lifetime depends on this chainstate manager.
// State transitions (e.g., processing blocks) will change the chain, so data retrieved
// from it is only consistent until new data is processed. It is the caller's responsibility
// to guard against these inconsistencies.
func (cm *ChainstateManager) GetActiveChain() *Chain {
	return &Chain{C.btck_chainstate_manager_get_active_chain((*C.btck_ChainstateManager)(cm.ptr))}
}

// GetBlockTreeEntryByHash retrieves a block tree entry by its block hash.
//
// Parameters:
//   - blockHash: Hash of the block to look up
//
// Returns nil if no block with this hash is found in the block index. The returned
// BlockTreeEntry is a non-owned pointer valid for the lifetime of this chainstate
// manager.
func (cm *ChainstateManager) GetBlockTreeEntryByHash(blockHash BlockHashLike) *BlockTreeEntry {
	ptr := C.btck_chainstate_manager_get_block_tree_entry_by_hash((*C.btck_ChainstateManager)(cm.ptr), blockHash.blockHashPtr())
	if ptr == nil {
		return nil
	}
	return &BlockTreeEntry{ptr: ptr}
}

// ImportBlocks triggers a reindex and/or imports block files from the filesystem.
//
// This starts a reindex if the wipe options were previously set via ChainstateManagerOptions.
// It can also import existing block files from the specified filesystem paths.
//
// Parameters:
//   - blockFilePaths: Array of full filesystem paths to block files to import (can be empty)
//
// Returns an error if the import fails. This is a long-running operation that can
// be interrupted via Context.Interrupt().
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
