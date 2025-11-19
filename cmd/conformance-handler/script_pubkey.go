package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleScriptPubkeyCreate creates a ScriptPubkey from hex and stores it in the registry
func handleScriptPubkeyCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		ScriptPubkeyHex string `json:"script_pubkey"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	scriptBytes, err := hex.DecodeString(params.ScriptPubkeyHex)
	if err != nil {
		return Response{}, fmt.Errorf("script_pubkey must be valid hex: %w", err)
	}

	spk := kernel.NewScriptPubkey(scriptBytes)
	registry.Store(req.Ref, spk)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleScriptPubkeyDestroy destroys a ScriptPubkey from the registry
func handleScriptPubkeyDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		ScriptPubkey RefObject `json:"script_pubkey"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.ScriptPubkey.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

// handleScriptPubkeyVerify verifies a script against a transaction
func handleScriptPubkeyVerify(registry *Registry, req Request) (Response, error) {
	var params struct {
		ScriptPubkey     RefObject       `json:"script_pubkey"`
		Amount           int64           `json:"amount"`
		TxTo             RefObject       `json:"tx_to"`
		InputIndex       uint            `json:"input_index"`
		Flags            json.RawMessage `json:"flags"`
		PrecomputedTxDat *RefObject      `json:"precomputed_txdata"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	scriptPubkey, err := registry.GetScriptPubkey(params.ScriptPubkey.Ref)
	if err != nil {
		return Response{}, err
	}

	tx, err := registry.GetTransaction(params.TxTo.Ref)
	if err != nil {
		return Response{}, err
	}

	flags, err := parseScriptFlags(params.Flags)
	if err != nil {
		return Response{}, fmt.Errorf("invalid flags: %w", err)
	}

	var precomputedTxData *kernel.PrecomputedTransactionData
	if params.PrecomputedTxDat != nil && params.PrecomputedTxDat.Ref != "" {
		precomputedTxData, err = registry.GetPrecomputedTransactionData(params.PrecomputedTxDat.Ref)
		if err != nil {
			return Response{}, err
		}
	}

	valid, err := scriptPubkey.Verify(params.Amount, tx, precomputedTxData, params.InputIndex, flags)
	if err != nil {
		var scriptVerifyError *kernel.ScriptVerifyError
		if errors.As(err, &scriptVerifyError) {
			switch {
			case errors.Is(err, kernel.ErrVerifyScriptVerifyInvalidFlagsCombination):
				return NewErrorResponse("btck_ScriptVerifyStatus", "ERROR_INVALID_FLAGS_COMBINATION"), nil
			case errors.Is(err, kernel.ErrVerifyScriptVerifySpentOutputsRequired):
				return NewErrorResponse("btck_ScriptVerifyStatus", "ERROR_SPENT_OUTPUTS_REQUIRED"), nil
			case errors.Is(err, kernel.ErrVerifyScriptVerifyTxInputIndex), errors.Is(err, kernel.ErrVerifyScriptVerifyInvalidFlags):
				return NewEmptyErrorResponse(), nil
			default:
				panic("scriptPubkey.Verify returned unhandled ScriptVerifyError (request ID: " + req.ID + "): " + err.Error())
			}
		}
		panic("scriptPubkey.Verify returned non-ScriptVerifyError (request ID: " + req.ID + "): " + err.Error())
	}

	return NewSuccessResponse(valid), nil
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
