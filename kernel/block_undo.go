package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &BlockUndo{}

// BlockUndo wraps the C kernel_BlockUndo
type BlockUndo struct {
	ptr *C.kernel_BlockUndo
}

// Size returns the number of transactions whose undo data is contained in block undo
func (bu *BlockUndo) Size() uint64 {
	checkReady(bu)
	return uint64(C.kernel_block_undo_size(bu.ptr))
}

// GetTransactionUndoSize returns the number of previous transaction outputs
// contained in the transaction undo data at the specified index
func (bu *BlockUndo) GetTransactionUndoSize(transactionUndoIndex uint64) uint64 {
	checkReady(bu)
	return uint64(C.kernel_get_transaction_undo_size(bu.ptr, C.uint64_t(transactionUndoIndex)))
}

// GetUndoOutputHeightByIndex returns the block height of the block that contains
// the output at output_index within the transaction undo data at the provided index
func (bu *BlockUndo) GetUndoOutputHeightByIndex(transactionUndoIndex, outputIndex uint64) uint32 {
	checkReady(bu)
	return uint32(C.kernel_get_undo_output_height_by_index(bu.ptr, C.uint64_t(transactionUndoIndex), C.uint64_t(outputIndex)))
}

// GetUndoOutputByIndex returns a transaction output contained in the transaction
// undo data at the specified indices
func (bu *BlockUndo) GetUndoOutputByIndex(transactionUndoIndex, outputIndex uint64) (*TransactionOutput, error) {
	checkReady(bu)

	ptr := C.kernel_get_undo_output_by_index(bu.ptr, C.uint64_t(transactionUndoIndex), C.uint64_t(outputIndex))
	if ptr == nil {
		return nil, ErrKernelGetUndoOutputByIndex
	}

	output := &TransactionOutput{ptr: ptr}
	runtime.SetFinalizer(output, (*TransactionOutput).destroy)
	return output, nil
}

func (bu *BlockUndo) destroy() {
	if bu.ptr != nil {
		C.kernel_block_undo_destroy(bu.ptr)
		bu.ptr = nil
	}
}

func (bu *BlockUndo) Destroy() {
	runtime.SetFinalizer(bu, nil)
	bu.destroy()
}

func (bu *BlockUndo) isReady() bool {
	return bu != nil && bu.ptr != nil
}

func (bu *BlockUndo) uninitializedError() error {
	return ErrBlockUndoUninitialized
}
