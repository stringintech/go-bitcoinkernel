package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"encoding/hex"
	"unsafe"
)

type txidCFuncs struct{}

func (txidCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_txid_destroy((*C.btck_Txid)(ptr))
}

func (txidCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_txid_copy((*C.btck_Txid)(ptr)))
}

type Txid struct {
	*handle
	txidApi
}

func newTxid(ptr *C.btck_Txid, fromOwned bool) *Txid {
	h := newHandle(unsafe.Pointer(ptr), txidCFuncs{}, fromOwned)
	return &Txid{handle: h, txidApi: txidApi{(*C.btck_Txid)(h.ptr)}}
}

type TxidView struct {
	txidApi
	ptr *C.btck_Txid
}

func newTxidView(ptr *C.btck_Txid) *TxidView {
	return &TxidView{
		txidApi: txidApi{ptr},
		ptr:     ptr,
	}
}

type txidApi struct {
	ptr *C.btck_Txid
}

// Copy creates a copy of the txid.
func (t *txidApi) Copy() *Txid {
	return newTxid(t.ptr, false)
}

// Equals checks if two txids are equal.
func (t *txidApi) Equals(other *Txid) bool {
	return C.btck_txid_equals(t.ptr, other.txidApi.ptr) != 0
}

// Bytes returns the 32-byte representation of the txid.
func (t *txidApi) Bytes() [32]byte {
	var output [32]C.uchar
	C.btck_txid_to_bytes(t.ptr, &output[0])
	return *(*[32]byte)(unsafe.Pointer(&output[0]))
}

// String returns the txid as a hex string in display order (reversed).
func (t *txidApi) String() string {
	hashBytes := t.Bytes()
	return hex.EncodeToString(ReverseBytes(hashBytes[:]))
}
