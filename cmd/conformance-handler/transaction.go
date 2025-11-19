package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleTransactionCreate creates a Transaction from raw hex and stores it in the registry
func handleTransactionCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		RawTransaction string `json:"raw_transaction"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	txBytes, err := hex.DecodeString(params.RawTransaction)
	if err != nil {
		return Response{}, fmt.Errorf("raw_transaction must be valid hex: %w", err)
	}

	tx, err := kernel.NewTransaction(txBytes)
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, tx)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionDestroy destroys a Transaction from the registry
func handleTransactionDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Transaction RefObject `json:"transaction"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.Transaction.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

// handleTransactionOutputCreate creates a TransactionOutput from a ScriptPubkey ref and amount
func handleTransactionOutputCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		ScriptPubkey RefObject `json:"script_pubkey"`
		Amount       int64     `json:"amount"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	spk, err := registry.GetScriptPubkey(params.ScriptPubkey.Ref)
	if err != nil {
		return Response{}, err
	}

	txOut := kernel.NewTransactionOutput(spk, params.Amount)
	registry.Store(req.Ref, txOut)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionOutputDestroy destroys a TransactionOutput from the registry
func handleTransactionOutputDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionOutput RefObject `json:"transaction_output"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.TransactionOutput.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}

// handlePrecomputedTransactionDataCreate creates PrecomputedTransactionData from a tx and spent outputs
func handlePrecomputedTransactionDataCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		TxTo         RefObject   `json:"tx_to"`
		SpentOutputs []RefObject `json:"spent_outputs"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	tx, err := registry.GetTransaction(params.TxTo.Ref)
	if err != nil {
		return Response{}, err
	}

	spentOutputs := make([]kernel.TransactionOutputLike, len(params.SpentOutputs))
	for i, soRef := range params.SpentOutputs {
		so, err := registry.GetTransactionOutput(soRef.Ref)
		if err != nil {
			return Response{}, err
		}
		spentOutputs[i] = so
	}

	ptd, err := kernel.NewPrecomputedTransactionData(tx, spentOutputs)
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, ptd)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handlePrecomputedTransactionDataDestroy destroys PrecomputedTransactionData from the registry
func handlePrecomputedTransactionDataDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		PrecomputedTxData RefObject `json:"precomputed_txdata"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.PrecomputedTxData.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}
