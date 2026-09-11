package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type transactionOutPointCFuncs struct{}

func (transactionOutPointCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_transaction_out_point_destroy((*C.btck_TransactionOutPoint)(ptr))
}

func (transactionOutPointCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_transaction_out_point_copy((*C.btck_TransactionOutPoint)(ptr)))
}

// TransactionOutPoint holds the txid and output index it is pointing to.
type TransactionOutPoint struct {
	*handle
	transactionOutPointApi
}

func newTransactionOutPoint(ptr *C.btck_TransactionOutPoint, fromOwned bool) *TransactionOutPoint {
	h := newHandle(unsafe.Pointer(ptr), transactionOutPointCFuncs{}, fromOwned)
	return &TransactionOutPoint{
		handle: h,
		transactionOutPointApi: transactionOutPointApi{
			ptr: func() *C.btck_TransactionOutPoint {
				return (*C.btck_TransactionOutPoint)(h.ptr)
			},
			owner: h,
		},
	}
}

// TransactionOutPointView holds the txid and output index it is pointing to.
type TransactionOutPointView struct {
	transactionOutPointApi
}

func newTransactionOutPointView(ptr *C.btck_TransactionOutPoint, owner cResourceHandle) *TransactionOutPointView {
	return &TransactionOutPointView{
		transactionOutPointApi: transactionOutPointApi{
			ptr: func() *C.btck_TransactionOutPoint {
				return ptr
			},
			owner: owner,
		},
	}
}

type transactionOutPointApi struct {
	ptr   func() *C.btck_TransactionOutPoint
	owner cResourceHandle
}

func (t *transactionOutPointApi) cPtr() *C.btck_TransactionOutPoint {
	return t.ptr()
}

// TransactionOutPointLike is implemented by *TransactionOutPoint and *TransactionOutPointView.
type TransactionOutPointLike interface {
	cPtr() *C.btck_TransactionOutPoint
	Copy() *TransactionOutPoint
	GetIndex() uint32
	GetTxid() *TxidView
}

var _ TransactionOutPointLike = (*TransactionOutPoint)(nil)
var _ TransactionOutPointLike = (*TransactionOutPointView)(nil)

// Copy creates a copy of the transaction out point.
func (t *transactionOutPointApi) Copy() *TransactionOutPoint {
	outPoint := newTransactionOutPoint(t.ptr(), false)
	runtime.KeepAlive(t)
	return outPoint
}

// GetIndex returns the output position from the transaction out point.
func (t *transactionOutPointApi) GetIndex() uint32 {
	index := uint32(C.btck_transaction_out_point_get_index(t.ptr()))
	runtime.KeepAlive(t)
	return index
}

// GetTxid returns the txid from the out point.
func (t *transactionOutPointApi) GetTxid() *TxidView {
	ptr := C.btck_transaction_out_point_get_txid(t.ptr())
	runtime.KeepAlive(t)
	return newTxidView(check(ptr), t.owner)
}
