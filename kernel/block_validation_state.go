package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"

type BlockValidationState struct {
	ptr *C.btck_BlockValidationState
}

func (bvs *BlockValidationState) ValidationMode() ValidationMode {
	mode := C.btck_block_validation_state_get_validation_mode(bvs.ptr)
	return ValidationMode(mode)
}

func (bvs *BlockValidationState) ValidationResult() BlockValidationResult {
	result := C.btck_block_validation_state_get_block_validation_result(bvs.ptr)
	return BlockValidationResult(result)
}

const (
	ValidationStateValid   = C.btck_ValidationMode_VALID
	ValidationStateInvalid = C.btck_ValidationMode_INVALID
	ValidationStateError   = C.btck_ValidationMode_INTERNAL_ERROR
)

type ValidationMode C.btck_ValidationMode

const (
	BlockResultUnset   = C.btck_BlockValidationResult_UNSET
	BlockConsensus     = C.btck_BlockValidationResult_CONSENSUS
	BlockCachedInvalid = C.btck_BlockValidationResult_CACHED_INVALID
	BlockInvalidHeader = C.btck_BlockValidationResult_INVALID_HEADER
	BlockMutated       = C.btck_BlockValidationResult_MUTATED
	BlockMissingPrev   = C.btck_BlockValidationResult_MISSING_PREV
	BlockInvalidPrev   = C.btck_BlockValidationResult_INVALID_PREV
	BlockTimeFuture    = C.btck_BlockValidationResult_TIME_FUTURE
	BlockHeaderLowWork = C.btck_BlockValidationResult_HEADER_LOW_WORK
)

type BlockValidationResult C.btck_BlockValidationResult
