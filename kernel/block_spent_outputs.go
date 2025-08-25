package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &BlockSpentOutputs{}

// BlockSpentOutputs wraps the C btck_BlockSpentOutputs
type BlockSpentOutputs struct {
	ptr *C.btck_BlockSpentOutputs
}

// Count returns the number of transaction spent outputs contained in block spent outputs
func (bso *BlockSpentOutputs) Count() uint64 {
	checkReady(bso)
	return uint64(C.btck_block_spent_outputs_count(bso.ptr))
}

// GetTransactionSpentOutputsAt returns the transaction spent outputs at the specified index
func (bso *BlockSpentOutputs) GetTransactionSpentOutputsAt(index uint64) (*TransactionSpentOutputs, error) {
	checkReady(bso)
	ptr := C.btck_block_spent_outputs_get_transaction_spent_outputs_at(bso.ptr, C.size_t(index))
	if ptr == nil {
		return nil, ErrKernelBlockSpentOutputsGetTransactionSpentOutputsAt
	}

	txSpentOutputs := &TransactionSpentOutputs{ptr: ptr}
	runtime.SetFinalizer(txSpentOutputs, (*TransactionSpentOutputs).destroy)
	return txSpentOutputs, nil
}

// Copy creates a copy of the block spent outputs
func (bso *BlockSpentOutputs) Copy() (*BlockSpentOutputs, error) {
	checkReady(bso)

	ptr := C.btck_block_spent_outputs_copy(bso.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockSpentOutputsCopy
	}

	outputs := &BlockSpentOutputs{ptr: ptr}
	runtime.SetFinalizer(outputs, (*BlockSpentOutputs).destroy)
	return outputs, nil
}

func (bso *BlockSpentOutputs) destroy() {
	if bso.ptr != nil {
		C.btck_block_spent_outputs_destroy(bso.ptr)
		bso.ptr = nil
	}
}

func (bso *BlockSpentOutputs) Destroy() {
	runtime.SetFinalizer(bso, nil)
	bso.destroy()
}

func (bso *BlockSpentOutputs) isReady() bool {
	return bso != nil && bso.ptr != nil
}

func (bso *BlockSpentOutputs) uninitializedError() error {
	return ErrBlockSpentOutputsUninitialized
}
