package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

// ValidationInterfaceCallbacks is implemented by types that want to receive kernel validation events.
//
// Note that these callbacks block any further validation execution when they are called.
type ValidationInterfaceCallbacks interface {
	BlockChecked(block *Block, state *BlockValidationStateView)
	PoWValidBlock(block *Block, entry *BlockTreeEntry)
	BlockConnected(block *Block, entry *BlockTreeEntry)
	BlockDisconnected(block *Block, entry *BlockTreeEntry)
}

//export go_validation_interface_block_checked_bridge
func go_validation_interface_block_checked_bridge(user_data unsafe.Pointer, block *C.btck_Block, state *C.btck_BlockValidationState) {
	callbacks := cgoHandleFromPointer(user_data).Value().(ValidationInterfaceCallbacks)
	callbacks.BlockChecked(newBlock(block, true), newBlockValidationStateView(state))
}

//export go_validation_interface_pow_valid_block_bridge
func go_validation_interface_pow_valid_block_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	callbacks := cgoHandleFromPointer(user_data).Value().(ValidationInterfaceCallbacks)
	callbacks.PoWValidBlock(newBlock(block, true), &BlockTreeEntry{ptr: entry})
}

//export go_validation_interface_block_connected_bridge
func go_validation_interface_block_connected_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	callbacks := cgoHandleFromPointer(user_data).Value().(ValidationInterfaceCallbacks)
	callbacks.BlockConnected(newBlock(block, true), &BlockTreeEntry{ptr: entry})
}

//export go_validation_interface_block_disconnected_bridge
func go_validation_interface_block_disconnected_bridge(user_data unsafe.Pointer, block *C.btck_Block, entry *C.btck_BlockTreeEntry) {
	callbacks := cgoHandleFromPointer(user_data).Value().(ValidationInterfaceCallbacks)
	callbacks.BlockDisconnected(newBlock(block, true), &BlockTreeEntry{ptr: entry})
}
