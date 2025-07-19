package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// ValidationInterfaceCallbacks contains all the Go callback function types for validation interface.
type ValidationInterfaceCallbacks struct {
	OnBlockChecked func(block *BlockPointer, state *BlockValidationState)
}

//export go_validation_interface_block_checked_bridge
func go_validation_interface_block_checked_bridge(user_data unsafe.Pointer, block *C.kernel_BlockPointer, state *C.kernel_BlockValidationState) {
	// Convert void* back to Handle - user_data contains Handle ID
	handle := cgo.Handle(user_data)
	// Retrieve original Go callback struct
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)

	if callbacks.OnBlockChecked != nil {
		// Note: BlockPointer and BlockValidationState from validation interface are const and owned by kernel library
		// We create wrappers but don't set finalizer since we don't own them
		goBlock := &BlockPointer{ptr: (*C.kernel_BlockPointer)(unsafe.Pointer(block))}
		goState := &BlockValidationState{ptr: (*C.kernel_BlockValidationState)(unsafe.Pointer(state))}
		callbacks.OnBlockChecked(goBlock, goState)
	}
}

// BlockPointer wraps the C kernel_BlockPointer for validation interface callbacks
type BlockPointer struct {
	ptr *C.kernel_BlockPointer
}

// GetHash returns the block hash
func (bp *BlockPointer) GetHash() (*BlockHash, error) {
	if bp.ptr == nil {
		return nil, ErrBlockUninitialized
	}

	hashPtr := C.kernel_block_pointer_get_hash(bp.ptr)
	if hashPtr == nil {
		return nil, ErrKernelBlockGetHash
	}

	return &BlockHash{ptr: hashPtr}, nil
}

// CopyData copies the block data into a byte array
func (bp *BlockPointer) CopyData() ([]byte, error) {
	if bp.ptr == nil {
		return nil, ErrBlockUninitialized
	}

	byteArray := C.kernel_block_pointer_copy_data(bp.ptr)
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
