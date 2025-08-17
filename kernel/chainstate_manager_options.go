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

var _ cManagedResource = &ChainstateManagerOptions{}

// ChainstateManagerOptions wraps the C btck_ChainstateManagerOptions
type ChainstateManagerOptions struct {
	ptr     *C.btck_ChainstateManagerOptions
	context *Context
}

// NewChainstateManagerOptions creates new chainstate manager options.
// The context must remain valid for the entire lifetime of the returned options.
func NewChainstateManagerOptions(context *Context, dataDir, blocksDir string) (*ChainstateManagerOptions, error) {
	if err := validateReady(context); err != nil {
		return nil, err
	}

	cDataDir := C.CString(dataDir)
	defer C.free(unsafe.Pointer(cDataDir))

	cBlocksDir := C.CString(blocksDir)
	defer C.free(unsafe.Pointer(cBlocksDir))

	ptr := C.btck_chainstate_manager_options_create(
		context.ptr,
		cDataDir,
		C.size_t(len(dataDir)),
		cBlocksDir,
		C.size_t(len(blocksDir)),
	)
	if ptr == nil {
		return nil, ErrKernelChainstateManagerOptionsCreate
	}

	opts := &ChainstateManagerOptions{
		ptr:     ptr,
		context: context,
	}
	runtime.SetFinalizer(opts, (*ChainstateManagerOptions).destroy)
	return opts, nil
}

// SetWorkerThreads sets the number of worker threads for validation
func (opts *ChainstateManagerOptions) SetWorkerThreads(threads int) {
	checkReady(opts)
	C.btck_chainstate_manager_options_set_worker_threads_num(opts.ptr, C.int(threads))
}

func (opts *ChainstateManagerOptions) SetWipeDBs(wipeBlockTree, wipeChainstate bool) bool {
	checkReady(opts)
	wipeBlockTreeInt := 0
	if wipeBlockTree {
		wipeBlockTreeInt = 1
	}
	wipeChainstateInt := 0
	if wipeChainstate {
		wipeChainstateInt = 1
	}
	return C.btck_chainstate_manager_options_set_wipe_dbs(
		opts.ptr,
		C.int(wipeBlockTreeInt),
		C.int(wipeChainstateInt),
	) != 0
}

func (opts *ChainstateManagerOptions) SetBlockTreeDBInMemory(inMemory bool) {
	checkReady(opts)
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_set_block_tree_db_in_memory(opts.ptr, C.int(inMemoryInt))
}

func (opts *ChainstateManagerOptions) SetChainstateDBInMemory(inMemory bool) {
	checkReady(opts)
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_set_chainstate_db_in_memory(opts.ptr, C.int(inMemoryInt))
}

func (opts *ChainstateManagerOptions) destroy() {
	if opts.ptr != nil {
		C.btck_chainstate_manager_options_destroy(opts.ptr)
		opts.ptr = nil
		opts.context = nil
	}
}

func (opts *ChainstateManagerOptions) Destroy() {
	runtime.SetFinalizer(opts, nil)
	opts.destroy()
}

func (opts *ChainstateManagerOptions) isReady() bool {
	return opts != nil && opts.ptr != nil
}

func (opts *ChainstateManagerOptions) uninitializedError() error {
	return ErrChainstateManagerOptionsUninitialized
}
