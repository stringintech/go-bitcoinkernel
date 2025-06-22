package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"errors"
	"runtime"
)

// Context wraps the C kernel_Context
// Once other validation objects are instantiated from it, the context needs to be kept in
// memory for the duration of their lifetimes.
//
// A constructed context can be safely used from multiple threads.
type Context struct {
	ptr *C.kernel_Context
}

// NewContext creates a new kernel context with the specified options.
// Kernel copies all necessary data from the options during context creation,
// so the caller can safely free the options object after this call returns.
func NewContext(options *ContextOptions) (*Context, error) {
	if options == nil {
		return nil, errors.New("context options cannot be nil")
	}

	ptr := C.kernel_context_create(options.ptr)
	if ptr == nil {
		return nil, ErrContextCreation
	}

	ctx := &Context{ptr: ptr}
	runtime.SetFinalizer(ctx, (*Context).destroy)
	return ctx, nil
}

// NewDefaultContext creates a new kernel context with default mainnet parameters.
// The defer statements are safe here because the Kernel copies all necessary
// data during context creation, so the caller can safely free the options and
// parameters objects immediately after the context is created.
func NewDefaultContext() (*Context, error) {
	opts, err := NewContextOptions()
	if err != nil {
		return nil, err
	}
	defer opts.Destroy()

	params, err := NewChainParameters(ChainTypeMainnet)
	if err != nil {
		return nil, err
	}
	defer params.Destroy()

	opts.SetChainParams(params)
	return NewContext(opts)
}

// Interrupt can be used to halt long-running validation functions
func (ctx *Context) Interrupt() bool {
	if ctx.ptr == nil {
		return false
	}
	return bool(C.kernel_context_interrupt(ctx.ptr))
}

func (ctx *Context) destroy() {
	if ctx.ptr != nil {
		C.kernel_context_destroy(ctx.ptr)
		ctx.ptr = nil
	}
}

func (ctx *Context) Destroy() {
	runtime.SetFinalizer(ctx, nil)
	ctx.destroy()
}
