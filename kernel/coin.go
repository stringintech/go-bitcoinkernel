package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
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
	return &Coin{
		handle: h,
		coinApi: coinApi{
			ptr: func() *C.btck_Coin {
				return (*C.btck_Coin)(h.ptr)
			},
			owner: h,
		},
	}
}

// CoinView holds information on a transaction output, including the height it was
// spent at and whether it is a coinbase output.
type CoinView struct {
	coinApi
}

func newCoinView(ptr *C.btck_Coin, owner cResourceHandle) *CoinView {
	return &CoinView{
		coinApi: coinApi{
			ptr: func() *C.btck_Coin {
				return ptr
			},
			owner: owner,
		},
	}
}

type coinApi struct {
	ptr   func() *C.btck_Coin
	owner cResourceHandle
}

func (c *coinApi) cPtr() *C.btck_Coin {
	return c.ptr()
}

// CoinLike is implemented by *Coin and *CoinView.
type CoinLike interface {
	cPtr() *C.btck_Coin
	Copy() *Coin
	GetOutput() *TransactionOutputView
	ConfirmationHeight() uint32
	IsCoinbase() bool
}

var _ CoinLike = (*Coin)(nil)
var _ CoinLike = (*CoinView)(nil)

// Copy creates a copy of the coin.
func (c *coinApi) Copy() *Coin {
	coin := newCoin(c.ptr(), false)
	runtime.KeepAlive(c)
	return coin
}

// GetOutput returns the transaction output contained in this coin.
//
// The returned TransactionOutputView is a non-owned pointer valid for the
// lifetime of this coin.
func (c *coinApi) GetOutput() *TransactionOutputView {
	ptr := C.btck_coin_get_output(c.ptr())
	runtime.KeepAlive(c)
	return newTransactionOutputView(check(ptr), c.owner)
}

// ConfirmationHeight returns the block height where the transaction that created this coin was included in.
func (c *coinApi) ConfirmationHeight() uint32 {
	height := uint32(C.btck_coin_confirmation_height(c.ptr()))
	runtime.KeepAlive(c)
	return height
}

// IsCoinbase returns true if this coin originates from a coinbase transaction.
func (c *coinApi) IsCoinbase() bool {
	coinbase := int(C.btck_coin_is_coinbase(c.ptr())) != 0
	runtime.KeepAlive(c)
	return coinbase
}
