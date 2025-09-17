package kernel

/*
#include "kernel/bitcoinkernel.h"
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

func NewLoggingConnection(callback LogCallback, options LoggingOptions) (*LoggingConnection, error) {
	cOptions := C.btck_LoggingOptions{
		log_timestamps:               boolToInt(options.LogTimestamps),
		log_time_micros:              boolToInt(options.LogTimeMicros),
		log_threadnames:              boolToInt(options.LogThreadNames),
		log_sourcelocations:          boolToInt(options.LogSourceLocations),
		always_print_category_levels: boolToInt(options.AlwaysPrintCategoryLevel),
	}

	callbackHandle := cgo.NewHandle(callback)
	ptr := C.btck_logging_connection_create((C.btck_LogCallback)(C.go_log_callback_bridge),
		unsafe.Pointer(callbackHandle), C.btck_DestroyCallback(C.go_delete_handle), cOptions)
	if ptr == nil {
		callbackHandle.Delete()
		return nil, &InternalError{"Failed to create logging connection"}
	}
	h := newUniqueHandle(unsafe.Pointer(ptr), loggingConnectionCFuncs{})
	return &LoggingConnection{uniqueHandle: h}, nil
}

// DisableLogging permanently disables the global internal logger.
// This function should only be called once and is not thread-safe
func DisableLogging() {
	C.btck_logging_disable()
}

// AddLogLevelCategory sets the log level for a specific category or all categories
func AddLogLevelCategory(category LogCategory, level LogLevel) {
	C.btck_logging_set_level_category(category.c(), level.c())
}

// EnableLogCategory enables logging for a specific category or all categories
func EnableLogCategory(category LogCategory) {
	C.btck_logging_enable_category(category.c())
}

// DisableLogCategory disables logging for a specific category or all categories
func DisableLogCategory(category LogCategory) {
	C.btck_logging_disable_category(category.c())
}

const (
	LogLevelTrace = C.btck_LogLevel_TRACE
	LogLevelDebug = C.btck_LogLevel_DEBUG
	LogLevelInfo  = C.btck_LogLevel_INFO
)

type LogLevel C.btck_LogLevel

func (l LogLevel) c() C.btck_LogLevel {
	switch l {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo:
		return C.btck_LogLevel(l)
	default:
		panic("Invalid log level")
	}
}

const (
	LogAll          = C.btck_LogCategory_ALL
	LogBench        = C.btck_LogCategory_BENCH
	LogBlockStorage = C.btck_LogCategory_BLOCKSTORAGE
	LogCoinDB       = C.btck_LogCategory_COINDB
	LogLevelDB      = C.btck_LogCategory_LEVELDB
	LogMempool      = C.btck_LogCategory_MEMPOOL
	LogPrune        = C.btck_LogCategory_PRUNE
	LogRand         = C.btck_LogCategory_RAND
	LogReindex      = C.btck_LogCategory_REINDEX
	LogValidation   = C.btck_LogCategory_VALIDATION
	LogKernel       = C.btck_LogCategory_KERNEL
)

type LogCategory C.btck_LogCategory

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
