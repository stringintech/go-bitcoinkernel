package main

import (
	"encoding/hex"
	"encoding/json"

	"github.com/stringintech/go-bitcoinkernel/kernel"
)

// handleBlockCreate creates a block from raw hex data
func handleBlockCreate(registry *Registry, req Request) Response {
	var params struct {
		RawBlock string `json:"raw_block"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	if req.Ref == "" {
		return NewInvalidParamsResponse(req.ID, "ref field is required")
	}

	// Decode hex to bytes
	blockBytes, err := hex.DecodeString(params.RawBlock)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, "raw_block must be valid hex")
	}

	// Create block
	block, err := kernel.NewBlock(blockBytes)
	if err != nil {
		return NewEmptyErrorResponse(req.ID)
	}

	registry.Store(req.Ref, block)

	return NewSuccessResponseWithRef(req.ID, req.Ref)
}

// handleBlockTreeEntryGetBlockHash gets the block hash from a block tree entry
func handleBlockTreeEntryGetBlockHash(registry *Registry, req Request) Response {
	var params struct {
		BlockTreeEntry RefObject `json:"block_tree_entry"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Get block tree entry from registry
	entry, err := registry.GetBlockTreeEntry(params.BlockTreeEntry.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get block hash and convert to string (handles display order conversion)
	hashView := entry.Hash()
	hashString := hashView.String()

	// Return hash as string
	return NewSuccessResponse(req.ID, hashString)
}
