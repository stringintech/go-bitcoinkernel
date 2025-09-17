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

type ChainstateManagerOptions struct {
	*uniqueHandle
}

func newChainstateManagerOptions(ptr *C.btck_ChainstateManagerOptions) *ChainstateManagerOptions {
	h := newUniqueHandle(unsafe.Pointer(ptr), chainstateManagerOptionsCFuncs{})
	return &ChainstateManagerOptions{uniqueHandle: h}
}

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

// SetWorkerThreads sets the number of worker threads for validation
func (opts *ChainstateManagerOptions) SetWorkerThreads(threads int) {
	C.btck_chainstate_manager_options_set_worker_threads_num((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(threads))
}

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

func (opts *ChainstateManagerOptions) SetBlockTreeDBInMemory(inMemory bool) {
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_set_block_tree_db_in_memory((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(inMemoryInt))
}

func (opts *ChainstateManagerOptions) SetChainstateDBInMemory(inMemory bool) {
	inMemoryInt := 0
	if inMemory {
		inMemoryInt = 1
	}
	C.btck_chainstate_manager_options_set_chainstate_db_in_memory((*C.btck_ChainstateManagerOptions)(opts.ptr), C.int(inMemoryInt))
}
