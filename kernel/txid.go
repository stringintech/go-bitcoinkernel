package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"encoding/hex"
	"runtime"
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
	return &Txid{
		handle: h,
		txidApi: txidApi{
			ptr: func() *C.btck_Txid {
				return (*C.btck_Txid)(h.ptr)
			},
			owner: h,
		},
	}
}

type TxidView struct {
	txidApi
}

func newTxidView(ptr *C.btck_Txid, owner cResourceHandle) *TxidView {
	return &TxidView{
		txidApi: txidApi{
			ptr: func() *C.btck_Txid {
				return ptr
			},
			owner: owner,
		},
	}
}

type txidApi struct {
	ptr   func() *C.btck_Txid
	owner cResourceHandle
}

func (t *txidApi) cPtr() *C.btck_Txid {
	return t.ptr()
}

// TxidLike is implemented by *Txid and *TxidView.
type TxidLike interface {
	cPtr() *C.btck_Txid
	Copy() *Txid
	Equals(TxidLike) bool
	Bytes() [32]byte
	String() string
}

var _ TxidLike = (*Txid)(nil)
var _ TxidLike = (*TxidView)(nil)

// Copy creates a copy of the txid.
func (t *txidApi) Copy() *Txid {
	txid := newTxid(t.ptr(), false)
	runtime.KeepAlive(t)
	return txid
}

// Equals checks if two txids are equal.
func (t *txidApi) Equals(other TxidLike) bool {
	equal := C.btck_txid_equals(t.ptr(), other.cPtr()) != 0
	runtime.KeepAlive(t)
	runtime.KeepAlive(other)
	return equal
}

// Bytes returns the 32-byte representation of the txid.
func (t *txidApi) Bytes() [32]byte {
	var output [32]C.uchar
	C.btck_txid_to_bytes(t.ptr(), &output[0])
	runtime.KeepAlive(t)
	return *(*[32]byte)(unsafe.Pointer(&output[0]))
}

// String returns the txid as a hex string in display order (reversed).
func (t *txidApi) String() string {
	hashBytes := t.Bytes()
	return hex.EncodeToString(ReverseBytes(hashBytes[:]))
}
