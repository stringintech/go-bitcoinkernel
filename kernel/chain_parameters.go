package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type chainParametersCFuncs struct{}

func (chainParametersCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_chain_parameters_destroy((*C.btck_ChainParameters)(ptr))
}

func (chainParametersCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_chain_parameters_copy((*C.btck_ChainParameters)(ptr)))
}

type ChainParameters struct {
	*handle
}

func newChainParameters(ptr *C.btck_ChainParameters, fromOwned bool) *ChainParameters {
	h := newHandle(unsafe.Pointer(ptr), chainParametersCFuncs{}, fromOwned)
	return &ChainParameters{handle: h}
}

func NewChainParameters(chainType ChainType) (*ChainParameters, error) {
	ptr := C.btck_chain_parameters_create(chainType.c())
	return newChainParameters(check(ptr), true), nil
}

func (cp *ChainParameters) Copy() *ChainParameters {
	return newChainParameters((*C.btck_ChainParameters)(cp.ptr), false)
}

const (
	ChainTypeMainnet  = C.btck_ChainType_MAINNET
	ChainTypeTestnet  = C.btck_ChainType_TESTNET
	ChainTypeTestnet4 = C.btck_ChainType_TESTNET_4
	ChainTypeSignet   = C.btck_ChainType_SIGNET
	ChainTypeRegtest  = C.btck_ChainType_REGTEST
)

type ChainType C.btck_ChainType

func (t ChainType) c() C.btck_ChainType {
	switch t {
	case ChainTypeMainnet, ChainTypeTestnet, ChainTypeTestnet4, ChainTypeSignet, ChainTypeRegtest:
		return C.btck_ChainType(t)
	default:
		panic("Invalid chain type")
	}
}
