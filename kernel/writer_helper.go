package kernel

// Writer implementation that minimizes memory copies where possible
// by using direct memory operations, though buffer growth requires copying

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
	buffer   []byte // Pre-allocated buffer
	position int    // Current write position
	err      error
}

//export go_writer_callback_bridge
func go_writer_callback_bridge(bytes unsafe.Pointer, size C.size_t, userdata unsafe.Pointer) C.int {
	// Convert void* back to Handle - userdata contains Handle ID
	handle := cgo.Handle(userdata)
	// Retrieve original Go callback data struct
	data := handle.Value().(*writerCallbackData)

	bytesSize := int(size)
	requiredCapacity := data.position + bytesSize

	// Grow buffer if needed (double capacity strategy)
	if requiredCapacity > len(data.buffer) {
		newCapacity := len(data.buffer) * 2
		if newCapacity < requiredCapacity {
			newCapacity = requiredCapacity
		}
		newBuffer := make([]byte, newCapacity)
		copy(newBuffer[:data.position], data.buffer[:data.position])
		data.buffer = newBuffer
	}

	if bytesSize > 0 {
		// Get pointer to destination in Go buffer
		dstPtr := unsafe.Pointer(&data.buffer[data.position])
		// Use C's memmove to copy directly
		C.memmove(dstPtr, bytes, size)
		data.position += bytesSize
	}

	return 0 // success
}

// writeToBytes is a helper function that uses a callback pattern to collect bytes
// It takes a function that calls the C API with the writer callback
func writeToBytes(writerFunc func(C.btck_WriteBytes, unsafe.Pointer) C.int) ([]byte, error) {
	// Pre-allocate buffer with reasonable initial capacity
	initialCapacity := 1024 // Start with 1KB, will grow as needed
	callbackData := &writerCallbackData{
		buffer:   make([]byte, initialCapacity),
		position: 0,
	}

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
	return callbackData.buffer[:callbackData.position], nil
}
