package kernel

/*
#include "bitcoinkernel.h"
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

// ChainParameters describes the properties of a blockchain network.
//
// These are placed into a kernel context through the kernel context options. The
// parameters may be instantiated for either mainnet, testnet, testnet4, signet, or
// regtest.
type ChainParameters struct {
	*handle
}

func newChainParameters(ptr *C.btck_ChainParameters, fromOwned bool) *ChainParameters {
	h := newHandle(unsafe.Pointer(ptr), chainParametersCFuncs{}, fromOwned)
	return &ChainParameters{handle: h}
}

// NewChainParameters creates chain parameters with default settings for the specified chain type.
//
// The chain parameters describe the properties of a blockchain network and are used
// when creating a kernel context.
//
// Parameters:
//   - chainType: One of ChainTypeMainnet, ChainTypeTestnet, ChainTypeTestnet4, ChainTypeSignet, or ChainTypeRegtest
func NewChainParameters(chainType ChainType) (*ChainParameters, error) {
	ptr := C.btck_chain_parameters_create(chainType.c())
	return newChainParameters(check(ptr), true), nil
}

// Copy creates a copy of the chain parameters.
func (cp *ChainParameters) Copy() *ChainParameters {
	return newChainParameters((*C.btck_ChainParameters)(cp.ptr), false)
}

type ChainType C.btck_ChainType

const (
	ChainTypeMainnet  ChainType = C.btck_ChainType_MAINNET
	ChainTypeTestnet  ChainType = C.btck_ChainType_TESTNET
	ChainTypeTestnet4 ChainType = C.btck_ChainType_TESTNET_4
	ChainTypeSignet   ChainType = C.btck_ChainType_SIGNET
	ChainTypeRegtest  ChainType = C.btck_ChainType_REGTEST
)

func (t ChainType) c() C.btck_ChainType {
	switch t {
	case ChainTypeMainnet, ChainTypeTestnet, ChainTypeTestnet4, ChainTypeSignet, ChainTypeRegtest:
		return C.btck_ChainType(t)
	default:
		panic("Invalid chain type")
	}
}
