package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import "runtime"

// ChainType represents the Bitcoin network type
type ChainType int

const (
	// ChainTypeMainnet represents the Bitcoin mainnet
	ChainTypeMainnet ChainType = iota
	// ChainTypeTestnet represents the Bitcoin testnet
	ChainTypeTestnet
	// ChainTypeTestnet4 represents the Bitcoin testnet4
	ChainTypeTestnet4
	// ChainTypeSignet represents the Bitcoin signet
	ChainTypeSignet
	// ChainTypeRegtest represents the Bitcoin regtest network
	ChainTypeRegtest
)

// ChainParameters wraps the C kernel_ChainParameters
type ChainParameters struct {
	ptr *C.kernel_ChainParameters
}

// NewChainParameters creates new chain parameters for the specified chain type
func NewChainParameters(chainType ChainType) (*ChainParameters, error) {
	var cChainType C.kernel_ChainType

	switch chainType {
	case ChainTypeMainnet:
		cChainType = C.kernel_CHAIN_TYPE_MAINNET
	case ChainTypeTestnet:
		cChainType = C.kernel_CHAIN_TYPE_TESTNET
	case ChainTypeTestnet4:
		cChainType = C.kernel_CHAIN_TYPE_TESTNET_4
	case ChainTypeSignet:
		cChainType = C.kernel_CHAIN_TYPE_SIGNET
	case ChainTypeRegtest:
		cChainType = C.kernel_CHAIN_TYPE_REGTEST
	default:
		return nil, ErrInvalidChainType
	}

	ptr := C.kernel_chain_parameters_create(cChainType)
	if ptr == nil {
		return nil, ErrChainParametersCreation
	}

	cp := &ChainParameters{ptr: ptr}
	runtime.SetFinalizer(cp, (*ChainParameters).destroy)
	return cp, nil
}

// destroy deallocates the chain parameters
func (cp *ChainParameters) destroy() {
	if cp.ptr != nil {
		C.kernel_chain_parameters_destroy(cp.ptr)
		cp.ptr = nil
	}
}

// Destroy explicitly destroys the chain parameters and removes the finalizer
func (cp *ChainParameters) Destroy() {
	runtime.SetFinalizer(cp, nil)
	cp.destroy()
}
