package main

import (
	"encoding/json"
	"fmt"
)

// handleChainGetHeight gets the current height of the chain
func handleChainGetHeight(registry *Registry, req Request) (Response, error) {
	var params struct {
		Chain RefObject `json:"chain"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(chain.GetHeight()), nil
}

// handleChainGetByHeight gets a block tree entry at the specified height
func handleChainGetByHeight(registry *Registry, req Request) (Response, error) {
	var params struct {
		Chain       RefObject `json:"chain"`
		BlockHeight int32     `json:"block_height"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	if req.Ref == "" {
		return Response{}, fmt.Errorf("ref field is required")
	}

	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return Response{}, err
	}

	entry := chain.GetByHeight(params.BlockHeight)
	if entry == nil {
		return NewEmptyErrorResponse(), nil
	}

	registry.Store(req.Ref, entry)

	return NewSuccessResponseWithRef(req.Ref), nil
}

// handleChainContains checks if a block tree entry is in the active chain
func handleChainContains(registry *Registry, req Request) (Response, error) {
	var params struct {
		Chain          RefObject `json:"chain"`
		BlockTreeEntry RefObject `json:"block_tree_entry"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return Response{}, fmt.Errorf("failed to parse params: %w", err)
	}

	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return Response{}, err
	}

	entry, err := registry.GetBlockTreeEntry(params.BlockTreeEntry.Ref)
	if err != nil {
		return Response{}, err
	}

	return NewSuccessResponse(chain.Contains(entry)), nil
}
