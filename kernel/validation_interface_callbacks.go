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
func go_validation_interface_block_checked_bridge(user_data unsafe.Pointer, block *C.btck_BlockPointer, state *C.btck_BlockValidationState) {
	// Convert void* back to Handle - user_data contains Handle ID
	handle := cgo.Handle(user_data)
	// Retrieve original Go callback struct
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)

	if callbacks.OnBlockChecked != nil {
		goBlock := &BlockPointer{ptr: (*C.btck_BlockPointer)(unsafe.Pointer(block))}
		goState := &BlockValidationState{ptr: (*C.btck_BlockValidationState)(unsafe.Pointer(state))}
		callbacks.OnBlockChecked(goBlock, goState)
	}
}

// BlockPointer wraps the C kernel_BlockPointer for validation interface callbacks
type BlockPointer struct {
	ptr *C.btck_BlockPointer
}

// GetHash returns the block hash
func (bp *BlockPointer) GetHash() (*BlockHash, error) {
	if bp.ptr == nil {
		return nil, ErrBlockUninitialized
	}

	hashPtr := C.btck_block_pointer_get_hash(bp.ptr)
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

	return writeToBytes(func(writer C.btck_WriteBytes, user_data unsafe.Pointer) C.int {
		return C.btck_block_pointer_to_bytes(bp.ptr, writer, user_data)
	})
}
