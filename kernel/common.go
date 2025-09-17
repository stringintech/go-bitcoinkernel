package kernel

import "C"
import (
	"runtime"
	"runtime/cgo"
	"unsafe"
)

type destroyableFuncs interface {
	destroy(ptr unsafe.Pointer)
}

type copyableFuncs interface {
	destroyableFuncs
	copy(ptr unsafe.Pointer) unsafe.Pointer
}

type uniqueHandle struct {
	ptr   unsafe.Pointer
	funcs destroyableFuncs
}

func newUniqueHandle(ptr unsafe.Pointer, funcs destroyableFuncs) *uniqueHandle {
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
	funcs copyableFuncs
}

func newHandle(ptr unsafe.Pointer, funcs copyableFuncs, fromOwned bool) *handle {
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

//export go_delete_handle
func go_delete_handle(handle unsafe.Pointer) {
	cgo.Handle(handle).Delete()
}
