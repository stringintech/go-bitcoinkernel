package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
#include <stdint.h>

// Bridge functions: exported Go functions that C library can call
// user_data contains the cgo.Handle ID as void* for callback identification
extern void go_notify_block_tip_bridge(void* user_data, btck_SynchronizationState state, btck_BlockTreeEntry* entry, double verification_progress);
extern void go_notify_header_tip_bridge(void* user_data, btck_SynchronizationState state, int64_t height, int64_t timestamp, int presync);
extern void go_notify_progress_bridge(void* user_data, const char* title, size_t title_len, int progress_percent, int resume_possible);
extern void go_notify_warning_set_bridge(void* user_data, btck_Warning warning, const char* message, size_t message_len);
extern void go_notify_warning_unset_bridge(void* user_data, btck_Warning warning);
extern void go_notify_flush_error_bridge(void* user_data, const char* message, size_t message_len);
extern void go_notify_fatal_error_bridge(void* user_data, const char* message, size_t message_len);
extern void go_validation_interface_block_checked_bridge(void* user_data, btck_Block* block, const btck_BlockValidationState* state);

extern void go_delete_handle(void* user_data);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

type contextOptionsCFuncs struct{}

func (contextOptionsCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_context_options_destroy((*C.btck_ContextOptions)(ptr))
}

type ContextOptions struct {
	*uniqueHandle
}

func newContextOptions(ptr *C.btck_ContextOptions) *ContextOptions {
	h := newUniqueHandle(unsafe.Pointer(ptr), contextOptionsCFuncs{})
	return &ContextOptions{uniqueHandle: h}
}

func NewContextOptions() (*ContextOptions, error) {
	ptr := C.btck_context_options_create()
	return newContextOptions(check(ptr)), nil
}

// SetChainParams sets the chain parameters for context options.
// Caller can safely free the chainParams object after this call returns.
func (opts *ContextOptions) SetChainParams(chainParams *ChainParameters) {
	C.btck_context_options_set_chainparams((*C.btck_ContextOptions)(opts.ptr), (*C.btck_ChainParameters)(chainParams.ptr))
}

func (opts *ContextOptions) SetNotifications(callbacks *NotificationCallbacks) error {
	notificationCallbacks := C.btck_NotificationInterfaceCallbacks{
		user_data:         unsafe.Pointer(cgo.NewHandle(callbacks)),
		user_data_destroy: C.btck_DestroyCallback(C.go_delete_handle),
		block_tip:         C.btck_NotifyBlockTip(C.go_notify_block_tip_bridge),
		header_tip:        C.btck_NotifyHeaderTip(C.go_notify_header_tip_bridge),
		progress:          C.btck_NotifyProgress(C.go_notify_progress_bridge),
		warning_set:       C.btck_NotifyWarningSet(C.go_notify_warning_set_bridge),
		warning_unset:     C.btck_NotifyWarningUnset(C.go_notify_warning_unset_bridge),
		flush_error:       C.btck_NotifyFlushError(C.go_notify_flush_error_bridge),
		fatal_error:       C.btck_NotifyFatalError(C.go_notify_fatal_error_bridge),
	}
	C.btck_context_options_set_notifications((*C.btck_ContextOptions)(opts.ptr), notificationCallbacks)
	return nil
}

func (opts *ContextOptions) SetValidationInterface(callbacks *ValidationInterfaceCallbacks) error {
	validationCallbacks := C.btck_ValidationInterfaceCallbacks{
		user_data:         unsafe.Pointer(cgo.NewHandle(callbacks)),
		user_data_destroy: C.btck_DestroyCallback(C.go_delete_handle),
		block_checked:     C.btck_ValidationInterfaceBlockChecked(C.go_validation_interface_block_checked_bridge),
	}
	C.btck_context_options_set_validation_interface((*C.btck_ContextOptions)(opts.ptr), validationCallbacks)
	return nil
}
