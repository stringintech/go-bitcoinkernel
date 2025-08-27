package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &Context{}

// Context wraps the C btck_Context
// Once other validation objects are instantiated from it, the context needs to be kept in
// memory for the duration of their lifetimes.
//
// A constructed context can be safely used from multiple threads.
type Context struct {
	ptr     *C.btck_Context
	options *ContextOptions
}

// NewContext creates a new kernel context with the specified options.
func NewContext(options *ContextOptions) (*Context, error) {
	if err := validateReady(options); err != nil {
		return nil, err
	}

	ptr := C.btck_context_create(options.ptr)
	if ptr == nil {
		return nil, ErrKernelContextCreate
	}

	ctx := &Context{
		ptr:     ptr,
		options: options,
	}
	runtime.SetFinalizer(ctx, (*Context).destroy)
	return ctx, nil
}

// NewDefaultContext creates a new kernel context with default mainnet parameters.
func NewDefaultContext() (*Context, error) {
	opts, err := NewContextOptions()
	if err != nil {
		return nil, err
	}

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
	checkReady(ctx)
	return C.btck_context_interrupt(ctx.ptr) != 0
}

func (ctx *Context) destroy() {
	if ctx.ptr != nil {
		C.btck_context_destroy(ctx.ptr)
		ctx.ptr = nil
	}
	if ctx.options != nil {
		ctx.options.destroy()
		ctx.options = nil
	}
}

func (ctx *Context) Destroy() {
	runtime.SetFinalizer(ctx, nil)
	ctx.destroy()
}

func (ctx *Context) isReady() bool {
	return ctx != nil && ctx.ptr != nil
}

func (ctx *Context) uninitializedError() error {
	return ErrContextUninitialized
}
