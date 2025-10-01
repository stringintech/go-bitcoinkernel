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

// NewScriptPubkey creates a new script pubkey from raw serialized script data.
//
// The script pubkey defines the conditions that must be met to spend a transaction output.
//
// Parameters:
//   - rawScriptPubkey: Serialized script pubkey data
func NewScriptPubkey(rawScriptPubkey []byte) *ScriptPubkey {
	ptr := C.btck_script_pubkey_create(unsafe.Pointer(&rawScriptPubkey[0]), C.size_t(len(rawScriptPubkey)))
	return newScriptPubkey(check(ptr), true)
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

// Copy creates a copy of the script pubkey.
func (s *scriptPubkeyApi) Copy() *ScriptPubkey {
	return newScriptPubkey(s.ptr, false)
}

// Bytes returns the serialized representation of the script pubkey.
//
// Returns an error if the serialization fails.
func (s *scriptPubkeyApi) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, user_data unsafe.Pointer) C.int {
		return C.btck_script_pubkey_to_bytes(s.ptr, writer, user_data)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize script pubkey"}
	}
	return bytes, nil
}

// Verify verifies if the input at inputIndex of txTo spends the script pubkey
// under the constraints specified by flags. If the witness flag is set in flags,
// the amount parameter is used. If the taproot flag is set, spentOutputs is used
// to validate taproot transactions.
//
// Parameters:
//   - amount: Amount of the script pubkey's associated output. May be zero if the witness flag is not set.
//   - txTo: Transaction spending the script pubkey.
//   - spentOutputs: Outputs spent by the transaction. May be nil if the taproot flag is not set.
//   - inputIndex: Index of the input in txTo spending the script pubkey.
//   - flags: ScriptFlags controlling validation constraints.
//
// Returns an error if verification fails.
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

// ScriptFlags represents script verification flags that may be composed with each other.
type ScriptFlags C.btck_ScriptVerificationFlags

const (
	ScriptFlagsVerifyNone                ScriptFlags = C.btck_ScriptVerificationFlags_NONE                // No verification flags
	ScriptFlagsVerifyP2SH                ScriptFlags = C.btck_ScriptVerificationFlags_P2SH                // Evaluate P2SH (BIP16) subscripts
	ScriptFlagsVerifyDERSig              ScriptFlags = C.btck_ScriptVerificationFlags_DERSIG              // Enforce strict DER (BIP66) compliance
	ScriptFlagsVerifyNullDummy           ScriptFlags = C.btck_ScriptVerificationFlags_NULLDUMMY           // Enforce NULLDUMMY (BIP147)
	ScriptFlagsVerifyCheckLockTimeVerify ScriptFlags = C.btck_ScriptVerificationFlags_CHECKLOCKTIMEVERIFY // Enable CHECKLOCKTIMEVERIFY (BIP65)
	ScriptFlagsVerifyCheckSequenceVerify ScriptFlags = C.btck_ScriptVerificationFlags_CHECKSEQUENCEVERIFY // Enable CHECKSEQUENCEVERIFY (BIP112)
	ScriptFlagsVerifyWitness             ScriptFlags = C.btck_ScriptVerificationFlags_WITNESS             // Enable WITNESS (BIP141)
	ScriptFlagsVerifyTaproot             ScriptFlags = C.btck_ScriptVerificationFlags_TAPROOT             // Enable TAPROOT (BIPs 341 & 342)
	ScriptFlagsVerifyAll                 ScriptFlags = C.btck_ScriptVerificationFlags_ALL                 // All verification flags combined
)

// ScriptVerifyStatus represents a collection of status codes that may be issued by the script verify function.
type ScriptVerifyStatus C.btck_ScriptVerifyStatus

const (
	ScriptVerifyOK                           ScriptVerifyStatus = C.btck_ScriptVerifyStatus_SCRIPT_VERIFY_OK                // Script verification succeeded
	ScriptVerifyErrorInvalidFlagsCombination ScriptVerifyStatus = C.btck_ScriptVerifyStatus_ERROR_INVALID_FLAGS_COMBINATION // The verification flags were combined in an invalid way
	ScriptVerifyErrorSpentOutputsRequired    ScriptVerifyStatus = C.btck_ScriptVerifyStatus_ERROR_SPENT_OUTPUTS_REQUIRED    // The taproot flag was set, so valid spent outputs must be provided
)
