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

// ChainstateManagerOptions wraps the C kernel_ChainstateManagerOptions
type ChainstateManagerOptions struct {
	ptr     *C.kernel_ChainstateManagerOptions
	context *Context
}

// NewChainstateManagerOptions creates new chainstate manager options.
// The context must remain valid for the entire lifetime of the returned options.
func NewChainstateManagerOptions(context *Context, dataDir, blocksDir string) (*ChainstateManagerOptions, error) {
	if context == nil || context.ptr == nil {
		return nil, ErrContextCreation
	}

	cDataDir := C.CString(dataDir)
	defer C.free(unsafe.Pointer(cDataDir))

	cBlocksDir := C.CString(blocksDir)
	defer C.free(unsafe.Pointer(cBlocksDir))

	ptr := C.kernel_chainstate_manager_options_create(
		context.ptr,
		cDataDir,
		C.size_t(len(dataDir)),
		cBlocksDir,
		C.size_t(len(blocksDir)),
	)
	if ptr == nil {
		return nil, ErrChainstateManagerOptionsCreation
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
	if opts.ptr != nil {
		C.kernel_chainstate_manager_options_set_worker_threads_num(opts.ptr, C.int(threads))
	}
}

func (opts *ChainstateManagerOptions) SetWipeDBs(wipeBlockTree, wipeChainstate bool) bool {
	if opts.ptr == nil {
		return false
	}
	return bool(C.kernel_chainstate_manager_options_set_wipe_dbs(
		opts.ptr,
		C.bool(wipeBlockTree),
		C.bool(wipeChainstate),
	))
}

func (opts *ChainstateManagerOptions) SetBlockTreeDBInMemory(inMemory bool) {
	if opts.ptr != nil {
		C.kernel_chainstate_manager_options_set_block_tree_db_in_memory(opts.ptr, C.bool(inMemory))
	}
}

func (opts *ChainstateManagerOptions) SetChainstateDBInMemory(inMemory bool) {
	if opts.ptr != nil {
		C.kernel_chainstate_manager_options_set_chainstate_db_in_memory(opts.ptr, C.bool(inMemory))
	}
}

func (opts *ChainstateManagerOptions) destroy() {
	if opts.ptr != nil {
		C.kernel_chainstate_manager_options_destroy(opts.ptr)
		opts.ptr = nil
		opts.context = nil
	}
}

func (opts *ChainstateManagerOptions) Destroy() {
	runtime.SetFinalizer(opts, nil)
	opts.destroy()
}
