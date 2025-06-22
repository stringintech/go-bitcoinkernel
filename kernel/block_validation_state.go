package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

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

// BlockValidationState wraps the C kernel_BlockValidationState
type BlockValidationState struct {
	ptr *C.kernel_BlockValidationState
}

func (bvs *BlockValidationState) ValidationMode() ValidationMode {
	if bvs.ptr == nil {
		return ValidationStateError
	}
	mode := C.kernel_get_validation_mode_from_block_validation_state(bvs.ptr)
	return ValidationMode(mode)
}

func (bvs *BlockValidationState) ValidationResult() BlockValidationResult {
	if bvs.ptr == nil {
		return BlockResultUnset
	}
	result := C.kernel_get_block_validation_result_from_block_validation_state(bvs.ptr)
	return BlockValidationResult(result)
}
