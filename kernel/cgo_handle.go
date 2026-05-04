package kernel

/*
#include <stdlib.h>
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

func newCgoHandlePointer(value any) unsafe.Pointer {
	userData := C.malloc(C.size_t(unsafe.Sizeof(cgo.Handle(0))))
	if userData == nil {
		panic("Failed to allocate callback handle storage")
	}
	handlePtr := (*cgo.Handle)(userData)
	*handlePtr = cgo.NewHandle(value)
	return userData
}

func cgoHandleFromPointer(userData unsafe.Pointer) cgo.Handle {
	handlePtr := (*cgo.Handle)(userData)
	return *handlePtr
}

func deleteCgoHandlePointer(userData unsafe.Pointer) {
	if userData != nil {
		cgoHandleFromPointer(userData).Delete()
		C.free(userData)
	}
}

//export go_delete_handle
func go_delete_handle(userData unsafe.Pointer) {
	deleteCgoHandlePointer(userData)
}
