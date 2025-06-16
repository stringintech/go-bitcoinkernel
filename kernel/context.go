package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"errors"
	"runtime"
)

// ContextOptions wraps the C kernel_ContextOptions
type ContextOptions struct {
	ptr *C.kernel_ContextOptions
}

// NewContextOptions creates new context options
func NewContextOptions() (*ContextOptions, error) {
	ptr := C.kernel_context_options_create()
	if ptr == nil {
		return nil, ErrContextOptionsCreation
	}

	opts := &ContextOptions{ptr: ptr}
	runtime.SetFinalizer(opts, (*ContextOptions).destroy)
	return opts, nil
}

// SetChainParams sets the chain parameters for these context options.
// The C++ kernel makes a copy of the chain parameters, so the caller can
// safely free the chainParams object after this call returns.
func (opts *ContextOptions) SetChainParams(chainParams *ChainParameters) {
	if opts.ptr != nil && chainParams.ptr != nil {
		C.kernel_context_options_set_chainparams(opts.ptr, chainParams.ptr)
	}
}

// destroy deallocates the context options
func (opts *ContextOptions) destroy() {
	if opts.ptr != nil {
		C.kernel_context_options_destroy(opts.ptr)
		opts.ptr = nil
	}
}

// Close explicitly destroys the context options and removes the finalizer
func (opts *ContextOptions) Close() {
	runtime.SetFinalizer(opts, nil)
	opts.destroy()
}

// Context wraps the C kernel_Context
type Context struct {
	ptr *C.kernel_Context
}

// NewContext creates a new kernel context with the specified options.
// The C++ kernel copies all necessary data from the options during context creation,
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
// The defer statements are safe here because the C++ kernel copies all necessary
// data during context creation, allowing us to free the temporary options and
// parameters objects immediately after the context is created.
func NewDefaultContext() (*Context, error) {
	// Create default options for mainnet
	opts, err := NewContextOptions()
	if err != nil {
		return nil, err
	}
	// Safe to defer: C++ kernel copies data from options during context creation
	defer opts.Close()

	// Create mainnet chain parameters
	params, err := NewChainParameters(ChainTypeMainnet)
	if err != nil {
		return nil, err
	}
	// Safe to defer: C++ kernel copies chain params when set on options
	defer params.Close()

	// Set chain parameters on options (C++ makes internal copy)
	opts.SetChainParams(params)

	// Create context with configured options (C++ copies all needed data)
	return NewContext(opts)
}

// Interrupt can be used to halt long-running validation functions
func (ctx *Context) Interrupt() bool {
	if ctx.ptr == nil {
		return false
	}
	return bool(C.kernel_context_interrupt(ctx.ptr))
}

// destroy deallocates the context
func (ctx *Context) destroy() {
	if ctx.ptr != nil {
		C.kernel_context_destroy(ctx.ptr)
		ctx.ptr = nil
	}
}

// Close explicitly destroys the context and removes the finalizer
func (ctx *Context) Close() {
	runtime.SetFinalizer(ctx, nil)
	ctx.destroy()
}

// IsValid returns true if the context is valid (non-nil)
func (ctx *Context) IsValid() bool {
	return ctx.ptr != nil
}
