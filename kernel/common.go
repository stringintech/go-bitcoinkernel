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
