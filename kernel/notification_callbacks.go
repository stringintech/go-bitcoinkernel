package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// NotificationCallbacks contains all the Go callback function types for notifications.
type NotificationCallbacks struct {
	OnBlockTip     func(state SynchronizationState, index *BlockIndex, progress float64)
	OnHeaderTip    func(state SynchronizationState, height int64, timestamp int64, presync bool)
	OnProgress     func(title string, percent int, resumable bool)
	OnWarningSet   func(warning Warning, message string)
	OnWarningUnset func(warning Warning)
	OnFlushError   func(message string)
	OnFatalError   func(message string)
}

//export go_notify_block_tip_bridge
func go_notify_block_tip_bridge(user_data unsafe.Pointer, state C.kernel_SynchronizationState, index *C.kernel_BlockIndex, verification_progress C.double) {
	// Convert void* back to Handle - user_data contains Handle ID
	handle := cgo.Handle(user_data)
	// Retrieve original Go callback struct
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnBlockTip != nil {
		goState := SynchronizationState(state)
		// Note: BlockIndex from notification is const and owned by kernel library
		// We create a wrapper but don't set finalizer since we don't own it
		goIndex := &BlockIndex{ptr: (*C.kernel_BlockIndex)(unsafe.Pointer(index))}
		callbacks.OnBlockTip(goState, goIndex, float64(verification_progress))
	}
}

//export go_notify_header_tip_bridge
func go_notify_header_tip_bridge(user_data unsafe.Pointer, state C.kernel_SynchronizationState, height C.int64_t, timestamp C.int64_t, presync C.bool) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnHeaderTip != nil {
		goState := SynchronizationState(state)
		callbacks.OnHeaderTip(goState, int64(height), int64(timestamp), bool(presync))
	}
}

//export go_notify_progress_bridge
func go_notify_progress_bridge(user_data unsafe.Pointer, title *C.char, title_len C.size_t, progress_percent C.int, resume_possible C.bool) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnProgress != nil {
		goTitle := C.GoStringN(title, C.int(title_len))
		callbacks.OnProgress(goTitle, int(progress_percent), bool(resume_possible))
	}
}

//export go_notify_warning_set_bridge
func go_notify_warning_set_bridge(user_data unsafe.Pointer, warning C.kernel_Warning, message *C.char, message_len C.size_t) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnWarningSet != nil {
		goWarning := Warning(warning)
		goMessage := C.GoStringN(message, C.int(message_len))
		callbacks.OnWarningSet(goWarning, goMessage)
	}
}

//export go_notify_warning_unset_bridge
func go_notify_warning_unset_bridge(user_data unsafe.Pointer, warning C.kernel_Warning) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnWarningUnset != nil {
		goWarning := Warning(warning)
		callbacks.OnWarningUnset(goWarning)
	}
}

//export go_notify_flush_error_bridge
func go_notify_flush_error_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnFlushError != nil {
		goMessage := C.GoStringN(message, C.int(message_len))
		callbacks.OnFlushError(goMessage)
	}
}

//export go_notify_fatal_error_bridge
func go_notify_fatal_error_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	handle := cgo.Handle(user_data)
	callbacks := handle.Value().(*NotificationCallbacks)

	if callbacks.OnFatalError != nil {
		goMessage := C.GoStringN(message, C.int(message_len))
		callbacks.OnFatalError(goMessage)
	}
}

// SynchronizationState represents the current sync state passed to tip changed callbacks
type SynchronizationState int

const (
	SyncStateInitReindex SynchronizationState = iota
	SyncStateInitDownload
	SyncStatePostInit
)

// Warning represents possible warning types issued by validation
type Warning int

const (
	WarningUnknownNewRulesActivated Warning = iota
	WarningLargeWorkInvalidChain
)
