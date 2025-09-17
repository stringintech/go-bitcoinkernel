package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"unsafe"
)

type scriptPubkeyCFuncs struct{}

func (scriptPubkeyCFuncs) destroy(ptr unsafe.Pointer) {
	C.btck_script_pubkey_destroy((*C.btck_ScriptPubkey)(ptr))
}

func (scriptPubkeyCFuncs) copy(ptr unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.btck_script_pubkey_copy((*C.btck_ScriptPubkey)(ptr)))
}

type ScriptPubkey struct {
	*handle
	scriptPubkeyApi
}

func newScriptPubkey(ptr *C.btck_ScriptPubkey, fromOwned bool) *ScriptPubkey {
	h := newHandle(unsafe.Pointer(ptr), scriptPubkeyCFuncs{}, fromOwned)
	return &ScriptPubkey{handle: h, scriptPubkeyApi: scriptPubkeyApi{(*C.btck_ScriptPubkey)(h.ptr)}}
}

// NewScriptPubkey creates a new script pubkey from raw serialized data
func NewScriptPubkey(rawScriptPubkey []byte) (*ScriptPubkey, error) {
	ptr := C.btck_script_pubkey_create(unsafe.Pointer(&rawScriptPubkey[0]), C.size_t(len(rawScriptPubkey)))
	if ptr == nil {
		return nil, &InternalError{"Failed to create script pubkey from bytes"}
	}
	return newScriptPubkey(ptr, true), nil
}

type ScriptPubkeyView struct {
	scriptPubkeyApi
	ptr *C.btck_ScriptPubkey
}

func newScriptPubkeyView(ptr *C.btck_ScriptPubkey) *ScriptPubkeyView {
	return &ScriptPubkeyView{
		scriptPubkeyApi: scriptPubkeyApi{ptr},
		ptr:             ptr,
	}
}

type scriptPubkeyApi struct {
	ptr *C.btck_ScriptPubkey
}

func (s *scriptPubkeyApi) Copy() *ScriptPubkey {
	return newScriptPubkey(s.ptr, false)
}

// Bytes returns the serialized script pubkey
func (s *scriptPubkeyApi) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, user_data unsafe.Pointer) C.int {
		return C.btck_script_pubkey_to_bytes(s.ptr, writer, user_data)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize script pubkey"}
	}
	return bytes, nil
}

// Verify verifies this script pubkey against a transaction input
func (s *scriptPubkeyApi) Verify(amount int64, txTo *Transaction, spentOutputs []*TransactionOutput, inputIndex uint, flags ScriptFlags) error {
	inputCount := txTo.CountInputs()
	if inputIndex >= uint(inputCount) {
		return ErrVerifyScriptVerifyTxInputIndex
	}

	if len(spentOutputs) > 0 && uint64(len(spentOutputs)) != inputCount {
		return ErrVerifyScriptVerifySpentOutputsMismatch
	}

	allFlags := ScriptFlagsVerifyAll
	if (flags & ^ScriptFlags(allFlags)) != 0 {
		return ErrVerifyScriptVerifyInvalidFlags
	}

	var cSpentOutputsPtr **C.btck_TransactionOutput
	if len(spentOutputs) > 0 {
		cSpentOutputs := make([]*C.btck_TransactionOutput, len(spentOutputs))
		for i, output := range spentOutputs {
			cSpentOutputs[i] = (*C.btck_TransactionOutput)(output.handle.ptr)
		}
		cSpentOutputsPtr = (**C.btck_TransactionOutput)(unsafe.Pointer(&cSpentOutputs[0]))
	}

	var cStatus C.btck_ScriptVerifyStatus
	result := C.btck_script_pubkey_verify(
		s.ptr,
		C.int64_t(amount),
		(*C.btck_Transaction)(txTo.handle.ptr),
		cSpentOutputsPtr,
		C.size_t(len(spentOutputs)),
		C.uint(inputIndex),
		C.btck_ScriptVerificationFlags(flags),
		&cStatus,
	)

	if result != 1 {
		status := ScriptVerifyStatus(cStatus)
		switch status {
		case ScriptVerifyErrorInvalidFlagsCombination:
			return ErrVerifyScriptVerifyInvalidFlagsCombination
		case ScriptVerifyErrorSpentOutputsRequired:
			return ErrVerifyScriptVerifySpentOutputsRequired
		default:
			return ErrVerifyScriptVerifyInvalid
		}
	}
	return nil
}

const (
	ScriptFlagsVerifyNone                = C.btck_ScriptVerificationFlags_NONE
	ScriptFlagsVerifyP2SH                = C.btck_ScriptVerificationFlags_P2SH
	ScriptFlagsVerifyDERSig              = C.btck_ScriptVerificationFlags_DERSIG
	ScriptFlagsVerifyNullDummy           = C.btck_ScriptVerificationFlags_NULLDUMMY
	ScriptFlagsVerifyCheckLockTimeVerify = C.btck_ScriptVerificationFlags_CHECKLOCKTIMEVERIFY
	ScriptFlagsVerifyCheckSequenceVerify = C.btck_ScriptVerificationFlags_CHECKSEQUENCEVERIFY
	ScriptFlagsVerifyWitness             = C.btck_ScriptVerificationFlags_WITNESS
	ScriptFlagsVerifyTaproot             = C.btck_ScriptVerificationFlags_TAPROOT
	ScriptFlagsVerifyAll                 = C.btck_ScriptVerificationFlags_ALL
)

type ScriptFlags C.btck_ScriptVerificationFlags

const (
	ScriptVerifyOK                           = C.btck_ScriptVerifyStatus_SCRIPT_VERIFY_OK
	ScriptVerifyErrorInvalidFlagsCombination = C.btck_ScriptVerifyStatus_ERROR_INVALID_FLAGS_COMBINATION
	ScriptVerifyErrorSpentOutputsRequired    = C.btck_ScriptVerifyStatus_ERROR_SPENT_OUTPUTS_REQUIRED
)

type ScriptVerifyStatus C.btck_ScriptVerifyStatus
