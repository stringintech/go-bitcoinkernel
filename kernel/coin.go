package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type coinCFuncs struct{}

func (coinCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_coin_destroy((*C.btck_Coin)(ptr))
}

func (coinCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_coin_copy((*C.btck_Coin)(ptr)))
}

// Coin holds information on a transaction output, including the height it was
// spent at and whether it is a coinbase output.
type Coin struct {
	*handle
	coinApi
}

func newCoin(ptr *C.btck_Coin, fromOwned bool) *Coin {
	h := newHandle(unsafe.Pointer(ptr), coinCFuncs{}, fromOwned)
	return &Coin{handle: h, coinApi: coinApi{(*C.btck_Coin)(h.ptr)}}
}

// CoinView holds information on a transaction output, including the height it was
// spent at and whether it is a coinbase output.
type CoinView struct {
	coinApi
	ptr *C.btck_Coin
}

func newCoinView(ptr *C.btck_Coin) *CoinView {
	return &CoinView{
		coinApi: coinApi{ptr},
		ptr:     ptr,
	}
}

type coinApi struct {
	ptr *C.btck_Coin
}

// Copy creates a copy of the coin.
func (c *coinApi) Copy() *Coin {
	return newCoin(c.ptr, false)
}

// GetOutput returns the transaction output contained in this coin.
//
// The returned TransactionOutputView is a non-owned pointer valid for the
// lifetime of this coin.
func (c *coinApi) GetOutput() *TransactionOutputView {
	ptr := C.btck_coin_get_output(c.ptr)
	return newTransactionOutputView(check(ptr))
}

// ConfirmationHeight returns the block height where the transaction that created this coin was included in.
func (c *coinApi) ConfirmationHeight() uint32 {
	return uint32(C.btck_coin_confirmation_height(c.ptr))
}

// IsCoinbase returns true if this coin originates from a coinbase transaction.
func (c *coinApi) IsCoinbase() bool {
	return int(C.btck_coin_is_coinbase(c.ptr)) != 0
}
