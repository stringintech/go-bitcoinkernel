package main

import (
	"encoding/json"
	"fmt"
)

// handleBlockTreeEntryGetHeight gets the height of a block tree entry.
func handleBlockTreeEntryGetHeight(registry *Registry, req Request) (Response, error) {
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

	return NewSuccessResponse(entry.Height()), nil
}

// handleBlockTreeEntryGetBlockHash gets the block hash from a block tree entry.
func handleBlockTreeEntryGetBlockHash(registry *Registry, req Request) (Response, error) {
	var params struct {
		BlockTreeEntry RefObject `json:"block_tree_entry"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	entry, err := registry.GetBlockTreeEntry(params.BlockTreeEntry.Ref)
	if err != nil {
		return Response{}, err
	}

	hash := entry.Hash()
	registry.Store(req.Ref, hash)
	return NewSuccessResponseWithRef(req.Ref), nil
}
