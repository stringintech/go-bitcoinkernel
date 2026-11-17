package kernel

import "C"
import (
	"runtime"
	"unsafe"
)

type cResourceFuncs interface {
	destroy(ptr unsafe.Pointer)
}

type copyableCResourceFuncs interface {
	cResourceFuncs
	copy(ptr unsafe.Pointer) unsafe.Pointer
}

// cResourceHandle is a finalizable handle that frees a C object.
//
// API types store one in their owner field. Owned types store their own handle,
// views store their parent's, so a reachable value keeps its owner reachable.
// Methods that pass a derived C pointer to C call runtime.KeepAlive on the
// owning value after the call, since the GC cannot see C's use of the pointer
// and could otherwise treat the value as unreachable mid-call.
type cResourceHandle interface {
	destroy()
}

type uniqueHandle struct {
	ptr   unsafe.Pointer
	funcs cResourceFuncs
}

func newUniqueHandle(ptr unsafe.Pointer, funcs cResourceFuncs) *uniqueHandle {
	if ptr == nil {
		panic("ptr must be provided to create handle")
	}
	if funcs == nil {
		panic("funcs must be provided to create handle")
	}
	h := &uniqueHandle{
		ptr:   ptr,
		funcs: funcs,
	}
	runtime.SetFinalizer(h, (*uniqueHandle).destroy)
	return h
}

func (h *uniqueHandle) destroy() {
	if h.ptr != nil {
		h.funcs.destroy(h.ptr)
		h.ptr = nil
	}
}

func (h *uniqueHandle) Destroy() {
	runtime.SetFinalizer(h, nil)
	h.destroy()
}

type handle struct {
	ptr   unsafe.Pointer
	funcs copyableCResourceFuncs
}

// newHandle wraps ptr in a finalizable handle. If fromOwned is true, the handle
// takes ownership of ptr, which must be a C object nobody else frees. If it is
// false, ptr is borrowed and newHandle stores a C copy of it. The caller must
// keep the source of ptr reachable until newHandle returns.
func newHandle(ptr unsafe.Pointer, funcs copyableCResourceFuncs, fromOwned bool) *handle {
	if ptr == nil {
		panic("ptr must be provided to create handle")
	}
	if funcs == nil {
		panic("funcs must be provided to create handle")
	}

	if !fromOwned {
		ptr = funcs.copy(ptr)
		if ptr == nil {
			panic(ErrKernelInstantiate)
		}
	}

	h := &handle{
		ptr:   ptr,
		funcs: funcs,
	}
	runtime.SetFinalizer(h, (*handle).destroy)
	return h
}

func (h *handle) destroy() {
	if h.ptr != nil {
		h.funcs.destroy(h.ptr)
		h.ptr = nil
	}
}

func (h *handle) Destroy() {
	runtime.SetFinalizer(h, nil)
	h.destroy()
}

func ReverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}
