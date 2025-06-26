package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// ScriptPubkey wraps the C kernel_ScriptPubkey
type ScriptPubkey struct {
	ptr *C.kernel_ScriptPubkey
}

// NewScriptPubkeyFromRaw creates a new script pubkey from raw serialized data
func NewScriptPubkeyFromRaw(rawScriptPubkey []byte) (*ScriptPubkey, error) {
	if len(rawScriptPubkey) == 0 {
		return nil, ErrInvalidScriptPubkeyData
	}

	ptr := C.kernel_script_pubkey_create((*C.uchar)(unsafe.Pointer(&rawScriptPubkey[0])), C.size_t(len(rawScriptPubkey)))
	if ptr == nil {
		return nil, ErrScriptPubkeyCreation
	}

	scriptPubkey := &ScriptPubkey{ptr: ptr}
	runtime.SetFinalizer(scriptPubkey, (*ScriptPubkey).destroy)
	return scriptPubkey, nil
}

// Data returns the serialized script pubkey data
func (s *ScriptPubkey) Data() ([]byte, error) {
	if s.ptr == nil {
		return nil, ErrInvalidScriptPubkey
	}

	byteArray := C.kernel_copy_script_pubkey_data(s.ptr)
	if byteArray == nil {
		return nil, ErrScriptPubkeyDataCopy
	}
	defer C.kernel_byte_array_destroy(byteArray)

	size := int(byteArray.size)
	if size == 0 {
		return nil, nil
	}

	// Copy the data to Go slice
	data := C.GoBytes(unsafe.Pointer(byteArray.data), C.int(size))
	return data, nil
}

func (s *ScriptPubkey) destroy() {
	if s.ptr != nil {
		C.kernel_script_pubkey_destroy(s.ptr)
		s.ptr = nil
	}
}

func (s *ScriptPubkey) Destroy() {
	runtime.SetFinalizer(s, nil)
	s.destroy()
}
