package kernel

/*
#include "kernel/bitcoinkernel.h"
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

type Coin struct {
	*handle
	coinApi
}

func newCoin(ptr *C.btck_Coin, fromOwned bool) *Coin {
	h := newHandle(unsafe.Pointer(ptr), coinCFuncs{}, fromOwned)
	return &Coin{handle: h, coinApi: coinApi{(*C.btck_Coin)(h.ptr)}}
}

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

func (c *coinApi) Copy() *Coin {
	return newCoin(c.ptr, false)
}

func (c *coinApi) GetOutput() *TransactionOutputView {
	ptr := C.btck_coin_get_output(c.ptr)
	return newTransactionOutputView(check(ptr))
}

func (c *coinApi) ConfirmationHeight() uint32 {
	return uint32(C.btck_coin_confirmation_height(c.ptr))
}

// IsCoinbase returns true if this coin is from a coinbase transaction
func (c *coinApi) IsCoinbase() bool {
	return int(C.btck_coin_is_coinbase(c.ptr)) != 0
}
