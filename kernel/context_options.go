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
extern void go_validation_interface_block_checked_bridge(void* user_data, const btck_BlockPointer* block, const btck_BlockValidationState* state);
*/
import "C"
import (
	"fmt"
	"runtime"
	"runtime/cgo"
	"unsafe"
)

var _ cResource = &ContextOptions{}

// ContextOptions wraps the C btck_ContextOptions.
// Once the options is set on a context, the context is responsible for its lifetime; otherwise it is garbage collected
type ContextOptions struct {
	ptr                *C.btck_ContextOptions
	notificationHandle cgo.Handle
	validationHandle   cgo.Handle
}

func NewContextOptions() (*ContextOptions, error) {
	ptr := C.btck_context_options_create()
	if ptr == nil {
		return nil, ErrKernelContextOptionsCreate
	}

	opts := &ContextOptions{ptr: ptr}
	runtime.SetFinalizer(opts, (*ContextOptions).finalize)
	return opts, nil
}

// SetChainParams sets the chain parameters for these context options.
// Kernel makes a copy of the chain parameters, so the caller can
// safely free the chainParams object after this call returns.
func (opts *ContextOptions) SetChainParams(chainParams *ChainParameters) {
	checkReady(opts)
	if chainParams == nil || chainParams.ptr == nil {
		panic(ErrChainParametersUninitialized)
	}
	C.btck_context_options_set_chainparams(opts.ptr, chainParams.ptr)
}

// SetNotifications sets the notification callbacks for these context options.
// The context created with these options will be configured with these notifications.
func (opts *ContextOptions) SetNotifications(callbacks *NotificationCallbacks) error {
	checkReady(opts)
	if callbacks == nil {
		return fmt.Errorf("nil notification callbacks")
	}
	if opts.notificationHandle != 0 {
		return fmt.Errorf("notification callbacks already set")
	}

	// Create a handle for the callbacks - this prevents garbage collection
	// and provides a stable ID that can be passed through C code safely
	opts.notificationHandle = cgo.NewHandle(callbacks)

	// Create notification callbacks struct and call C library directly
	notificationCallbacks := C.btck_NotificationInterfaceCallbacks{
		user_data:         unsafe.Pointer(opts.notificationHandle),
		user_data_destroy: nil, // Go handles memory management via cgo.Handle
		block_tip:         C.btck_NotifyBlockTip(C.go_notify_block_tip_bridge),
		header_tip:        C.btck_NotifyHeaderTip(C.go_notify_header_tip_bridge),
		progress:          C.btck_NotifyProgress(C.go_notify_progress_bridge),
		warning_set:       C.btck_NotifyWarningSet(C.go_notify_warning_set_bridge),
		warning_unset:     C.btck_NotifyWarningUnset(C.go_notify_warning_unset_bridge),
		flush_error:       C.btck_NotifyFlushError(C.go_notify_flush_error_bridge),
		fatal_error:       C.btck_NotifyFatalError(C.go_notify_fatal_error_bridge),
	}
	C.btck_context_options_set_notifications(opts.ptr, notificationCallbacks)
	return nil
}

// SetValidationInterface sets the validation interface callbacks for these context options.
// The context created with these options will be configured with these validation callbacks.
func (opts *ContextOptions) SetValidationInterface(callbacks *ValidationInterfaceCallbacks) error {
	checkReady(opts)
	if callbacks == nil {
		return fmt.Errorf("nil validation interface callbacks")
	}
	if opts.validationHandle != 0 {
		return fmt.Errorf("validation interface callbacks already set")
	}

	// Create a handle for the callbacks - this prevents garbage collection
	// and provides a stable ID that can be passed through C code safely
	opts.validationHandle = cgo.NewHandle(callbacks)

	// Create validation callbacks struct and call C library directly
	validationCallbacks := C.btck_ValidationInterfaceCallbacks{
		user_data:         unsafe.Pointer(opts.validationHandle),
		user_data_destroy: nil, // Go handles memory management via cgo.Handle
		block_checked:     C.btck_ValidationInterfaceBlockChecked(C.go_validation_interface_block_checked_bridge),
	}
	C.btck_context_options_set_validation_interface(opts.ptr, validationCallbacks)
	return nil
}

func (opts *ContextOptions) destroy() {
	opts.finalize()
	runtime.SetFinalizer(opts, nil)
}

func (opts *ContextOptions) finalize() {
	if opts.ptr != nil {
		C.btck_context_options_destroy(opts.ptr)
		opts.ptr = nil
	}
	if opts.notificationHandle != 0 {
		// Delete exposes notification callbacks to garbage collection
		opts.notificationHandle.Delete()
		opts.notificationHandle = 0
	}
	if opts.validationHandle != 0 {
		// Delete exposes validation callbacks to garbage collection
		opts.validationHandle.Delete()
		opts.validationHandle = 0
	}
}

func (opts *ContextOptions) isReady() bool {
	return opts != nil && opts.ptr != nil
}

func (opts *ContextOptions) uninitializedError() error {
	return ErrContextOptionsUninitialized
}
