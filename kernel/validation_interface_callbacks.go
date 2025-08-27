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
	OnBlockChecked func(block *Block, state *BlockValidationState)
}

//export go_validation_interface_block_checked_bridge
func go_validation_interface_block_checked_bridge(user_data unsafe.Pointer, block *C.btck_Block, state *C.btck_BlockValidationState) {
	// Convert void* back to Handle - user_data contains Handle ID
	handle := cgo.Handle(user_data)
	// Retrieve original Go callback struct
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)

	if callbacks.OnBlockChecked != nil {
		callbacks.OnBlockChecked(newBlockFromPtr(block), &BlockValidationState{ptr: state})
	}
}
