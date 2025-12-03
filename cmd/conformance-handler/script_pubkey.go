package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleScriptPubkeyVerify verifies a script against a transaction
func handleScriptPubkeyVerify(req Request) Response {
	var params struct {
		ScriptPubkeyHex string          `json:"script_pubkey"`
		Amount          int64           `json:"amount"`
		TxToHex         string          `json:"tx_to"`
		InputIndex      uint            `json:"input_index"`
		Flags           json.RawMessage `json:"flags"`
		SpentOutputs    []SpentOutput   `json:"spent_outputs"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "")
	}

	// Decode script pubkey
	var scriptBytes []byte
	var err error
	if params.ScriptPubkeyHex != "" {
		scriptBytes, err = hex.DecodeString(params.ScriptPubkeyHex)
		if err != nil {
			return NewInvalidParamsResponse(req.ID, "script pubkey hex")
		}
	}

	// Decode transaction
	txBytes, err := hex.DecodeString(params.TxToHex)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, "transaction hex")
	}

	// Parse flags
	flags, err := parseScriptFlags(params.Flags)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, "flags")
	}

	// Parse spent outputs
	spentOutputs, err := parseSpentOutputs(params.SpentOutputs)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, "spent outputs")
	}
	defer func() {
		for _, so := range spentOutputs {
			so.Destroy()
		}
	}()

	// Create script pubkey and transaction
	scriptPubkey := kernel.NewScriptPubkey(scriptBytes)
	defer scriptPubkey.Destroy()

	tx, err := kernel.NewTransaction(txBytes)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, "transaction decode")
	}
	defer tx.Destroy()

	// Verify script
	valid, err := scriptPubkey.Verify(params.Amount, tx, spentOutputs, params.InputIndex, flags)
	if err != nil {
		var scriptVerifyError *kernel.ScriptVerifyError
		if errors.As(err, &scriptVerifyError) {
			switch {
			case errors.Is(err, kernel.ErrVerifyScriptVerifyInvalidFlagsCombination):
				return NewErrorResponse(req.ID, "btck_ScriptVerifyStatus", "ERROR_INVALID_FLAGS_COMBINATION")
			case errors.Is(err, kernel.ErrVerifyScriptVerifySpentOutputsRequired):
				return NewErrorResponse(req.ID, "btck_ScriptVerifyStatus", "ERROR_SPENT_OUTPUTS_REQUIRED")
			default:
				panic("scriptPubkey.Verify returned unhandled ScriptVerifyError (request ID: " + req.ID + "): " + err.Error())
			}
		}
		panic("scriptPubkey.Verify returned non-ScriptVerifyError (request ID: " + req.ID + "): " + err.Error())
	}

	return NewSuccessResponse(req.ID, valid)
}

// parseScriptFlags parses flags from array or numeric format
func parseScriptFlags(flagsJSON json.RawMessage) (kernel.ScriptFlags, error) {
	// Try array format first
	var flagsArray []string
	if err := json.Unmarshal(flagsJSON, &flagsArray); err == nil {
		var result kernel.ScriptFlags
		for _, flagStr := range flagsArray {
			flag, err := parseSingleFlag(flagStr)
			if err != nil {
				return 0, err
			}
			result |= flag
		}
		return result, nil
	}

	// Numeric flags
	var numFlags uint32
	if err := json.Unmarshal(flagsJSON, &numFlags); err != nil {
		return 0, errors.New("invalid flags format: must be array or number")
	}
	return kernel.ScriptFlags(numFlags), nil
}

// parseSingleFlag maps a flag string to its kernel constant
func parseSingleFlag(flagStr string) (kernel.ScriptFlags, error) {
	switch flagStr {
	case "btck_ScriptVerificationFlags_NONE":
		return kernel.ScriptFlagsVerifyNone, nil
	case "btck_ScriptVerificationFlags_P2SH":
		return kernel.ScriptFlagsVerifyP2SH, nil
	case "btck_ScriptVerificationFlags_DERSIG":
		return kernel.ScriptFlagsVerifyDERSig, nil
	case "btck_ScriptVerificationFlags_NULLDUMMY":
		return kernel.ScriptFlagsVerifyNullDummy, nil
	case "btck_ScriptVerificationFlags_CHECKLOCKTIMEVERIFY":
		return kernel.ScriptFlagsVerifyCheckLockTimeVerify, nil
	case "btck_ScriptVerificationFlags_CHECKSEQUENCEVERIFY":
		return kernel.ScriptFlagsVerifyCheckSequenceVerify, nil
	case "btck_ScriptVerificationFlags_WITNESS":
		return kernel.ScriptFlagsVerifyWitness, nil
	case "btck_ScriptVerificationFlags_TAPROOT":
		return kernel.ScriptFlagsVerifyTaproot, nil
	case "btck_ScriptVerificationFlags_ALL":
		return kernel.ScriptFlagsVerifyAll, nil
	default:
		return 0, errors.New("unknown flag: " + flagStr)
	}
}

type SpentOutput struct {
	ScriptPubkeyHex string `json:"script_pubkey"`
	Amount          int64  `json:"amount"`
}

// parseSpentOutputs parses spent outputs
func parseSpentOutputs(spentOutputParams []SpentOutput) (spentOutputs []*kernel.TransactionOutput, err error) {
	defer func() {
		// Clean up already created outputs on error
		if err != nil {
			for _, so := range spentOutputs {
				if so != nil {
					so.Destroy()
				}
			}
		}
	}()
	for _, so := range spentOutputParams {
		var scriptBytes []byte
		scriptBytes, err = hex.DecodeString(so.ScriptPubkeyHex)
		if err != nil {
			return
		}
		scriptPubkeyOut := kernel.NewScriptPubkey(scriptBytes)
		spentOutputs = append(spentOutputs, kernel.NewTransactionOutput(scriptPubkeyOut, so.Amount))
	}
	return spentOutputs, nil
}
