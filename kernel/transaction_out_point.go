package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
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
	return &TransactionOutPoint{handle: h, transactionOutPointApi: transactionOutPointApi{(*C.btck_TransactionOutPoint)(h.ptr)}}
}

// TransactionOutPointView holds the txid and output index it is pointing to.
type TransactionOutPointView struct {
	transactionOutPointApi
	ptr *C.btck_TransactionOutPoint
}

func newTransactionOutPointView(ptr *C.btck_TransactionOutPoint) *TransactionOutPointView {
	return &TransactionOutPointView{
		transactionOutPointApi: transactionOutPointApi{ptr},
		ptr:                    ptr,
	}
}

type transactionOutPointApi struct {
	ptr *C.btck_TransactionOutPoint
}

// Copy creates a copy of the transaction out point.
func (t *transactionOutPointApi) Copy() *TransactionOutPoint {
	return newTransactionOutPoint(t.ptr, false)
}

// GetIndex returns the output position from the transaction out point.
func (t *transactionOutPointApi) GetIndex() uint32 {
	return uint32(C.btck_transaction_out_point_get_index(t.ptr))
}

// GetTxid returns the txid from the out point.
func (t *transactionOutPointApi) GetTxid() *TxidView {
	ptr := C.btck_transaction_out_point_get_txid(t.ptr)
	return newTxidView(check(ptr))
}
