package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import "runtime"

var _ cManagedResource = &ChainParameters{}

// ChainParameters wraps the C btck_ChainParameters
type ChainParameters struct {
	ptr *C.btck_ChainParameters
}

func NewChainParameters(chainType ChainType) (*ChainParameters, error) {
	cType, err := chainType.c()
	if err != nil {
		return nil, err
	}
	ptr := C.btck_chain_parameters_create(cType)
	if ptr == nil {
		return nil, ErrKernelChainParametersCreate
	}

	cp := &ChainParameters{ptr: ptr}
	runtime.SetFinalizer(cp, (*ChainParameters).destroy)
	return cp, nil
}

func (cp *ChainParameters) destroy() {
	if cp.ptr != nil {
		C.btck_chain_parameters_destroy(cp.ptr)
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

const (
	ChainTypeMainnet  = C.btck_ChainType_MAINNET
	ChainTypeTestnet  = C.btck_ChainType_TESTNET
	ChainTypeTestnet4 = C.btck_ChainType_TESTNET_4
	ChainTypeSignet   = C.btck_ChainType_SIGNET
	ChainTypeRegtest  = C.btck_ChainType_REGTEST
)

type ChainType C.btck_ChainType

func (t ChainType) c() (C.btck_ChainType, error) {
	switch t {
	case ChainTypeMainnet, ChainTypeTestnet, ChainTypeTestnet4, ChainTypeSignet, ChainTypeRegtest:
		return C.btck_ChainType(t), nil
	default:
		return 0, ErrInvalidChainType
	}
}
