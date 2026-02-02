package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import "unsafe"

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
//   - rawHeader: Serialized block header data
//
// Returns an error if the header data is malformed, invalid length, or cannot be parsed.
func NewBlockHeader(rawHeader []byte) (*BlockHeader, error) {
	ptr := C.btck_block_header_create(unsafe.Pointer(unsafe.SliceData(rawHeader)), C.size_t(len(rawHeader)))
	if ptr == nil {
		return nil, &InternalError{"Failed to create block header from bytes"}
	}
	return newBlockHeader(ptr, true), nil
}

// Copy creates a copy of the block header.
func (bh *BlockHeader) Copy() *BlockHeader {
	return newBlockHeader((*C.btck_BlockHeader)(bh.ptr), false)
}

// Hash calculates and returns the hash of this block header.
func (bh *BlockHeader) Hash() *BlockHash {
	return newBlockHash(C.btck_block_header_get_hash((*C.btck_BlockHeader)(bh.ptr)), true)
}

// PrevHash returns the hash of the previous block in the chain.
func (bh *BlockHeader) PrevHash() *BlockHashView {
	ptr := C.btck_block_header_get_prev_hash((*C.btck_BlockHeader)(bh.ptr))
	return newBlockHashView(check(ptr))
}

// Timestamp returns the block timestamp.
func (bh *BlockHeader) Timestamp() uint32 {
	return uint32(C.btck_block_header_get_timestamp((*C.btck_BlockHeader)(bh.ptr)))
}

// Bits returns the nBits difficulty target.
func (bh *BlockHeader) Bits() uint32 {
	return uint32(C.btck_block_header_get_bits((*C.btck_BlockHeader)(bh.ptr)))
}

// Version returns the block version.
func (bh *BlockHeader) Version() int32 {
	return int32(C.btck_block_header_get_version((*C.btck_BlockHeader)(bh.ptr)))
}

// Nonce returns the nonce.
func (bh *BlockHeader) Nonce() uint32 {
	return uint32(C.btck_block_header_get_nonce((*C.btck_BlockHeader)(bh.ptr)))
}
