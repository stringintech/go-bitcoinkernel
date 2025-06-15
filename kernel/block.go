package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// BlockHash represents a Bitcoin block hash
type BlockHash struct {
	ptr *C.kernel_BlockHash
}

// Block wraps the C kernel_Block
type Block struct {
	ptr *C.kernel_Block
}

// BlockIndex wraps the C kernel_BlockIndex
type BlockIndex struct {
	ptr *C.kernel_BlockIndex
}

// ValidationMode represents the validation state mode
type ValidationMode int

const (
	ValidationStateValid ValidationMode = iota
	ValidationStateInvalid
	ValidationStateError
)

// BlockValidationResult represents the validation result for a block
type BlockValidationResult int

const (
	BlockResultUnset BlockValidationResult = iota
	BlockConsensus
	BlockCachedInvalid
	BlockInvalidHeader
	BlockMutated
	BlockMissingPrev
	BlockInvalidPrev
	BlockTimeFuture
	BlockHeaderLowWork
)

// BlockValidationState wraps the C kernel_BlockValidationState
type BlockValidationState struct {
	ptr *C.kernel_BlockValidationState
}

// NewBlockFromRaw creates a new block from raw serialized data
func NewBlockFromRaw(rawBlock []byte) (*Block, error) {
	if len(rawBlock) == 0 {
		return nil, ErrInvalidBlockData
	}

	ptr := C.kernel_block_create((*C.uchar)(unsafe.Pointer(&rawBlock[0])), C.size_t(len(rawBlock)))
	if ptr == nil {
		return nil, ErrBlockCreation
	}

	block := &Block{ptr: ptr}
	runtime.SetFinalizer(block, (*Block).destroy)
	return block, nil
}

// Hash returns the hash of the block
func (b *Block) Hash() (*BlockHash, error) {
	if b.ptr == nil {
		return nil, ErrInvalidBlock
	}

	ptr := C.kernel_block_get_hash(b.ptr)
	if ptr == nil {
		return nil, ErrHashCalculation
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

// Data returns the serialized block data
func (b *Block) Data() ([]byte, error) {
	if b.ptr == nil {
		return nil, ErrInvalidBlock
	}

	byteArray := C.kernel_copy_block_data(b.ptr)
	if byteArray == nil {
		return nil, ErrBlockDataCopy
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

// destroy deallocates the block
func (b *Block) destroy() {
	if b.ptr != nil {
		C.kernel_block_destroy(b.ptr)
		b.ptr = nil
	}
}

// Close explicitly destroys the block and removes the finalizer
func (b *Block) Close() {
	runtime.SetFinalizer(b, nil)
	b.destroy()
}

// Height returns the height of the block index
func (bi *BlockIndex) Height() int32 {
	if bi.ptr == nil {
		return -1
	}
	return int32(C.kernel_block_index_get_height(bi.ptr))
}

// Hash returns the block hash associated with this block index
func (bi *BlockIndex) Hash() (*BlockHash, error) {
	if bi.ptr == nil {
		return nil, ErrInvalidBlockIndex
	}

	ptr := C.kernel_block_index_get_block_hash(bi.ptr)
	if ptr == nil {
		return nil, ErrHashCalculation
	}

	hash := &BlockHash{ptr: ptr}
	runtime.SetFinalizer(hash, (*BlockHash).destroy)
	return hash, nil
}

// Previous returns the previous block index in the chain
func (bi *BlockIndex) Previous() *BlockIndex {
	if bi.ptr == nil {
		return nil
	}

	ptr := C.kernel_get_previous_block_index(bi.ptr)
	if ptr == nil {
		return nil
	}

	prevIndex := &BlockIndex{ptr: ptr}
	runtime.SetFinalizer(prevIndex, (*BlockIndex).destroy)
	return prevIndex
}

// destroy deallocates the block index
func (bi *BlockIndex) destroy() {
	if bi.ptr != nil {
		C.kernel_block_index_destroy(bi.ptr)
		bi.ptr = nil
	}
}

// Close explicitly destroys the block index and removes the finalizer
func (bi *BlockIndex) Close() {
	runtime.SetFinalizer(bi, nil)
	bi.destroy()
}

// destroy deallocates the block hash
func (bh *BlockHash) destroy() {
	if bh.ptr != nil {
		C.kernel_block_hash_destroy(bh.ptr)
		bh.ptr = nil
	}
}

// Close explicitly destroys the block hash and removes the finalizer
func (bh *BlockHash) Close() {
	runtime.SetFinalizer(bh, nil)
	bh.destroy()
}

// Bytes returns the raw hash bytes
func (bh *BlockHash) Bytes() []byte {
	if bh.ptr == nil {
		return nil
	}
	// BlockHash is a 32-byte array in the C struct
	return C.GoBytes(unsafe.Pointer(&bh.ptr.hash[0]), 32)
}

// ValidationMode returns the validation mode from the block validation state
func (bvs *BlockValidationState) ValidationMode() ValidationMode {
	if bvs.ptr == nil {
		return ValidationStateError
	}
	mode := C.kernel_get_validation_mode_from_block_validation_state(bvs.ptr)
	return ValidationMode(mode)
}

// ValidationResult returns the validation result from the block validation state
func (bvs *BlockValidationState) ValidationResult() BlockValidationResult {
	if bvs.ptr == nil {
		return BlockResultUnset
	}
	result := C.kernel_get_block_validation_result_from_block_validation_state(bvs.ptr)
	return BlockValidationResult(result)
}
