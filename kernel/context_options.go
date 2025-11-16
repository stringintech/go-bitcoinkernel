package kernel

/*
#include "bitcoinkernel.h"
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
extern void go_validation_interface_pow_valid_block_bridge(void* user_data, const btck_BlockTreeEntry* entry, btck_Block* block);
extern void go_validation_interface_block_connected_bridge(void* user_data, btck_Block* block, const btck_BlockTreeEntry* entry);
extern void go_validation_interface_block_disconnected_bridge(void* user_data, btck_Block* block, const btck_BlockTreeEntry* entry);

extern void go_delete_handle(void* user_data);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// ContextOption is a functional option for configuring context options.
type ContextOption func(*C.btck_ContextOptions) error

// WithChainType returns a ContextOption that sets the chain parameters for the context.
// The context will be configured for these chain parameters.
//
// Parameters:
//   - chainType: The type of chain (ChainTypeMainnet, ChainTypeTestnet, ChainTypeRegtest, etc.)
func WithChainType(chainType ChainType) ContextOption {
	return func(opts *C.btck_ContextOptions) error {
		chainParams, err := NewChainParameters(chainType)
		if err != nil {
			return err
		}
		defer chainParams.Destroy()
		C.btck_context_options_set_chainparams(opts, (*C.btck_ChainParameters)(chainParams.ptr))
		return nil
	}
}

// WithNotifications returns a ContextOption that sets the kernel notifications for the context.
// The context will be configured with these notifications.
//
// Parameters:
//   - callbacks: Notification callbacks to set
func WithNotifications(callbacks *NotificationCallbacks) ContextOption {
	return func(opts *C.btck_ContextOptions) error {
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
		C.btck_context_options_set_notifications(opts, notificationCallbacks)
		return nil
	}
}

// WithValidationInterface returns a ContextOption that sets the validation interface callbacks.
// The callbacks will be triggered from validation events issued by the chainstate manager
// created from the same context.
//
// Parameters:
//   - callbacks: The callbacks used for passing validation information to the user
func WithValidationInterface(callbacks *ValidationInterfaceCallbacks) ContextOption {
	return func(opts *C.btck_ContextOptions) error {
		validationCallbacks := C.btck_ValidationInterfaceCallbacks{
			user_data:          unsafe.Pointer(cgo.NewHandle(callbacks)),
			user_data_destroy:  C.btck_DestroyCallback(C.go_delete_handle),
			block_checked:      C.btck_ValidationInterfaceBlockChecked(C.go_validation_interface_block_checked_bridge),
			pow_valid_block:    C.btck_ValidationInterfacePoWValidBlock(C.go_validation_interface_pow_valid_block_bridge),
			block_connected:    C.btck_ValidationInterfaceBlockConnected(C.go_validation_interface_block_connected_bridge),
			block_disconnected: C.btck_ValidationInterfaceBlockDisconnected(C.go_validation_interface_block_disconnected_bridge),
		}
		C.btck_context_options_set_validation_interface(opts, validationCallbacks)
		return nil
	}
}
