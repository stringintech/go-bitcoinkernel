package kernel

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
#include <stdint.h>

// Bridge function: exported Go function that C library can call
// user_data contains the cgo.Handle ID as void* for callback identification
extern void go_log_callback_bridge(void* user_data, char* message, size_t message_len);
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

// LoggingConnection wraps the C btck_LoggingConnection.
// Functions changing the logging settings are global and change
// the settings for all existing btck_LoggingConnection instances.
type LoggingConnection struct {
	ptr    *C.btck_LoggingConnection
	handle cgo.Handle // Prevents callback GC until Delete() called
}

// boolToInt converts Go bool to C int (0 for false, 1 for true)
func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
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

	connection := &LoggingConnection{
		ptr:    nil,
		handle: 0,
	}

	// Create a handle for the callback - this prevents garbage collection
	// and provides a stable ID that can be passed through C code safely
	connection.handle = cgo.NewHandle(callback)

	cOptions := C.btck_LoggingOptions{
		log_timestamps:               boolToInt(options.LogTimestamps),
		log_time_micros:              boolToInt(options.LogTimeMicros),
		log_threadnames:              boolToInt(options.LogThreadNames),
		log_sourcelocations:          boolToInt(options.LogSourceLocations),
		always_print_category_levels: boolToInt(options.AlwaysPrintCategoryLevel),
	}

	connection.ptr = C.btck_logging_connection_create((C.btck_LogCallback)(C.go_log_callback_bridge),
		unsafe.Pointer(connection.handle), nil, cOptions)
	if connection.ptr == nil {
		connection.handle.Delete()
		return nil, ErrKernelLoggingConnectionCreate
	}

	runtime.SetFinalizer(connection, (*LoggingConnection).destroy)
	return connection, nil
}

func (lc *LoggingConnection) destroy() {
	if lc.ptr != nil {
		C.btck_logging_connection_destroy(lc.ptr)
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
	C.btck_logging_disable()
}

// Global mutex for thread-safe category management
var loggingMutex = sync.RWMutex{}

// AddLogLevelCategory sets the log level for a specific category or all categories
func AddLogLevelCategory(category LogCategory, level LogLevel) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.btck_logging_set_level_category(category.mustC(), level.mustC())
}

// EnableLogCategory enables logging for a specific category or all categories
func EnableLogCategory(category LogCategory) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.btck_logging_enable_category(category.mustC())
}

// DisableLogCategory disables logging for a specific category or all categories
func DisableLogCategory(category LogCategory) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	C.btck_logging_disable_category(category.mustC())
}

const (
	LogLevelTrace = C.btck_LogLevel_TRACE
	LogLevelDebug = C.btck_LogLevel_DEBUG
	LogLevelInfo  = C.btck_LogLevel_INFO
)

type LogLevel C.btck_LogLevel

func (l LogLevel) mustC() C.btck_LogLevel {
	c, err := l.c()
	if err != nil {
		panic(err)
	}
	return c
}

func (l LogLevel) c() (C.btck_LogLevel, error) {
	switch l {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo:
		return C.btck_LogLevel(l), nil
	default:
		return 0, ErrInvalidLogLevel
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

func (c LogCategory) mustC() C.btck_LogCategory {
	cType, err := c.c()
	if err != nil {
		panic(err)
	}
	return cType
}

func (c LogCategory) c() (C.btck_LogCategory, error) {
	switch c {
	case LogAll, LogBench, LogBlockStorage, LogCoinDB, LogLevelDB, LogMempool, LogPrune, LogRand, LogReindex, LogValidation, LogKernel:
		return C.btck_LogCategory(c), nil
	default:
		return 0, ErrInvalidLogCategory
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
