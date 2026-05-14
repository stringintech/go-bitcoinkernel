package main

import (
	"encoding/json"
	"fmt"
)

// handleTransactionInputGetOutPoint retrieves the out point view from a transaction input and stores it in the registry.
func handleTransactionInputGetOutPoint(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionInput RefObject `json:"transaction_input"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	ti, err := registry.GetTransactionInput(params.TransactionInput.Ref)
	if err != nil {
		return Response{}, err
	}

	registry.Store(req.Ref, ti.GetOutPoint())

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionInputCopy copies a transaction input and stores it in the registry
func handleTransactionInputCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionInput RefObject `json:"transaction_input"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	ti, err := registry.GetTransactionInput(params.TransactionInput.Ref)
	if err != nil {
		return Response{}, err
	}

	tiCopy := ti.Copy()
	registry.Store(req.Ref, tiCopy)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionOutPointGetTxid returns the txid view of an out point.
// The request must include ref so the handler can return the txid as a view ref.
func handleTransactionOutPointGetTxid(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionOutPoint RefObject `json:"transaction_out_point"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	op, err := registry.GetTransactionOutPoint(params.TransactionOutPoint.Ref)
	if err != nil {
		return Response{}, err
	}

	txidView := op.GetTxid()
	registry.Store(req.Ref, txidView)
	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionOutPointGetIndex returns the output index of an out point
func handleTransactionOutPointGetIndex(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionOutPoint RefObject `json:"transaction_out_point"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	op, err := registry.GetTransactionOutPoint(params.TransactionOutPoint.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(op.GetIndex()), nil
}

// handleTransactionOutPointCopy copies an out point and stores it in the registry
func handleTransactionOutPointCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionOutPoint RefObject `json:"transaction_out_point"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	op, err := registry.GetTransactionOutPoint(params.TransactionOutPoint.Ref)
	if err != nil {
		return Response{}, err
	}

	opCopy := op.Copy()
	registry.Store(req.Ref, opCopy)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTransactionOutPointDestroy destroys an out point from the registry
func handleTransactionOutPointDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		TransactionOutPoint RefObject `json:"transaction_out_point"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.TransactionOutPoint.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}
