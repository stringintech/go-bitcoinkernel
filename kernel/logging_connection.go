package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
#include <stdint.h>

// Bridge function: exported Go function that C library can call
// user_data contains the cgo.Handle ID as void* for callback identification
extern void go_log_callback_bridge(void* user_data, char* message, size_t message_len);

// Wrapper function: C helper to create logging connection with Go callback
// Converts Handle ID to void* and passes to C library
static inline kernel_LoggingConnection* create_logging_connection_wrapper(uintptr_t context, kernel_LoggingOptions options) {
    return kernel_logging_connection_create((kernel_LogCallback)go_log_callback_bridge, (void*)context, options);
}
*/
import "C"
import (
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"
)

// LogCallback is the Go callback function type for log messages.
type LogCallback func(message string)

var _ cManagedResource = &LoggingConnection{}

// LoggingConnection wraps the C kernel_LoggingConnection.
// Functions changing the logging settings are global and change
// the settings for all existing kernel_LoggingConnection instances.
type LoggingConnection struct {
	ptr    *C.kernel_LoggingConnection
	handle cgo.Handle // Prevents callback GC until Delete() called
}

//export go_log_callback_bridge
func go_log_callback_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	// Convert void* back to Handle - user_data contains Handle ID
	handle := cgo.Handle(user_data)
	// Retrieve original Go callback
	callback := handle.Value().(LogCallback)

	goMessage := C.GoStringN(message, C.int(message_len))

	// Call the actual Go callback
	callback(goMessage)
}

func NewLoggingConnection(callback LogCallback, options LoggingOptions) (*LoggingConnection, error) {
	if callback == nil {
		return nil, ErrLoggingConnectionUninitialized //FIXME
	}

	// Create a handle for the callback - this prevents garbage collection
	// and provides a stable ID that can be passed through C code safely
	handle := cgo.NewHandle(callback)

	cOptions := C.kernel_LoggingOptions{
		log_timestamps:               C.bool(options.LogTimestamps),
		log_time_micros:              C.bool(options.LogTimeMicros),
		log_threadnames:              C.bool(options.LogThreadNames),
		log_sourcelocations:          C.bool(options.LogSourceLocations),
		always_print_category_levels: C.bool(options.AlwaysPrintCategoryLevel),
	}

	ptr := C.create_logging_connection_wrapper(C.uintptr_t(handle), cOptions)
	if ptr == nil {
		handle.Delete()
		return nil, ErrKernelLoggingConnectionCreate
	}

	connection := &LoggingConnection{
		ptr:    ptr,
		handle: handle,
	}

	runtime.SetFinalizer(connection, (*LoggingConnection).destroy)
	return connection, nil
}

func (lc *LoggingConnection) destroy() {
	if lc.ptr != nil {
		C.kernel_logging_connection_destroy(lc.ptr)
		lc.ptr = nil
	}
	if lc.handle != 0 {
		// Delete exposes callback to garbage collection
		lc.handle.Delete()
		lc.handle = 0
	}
}

func (lc *LoggingConnection) Destroy() {
	runtime.SetFinalizer(lc, nil)
	lc.destroy()
}

func (lc *LoggingConnection) isReady() bool {
	return lc != nil && lc.ptr != nil && lc.handle != 0
}

func (lc *LoggingConnection) uninitializedError() error {
	return ErrLoggingConnectionUninitialized
}

// DisableLogging permanently disables the global internal logger.
// This function should only be called once and is not thread-safe
func DisableLogging() {
	C.kernel_disable_logging()
}

// Global mutex for thread-safe category management
var loggingMutex = sync.RWMutex{}

// AddLogLevelCategory sets the log level for a specific category or all categories
func AddLogLevelCategory(category LogCategory, level LogLevel) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.kernel_add_log_level_category(category.mustC(), level.mustC())
}

// EnableLogCategory enables logging for a specific category or all categories
func EnableLogCategory(category LogCategory) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.kernel_enable_log_category(category.mustC())
}

// DisableLogCategory disables logging for a specific category or all categories
func DisableLogCategory(category LogCategory) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.kernel_disable_log_category(category.mustC())
}

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
)

func (l LogLevel) mustC() C.kernel_LogLevel {
	c, err := l.c()
	if err != nil {
		panic(err)
	}
	return c
}

func (l LogLevel) c() (C.kernel_LogLevel, error) {
	if l < LogLevelTrace || l > LogLevelInfo {
		return 0, ErrInvalidLogLevel
	}
	return C.kernel_LogLevel(l), nil
}

// LogCategory represents a logging category
type LogCategory int

const (
	LogAll LogCategory = iota
	LogBench
	LogBlockStorage
	LogCoinDB
	LogLevelDB
	LogMempool
	LogPrune
	LogRand
	LogReindex
	LogValidation
	LogKernel
)

func (c LogCategory) mustC() C.kernel_LogCategory {
	cType, err := c.c()
	if err != nil {
		panic(err)
	}
	return cType
}

func (c LogCategory) c() (C.kernel_LogCategory, error) {
	if c < LogAll || c > LogKernel {
		return 0, ErrInvalidLogCategory
	}
	return C.kernel_LogCategory(c), nil
}

// LoggingOptions configures the format of log messages
type LoggingOptions struct {
	LogTimestamps            bool // Prepend a timestamp to log messages
	LogTimeMicros            bool // Log timestamps in microsecond precision
	LogThreadNames           bool // Prepend the name of the thread to log messages
	LogSourceLocations       bool // Prepend the source location to log messages
	AlwaysPrintCategoryLevel bool // Prepend the log category and level to log messages
}
