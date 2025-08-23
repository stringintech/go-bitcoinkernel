package kernel

// Writer implementation that minimizes memory copies

/*
#include "kernel/bitcoinkernel.h"
#include <stdlib.h>
#include <string.h>

// Bridge function: exported Go function that C library can call
extern int go_writer_callback_bridge(void* bytes, size_t size, void* userdata);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// writerCallbackData holds the growing buffer that collects written bytes
// Uses a capacity-based growth strategy to reduce reallocations
type writerCallbackData struct {
	buffer []byte // Pre-allocated buffer
	err    error
}

//export go_writer_callback_bridge
func go_writer_callback_bridge(bytes unsafe.Pointer, size C.size_t, userdata unsafe.Pointer) C.int {
	// Convert void* back to Handle - userdata contains Handle ID
	handle := cgo.Handle(userdata)
	// Retrieve original Go callback data struct
	data := handle.Value().(*writerCallbackData)

	if size > 0 {
		// Create a Go slice view of the C memory
		cBytes := unsafe.Slice((*byte)(bytes), int(size))
		data.buffer = append(data.buffer, cBytes...)
	}
	return 0
}

// writeToBytes is a helper function that uses a callback pattern to collect bytes
// It takes a function that calls the C API with the writer callback
func writeToBytes(writerFunc func(C.btck_WriteBytes, unsafe.Pointer) C.int) ([]byte, error) {
	callbackData := &writerCallbackData{}

	// Create cgo handle for the callback data
	handle := cgo.NewHandle(callbackData)
	defer handle.Delete()

	// Call the C function with our callback
	result := writerFunc((C.btck_WriteBytes)(C.go_writer_callback_bridge), unsafe.Pointer(handle))
	if result != 0 {
		return nil, &KernelError{Operation: "writer_callback", Detail: "serialization failed"}
	}

	if callbackData.err != nil {
		return nil, callbackData.err
	}

	// Return exactly the bytes that were written (slice the buffer to actual size)
	return callbackData.buffer, nil
}
