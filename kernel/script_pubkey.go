package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ cManagedResource = &ScriptPubkey{}

// ScriptPubkey wraps the C btck_ScriptPubkey
type ScriptPubkey struct {
	ptr *C.btck_ScriptPubkey
}

// NewScriptPubkeyFromRaw creates a new script pubkey from raw serialized data
func NewScriptPubkeyFromRaw(rawScriptPubkey []byte) (*ScriptPubkey, error) {
	if len(rawScriptPubkey) == 0 {
		return nil, ErrEmptyScriptPubkeyData
	}
	ptr := C.btck_script_pubkey_create(unsafe.Pointer(&rawScriptPubkey[0]), C.size_t(len(rawScriptPubkey)))
	if ptr == nil {
		return nil, ErrKernelScriptPubkeyCreate
	}

	scriptPubkey := &ScriptPubkey{ptr: ptr}
	runtime.SetFinalizer(scriptPubkey, (*ScriptPubkey).destroy)
	return scriptPubkey, nil
}

// Data returns the serialized script pubkey data
func (s *ScriptPubkey) Data() ([]byte, error) {
	checkReady(s)

	return writeToBytes(func(writer C.btck_WriteBytes, user_data unsafe.Pointer) C.int {
		return C.btck_script_pubkey_to_bytes(s.ptr, writer, user_data)
	})
}

// Copy creates a copy of the script pubkey
func (s *ScriptPubkey) Copy() (*ScriptPubkey, error) {
	checkReady(s)

	ptr := C.btck_script_pubkey_copy(s.ptr)
	if ptr == nil {
		return nil, ErrKernelScriptPubkeyCopy
	}

	scriptPubkey := &ScriptPubkey{ptr: ptr}
	runtime.SetFinalizer(scriptPubkey, (*ScriptPubkey).destroy)
	return scriptPubkey, nil
}

func (s *ScriptPubkey) destroy() {
	if s.ptr != nil {
		C.btck_script_pubkey_destroy(s.ptr)
		s.ptr = nil
	}
}

func (s *ScriptPubkey) Destroy() {
	runtime.SetFinalizer(s, nil)
	s.destroy()
}

func (s *ScriptPubkey) isReady() bool {
	return s != nil && s.ptr != nil
}

func (s *ScriptPubkey) uninitializedError() error {
	return ErrScriptPubkeyUninitialized
}
