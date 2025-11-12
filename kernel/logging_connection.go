package kernel

/*
#include "bitcoinkernel.h"
#include <stdlib.h>
#include <stdint.h>

// Bridge function: exported Go function that C library can call
// user_data contains the cgo.Handle ID as void* for callback identification
extern void go_log_callback_bridge(void* user_data, char* message, size_t message_len);

extern void go_delete_handle(void* user_data);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// LogCallback is the Go callback function type for log messages.
type LogCallback func(message string)

type loggingConnectionCFuncs struct{}

func (loggingConnectionCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_logging_connection_destroy((*C.btck_LoggingConnection)(ptr))
}

// LoggingConnection can be used to manually stop logging.
//
// Messages that were logged before a connection is created are buffered in a 1MB
// buffer. Logging can alternatively be permanently disabled by calling DisableLogging.
// Functions changing the logging settings are global and change the settings for all
// existing LoggingConnection instances.
type LoggingConnection struct {
	*uniqueHandle
}

//export go_log_callback_bridge
func go_log_callback_bridge(user_data unsafe.Pointer, message *C.char, message_len C.size_t) {
	handle := cgo.Handle(user_data)
	callback := handle.Value().(LogCallback)
	goMessage := C.GoStringN(message, C.int(message_len))
	callback(goMessage)
}

// NewLoggingConnection creates a logging connection that routes kernel logs to a Go callback.
//
// Messages logged before the first connection is created are buffered in a 1MB buffer and
// delivered immediately when this function is called.
//
// Parameters:
//   - callback: Function that will receive log messages
//
// Returns an error if the logging connection cannot be created.
func NewLoggingConnection(callback LogCallback) (*LoggingConnection, error) {
	callbackHandle := cgo.NewHandle(callback)
	ptr := C.btck_logging_connection_create((C.btck_LogCallback)(C.go_log_callback_bridge),
		unsafe.Pointer(callbackHandle), C.btck_DestroyCallback(C.go_delete_handle))
	if ptr == nil {
		callbackHandle.Delete()
		return nil, &InternalError{"Failed to create logging connection"}
	}
	h := newUniqueHandle(unsafe.Pointer(ptr), loggingConnectionCFuncs{})
	return &LoggingConnection{uniqueHandle: h}, nil
}

// SetLoggingOptions configures the formatting options for the global internal logger.
//
// This changes global settings and affects all existing LoggingConnection instances.
//
// Parameters:
//   - options: Formatting options for log messages (timestamps, thread names, etc.)
func SetLoggingOptions(options LoggingOptions) {
	cOptions := C.btck_LoggingOptions{
		log_timestamps:               boolToInt(options.LogTimestamps),
		log_time_micros:              boolToInt(options.LogTimeMicros),
		log_threadnames:              boolToInt(options.LogThreadNames),
		log_sourcelocations:          boolToInt(options.LogSourceLocations),
		always_print_category_levels: boolToInt(options.AlwaysPrintCategoryLevel),
	}
	C.btck_logging_set_options(cOptions)
}

// DisableLogging permanently disables the global internal logger.
//
// Log messages will be buffered until this function is called, or a logging connection
// is created. This must not be called while a logging connection already exists.
// This function should only be called once and is not thread-safe or re-entry safe.
func DisableLogging() {
	C.btck_logging_disable()
}

// AddLogLevelCategory sets the log level for a specific category.
//
// This does not enable the selected categories. Use EnableLogCategory to
// start logging from a specific, or all categories. This changes a global
// setting and will override settings for all existing LoggingConnection instances.
//
// Parameters:
//   - category: If LogAll is chosen, sets both the global fallback log level
//     used by all categories that don't have a specific level set, and also
//     sets the log level for messages logged with the LogAll category itself.
//     For any other category, sets a category-specific log level that overrides
//     the global fallback for that category only.
//   - level: Minimum log level (Trace, Debug, or Info)
//
// Messages at the specified level and above will be logged for the category.
func AddLogLevelCategory(category LogCategory, level LogLevel) {
	C.btck_logging_set_level_category(category.c(), level.c())
}

// EnableLogCategory enables logging for a specific category.
//
// This changes a global setting and affects all existing LoggingConnection instances.
//
// Parameters:
//   - category: Log category to enable (or LogAll to enable all categories)
func EnableLogCategory(category LogCategory) {
	C.btck_logging_enable_category(category.c())
}

// DisableLogCategory disables logging for a specific category.
//
// This changes a global setting and affects all existing LoggingConnection instances.
//
// Parameters:
//   - category: Log category to disable (or LogAll to disable all categories)
func DisableLogCategory(category LogCategory) {
	C.btck_logging_disable_category(category.c())
}

// LogLevel represents the level at which logs should be produced.
type LogLevel C.btck_LogLevel

const (
	LogLevelTrace LogLevel = C.btck_LogLevel_TRACE
	LogLevelDebug LogLevel = C.btck_LogLevel_DEBUG
	LogLevelInfo  LogLevel = C.btck_LogLevel_INFO
)

func (l LogLevel) c() C.btck_LogLevel {
	switch l {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo:
		return C.btck_LogLevel(l)
	default:
		panic("Invalid log level")
	}
}

// LogCategory represents a collection of logging categories that may be encountered by kernel code.
type LogCategory C.btck_LogCategory

const (
	LogAll          LogCategory = C.btck_LogCategory_ALL
	LogBench        LogCategory = C.btck_LogCategory_BENCH
	LogBlockStorage LogCategory = C.btck_LogCategory_BLOCKSTORAGE
	LogCoinDB       LogCategory = C.btck_LogCategory_COINDB
	LogLevelDB      LogCategory = C.btck_LogCategory_LEVELDB
	LogMempool      LogCategory = C.btck_LogCategory_MEMPOOL
	LogPrune        LogCategory = C.btck_LogCategory_PRUNE
	LogRand         LogCategory = C.btck_LogCategory_RAND
	LogReindex      LogCategory = C.btck_LogCategory_REINDEX
	LogValidation   LogCategory = C.btck_LogCategory_VALIDATION
	LogKernel       LogCategory = C.btck_LogCategory_KERNEL
)

func (c LogCategory) c() C.btck_LogCategory {
	switch c {
	case LogAll, LogBench, LogBlockStorage, LogCoinDB, LogLevelDB, LogMempool, LogPrune, LogRand, LogReindex, LogValidation, LogKernel:
		return C.btck_LogCategory(c)
	default:
		panic("Invalid log category")
	}
}

// LoggingOptions configures the format of log messages
type LoggingOptions struct {
	LogTimestamps            bool // Prepend a timestamp to log messages
	LogTimeMicros            bool // Log timestamps in microsecond precision
	LogThreadNames           bool // Prepend the name of the thread to log messages
	LogSourceLocations       bool // Prepend the source location to log messages
	AlwaysPrintCategoryLevel bool // Prepend the log category and level to log messages
}

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
