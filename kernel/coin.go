package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
)

var _ cManagedResource = &Coin{}

// Coin wraps the C btck_Coin
type Coin struct {
	ptr *C.btck_Coin
}

// GetOutput returns the transaction output held within the coin
func (c *Coin) GetOutput() (*TransactionOutput, error) {
	checkReady(c)
	ptr := C.btck_coin_get_output(c.ptr)
	if ptr == nil {
		return nil, ErrKernelTransactionOutputCreate
	}

	output := &TransactionOutput{ptr: ptr}
	runtime.SetFinalizer(output, (*TransactionOutput).destroy)
	return output, nil
}

// ConfirmationHeight returns the block height at which this coin was confirmed
func (c *Coin) ConfirmationHeight() uint32 {
	checkReady(c)
	return uint32(C.btck_coin_confirmation_height(c.ptr))
}

// IsCoinbase returns true if this coin is from a coinbase transaction
func (c *Coin) IsCoinbase() bool {
	checkReady(c)
	return int(C.btck_coin_is_coinbase(c.ptr)) != 0
}

func (c *Coin) destroy() {
	if c.ptr != nil {
		C.btck_coin_destroy(c.ptr)
		c.ptr = nil
	}
}

func (c *Coin) Destroy() {
	runtime.SetFinalizer(c, nil)
	c.destroy()
}

func (c *Coin) isReady() bool {
	return c != nil && c.ptr != nil
}

func (c *Coin) uninitializedError() error {
	return ErrCoinUninitialized
}

// Copy creates a copy of the coin
func (c *Coin) Copy() (*Coin, error) {
	if !c.isReady() {
		return nil, c.uninitializedError()
	}

	ptr := C.btck_coin_copy(c.ptr)
	if ptr == nil {
		return nil, ErrKernelCoinCopy
	}

	coin := &Coin{ptr: ptr}
	runtime.SetFinalizer(coin, (*Coin).destroy)
	return coin, nil
}
