package kernel

/*
#include "bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

// NotificationCallbacks is implemented by types that want to receive kernel notification events.
type NotificationCallbacks interface {
	BlockTip(state SynchronizationState, entry *BlockTreeEntry, progress float64)
	HeaderTip(state SynchronizationState, height int64, timestamp int64, presync bool)
	Progress(title string, percent int, resumable bool)
	WarningSet(warning Warning, message string)
	WarningUnset(warning Warning)
	FlushError(message string)
	FatalError(message string)
}

// SynchronizationState represents the current sync state passed to tip changed callbacks.
type SynchronizationState C.btck_SynchronizationState

const (
	SyncStateInitReindex  SynchronizationState = C.btck_SynchronizationState_INIT_REINDEX
	SyncStateInitDownload SynchronizationState = C.btck_SynchronizationState_INIT_DOWNLOAD
	SyncStatePostInit     SynchronizationState = C.btck_SynchronizationState_POST_INIT
)

// Warning represents possible warning types issued by validation.
type Warning C.btck_Warning

const (
	WarningUnknownNewRulesActivated Warning = C.btck_Warning_UNKNOWN_NEW_RULES_ACTIVATED
	WarningLargeWorkInvalidChain    Warning = C.btck_Warning_LARGE_WORK_INVALID_CHAIN
)

//export go_notify_block_tip_bridge
func go_notify_block_tip_bridge(user_data unsafe.Pointer, state C.btck_SynchronizationState, entry *C.btck_BlockTreeEntry, verification_progress C.double) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.BlockTip(SynchronizationState(state), &BlockTreeEntry{ptr: (*C.btck_BlockTreeEntry)(unsafe.Pointer(entry))}, float64(verification_progress))
}

//export go_notify_header_tip_bridge
func go_notify_header_tip_bridge(user_data unsafe.Pointer, state C.btck_SynchronizationState, height C.int64_t, timestamp C.int64_t, presync C.int) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.HeaderTip(SynchronizationState(state), int64(height), int64(timestamp), presync != 0)
}

//export go_notify_progress_bridge
func go_notify_progress_bridge(user_data unsafe.Pointer, title *C.char, title_len C.size_t, progress_percent C.int, resume_possible C.int) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.Progress(C.GoStringN(title, C.int(title_len)), int(progress_percent), resume_possible != 0)
}

//export go_notify_warning_set_bridge
func go_notify_warning_set_bridge(user_data unsafe.Pointer, warning C.btck_Warning, message *C.char, message_len C.size_t) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.WarningSet(Warning(warning), C.GoStringN(message, C.int(message_len)))
}

//export go_notify_warning_unset_bridge
func go_notify_warning_unset_bridge(user_data unsafe.Pointer, warning C.btck_Warning) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.WarningUnset(Warning(warning))
}

//export go_notify_flush_error_bridge
func go_notify_flush_error_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.FlushError(C.GoStringN(message, C.int(message_len)))
}

//export go_notify_fatal_error_bridge
func go_notify_fatal_error_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	callbacks := cgoHandleFromPointer(user_data).Value().(NotificationCallbacks)
	callbacks.FatalError(C.GoStringN(message, C.int(message_len)))
}
