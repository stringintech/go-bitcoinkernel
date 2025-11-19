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

// handleBlockTreeEntryGetBlockHash gets the block hash from a block tree entry
func handleBlockTreeEntryGetBlockHash(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockTreeEntry RefObject `json:"block_tree_entry"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	entry, err := registry.GetBlockTreeEntry(params.BlockTreeEntry.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(entry.Hash().String()), nil
}
