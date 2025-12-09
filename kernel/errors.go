package kernel

var (
	ErrKernelInstantiate = &kernelError{"Failed to instantiate btck object"}

	ErrKernelIndexOutOfBounds = &kernelError{"Index out of bounds"}

	ErrVerifyScriptVerifyTxInputIndex            = &ScriptVerifyError{"Transaction input index out of range"}
	ErrVerifyScriptVerifyInvalidFlags            = &ScriptVerifyError{"Invalid script verification flags"}
	ErrVerifyScriptVerifyInvalidFlagsCombination = &ScriptVerifyError{"Invalid combination of script verification flags"}
	ErrVerifyScriptVerifySpentOutputsMismatch    = &ScriptVerifyError{"Spent outputs count mismatch"}
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
