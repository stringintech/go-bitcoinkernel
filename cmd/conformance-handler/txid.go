package main

import (
	"encoding/json"
	"fmt"
)

// handleTxidToBytes returns the raw stored 32-byte txid as a hex string.
func handleTxidToBytes(registry *Registry, req Request) (Response, error) {
	var params struct {
		Txid RefObject `json:"txid"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	txid, err := registry.GetTxid(params.Txid.Ref)
	if err != nil {
		return Response{}, err
	}

	raw := txid.Bytes()
	return NewSuccessResponse(fmt.Sprintf("%x", raw)), nil
}

// handleTxidEquals checks if two txids are equal
func handleTxidEquals(registry *Registry, req Request) (Response, error) {
	var params struct {
		Txid1 RefObject `json:"txid1"`
		Txid2 RefObject `json:"txid2"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	txid, err := registry.GetTxid(params.Txid1.Ref)
	if err != nil {
		return Response{}, err
	}

	txid2, err := registry.GetTxid(params.Txid2.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(txid.Bytes() == txid2.Bytes()), nil
}

// handleTxidCopy copies a txid and stores the copy in the registry
func handleTxidCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Txid RefObject `json:"txid"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	txid, err := registry.GetTxid(params.Txid.Ref)
	if err != nil {
		return Response{}, err
	}

	txidCopy := txid.Copy()
	registry.Store(req.Ref, txidCopy)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleTxidDestroy destroys a txid from the registry
func handleTxidDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Txid RefObject `json:"txid"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.Txid.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}
