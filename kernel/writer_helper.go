package kernel

/*
#include "bitcoinkernel.h"
#include <stdlib.h>
#include <string.h>

extern int go_writer_callback_bridge(void* bytes, size_t size, void* userdata);
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// writerCallbackData holds the growing buffer that collects written bytes
type writerCallbackData struct {
	buffer []byte
}

//export go_writer_callback_bridge
func go_writer_callback_bridge(bytes unsafe.Pointer, size C.size_t, userdata unsafe.Pointer) C.int {
	if size > 0 {
		data := cgo.Handle(userdata).Value().(*writerCallbackData)
		// Create a Go slice view of the C memory
		cBytes := unsafe.Slice((*byte)(bytes), int(size))
		data.buffer = append(data.buffer, cBytes...)
	}
	return 0
}

// writeToBytes is a helper function that uses a callback pattern to collect bytes
// It takes a function that calls the C API with the writer callback
func writeToBytes(writerFunc func(C.btck_WriteBytes, unsafe.Pointer) C.int) (bytes []byte, ok bool) {
	callbackData := &writerCallbackData{}
	handle := cgo.NewHandle(callbackData)
	defer handle.Delete()

	result := writerFunc((C.btck_WriteBytes)(C.go_writer_callback_bridge), unsafe.Pointer(handle))
	if result != 0 {
		return nil, false
	}
	return callbackData.buffer, true
}
