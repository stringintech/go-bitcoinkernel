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
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)
	if callbacks.OnBlockChecked != nil {
		callbacks.OnBlockChecked(newBlock(block, true), &BlockValidationState{ptr: state})
	}
}
