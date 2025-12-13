package main

import (
	"encoding/json"
)

// handleChainGetHeight gets the current height of the chain
func handleChainGetHeight(registry *Registry, req Request) Response {
	var params struct {
		Chain RefObject `json:"chain"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Get chain from registry
	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get height
	height := chain.GetHeight()

	// Return height as integer
	return NewSuccessResponse(req.ID, height)
}

// handleChainGetByHeight gets a block tree entry at the specified height
func handleChainGetByHeight(registry *Registry, req Request) Response {
	var params struct {
		Chain       RefObject `json:"chain"`
		BlockHeight int32     `json:"block_height"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	if req.Ref == "" {
		return NewInvalidParamsResponse(req.ID, "ref field is required")
	}

	// Get chain from registry
	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get block tree entry at height
	entry := chain.GetByHeight(params.BlockHeight)
	if entry == nil {
		return NewEmptyErrorResponse(req.ID)
	}

	registry.Store(req.Ref, entry)

	return NewSuccessResponseWithRef(req.ID, req.Ref)
}

// handleChainContains checks if a block tree entry is in the active chain
func handleChainContains(registry *Registry, req Request) Response {
	var params struct {
		Chain          RefObject `json:"chain"`
		BlockTreeEntry RefObject `json:"block_tree_entry"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewInvalidParamsResponse(req.ID, "failed to parse params")
	}

	// Get chain from registry
	chain, err := registry.GetChain(params.Chain.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Get block tree entry from registry
	entry, err := registry.GetBlockTreeEntry(params.BlockTreeEntry.Ref)
	if err != nil {
		return NewInvalidParamsResponse(req.ID, err.Error())
	}

	// Check if chain contains the entry
	contains := chain.Contains(entry)

	// Return boolean result
	return NewSuccessResponse(req.ID, contains)
}
