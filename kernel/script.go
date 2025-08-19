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

	var cSpentOutputs **C.btck_TransactionOutput
	spentOutputsLen := C.size_t(len(spentOutputs))
	if len(spentOutputs) > 0 {
		for i, output := range spentOutputs {
			if err := validateReady(output); err != nil {
				return fmt.Errorf("invalid transaction output at index %d: %w", i, err)
			}
		}
		cPtrs := make([]*C.btck_TransactionOutput, len(spentOutputs))
		for i, output := range spentOutputs {
			cPtrs[i] = output.ptr
		}
		cSpentOutputs = &cPtrs[0]
	}

	var cStatus C.btck_ScriptVerifyStatus
	success := C.btck_script_pubkey_verify(
		scriptPubkey.ptr,
		C.int64_t(amount),
		txTo.ptr,
		cSpentOutputs,
		spentOutputsLen,
		C.uint(inputIndex),
		flags.c(),
		&cStatus,
	)
	if success == 0 {
		status := scriptVerifyStatusFromC(cStatus)
		return status.err()
	}
	return nil
}

// ScriptFlags represents script verification flags
type ScriptFlags uint

const (
	ScriptFlagsVerifyNone                = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_NONE)
	ScriptFlagsVerifyP2SH                = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_P2SH)
	ScriptFlagsVerifyDERSig              = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_DERSIG)
	ScriptFlagsVerifyNullDummy           = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_NULLDUMMY)
	ScriptFlagsVerifyCheckLockTimeVerify = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_CHECKLOCKTIMEVERIFY)
	ScriptFlagsVerifyCheckSequenceVerify = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_CHECKSEQUENCEVERIFY)
	ScriptFlagsVerifyWitness             = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_WITNESS)
	ScriptFlagsVerifyTaproot             = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_TAPROOT)
	ScriptFlagsVerifyAll                 = ScriptFlags(C.btck_SCRIPT_FLAGS_VERIFY_ALL)
)

func (s ScriptFlags) c() C.uint {
	return C.uint(s)
}

// ScriptVerifyStatus represents the status of script verification
type ScriptVerifyStatus int

const (
	ScriptVerifyOK ScriptVerifyStatus = iota
	ScriptVerifyErrorInvalidFlags
	ScriptVerifyErrorInvalidFlagsCombination
	ScriptVerifyErrorSpentOutputsRequired
)

func (s ScriptVerifyStatus) err() error {
	switch s {
	case ScriptVerifyErrorInvalidFlags:
		return ErrKernelScriptVerifyInvalidFlags
	case ScriptVerifyErrorInvalidFlagsCombination:
		return ErrKernelScriptVerifyInvalidFlagsCombination
	case ScriptVerifyErrorSpentOutputsRequired:
		return ErrKernelScriptVerifySpentOutputsRequired
	default:
		return ErrKernelScriptVerify
	}
}

func scriptVerifyStatusFromC(status C.btck_ScriptVerifyStatus) ScriptVerifyStatus {
	s := ScriptVerifyStatus(status)
	if s < ScriptVerifyOK || s > ScriptVerifyErrorSpentOutputsRequired {
		panic(ErrInvalidScriptVerifyStatus)
	}
	return s
}
