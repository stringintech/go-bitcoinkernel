package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

// ContextOptions wraps the C kernel_ContextOptions
type ContextOptions struct {
	ptr *C.kernel_ContextOptions
}

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
// Kernel makes a copy of the chain parameters, so the caller can
// safely free the chainParams object after this call returns.
func (opts *ContextOptions) SetChainParams(chainParams *ChainParameters) {
	if opts.ptr == nil {
		panic("ContextOptions pointer is nil")
	}
	if chainParams.ptr == nil {
		panic("ChainParameters pointer is nil")
	}
	C.kernel_context_options_set_chainparams(opts.ptr, chainParams.ptr)
}

func (opts *ContextOptions) destroy() {
	if opts.ptr != nil {
		C.kernel_context_options_destroy(opts.ptr)
		opts.ptr = nil
	}
}

func (opts *ContextOptions) Destroy() {
	runtime.SetFinalizer(opts, nil)
	opts.destroy()
}
