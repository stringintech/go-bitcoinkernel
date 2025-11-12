package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type contextCFuncs struct{}

func (contextCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_context_destroy((*C.btck_Context)(ptr))
}

func (contextCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_context_copy((*C.btck_Context)(ptr)))
}

// Context is a central structure that holds chain-specific parameters and
// callbacks for handling error and validation events.
//
// Once other validation objects are instantiated from it, the context is kept in
// memory for the duration of their lifetimes.
//
// A constructed context can be safely used from multiple threads.
type Context struct {
	*handle
}

func newContext(ptr *C.btck_Context, fromOwned bool) *Context {
	h := newHandle(unsafe.Pointer(ptr), contextCFuncs{}, fromOwned)
	return &Context{handle: h}
}

// NewContext creates a new kernel context.
//
// Parameters:
//   - options: Context configuration created by NewContextOptions (can be nil for defaults)
//
// Returns an error if the context cannot be created. If options is nil or not configured,
// the context assumes mainnet chain parameters and no callbacks.
func NewContext(options *ContextOptions) (*Context, error) {
	ptr := C.btck_context_create((*C.btck_ContextOptions)(options.ptr))
	if ptr == nil {
		return nil, &InternalError{"Failed to create context"}
	}
	return newContext(ptr, true), nil
}

// Interrupt halts long-running validation functions like reindexing or block import.
//
// Returns an error if the interrupt signal cannot be delivered.
func (ctx *Context) Interrupt() error {
	result := C.btck_context_interrupt((*C.btck_Context)(ctx.ptr))
	if result != 0 {
		return &InternalError{"Context interrupt failed"}
	}
	return nil
}

// Copy creates a shallow copy of the context by incrementing its reference count.
//
// The context is reference-counted internally, so this operation is efficient and does
// not duplicate the underlying data.
func (ctx *Context) Copy() *Context {
	return newContext((*C.btck_Context)(ctx.ptr), false)
}
