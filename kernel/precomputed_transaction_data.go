package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type precomputedTransactionDataCFuncs struct{}

func (precomputedTransactionDataCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_precomputed_transaction_data_destroy((*C.btck_PrecomputedTransactionData)(ptr))
}

func (precomputedTransactionDataCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_precomputed_transaction_data_copy((*C.btck_PrecomputedTransactionData)(ptr)))
}

type PrecomputedTransactionData struct {
	*handle
}

func newPrecomputedTransactionData(ptr *C.btck_PrecomputedTransactionData, fromOwned bool) *PrecomputedTransactionData {
	h := newHandle(unsafe.Pointer(ptr), precomputedTransactionDataCFuncs{}, fromOwned)
	return &PrecomputedTransactionData{handle: h}
}

// NewPrecomputedTransactionData creates precomputed transaction data for script verification.
//
// Precomputed data is reusable when verifying multiple inputs of the same transaction.
// This avoids recomputing transaction hashes for each input.
//
// Required when verifying a taproot input.
//
// The underlying C implementation copies the spent outputs data, so the spentOutputs slice
// and its elements can be safely freed after this function returns.
//
// Parameters:
//   - txTo: The transaction to precompute data for
//   - spentOutputs: Outputs spent by the transaction. May be nil for non-taproot verification.
//
// Returns an error if the precomputation fails.
func NewPrecomputedTransactionData(txTo TransactionLike, spentOutputs []TransactionOutputLike) (*PrecomputedTransactionData, error) {
	var cSpentOutputsPtr **C.btck_TransactionOutput
	if len(spentOutputs) > 0 {
		cSpentOutputs := make([]*C.btck_TransactionOutput, len(spentOutputs))
		for i, output := range spentOutputs {
			cSpentOutputs[i] = output.cPtr()
		}
		cSpentOutputsPtr = (**C.btck_TransactionOutput)(unsafe.Pointer(&cSpentOutputs[0]))
	}

	ptr := C.btck_precomputed_transaction_data_create(
		txTo.cPtr(),
		cSpentOutputsPtr,
		C.size_t(len(spentOutputs)),
	)
	if ptr == nil {
		return nil, &InternalError{"Failed to create precomputed transaction data"}
	}
	return newPrecomputedTransactionData(ptr, true), nil
}

// Copy creates a copy of the precomputed transaction data.
func (p *PrecomputedTransactionData) Copy() *PrecomputedTransactionData {
	return newPrecomputedTransactionData((*C.btck_PrecomputedTransactionData)(p.ptr), false)
}
