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
		ScriptPubkeyHex string          `json:"script_pubkey_hex"`
		Amount          int64           `json:"amount"`
		TxHex           string          `json:"tx_hex"`
		InputIndex      uint            `json:"input_index"`
		Flags           json.RawMessage `json:"flags"`
		SpentOutputs    []SpentOutput   `json:"spent_outputs"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, "InvalidParams", err.Error())
	}

	// Decode script pubkey
	var scriptBytes []byte
	var err error
	if params.ScriptPubkeyHex != "" {
		scriptBytes, err = hex.DecodeString(params.ScriptPubkeyHex)
		if err != nil {
			return NewErrorResponse(req.ID, "InvalidParams", "invalid script pubkey hex")
		}
	}

	// Decode transaction
	txBytes, err := hex.DecodeString(params.TxHex)
	if err != nil {
		return NewErrorResponse(req.ID, "InvalidParams", "invalid transaction hex")
	}

	// Parse flags
	flags, err := parseScriptFlags(params.Flags)
	if err != nil {
		return NewErrorResponse(req.ID, "InvalidParams", err.Error())
	}

	// Parse spent outputs
	spentOutputs, err := parseSpentOutputs(params.SpentOutputs)
	if err != nil {
		return NewErrorResponse(req.ID, "InvalidParams", err.Error())
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
		return NewErrorResponse(req.ID, "TransactionDecode", err.Error())
	}
	defer tx.Destroy()

	// Verify script
	err = scriptPubkey.Verify(params.Amount, tx, spentOutputs, params.InputIndex, flags)
	if err != nil {
		return newScriptVerifyErrorResponse(req.ID, err)
	}

	return NewEmptySuccessResponse(req.ID)
}

// parseScriptFlags parses flags from string or numeric format
func parseScriptFlags(flagsJSON json.RawMessage) (kernel.ScriptFlags, error) {
	var flagsStr string
	if err := json.Unmarshal(flagsJSON, &flagsStr); err == nil {
		// String flags
		switch flagsStr {
		case "VERIFY_ALL_PRE_TAPROOT":
			return kernel.ScriptFlags(kernel.ScriptFlagsVerifyAll &^ kernel.ScriptFlagsVerifyTaproot), nil
		case "VERIFY_ALL":
			return kernel.ScriptFlagsVerifyAll, nil
		case "VERIFY_NONE":
			return kernel.ScriptFlagsVerifyNone, nil
		case "VERIFY_WITNESS":
			return kernel.ScriptFlagsVerifyWitness, nil
		case "VERIFY_TAPROOT":
			return kernel.ScriptFlagsVerifyTaproot, nil
		default:
			return 0, errors.New("unknown flags")
		}
	}

	// Numeric flags
	var numFlags uint32
	if err := json.Unmarshal(flagsJSON, &numFlags); err != nil {
		return 0, errors.New("invalid flags format")
	}
	return kernel.ScriptFlags(numFlags), nil
}

type SpentOutput struct {
	ScriptPubkeyHex string `json:"script_pubkey_hex"`
	Value           int64  `json:"value"`
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
		spentOutputs = append(spentOutputs, kernel.NewTransactionOutput(scriptPubkeyOut, so.Value))
	}
	return spentOutputs, nil
}

func newScriptVerifyErrorResponse(requestID string, err error) Response {
	var scriptVerifyError *kernel.ScriptVerifyError
	if errors.As(err, &scriptVerifyError) {
		switch {
		case errors.Is(err, kernel.ErrVerifyScriptVerifyTxInputIndex):
			return NewErrorResponse(requestID, "ScriptVerify", "TxInputIndex")
		case errors.Is(err, kernel.ErrVerifyScriptVerifyInvalidFlags):
			return NewErrorResponse(requestID, "ScriptVerify", "InvalidFlags")
		case errors.Is(err, kernel.ErrVerifyScriptVerifyInvalidFlagsCombination):
			return NewErrorResponse(requestID, "ScriptVerify", "InvalidFlagsCombination")
		case errors.Is(err, kernel.ErrVerifyScriptVerifySpentOutputsMismatch):
			return NewErrorResponse(requestID, "ScriptVerify", "SpentOutputsMismatch")
		case errors.Is(err, kernel.ErrVerifyScriptVerifySpentOutputsRequired):
			return NewErrorResponse(requestID, "ScriptVerify", "SpentOutputsRequired")
		case errors.Is(err, kernel.ErrVerifyScriptVerifyInvalid):
			return NewErrorResponse(requestID, "ScriptVerify", "Invalid")
		default:
			return NewErrorResponse(requestID, "ScriptVerify", "Invalid")
		}
	}
	return NewErrorResponse(requestID, "KernelError", err.Error())
}
