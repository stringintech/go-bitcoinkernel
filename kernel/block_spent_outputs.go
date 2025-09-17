package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type blockSpentOutputsCFuncs struct{}

func (blockSpentOutputsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_block_spent_outputs_destroy((*C.btck_BlockSpentOutputs)(ptr))
}

func (blockSpentOutputsCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_block_spent_outputs_copy((*C.btck_BlockSpentOutputs)(ptr)))
}

type BlockSpentOutputs struct {
	*handle
}

func newBlockSpentOutputs(ptr *C.btck_BlockSpentOutputs, fromOwned bool) *BlockSpentOutputs {
	h := newHandle(unsafe.Pointer(ptr), blockSpentOutputsCFuncs{}, fromOwned)
	return &BlockSpentOutputs{handle: h}
}

func (bso *BlockSpentOutputs) Count() uint64 {
	return uint64(C.btck_block_spent_outputs_count((*C.btck_BlockSpentOutputs)(bso.ptr)))
}

// GetTransactionSpentOutputsAt returns the transaction spent outputs at the specified index
func (bso *BlockSpentOutputs) GetTransactionSpentOutputsAt(index uint64) (*TransactionSpentOutputsView, error) {
	if index >= bso.Count() {
		return nil, ErrKernelIndexOutOfBounds
	}
	ptr := C.btck_block_spent_outputs_get_transaction_spent_outputs_at((*C.btck_BlockSpentOutputs)(bso.ptr), C.size_t(index))
	return newTransactionSpentOutputsView(check(ptr)), nil
}

func (bso *BlockSpentOutputs) Copy() *BlockSpentOutputs {
	return newBlockSpentOutputs((*C.btck_BlockSpentOutputs)(bso.ptr), false)
}
