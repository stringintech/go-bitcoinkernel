package kernel

/*
#include "kernel/bitcoinkernel.h"
*/
import "C"
import (
	"fmt"
)

// VerifyScript verifies a script pubkey against a transaction input
func VerifyScript(scriptPubkey *ScriptPubkey, amount int64, txTo *Transaction,
	spentOutputs []*TransactionOutput, inputIndex uint, flags ScriptFlags) error {
	if err := validateReady(scriptPubkey); err != nil {
		return err
	}
	if err := validateReady(txTo); err != nil {
		return err
	}

	var cSpentOutputs **C.kernel_TransactionOutput
	spentOutputsLen := C.size_t(len(spentOutputs))
	if len(spentOutputs) > 0 {
		for i, output := range spentOutputs {
			if err := validateReady(output); err != nil {
				return fmt.Errorf("invalid transaction output at index %d: %w", i, err)
			}
		}
		cPtrs := make([]*C.kernel_TransactionOutput, len(spentOutputs))
		for i, output := range spentOutputs {
			cPtrs[i] = output.ptr
		}
		cSpentOutputs = &cPtrs[0]
	}

	var cStatus C.kernel_ScriptVerifyStatus
	success := C.kernel_verify_script(
		scriptPubkey.ptr,
		C.int64_t(amount),
		txTo.ptr,
		cSpentOutputs,
		spentOutputsLen,
		C.uint(inputIndex),
		flags.c(),
		&cStatus,
	)
	if !success {
		status := scriptVerifyStatusFromC(cStatus)
		return status.err()
	}
	return nil
}

// ScriptFlags represents script verification flags
type ScriptFlags uint

const (
	ScriptFlagsVerifyNone                = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_NONE)
	ScriptFlagsVerifyP2SH                = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_P2SH)
	ScriptFlagsVerifyDERSig              = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_DERSIG)
	ScriptFlagsVerifyNullDummy           = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_NULLDUMMY)
	ScriptFlagsVerifyCheckLockTimeVerify = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_CHECKLOCKTIMEVERIFY)
	ScriptFlagsVerifyCheckSequenceVerify = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_CHECKSEQUENCEVERIFY)
	ScriptFlagsVerifyWitness             = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_WITNESS)
	ScriptFlagsVerifyTaproot             = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_TAPROOT)
	ScriptFlagsVerifyAll                 = ScriptFlags(C.kernel_SCRIPT_FLAGS_VERIFY_ALL)
)

func (s ScriptFlags) c() C.uint {
	return C.uint(s)
}

// ScriptVerifyStatus represents the status of script verification
type ScriptVerifyStatus int

const (
	ScriptVerifyOK ScriptVerifyStatus = iota
	ScriptVerifyErrorTxInputIndex
	ScriptVerifyErrorInvalidFlags
	ScriptVerifyErrorInvalidFlagsCombination
	ScriptVerifyErrorSpentOutputsRequired
	ScriptVerifyErrorSpentOutputsMismatch
)

func (s ScriptVerifyStatus) err() error {
	switch s {
	case ScriptVerifyErrorTxInputIndex:
		return ErrScriptVerifyTxInputIndex
	case ScriptVerifyErrorInvalidFlags:
		return ErrScriptVerifyInvalidFlags
	case ScriptVerifyErrorInvalidFlagsCombination:
		return ErrScriptVerifyInvalidFlagsCombination
	case ScriptVerifyErrorSpentOutputsRequired:
		return ErrScriptVerifySpentOutputsRequired
	case ScriptVerifyErrorSpentOutputsMismatch:
		return ErrScriptVerifySpentOutputsMismatch
	default:
		return ErrScriptVerify
	}
}

func scriptVerifyStatusFromC(status C.kernel_ScriptVerifyStatus) ScriptVerifyStatus {
	s := ScriptVerifyStatus(status)
	if s < ScriptVerifyOK || s > ScriptVerifyErrorSpentOutputsMismatch {
		panic(ErrInvalidScriptVerifyStatus)
	}
	return s
}
