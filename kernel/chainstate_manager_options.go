package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type chainstateManagerOptionsCFuncs struct{}

func (chainstateManagerOptionsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_chainstate_manager_options_destroy((*C.btck_ChainstateManagerOptions)(ptr))
}

// ChainstateManagerOptions holds configuration options for creating a new chainstate manager.
//
// Options are initialized with sensible defaults and can be customized using the
// setter methods before creating the chainstate manager.
type ChainstateManagerOptions struct {
	*uniqueHandle
}

func newChainstateManagerOptions(ptr *C.btck_ChainstateManagerOptions) *ChainstateManagerOptions {
	h := newUniqueHandle(unsafe.Pointer(ptr), chainstateManagerOptionsCFuncs{})
	return &ChainstateManagerOptions{uniqueHandle: h}
}

// NewChainstateManagerOptions creates options for configuring a chainstate manager.
//
// The options associate with the provided kernel context and specify the data and block
// directories. If the directories do not exist, they will be created.
//
// Parameters:
//   - context: Kernel context that the chainstate manager will associate with
//   - dataDir: Path to the directory containing chainstate data
//   - blocksDir: Path to the directory containing block data
//
// Returns an error if the options cannot be created.
func NewChainstateManagerOptions(context *Context, dataDir, blocksDir string) (*ChainstateManagerOptions, error) {
	cDataDir := C.CString(dataDir)
	defer C.free(unsafe.Pointer(cDataDir))

	cBlocksDir := C.CString(blocksDir)
	defer C.free(unsafe.Pointer(cBlocksDir))

	ptr := C.btck_chainstate_manager_options_create((*C.btck_Context)(context.ptr), cDataDir, C.size_t(len(dataDir)),
		cBlocksDir, C.size_t(len(blocksDir)))
	if ptr == nil {
		return nil, &InternalError{"Failed to create chainstate manager options"}
	}
	return newChainstateManagerOptions(ptr), nil
}

// SetWorkerThreads configures the number of worker threads for parallel validation.
//
// Parameters:
//   - threads: Number of worker threads (0 disables parallel verification, max is clamped to 15)
func (opts *ChainstateManagerOptions) SetWorkerThreads(threads int) {
	C.btck_chainstate_manager_options_set_worker_threads_num((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(threads))
}

// SetWipeDBs configures which databases to wipe on startup.
//
// When combined with ImportBlocks, this triggers a full reindex (if wipeBlockTree is true)
// or chainstate-only reindex (if only wipeChainstate is true).
//
// Parameters:
//   - wipeBlockTree: Whether to wipe the block tree database (requires wipeChainstate to also be true)
//   - wipeChainstate: Whether to wipe the chainstate database
//
// Returns an error if wipeBlockTree is true but wipeChainstate is false.
func (opts *ChainstateManagerOptions) SetWipeDBs(wipeBlockTree, wipeChainstate bool) error {
	wipeBlockTreeInt := 0
	if wipeBlockTree {
		wipeBlockTreeInt = 1
	}
	wipeChainstateInt := 0
	if wipeChainstate {
		wipeChainstateInt = 1
	}
	result := C.btck_chainstate_manager_options_set_wipe_dbs((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(wipeBlockTreeInt), C.int(wipeChainstateInt))
	if result != 0 {
		return &InternalError{"Failed to set wipe db"}
	}
	return nil
}

// UpdateBlockTreeDBInMemory configures whether the block tree database is stored in memory.
//
// Parameters:
//   - inMemory: If true, the block tree database will be kept entirely in memory
func (opts *ChainstateManagerOptions) UpdateBlockTreeDBInMemory(inMemory bool) {
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_update_block_tree_db_in_memory((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(inMemoryInt))
}

// UpdateChainstateDBInMemory configures whether the chainstate database is stored in memory.
//
// Parameters:
//   - inMemory: If true, the chainstate database will be kept entirely in memory
func (opts *ChainstateManagerOptions) UpdateChainstateDBInMemory(inMemory bool) {
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_update_chainstate_db_in_memory((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(inMemoryInt))
}
