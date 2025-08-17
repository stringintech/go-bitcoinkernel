package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

var _ cResource = &BlockValidationState{}

// BlockValidationState wraps the C btck_BlockValidationState
type BlockValidationState struct {
	ptr *C.btck_BlockValidationState
}

func (bvs *BlockValidationState) ValidationMode() ValidationMode {
	checkReady(bvs)
	mode := C.btck_block_validation_state_get_validation_mode(bvs.ptr)
	return ValidationMode(mode)
}

func (bvs *BlockValidationState) ValidationResult() BlockValidationResult {
	checkReady(bvs)
	result := C.btck_block_validation_state_get_block_validation_result(bvs.ptr)
	return BlockValidationResult(result)
}

func (bvs *BlockValidationState) isReady() bool {
	return bvs != nil && bvs.ptr != nil
}

func (bvs *BlockValidationState) uninitializedError() error {
	return ErrBlockValidationStateUninitialized
}

// ValidationMode represents the validation state mode
type ValidationMode int

const (
	ValidationStateValid ValidationMode = iota
	ValidationStateInvalid
	ValidationStateError
)

// BlockValidationResult represents the validation result for a block
type BlockValidationResult int

const (
	BlockResultUnset BlockValidationResult = iota
	BlockConsensus
	BlockCachedInvalid
	BlockInvalidHeader
	BlockMutated
	BlockMissingPrev
	BlockInvalidPrev
	BlockTimeFuture
	BlockHeaderLowWork
)
