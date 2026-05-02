package kernel

/*
#include "bitcoinkernel.h"
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
	return &ScriptPubkey{
		handle: h,
		scriptPubkeyApi: scriptPubkeyApi{
			ptr: func() *C.btck_ScriptPubkey {
				return (*C.btck_ScriptPubkey)(h.ptr)
			},
		},
	}
}

// NewScriptPubkey creates a new script pubkey from raw serialized script data.
//
// The script pubkey defines the conditions that must be met to spend a transaction output.
//
// Parameters:
//   - rawScriptPubkey: Serialized script pubkey data
func NewScriptPubkey(rawScriptPubkey []byte) *ScriptPubkey {
	ptr := C.btck_script_pubkey_create(unsafe.Pointer(unsafe.SliceData(rawScriptPubkey)), C.size_t(len(rawScriptPubkey)))
	return newScriptPubkey(check(ptr), true)
}

type ScriptPubkeyView struct {
	scriptPubkeyApi
}

func newScriptPubkeyView(ptr *C.btck_ScriptPubkey) *ScriptPubkeyView {
	return &ScriptPubkeyView{
		scriptPubkeyApi: scriptPubkeyApi{
			ptr: func() *C.btck_ScriptPubkey {
				return ptr
			},
		},
	}
}

type scriptPubkeyApi struct {
	ptr func() *C.btck_ScriptPubkey
}

func (s *scriptPubkeyApi) cPtr() *C.btck_ScriptPubkey {
	return s.ptr()
}

// ScriptPubkeyLike is implemented by *ScriptPubkey and *ScriptPubkeyView.
type ScriptPubkeyLike interface {
	cPtr() *C.btck_ScriptPubkey
	Copy() *ScriptPubkey
	Bytes() ([]byte, error)
	Verify(int64, TransactionLike, *PrecomputedTransactionData, uint, ScriptFlags) (bool, error)
}

var _ ScriptPubkeyLike = (*ScriptPubkey)(nil)
var _ ScriptPubkeyLike = (*ScriptPubkeyView)(nil)

// Copy creates a copy of the script pubkey.
func (s *scriptPubkeyApi) Copy() *ScriptPubkey {
	return newScriptPubkey(s.ptr(), false)
}

// Bytes returns the serialized representation of the script pubkey.
//
// Returns an error if the serialization fails.
func (s *scriptPubkeyApi) Bytes() ([]byte, error) {
	bytes, ok := writeToBytes(func(writer C.btck_WriteBytes, user_data unsafe.Pointer) C.int {
		return C.btck_script_pubkey_to_bytes(s.ptr(), writer, user_data)
	})
	if !ok {
		return nil, &SerializationError{"Failed to serialize script pubkey"}
	}
	return bytes, nil
}

// Verify verifies if the input at inputIndex of txTo spends the script pubkey
// under the constraints specified by flags. If the witness flag is set in flags,
// the amount parameter is used. If the taproot flag is set, precomputedTxData must
// contain the spent outputs.
//
// Parameters:
//   - amount: Amount of the script pubkey's associated output. May be zero if the witness flag is not set.
//   - txTo: Transaction spending the script pubkey.
//   - precomputedTxData: Precomputed transaction data. May be nil if the taproot flag is not set.
//   - inputIndex: Index of the input in txTo spending the script pubkey.
//   - flags: ScriptFlags controlling validation constraints.
//
// Returns:
//   - bool: true if the script is valid, false if invalid (only meaningful when error is nil)
//   - error: non-nil if verification could not be performed due to malformed input;
//     nil if verification completed successfully (check bool for validity result)
func (s *scriptPubkeyApi) Verify(amount int64, txTo TransactionLike, precomputedTxData *PrecomputedTransactionData, inputIndex uint, flags ScriptFlags) (bool, error) {
	inputCount := txTo.CountInputs()
	if inputIndex >= uint(inputCount) {
		return false, ErrVerifyScriptVerifyTxInputIndex
	}

	allFlags := ScriptFlagsVerifyAll
	if (flags & ^ScriptFlags(allFlags)) != 0 {
		return false, ErrVerifyScriptVerifyInvalidFlags
	}

	var cPrecomputedTxData *C.btck_PrecomputedTransactionData
	if precomputedTxData != nil {
		cPrecomputedTxData = (*C.btck_PrecomputedTransactionData)(precomputedTxData.ptr)
	}

	var cStatus C.btck_ScriptVerifyStatus
	result := C.btck_script_pubkey_verify(
		s.ptr(),
		C.int64_t(amount),
		txTo.cPtr(),
		cPrecomputedTxData,
		C.uint(inputIndex),
		C.btck_ScriptVerificationFlags(flags),
		&cStatus,
	)

	// Check for errors that prevented verification
	if cStatus == C.btck_ScriptVerifyStatus_ERROR_INVALID_FLAGS_COMBINATION {
		return false, ErrVerifyScriptVerifyInvalidFlagsCombination
	}
	if cStatus == C.btck_ScriptVerifyStatus_ERROR_SPENT_OUTPUTS_REQUIRED {
		return false, ErrVerifyScriptVerifySpentOutputsRequired
	}

	// Verification completed: result indicates validity
	// result == 1: script is valid
	// result != 1: script is invalid
	return result == 1, nil
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
