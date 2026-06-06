package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleBlockCreate creates a block from raw hex data
func handleBlockCreate(registry *Registry, req Request) (Response, error) {
	var params struct {
		RawBlock string `json:"raw_block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	blockBytes, err := hex.DecodeString(params.RawBlock)
	if err != nil {
		return Response{}, fmt.Errorf("raw_block must be valid hex: %w", err)
	}

	block, err := kernel.NewBlock(blockBytes)
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, block)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockGetHash gets the hash of a block and stores it in the registry.
func handleBlockGetHash(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	hash := block.Hash()
	registry.Store(req.Ref, hash)
	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockGetHeader extracts the block header and stores it in the registry
func handleBlockGetHeader(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	header := block.GetHeader()
	registry.Store(req.Ref, header)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockCountTransactions returns the number of transactions in the block
func handleBlockCountTransactions(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(block.CountTransactions()), nil
}

// handleBlockGetTransactionAt retrieves the transaction at the given index and stores it in the registry
func handleBlockGetTransactionAt(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block            RefObject `json:"block"`
		TransactionIndex uint64    `json:"transaction_index"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	txView, err := block.GetTransactionAt(params.TransactionIndex)
	if err != nil {
		return Response{}, err
	}

	registry.Store(req.Ref, txView)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockToBytes returns the consensus-serialized block as a hex string
func handleBlockToBytes(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	data, err := block.Bytes()
	if err != nil {
		return NewEmptyErrorResponse(), nil
	}

	return NewSuccessResponse(hex.EncodeToString(data)), nil
}

// handleBlockCopy copies a block and stores the copy in the registry
func handleBlockCopy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	block, err := registry.GetBlock(params.Block.Ref)
	if err != nil {
		return Response{}, err
	}

	blockCopy := block.Copy()
	registry.Store(req.Ref, blockCopy)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleBlockDestroy destroys a block from the registry
func handleBlockDestroy(registry *Registry, req Request) (Response, error) {
	var params struct {
		Block RefObject `json:"block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if err := registry.Destroy(params.Block.Ref); err != nil {
		return Response{}, err
	}

	return NewEmptySuccessResponse(), nil
}
