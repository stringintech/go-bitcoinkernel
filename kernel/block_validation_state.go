package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"

// BlockValidationState holds the state of a block during validation.
//
// Contains information about whether validation was successful, and if not,
// which step during block validation failed. This is typically provided through
// validation interface callbacks.
type BlockValidationState struct {
	ptr *C.btck_BlockValidationState
}

// ValidationMode returns whether the block is valid, invalid, or encountered an error.
//
// Returns one of:
//   - ValidationStateValid: Block passed validation
//   - ValidationStateInvalid: Block failed validation
//   - ValidationStateError: Internal error during validation
func (bvs *BlockValidationState) ValidationMode() ValidationMode {
	mode := C.btck_block_validation_state_get_validation_mode(bvs.ptr)
	return ValidationMode(mode)
}

// ValidationResult returns a granular reason for why the block was invalid.
//
// This provides detailed information about the specific validation failure, such as
// consensus violations, invalid headers, or missing previous blocks.
func (bvs *BlockValidationState) ValidationResult() BlockValidationResult {
	result := C.btck_block_validation_state_get_block_validation_result(bvs.ptr)
	return BlockValidationResult(result)
}

// ValidationMode indicates whether a validated data structure is valid, invalid,
// or an error was encountered during processing.
type ValidationMode C.btck_ValidationMode

const (
	ValidationStateValid   ValidationMode = C.btck_ValidationMode_VALID          // Data structure passed validation
	ValidationStateInvalid ValidationMode = C.btck_ValidationMode_INVALID        // Data structure failed validation
	ValidationStateError   ValidationMode = C.btck_ValidationMode_INTERNAL_ERROR // Internal error during validation
)

// BlockValidationResult provides a granular reason why a block was invalid.
type BlockValidationResult C.btck_BlockValidationResult

const (
	BlockResultUnset   BlockValidationResult = C.btck_BlockValidationResult_UNSET           // Initial value, block has not yet been rejected
	BlockConsensus     BlockValidationResult = C.btck_BlockValidationResult_CONSENSUS       // Invalid by consensus rules (excluding specific reasons below)
	BlockCachedInvalid BlockValidationResult = C.btck_BlockValidationResult_CACHED_INVALID  // Block was cached as invalid without storing the reason
	BlockInvalidHeader BlockValidationResult = C.btck_BlockValidationResult_INVALID_HEADER  // Invalid proof of work or timestamp too old
	BlockMutated       BlockValidationResult = C.btck_BlockValidationResult_MUTATED         // Block data didn't match the data committed to by PoW
	BlockMissingPrev   BlockValidationResult = C.btck_BlockValidationResult_MISSING_PREV    // Previous block is missing
	BlockInvalidPrev   BlockValidationResult = C.btck_BlockValidationResult_INVALID_PREV    // Previous block is invalid
	BlockTimeFuture    BlockValidationResult = C.btck_BlockValidationResult_TIME_FUTURE     // Block timestamp was >2 hours in the future
	BlockHeaderLowWork BlockValidationResult = C.btck_BlockValidationResult_HEADER_LOW_WORK // Block header may be on a too-little-work chain
)
