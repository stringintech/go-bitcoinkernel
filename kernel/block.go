package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ cManagedResource = &Block{}

// Block wraps the C kernel_Block
type Block struct {
	ptr *C.kernel_Block
}

// NewBlockFromRaw creates a new block from raw serialized data
func NewBlockFromRaw(rawBlock []byte) (*Block, error) {
	if len(rawBlock) == 0 {
		return nil, ErrEmptyBlockData
	}
	ptr := C.kernel_block_create((*C.uchar)(unsafe.Pointer(&rawBlock[0])), C.size_t(len(rawBlock)))
	if ptr == nil {
		return nil, ErrKernelBlockCreate
	}
	return newBlockFromPtr(ptr), nil
}

func newBlockFromPtr(ptr *C.kernel_Block) *Block {
	block := &Block{ptr: ptr}
	runtime.SetFinalizer(block, (*Block).destroy)
	return block
}

func (b *Block) Hash() (*BlockHash, error) {
	checkReady(b)

	ptr := C.kernel_block_get_hash(b.ptr)
	if ptr == nil {
		return nil, ErrKernelBlockGetHash
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

// Data returns the serialized block data
func (b *Block) Data() ([]byte, error) {
	checkReady(b)

	byteArray := C.kernel_copy_block_data(b.ptr)
	if byteArray == nil {
		return nil, ErrKernelCopyBlockData
	}
	defer C.kernel_byte_array_destroy(byteArray)

	size := int(byteArray.size)
	if size == 0 {
		return nil, nil
	}

	// Copy the data to Go slice
	data := C.GoBytes(unsafe.Pointer(byteArray.data), C.int(size))
	return data, nil
}

func (b *Block) destroy() {
	if b.ptr != nil {
		C.kernel_block_destroy(b.ptr)
		b.ptr = nil
	}
}

func (b *Block) Destroy() {
	runtime.SetFinalizer(b, nil)
	b.destroy()
}

func (b *Block) isReady() bool {
	return b != nil && b.ptr != nil
}

func (b *Block) uninitializedError() error {
	return ErrBlockUninitialized
}
