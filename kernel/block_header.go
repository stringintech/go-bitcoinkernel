package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type blockHeaderCFuncs struct{}

func (blockHeaderCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_block_header_destroy((*C.btck_BlockHeader)(ptr))
}

func (blockHeaderCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_block_header_copy((*C.btck_BlockHeader)(ptr)))
}

// BlockHeader represents a Bitcoin block header containing metadata about a block.
//
// A block header is 80 bytes and contains:
//   - Version
//   - Previous block hash
//   - Merkle root
//   - Timestamp
//   - Difficulty target (nBits)
//   - Nonce
type BlockHeader struct {
	*handle
}

func newBlockHeader(ptr *C.btck_BlockHeader, fromOwned bool) *BlockHeader {
	h := newHandle(unsafe.Pointer(ptr), blockHeaderCFuncs{}, fromOwned)
	return &BlockHeader{handle: h}
}

// NewBlockHeader creates a new block header from raw serialized header data.
//
// Parameters:
//   - rawHeader: Non-nil pointer to 80 bytes of serialized block header data
//
// Returns an error if deserialization fails.
func NewBlockHeader(rawHeader *[80]byte) (*BlockHeader, error) {
	if rawHeader == nil {
		panic("rawHeader must not be nil")
	}

	ptr := C.btck_block_header_create(unsafe.Pointer(rawHeader), C.size_t(80))
	if ptr == nil {
		return nil, &InternalError{"Failed to create block header from bytes"}
	}
	return newBlockHeader(ptr, true), nil
}

// Copy creates a copy of the block header.
func (bh *BlockHeader) Copy() *BlockHeader {
	header := newBlockHeader((*C.btck_BlockHeader)(bh.ptr), false)
	runtime.KeepAlive(bh)
	return header
}

// Hash calculates and returns the hash of this block header.
func (bh *BlockHeader) Hash() *BlockHash {
	hash := newBlockHash(C.btck_block_header_get_hash((*C.btck_BlockHeader)(bh.ptr)), true)
	runtime.KeepAlive(bh)
	return hash
}

// PrevHash returns the hash of the previous block in the chain.
func (bh *BlockHeader) PrevHash() *BlockHashView {
	ptr := C.btck_block_header_get_prev_hash((*C.btck_BlockHeader)(bh.ptr))
	runtime.KeepAlive(bh)
	return newBlockHashView(check(ptr), bh.handle)
}

// Timestamp returns the block timestamp.
func (bh *BlockHeader) Timestamp() uint32 {
	timestamp := uint32(C.btck_block_header_get_timestamp((*C.btck_BlockHeader)(bh.ptr)))
	runtime.KeepAlive(bh)
	return timestamp
}

// Bits returns the nBits difficulty target.
func (bh *BlockHeader) Bits() uint32 {
	bits := uint32(C.btck_block_header_get_bits((*C.btck_BlockHeader)(bh.ptr)))
	runtime.KeepAlive(bh)
	return bits
}

// Version returns the block version.
func (bh *BlockHeader) Version() int32 {
	version := int32(C.btck_block_header_get_version((*C.btck_BlockHeader)(bh.ptr)))
	runtime.KeepAlive(bh)
	return version
}

// Nonce returns the nonce.
func (bh *BlockHeader) Nonce() uint32 {
	nonce := uint32(C.btck_block_header_get_nonce((*C.btck_BlockHeader)(bh.ptr)))
	runtime.KeepAlive(bh)
	return nonce
}

// Bytes returns the consensus serialized representation of the block header.
//
// Returns an error if the serialization fails.
func (bh *BlockHeader) Bytes() ([80]byte, error) {
	var buf [80]byte
	result := C.btck_block_header_to_bytes(
		(*C.btck_BlockHeader)(bh.ptr),
		(*C.uchar)(unsafe.Pointer(&buf[0])),
	)
	runtime.KeepAlive(bh)
	if result != 0 {
		return [80]byte{}, &SerializationError{"Failed to serialize block header"}
	}
	return buf, nil
}
