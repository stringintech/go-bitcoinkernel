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

// Context wraps the C btck_Context
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

func NewContext(options *ContextOptions) (*Context, error) {
	ptr := C.btck_context_create((*C.btck_ContextOptions)(options.ptr))
	if ptr == nil {
		return nil, &InternalError{"Failed to create context"}
	}
	return newContext(ptr, true), nil
}

// Interrupt can be used to halt long-running validation functions
func (ctx *Context) Interrupt() error {
	result := C.btck_context_interrupt((*C.btck_Context)(ctx.handle.ptr))
	if result != 0 {
		return &InternalError{"Context interrupt failed"}
	}
	return nil
}

func (ctx *Context) Copy() *Context {
	return newContext((*C.btck_Context)(ctx.handle.ptr), false)
}
