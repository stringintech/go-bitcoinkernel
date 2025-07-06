package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
#include <stdint.h>

// Bridge functions: exported Go functions that C library can call
// user_data contains the cgo.Handle ID as void* for callback identification
extern void go_notify_block_tip_bridge(void* user_data, kernel_SynchronizationState state, kernel_BlockIndex* index, double verification_progress);
extern void go_notify_header_tip_bridge(void* user_data, kernel_SynchronizationState state, int64_t height, int64_t timestamp, bool presync);
extern void go_notify_progress_bridge(void* user_data, char* title, size_t title_len, int progress_percent, bool resume_possible);
extern void go_notify_warning_set_bridge(void* user_data, kernel_Warning warning, char* message, size_t message_len);
extern void go_notify_warning_unset_bridge(void* user_data, kernel_Warning warning);
extern void go_notify_flush_error_bridge(void* user_data, char* message, size_t message_len);
extern void go_notify_fatal_error_bridge(void* user_data, char* message, size_t message_len);

// Wrapper function: C helper to set notifications with Go callbacks
// Converts Handle ID to void* and passes to C library
static inline void set_notifications_wrapper(kernel_ContextOptions* opts, uintptr_t handle) {
    kernel_NotificationInterfaceCallbacks callbacks = {
        .user_data = (void*)handle,
        .block_tip = (kernel_NotifyBlockTip)go_notify_block_tip_bridge,
        .header_tip = (kernel_NotifyHeaderTip)go_notify_header_tip_bridge,
        .progress = (kernel_NotifyProgress)go_notify_progress_bridge,
        .warning_set = (kernel_NotifyWarningSet)go_notify_warning_set_bridge,
        .warning_unset = (kernel_NotifyWarningUnset)go_notify_warning_unset_bridge,
        .flush_error = (kernel_NotifyFlushError)go_notify_flush_error_bridge,
        .fatal_error = (kernel_NotifyFatalError)go_notify_fatal_error_bridge,
    };
    kernel_context_options_set_notifications(opts, callbacks);
}
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
)

var _ cManagedResource = &ContextOptions{}

// ContextOptions wraps the C kernel_ContextOptions
type ContextOptions struct {
	ptr                *C.kernel_ContextOptions
	notificationHandle cgo.Handle // Prevents notification callbacks GC until Destroy() called
}

func NewContextOptions() (*ContextOptions, error) {
	ptr := C.kernel_context_options_create()
	if ptr == nil {
		return nil, ErrKernelContextOptionsCreate
	}

	opts := &ContextOptions{ptr: ptr}
	runtime.SetFinalizer(opts, (*ContextOptions).destroy)
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
	C.kernel_context_options_set_chainparams(opts.ptr, chainParams.ptr)
}

// SetNotifications sets the notification callbacks for these context options.
// The context created with these options will be configured with these notifications.
func (opts *ContextOptions) SetNotifications(callbacks *NotificationCallbacks) error {
	checkReady(opts)
	if callbacks == nil {
		return ErrNilNotificationCallbacks
	}

	// Create a handle for the callbacks - this prevents garbage collection
	// and provides a stable ID that can be passed through C code safely
	handle := cgo.NewHandle(callbacks)

	// Call the C wrapper function to set all notification callbacks
	C.set_notifications_wrapper(opts.ptr, C.uintptr_t(handle))

	// Store the handle to prevent GC and allow cleanup
	opts.notificationHandle = handle
	return nil
}

func (opts *ContextOptions) destroy() {
	if opts.ptr != nil {
		C.kernel_context_options_destroy(opts.ptr)
		opts.ptr = nil
	}
	if opts.notificationHandle != 0 {
		// Delete exposes notification callbacks to garbage collection
		opts.notificationHandle.Delete()
		opts.notificationHandle = 0
	}
}

func (opts *ContextOptions) Destroy() {
	runtime.SetFinalizer(opts, nil)
	opts.destroy()
}

func (opts *ContextOptions) isReady() bool {
	return opts != nil && opts.ptr != nil
}

func (opts *ContextOptions) uninitializedError() error {
	return ErrContextOptionsUninitialized
}
