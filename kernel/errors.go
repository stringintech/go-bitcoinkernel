package kernel

var (
	ErrKernelInstantiate = &kernelError{"Failed to instantiate btck object"}

	ErrKernelIndexOutOfBounds = &kernelError{"Index out of bounds"}

	ErrVerifyScriptVerifyTxInputIndex            = &ScriptVerifyError{"Transaction input index out of range"}
	ErrVerifyScriptVerifyInvalidFlags            = &ScriptVerifyError{"Invalid script verification flags"}
	ErrVerifyScriptVerifyInvalidFlagsCombination = &ScriptVerifyError{"Invalid combination of script verification flags"}
	ErrVerifyScriptVerifySpentOutputsRequired    = &ScriptVerifyError{"Spent outputs required for verification"}
)

// check panics if ptr is nil, otherwise returns ptr unchanged; used when C calls are not expected to return null
func check[T any](ptr T) T {
	if any(ptr) == nil {
		panic(ErrKernelInstantiate)
	}
	return ptr
}

type KernelError interface {
	Error() string
	isKernelError()
}

type kernelError struct {
	Msg string
}

func (e *kernelError) Error() string {
	return e.Msg
}

func (e *kernelError) isKernelError() {}

// InternalError is returned when a call to the underlying library fails.
type InternalError struct {
	Msg string
}

func (e *InternalError) Error() string {
	return e.Msg
}

func (e *InternalError) isKernelError() {}

type SerializationError struct {
	Msg string
}

func (e *SerializationError) Error() string {
	return e.Msg
}

func (e *SerializationError) isKernelError() {}

// ScriptVerifyError represents errors that prevent script verification from executing.
type ScriptVerifyError struct {
	Msg string
}

func (e *ScriptVerifyError) Error() string {
	return "Script verification failed: " + e.Msg
}

func (e *ScriptVerifyError) isKernelError() {}

// BlockValidationError represents a block or header validation failure with detailed state information.
type BlockValidationError struct {
	State *BlockValidationState
}

func (e *BlockValidationError) Error() string {
	if e.State == nil {
		return "Block validation failed"
	}
	mode := e.State.ValidationMode()
	result := e.State.ValidationResult()
	return "Block validation failed: mode=" + validationModeString(mode) + ", result=" + blockValidationResultString(result)
}

func (e *BlockValidationError) isKernelError() {}

// GetState returns the validation state containing detailed error information.
func (e *BlockValidationError) GetState() *BlockValidationState {
	return e.State
}

func validationModeString(mode ValidationMode) string {
	switch mode {
	case ValidationStateValid:
		return "valid"
	case ValidationStateInvalid:
		return "invalid"
	case ValidationStateError:
		return "error"
	default:
		return "unknown"
	}
}

func blockValidationResultString(result BlockValidationResult) string {
	switch result {
	case BlockResultUnset:
		return "unset"
	case BlockConsensus:
		return "consensus"
	case BlockCachedInvalid:
		return "cached_invalid"
	case BlockInvalidHeader:
		return "invalid_header"
	case BlockMutated:
		return "mutated"
	case BlockMissingPrev:
		return "missing_prev"
	case BlockInvalidPrev:
		return "invalid_prev"
	case BlockTimeFuture:
		return "time_future"
	case BlockHeaderLowWork:
		return "header_low_work"
	default:
		return "unknown"
	}
}
