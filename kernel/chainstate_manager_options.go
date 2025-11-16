package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"

// ChainstateManagerOption is a functional option for configuring chainstate manager.
type ChainstateManagerOption func(*C.btck_ChainstateManagerOptions) error

// WithWorkerThreads returns a ChainstateManagerOption that configures the number of worker threads for parallel validation.
//
// Parameters:
//   - threads: Number of worker threads (0 disables parallel verification, max is clamped to 15)
func WithWorkerThreads(threads int) ChainstateManagerOption {
	return func(opts *C.btck_ChainstateManagerOptions) error {
		C.btck_chainstate_manager_options_set_worker_threads_num(opts, C.int(threads))
		return nil
	}
}

// WithWipeDBs returns a ChainstateManagerOption that configures which databases to wipe on startup.
//
// When combined with ImportBlocks, this triggers a full reindex (if wipeBlockTree is true)
// or chainstate-only reindex (if only wipeChainstate is true).
//
// Parameters:
//   - wipeBlockTree: Whether to wipe the block tree database (requires wipeChainstate to also be true)
//   - wipeChainstate: Whether to wipe the chainstate database
//
// Returns an error if wipeBlockTree is true but wipeChainstate is false.
func WithWipeDBs(wipeBlockTree, wipeChainstate bool) ChainstateManagerOption {
	return func(opts *C.btck_ChainstateManagerOptions) error {
		wipeBlockTreeInt := 0
		if wipeBlockTree {
			wipeBlockTreeInt = 1
		}
		wipeChainstateInt := 0
		if wipeChainstate {
			wipeChainstateInt = 1
		}
		result := C.btck_chainstate_manager_options_set_wipe_dbs(opts, C.int(wipeBlockTreeInt), C.int(wipeChainstateInt))
		if result != 0 {
			return &InternalError{"Failed to set wipe db"}
		}
		return nil
	}
}

// WithBlockTreeDBInMemory returns a ChainstateManagerOption that configures
// the block tree database to be stored in memory.
func WithBlockTreeDBInMemory() ChainstateManagerOption {
	return func(opts *C.btck_ChainstateManagerOptions) error {
		C.btck_chainstate_manager_options_update_block_tree_db_in_memory(opts, C.int(1))
		return nil
	}
}

// WithChainstateDBInMemory returns a ChainstateManagerOption that configures
// the chainstate database to be stored in memory.
func WithChainstateDBInMemory() ChainstateManagerOption {
	return func(opts *C.btck_ChainstateManagerOptions) error {
		C.btck_chainstate_manager_options_update_chainstate_db_in_memory(opts, C.int(1))
		return nil
	}
}
