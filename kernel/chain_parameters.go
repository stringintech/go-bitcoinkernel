package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import "runtime"

var _ cManagedResource = &ChainParameters{}

// ChainParameters wraps the C kernel_ChainParameters
type ChainParameters struct {
	ptr *C.kernel_ChainParameters
}

func NewChainParameters(chainType ChainType) (*ChainParameters, error) {
	cType, err := chainType.c()
	if err != nil {
		return nil, err
	}
	ptr := C.kernel_chain_parameters_create(cType)
	if ptr == nil {
		return nil, ErrKernelChainParametersCreate
	}

	cp := &ChainParameters{ptr: ptr}
	runtime.SetFinalizer(cp, (*ChainParameters).destroy)
	return cp, nil
}

func (cp *ChainParameters) destroy() {
	if cp.ptr != nil {
		C.kernel_chain_parameters_destroy(cp.ptr)
		cp.ptr = nil
	}
}

func (cp *ChainParameters) Destroy() {
	runtime.SetFinalizer(cp, nil)
	cp.destroy()
}

func (cp *ChainParameters) isReady() bool {
	return cp != nil && cp.ptr != nil
}

func (cp *ChainParameters) uninitializedError() error {
	return ErrChainParametersUninitialized
}

// ChainType represents the Bitcoin network type
type ChainType int

const (
	ChainTypeMainnet ChainType = iota
	ChainTypeTestnet
	ChainTypeTestnet4
	ChainTypeSignet
	ChainTypeRegtest
)

func (t ChainType) c() (C.kernel_ChainType, error) {
	if t < ChainTypeMainnet || t > ChainTypeRegtest {
		return 0, ErrInvalidChainType
	}
	return C.kernel_ChainType(t), nil
}
