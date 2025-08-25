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

func (s ScriptFlags) c() C.uint {
	return C.uint(s)
}

const (
	ScriptVerifyOK                           = C.btck_ScriptVerifyStatus_SCRIPT_VERIFY_OK
	ScriptVerifyErrorInvalidFlagsCombination = C.btck_ScriptVerifyStatus_ERROR_INVALID_FLAGS_COMBINATION
	ScriptVerifyErrorSpentOutputsRequired    = C.btck_ScriptVerifyStatus_ERROR_SPENT_OUTPUTS_REQUIRED
)

type ScriptVerifyStatus C.btck_ScriptVerifyStatus

func (s ScriptVerifyStatus) err() error {
	switch s {
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
