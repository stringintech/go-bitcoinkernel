package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// ValidationInterfaceCallbacks holds the validation interface callbacks.
//
// Note that these callbacks block any further validation execution when they are called.
type ValidationInterfaceCallbacks struct {
	OnBlockChecked      func(block *Block, state *BlockValidationState) // Called when a new block has been fully validated. Contains the result of its validation.
	OnPoWValidBlock     func(block *Block, entry *BlockTreeEntry)       // Called when a new block extends the header chain and has a valid transaction and segwit merkle root.
	OnBlockConnected    func(block *Block, entry *BlockTreeEntry)       // Called when a block is valid and has now been connected to the best chain.
	OnBlockDisconnected func(block *Block, entry *BlockTreeEntry)       // Called during a re-org when a block has been removed from the best chain.
}

//export go_validation_interface_block_checked_bridge
func go_validation_interface_block_checked_bridge(user_data unsafe.Pointer, block *C.btck_Block, state *C.btck_BlockValidationState) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)
	if callbacks.OnBlockChecked != nil {
		callbacks.OnBlockChecked(newBlock(block, true), &BlockValidationState{ptr: state})
	}
}

//export go_validation_interface_pow_valid_block_bridge
func go_validation_interface_pow_valid_block_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)
	if callbacks.OnPoWValidBlock != nil {
		callbacks.OnPoWValidBlock(newBlock(block, true), &BlockTreeEntry{ptr: entry})
	}
}

//export go_validation_interface_block_connected_bridge
func go_validation_interface_block_connected_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)
	if callbacks.OnBlockConnected != nil {
		callbacks.OnBlockConnected(newBlock(block, true), &BlockTreeEntry{ptr: entry})
	}
}

//export go_validation_interface_block_disconnected_bridge
func go_validation_interface_block_disconnected_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*ValidationInterfaceCallbacks)
	if callbacks.OnBlockDisconnected != nil {
		callbacks.OnBlockDisconnected(newBlock(block, true), &BlockTreeEntry{ptr: entry})
	}
}
